package orchestrator

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

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
)

const (
	NEXT_BLOCK_WAIT_TIME = 1500
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
	CoinType         uint32         `json:"coin_type" yaml:"coin_type"`

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
	stdlog.Println("Here1")
	custodialAddress, err := AccAddressFromBech32(custAddr, c.AccountPrefix)
	if err != nil {
		stdlog.Println("Here1.5")
		return err
	}
	stdlog.Println("Here2")
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

	stdlog.Println(rpcClient, "cosmos rpcClient")

	if c.logger == nil {
		c.logger = defaultChainLogger()
	}

	return nil

}

func defaultChainLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

func (c *CosmosChain) UseSDKContext() func() {

	sdkConf := sdk.GetConfig()
	sdkConf.SetBech32PrefixForAccount(c.AccountPrefix, c.AccountPrefix+"pub")
	sdkConf.SetBech32PrefixForValidator(c.AccountPrefix+"valoper", c.AccountPrefix+"valoperpub")
	sdkConf.SetBech32PrefixForConsensusNode(c.AccountPrefix+"valcons", c.AccountPrefix+"valconspub")

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

func StartListeningCosmosEvent(valAddr string, orcSeeds []string, nativeCliCtx client.Context, clientCtx client.Context, chain *CosmosChain, native *NativeChain, codec *codec.ProtoCodec) {
	ctx := context.Background()

	var cHeight uint64
	stdlog.Println(chain, "printing chain")
	abciInfoCosmos, err := chain.Client.ABCIInfo(ctx)
	if err != nil {
		stdlog.Println("error getting abci info", err)
		stdlog.Println("error getting cosmos abci info", err)
	}
	cHeight = uint64(abciInfoCosmos.Response.LastBlockHeight)

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

	EventListForBlock, err := rpcClient.Subscribe(ctx, "orchestrator-client", query)
	if err != nil {
		stdlog.Println("error subscribing, to the Event", err)
	}

	for e := range EventListForBlock {
		stdlog.Println("listening to events on Cosmos side")
		slashSlice := e.Events["slash.address"]

		if slashSlice != nil {
			for _, slashAddr := range slashSlice {
				err := chain.SlashingHandler(slashAddr, orcSeeds, valAddr, nativeCliCtx, native, chain, int64(cHeight))
				if err != nil {
					stdlog.Println("proposal handling error")
					panic(err)
				}
			}
		}

		propIDSlice := e.Events["active_proposal.proposal_id"]
		if propIDSlice != nil {
			for _, propId := range propIDSlice {
				err := chain.ProposalHandler(propId, orcSeeds, nativeCliCtx, native, chain, int64(cHeight))
				if err != nil {
					stdlog.Println("proposal handling error")
					panic(err)
				}
			}
		}

		//TODO : Add undelegate success handler

	}
}
func StartListeningCosmosDeposit(valAddr string, orcSeeds []string, nativeCliCtx client.Context, clientCtx client.Context, chain *CosmosChain, native *NativeChain, codec *codec.ProtoCodec) {
	ctx := context.Background()
	var cHeight, nHeight uint64

	if _, err := os.Stat(filepath.Join(chain.HomePath, "status.json")); err == nil {
		cHeight, nHeight = GetHeightStatus(chain.HomePath)
	} else if errors.Is(err, os.ErrNotExist) {
		abciInfoCosmos, err := chain.Client.ABCIInfo(ctx)
		if err != nil {
			stdlog.Println("error getting abci info", err)
			stdlog.Println("error getting cosmos abci info", err)
		}
		cHeight = uint64(abciInfoCosmos.Response.LastBlockHeight)
		stdlog.Println("cosmos Block height- ", cHeight)

		abciInfoNative, err := native.Client.ABCIInfo(ctx)
		if err != nil {
			stdlog.Println("error getting abci info", err)
			stdlog.Println("error getting native abci info", err)
		}
		nHeight = uint64(abciInfoNative.Response.LastBlockHeight)
		stdlog.Println("native Block height- ", nHeight)

		SetStatus(chain.HomePath, cHeight, nHeight)

	}
	for cHeight > 0 && nHeight > 0 {
		stdlog.Println("cosmos Block height ", cHeight)
		err := chain.DepositHandler(valAddr, orcSeeds, nativeCliCtx, clientCtx, native, int64(cHeight), codec)
		if err != nil {
			stdlog.Fatalln()
		}

		_, err = chain.Client.ABCIInfo(ctx)
		if err != nil {
			stdlog.Println("error getting cosmos abci info", err)

		}

		time.Sleep(NEXT_BLOCK_WAIT_TIME * time.Millisecond)
		cHeight++

		SetStatus(chain.HomePath, cHeight, nHeight)

	}
}
