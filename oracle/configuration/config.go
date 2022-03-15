package configuration

import (
	"github.com/persistenceOne/pStake-native/oracle/constants"
	"github.com/spf13/cobra"
)

type Config struct {
	ValAddress   string       `json:"val_address"`
	CosmosConfig CosmosConfig `json:"cosmos_config"`
	NativeConfig NativeConfig `json:"native_config"`
}

var config = Default()

func GetConfig() Config {
	return config
}

func SetConfig(cmd *cobra.Command) Config {
	//  val flag
	valAddress, err := cmd.Flags().GetString(constants.FlagValAddress)
	//	cosmos Config
	cosmosChainID, err := cmd.Flags().GetString(constants.FlagCosmosChainID)
	if err != nil {
		panic(err)
	}
	cosmosCustodialAddr, err := cmd.Flags().GetString(constants.FlagCosmosCustodialAddr)
	if err != nil {
		panic(err)
	}
	cosmosDenom, err := cmd.Flags().GetString(constants.FlagCosmosDenom)
	if err != nil {
		panic(err)
	}
	cosmosRPCAddr, err := cmd.Flags().GetString(constants.FlagCosmosRPCAddr)
	if err != nil {
		panic(err)
	}
	cosmosAccountPrefix, err := cmd.Flags().GetString(constants.FlagCosmosAccountPrefix)
	if err != nil {
		panic(err)
	}
	cosmosGasAdjustment, err := cmd.Flags().GetFloat64(constants.FlagCosmosGasAdjustment)
	if err != nil {
		panic(err)
	}
	cosmosGasPrice, err := cmd.Flags().GetString(constants.FlagCosmosGasPrice)
	if err != nil {
		panic(err)
	}

	cosmosConfig := NewCosmosConfig(
		cosmosChainID,
		cosmosCustodialAddr,
		cosmosDenom,
		cosmosRPCAddr,
		cosmosAccountPrefix,
		cosmosGasAdjustment,
		cosmosGasPrice,
	)
	//	native Config

	nativeChainID, err := cmd.Flags().GetString(constants.FlagNativeChainID)
	if err != nil {
		panic(err)
	}
	nativeModuleName, err := cmd.Flags().GetString(constants.FlagNativeModuleName)
	if err != nil {
		panic(err)
	}
	nativeDenom, err := cmd.Flags().GetString(constants.FlagNativeDenom)
	if err != nil {
		panic(err)
	}
	nativeRPCAddr, err := cmd.Flags().GetString(constants.FlagNativeRPCAddr)
	if err != nil {
		panic(err)
	}
	nativeAccountPrefix, err := cmd.Flags().GetString(constants.FlagNativeAccountPrefix)
	if err != nil {
		panic(err)
	}
	nativeGasAdjustment, err := cmd.Flags().GetFloat64(constants.FlagNativeGasAdjustment)
	if err != nil {
		panic(err)
	}
	nativeGasPrice, err := cmd.Flags().GetString(constants.FlagNativeGasPrice)
	if err != nil {
		panic(err)
	}

	nativeConfig := NewNativeConfig(
		nativeChainID,
		nativeModuleName,
		nativeDenom,
		nativeRPCAddr,
		nativeAccountPrefix,
		nativeGasAdjustment,
		nativeGasPrice,
	)
	return NewConfig(
		valAddress,
		*cosmosConfig,
		*nativeConfig,
	)

}

func (c *Config) GetCosmosConfig() CosmosConfig {
	return c.CosmosConfig
}

func (c *Config) GetNativeConfig() NativeConfig {
	return c.NativeConfig
}
func NewConfig(val string, cosmosConfig CosmosConfig, nativeConfig NativeConfig) Config {
	return Config{
		ValAddress:   val,
		CosmosConfig: cosmosConfig,
		NativeConfig: nativeConfig,
	}
}

func Default() Config {
	return Config{
		CosmosConfig: *DefaultCosmosConfig(),
		NativeConfig: *DefaultNativeConfig(),
	}
}

type CosmosConfig struct {
	ChainID       string  `json:"chain_id"`
	CustodialAddr string  `json:"custodial_addr"`
	Denom         string  `json:"denom"`
	RPCAddr       string  `json:"rpc_addr"`
	AccountPrefix string  `json:"account_prefix"`
	GasAdjustment float64 `json:"gas_adjustment"`
	GasPrices     string  `json:"gas_price"`
}

type NativeConfig struct {
	ChainID       string  `json:"chain_id"`
	ModuleName    string  `json:"module_name"`
	Denom         string  `json:"denom"`
	RPCAddr       string  `json:"rpc_addr"`
	AccountPrefix string  `json:"account_prefix"`
	GasAdjustment float64 `json:"gas_adjustment"`
	GasPrices     string  `json:"gas_price"`
}

func DefaultCosmosConfig() *CosmosConfig {
	return &CosmosConfig{
		ChainID:       constants.CosmosChainID,
		CustodialAddr: constants.CosmosCustodialAddr,
		Denom:         constants.CosmosDenom,
		RPCAddr:       constants.CosmosRPCAddr,
		AccountPrefix: constants.CosmosAccountPrefix,
		GasAdjustment: constants.CosmosGasAdjustment,
		GasPrices:     constants.CosmosGasPrice,
	}
}

func DefaultNativeConfig() *NativeConfig {
	return &NativeConfig{
		ChainID:       constants.NativeChainID,
		ModuleName:    constants.ModuleName,
		Denom:         constants.NativeDenom,
		RPCAddr:       constants.NativeRPCAddr,
		AccountPrefix: constants.NativeAccountPrefix,
		GasAdjustment: constants.NativeGasAdjustment,
		GasPrices:     constants.NativeGasPrice,
	}
}

func NewCosmosConfig(chainID string, custodial_addr string, denom string, rpcaddr string, accountPrefix string, gasAdjustment float64, gasPrice string) *CosmosConfig {
	return &CosmosConfig{
		ChainID:       chainID,
		CustodialAddr: custodial_addr,
		Denom:         denom,
		RPCAddr:       rpcaddr,
		AccountPrefix: accountPrefix,
		GasAdjustment: gasAdjustment,
		GasPrices:     gasPrice,
	}
}

func NewNativeConfig(chainID string, moduleName string, denom string, rpcaddr string, accountPrefix string, gasAdjustment float64, gasPrice string) *NativeConfig {
	return &NativeConfig{
		ChainID:       chainID,
		ModuleName:    moduleName,
		Denom:         denom,
		RPCAddr:       rpcaddr,
		AccountPrefix: accountPrefix,
		GasAdjustment: gasAdjustment,
		GasPrices:     gasPrice,
	}
}
