package oracle

//func InitCosmosChain(timeout, homePath string) (*oracle.Chain, error) {
//	chain := &oracle.Chain{}
//	chain.Key = "cosmos"
//	chain.ChainID = configuration.GetConfig().CosmosConfig.ChainID
//	chain.RPCAddr = configuration.GetConfig().CosmosConfig.RPCAddr
//	chain.AccountPrefix = configuration.GetConfig().CosmosConfig.AccountPrefix
//	chain.GasAdjustment = configuration.GetConfig().CosmosConfig.GasAdjustment
//	chain.GasPrices = configuration.GetConfig().CosmosConfig.GasPrice
//	chain.TrustingPeriod = "21h"
//
//	to, err := time.ParseDuration(timeout)
//	if err != nil {
//		return nil, err
//	}
//	err = chain.Init(homePath, to, nil, true)
//
//	if chain.KeyExists(chain.Key) {
//		println("Key Exists")
//		err = chain.Keybase.Delete(chain.Key)
//		if err != nil {
//			return chain, err
//		}
//	}
//
//	_, err = helpers.KeyAddOrRestore(chain, chain.Key, constants.CosmosCoinType)
//	if err != nil {
//		return chain, err
//	}
//
//	if err = chain.Start(); err != nil {
//		if err != tendermintService.ErrAlreadyStarted {
//			chain.Error(err)
//			return chain, err
//		}
//	}
//	return chain, nil
//
//}

import (
	"context"
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
	"os"
	"path"
	"sync"
	"time"
)

type CosmosChain struct {
	Key              string         `json:"key" yaml:"key"`
	ChainID          string         `json:"chain_id" yaml:"chain_id"`
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
	client, err := newRPCClient(c.RPCAddr, timeout)
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
	c.Client = client
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

func StartGettingDepositTx(nativeCliCtx client.Context, ClientCtx client.Context, chain *CosmosChain, native *NativeChain, codec *codec.ProtoCodec) {
	ctx := context.Background()
	for {
		abciInfoCosmos, err := chain.Client.ABCIInfo(ctx)
		if err != nil {
			fmt.Println("error getting abci info", err)
			time.Sleep(time.Second * 5)
			continue
		}

		block_height := abciInfoCosmos.Response.LastBlockHeight
		block_height -= 4
		fmt.Println("block height", block_height)

		err = chain.DepositHandler(nativeCliCtx, ClientCtx, native, block_height, codec)

	}
}
