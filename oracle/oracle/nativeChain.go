package oracle

//func InitNativeChain(timeout, homePath string) (*oracle.Chain, error) {
//	chain := &oracle.Chain{}
//	chain.Key = "native"
//	chain.ChainID = configuration.GetConfig().NativeConfig.ChainID
//	chain.RPCAddr = configuration.GetConfig().NativeConfig.RPCAddr
//	chain.AccountPrefix = configuration.GetConfig().NativeConfig.AccountPrefix
//	chain.GasAdjustment = configuration.GetConfig().NativeConfig.GasAdjustment
//	chain.GasPrices = configuration.GetConfig().NativeConfig.GasPrice
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
//	_, err = helpers.KeyAddOrRestore(chain, chain.Key, constants.NativeCoinType)
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
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	keys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pStake-native/oracle/utils"
	"github.com/tendermint/tendermint/libs/log"
	provtypes "github.com/tendermint/tendermint/light/provider"
	prov "github.com/tendermint/tendermint/light/provider/http"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	logg "log"
	"os"
	"path/filepath"
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

	if c.logger == nil {
		c.logger = defaultChainLogger()
	}

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
	var cHeight, nHeight uint64

	if _, err := os.Stat(filepath.Join(chain.HomePath, "status.json")); err == nil {
		cHeight, nHeight = utils.GetHeightStatus(chain.HomePath)
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
		cHeight = uint64(abciInfoNative.Response.LastBlockHeight)

		utils.NewStatusJSON(chain.HomePath, cHeight, nHeight)

	}
	for cHeight > 0 && nHeight > 0 {
		fmt.Println("cosmos Block height- ", cHeight)
		fmt.Println("native Block Height", nHeight)

		err := native.StakingHandler(valAddr, orcSeeds, nativeCliCtx, ClientCtx, native, int64(nHeight), codec)
		if err != nil {
			logg.Fatalln()
		}
		nHeight += 1
		utils.NewStatusJSON(chain.HomePath, cHeight, nHeight)

	}

}
