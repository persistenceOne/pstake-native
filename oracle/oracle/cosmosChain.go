package oracle

import (
	"context"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	keys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	provtypes "github.com/tendermint/tendermint/light/provider"
	prov "github.com/tendermint/tendermint/light/provider/http"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	logg "log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type CosmosChain struct {
	Key              string         `json:"key" yaml:"key"`
	ChainID          string         `json:"chain_id" yaml:"chain_id"`
	GRPCAddr         string         `json:"grpc_addr" yaml:"grpc_addr"`
	RPCAddr          string         `json:"rpc_addr" yaml:"rpc_addr"`
	AccountPrefix    string         `json:"account_prefix" yaml:"account_prefix"`
	GasAdjustment    float64        `json:"gas_adjustment" yaml:"gas_adjustment"`
	GasPrices        string         `json:"gas_prices" yaml:"gas_prices"`
	CustodialAddress sdk.AccAddress `json:"custodial_address" yaml:"custodial_address"`

	HomePath string                `json:"home_path" yaml:"home_path"`
	KeyBase  keys.Keyring          `json:"-" yaml:""`
	Client   rpcclient.Client      `json:"-" yaml:""`
	Encoding params.EncodingConfig `json:"-" yaml:""`
	Provider provtypes.Provider    `json:"-" yaml:""`

	address sdk.AccAddress
	logger  log.Logger
	timeout time.Duration
	debug   bool
}

func (c *CosmosChain) Init(custAddr string, homepath string, timeout time.Duration, logger log.Logger, debug bool) error {
	keybase, err := keys.New(c.ChainID, "test", keysDir(homepath, c.ChainID), nil)
	if err != nil {
		return err
	}
	rpcClient, err := newRPCClient(c.RPCAddr, timeout)
	if err != nil {
		return err
	}
	liteprovider, err := prov.New(c.ChainID, c.RPCAddr)
	if err != nil {
		return err
	}
	_, err = sdk.ParseDecCoins(c.GasPrices)
	if err != nil {
		return fmt.Errorf("failed to parse gas prices (%s) for chain %s", c.GasPrices, c.ChainID)
	}
	custodialAddress, err := sdk.AccAddressFromBech32(custAddr)

	if err != nil {
		return err
	}

	encodingConfig := c.MakeEncodingConfig()

	c.KeyBase = keybase
	c.Client = rpcClient
	c.HomePath = homepath
	c.Encoding = encodingConfig
	c.logger = logger
	c.timeout = timeout
	c.debug = debug
	c.Provider = liteprovider
	c.CustodialAddress = custodialAddress

	if c.logger == nil {
		c.logger = defaultChainLogger()
	}

	return nil

}

func defaultChainLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

func (c *CosmosChain) UseSDKContext() func() {
	//fmt.Println("SDKCONTEXT LOCK")
	//sdkContextMutex.Lock()

	sdkConf := sdk.GetConfig()
	sdkConf.SetBech32PrefixForAccount(c.AccountPrefix, c.AccountPrefix+"pub")
	sdkConf.SetBech32PrefixForValidator(c.AccountPrefix+"valoper", c.AccountPrefix+"valoperpub")
	sdkConf.SetBech32PrefixForConsensusNode(c.AccountPrefix+"valcons", c.AccountPrefix+"valconspub")

	//return sdkContextMutex.Unlock
	return func() {}
}

func keysDir(home, chainID string) string {
	return path.Join(home, "keys", chainID)
}

func newRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}

	httpClient.Timeout = timeout
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	return rpcClient, nil
}
func (c *CosmosChain) KeyExists(name string) bool {
	k, err := c.KeyBase.Key(name)
	if err != nil {
		return false
	}

	return k.GetName() == name
}
func (c *CosmosChain) Start() error {
	return c.Client.Start()
}

var sdkContextMutex sync.Mutex

func StartListeningCosmosSideActions(valAddr string, orcSeeds []string, nativeCliCtx client.Context, ClientCtx client.Context, chain *CosmosChain, native *NativeChain, codec *codec.ProtoCodec) {
	ctx := context.Background()
	var cHeight, nHeight uint64

	rpcClient, err := rpchttp.New(chain.RPCAddr, "/websocket")
	if err != nil {
		_ = fmt.Errorf("RPC address invalid %v", err)
		return
	}
	err = rpcClient.Start()

	if err != nil {
		_ = fmt.Errorf("unable to reach RPC %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "tm.event = 'NewBlock'"

	if _, err := os.Stat(filepath.Join(chain.HomePath, "status.json")); err == nil {
		cHeight, nHeight = GetHeightStatus(chain.HomePath)
	} else if errors.Is(err, os.ErrNotExist) {
		abciInfoCosmos, err := chain.Client.ABCIInfo(ctx)
		if err != nil {
			fmt.Println("error getting abci info", err)
			logg.Println("error getting cosmos abci info", err)
		}
		cHeight = uint64(abciInfoCosmos.Response.LastBlockHeight)

		abciInfoNative, err := native.Client.ABCIInfo(ctx)
		if err != nil {
			fmt.Println("error getting abci info", err)
			logg.Println("error getting native abci info", err)
		}
		nHeight = uint64(abciInfoNative.Response.LastBlockHeight)

		NewStatusJSON(chain.HomePath, cHeight, nHeight)

	}
	for cHeight > 0 && nHeight > 0 {
		fmt.Println("cosmos Block height- ", cHeight)
		fmt.Println("native Block Height", nHeight)
		err := chain.DepositHandler(valAddr, orcSeeds, nativeCliCtx, ClientCtx, native, int64(nHeight), codec)
		if err != nil {
			logg.Fatalln()
		}

		EventListForBlock, err := rpcClient.Subscribe(ctx, "oracle-client", query)
		if err != nil {
			logg.Println("error subscribing, to the Event", err)
		}

		for e := range EventListForBlock {
			logg.Println("listening to events on Cosmos side")
			propIDSlice := e.Events["active_proposal.proposal_id"]
			if propIDSlice == nil {
				continue
			}
			for _, propId := range propIDSlice {
				err := chain.ProposalHandler(propId, orcSeeds, nativeCliCtx, native, chain, int64(cHeight))
				if err != nil {
					logg.Println("proposal handling error")
					return
				}
			}
		}

		cHeight += 1
		NewStatusJSON(chain.HomePath, cHeight, nHeight)

	}

}
