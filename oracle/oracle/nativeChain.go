package oracle

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
	logg "log"
	"time"
)

type NativeChain struct {
	Key           string  `json:"key" yaml:"key"`
	ChainID       string  `json:"chain_id" yaml:"chain_id"`
	RPCAddr       string  `json:"rpc_addr" yaml:"rpc_addr"`
	AccountPrefix string  `json:"account_prefix" yaml:"account_prefix"`
	GasAdjustment float64 `json:"gas_adjustment" yaml:"gas_adjustment"`
	GasPrices     string  `json:"gas_prices" yaml:"gas_prices"`
	GRPCAddr      string  `json:"grpc_addr" yaml:"grpc_addr"`
	CoinType      uint32  `json:"coin_type" yaml:"coin_type"`
	//CustodialAddress sdk.AccAddress `json:"custodial_address" yaml:"custodial_address"`

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

func (c *NativeChain) Init(homepath string, timeout time.Duration, logger log.Logger, debug bool) error {
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

	if c.logger == nil {
		c.logger = defaultChainLogger()
	}
	fmt.Println(c, "chain-")
	fmt.Println(homepath)
	return nil

}

func (c *NativeChain) UseSDKContext() func() {

	sdkContextMutex.Lock()

	sdkConf := sdk.GetConfig()
	sdkConf.SetBech32PrefixForAccount(c.AccountPrefix, c.AccountPrefix+"pub")
	sdkConf.SetBech32PrefixForValidator(c.AccountPrefix+"valoper", c.AccountPrefix+"valoperpub")
	sdkConf.SetBech32PrefixForConsensusNode(c.AccountPrefix+"valcons", c.AccountPrefix+"valconspub")

	return sdkContextMutex.Unlock
}

func (c *NativeChain) KeyExists(name string) bool {
	k, err := c.KeyBase.Key(name)
	if err != nil {
		return false
	}

	return k.GetName() == name
}
func (c *NativeChain) Start() error {
	return c.Client.Start()
}

func StartListeningNativeSideActions(valAddr string, orcSeeds []string, nativeCliCtx client.Context, ClientCtx client.Context, chain *CosmosChain, native *NativeChain, codec *codec.ProtoCodec) {
	ctx := context.Background()

	fmt.Println("Listening to native side events")

	var nHeight uint64

	abciInfoNative, err := native.Client.ABCIInfo(ctx)
	if err != nil {
		fmt.Println("error getting abci info", err)
		logg.Println("error getting cosmos abci info", err)
	}
	nHeight = uint64(abciInfoNative.Response.LastBlockHeight)

	fmt.Println(nHeight)

	rpcClient, err := rpchttp.New(native.RPCAddr, "/websocket")
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

	EventListForBlock, err := rpcClient.Subscribe(ctx, "oracle-client", query)
	if err != nil {
		logg.Println("error subscribing, to the Event", err)
	}

	for e := range EventListForBlock {
		logg.Println("listening to events on native side")
		txIdSlice := e.Events["outgoing_txn.outgoing_tx_id"]

		if txIdSlice == nil {
			continue
		}
		for _, txId := range txIdSlice {

			err := native.OutgoingTxHandler(txId, valAddr, orcSeeds, nativeCliCtx, ClientCtx, native, chain)
			if err != nil {

				panic(err)
				continue
			}
		}

		txSignIdSlice := e.Events["signed_tx.outgoing_tx_id"]
		if txIdSlice == nil {
			continue
		}

		for _, txID := range txSignIdSlice {
			err := native.SignedOutgoingTxHandler(txID, valAddr, orcSeeds, nativeCliCtx, ClientCtx, native, chain)
			if err != nil {
				logg.Println("signed outgoing tx handling error")
				panic(err)
				return

			}
		}

	}

}
