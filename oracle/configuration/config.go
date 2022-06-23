package configuration

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/pstake-native/oracle/constants"
	"github.com/spf13/cobra"
	"path/filepath"
)

type Config struct {
	ValAddress   string       `json:"val_address"`
	OrcSeeds     []string     `json:"seeds"`
	CosmosConfig CosmosConfig `json:"cosmos_config"`
	NativeConfig NativeConfig `json:"native_config"`
}

var orcConfig = newConfig()

func GetConfig() Config {
	return orcConfig
}

func SetConfig(cmd *cobra.Command) Config {
	//  val flag

	ValAddress, err := cmd.Flags().GetString(constants.FlagValAddress)
	orcConfig.ValAddress = ValAddress
	//	cosmos Config
	orcConfig.CosmosConfig.ChainID, err = cmd.Flags().GetString(constants.FlagCosmosChainID)
	if err != nil {
		panic(err)
	}
	orcConfig.CosmosConfig.CustodialAddr, err = cmd.Flags().GetString(constants.FlagCosmosCustodialAddr)
	if err != nil {
		panic(err)
	}
	orcConfig.CosmosConfig.Denom, err = cmd.Flags().GetString(constants.FlagCosmosDenom)
	if err != nil {
		panic(err)
	}
	orcConfig.CosmosConfig.GRPCAddr, err = cmd.Flags().GetString(constants.FlagCosmosGRPCAddr)
	if err != nil {
		panic(err)
	}
	orcConfig.CosmosConfig.RPCAddr, err = cmd.Flags().GetString(constants.FlagCosmosRPCAddr)
	if err != nil {
		panic(err)
	}
	orcConfig.CosmosConfig.AccountPrefix, err = cmd.Flags().GetString(constants.FlagCosmosAccountPrefix)
	if err != nil {
		panic(err)
	}
	orcConfig.CosmosConfig.GasAdjustment, err = cmd.Flags().GetFloat64(constants.FlagCosmosGasAdjustment)
	if err != nil {
		panic(err)
	}
	orcConfig.CosmosConfig.GasPrice, err = cmd.Flags().GetString(constants.FlagCosmosGasPrice)
	if err != nil {
		panic(err)
	}
	orcConfig.CosmosConfig.CoinType, err = cmd.Flags().GetUint32(constants.FlagCosmosCoinType)

	//	native Config

	orcConfig.NativeConfig.ChainID, err = cmd.Flags().GetString(constants.FlagNativeChainID)
	if err != nil {
		panic(err)
	}
	orcConfig.NativeConfig.Denom, err = cmd.Flags().GetString(constants.FlagNativeDenom)
	if err != nil {
		panic(err)
	}
	orcConfig.NativeConfig.GRPCAddr, err = cmd.Flags().GetString(constants.FlagNativeGRPCAddr)
	if err != nil {
		panic(err)
	}
	orcConfig.NativeConfig.RPCAddr, err = cmd.Flags().GetString(constants.FlagNativeRPCAddr)
	if err != nil {
		panic(err)
	}
	orcConfig.NativeConfig.AccountPrefix, err = cmd.Flags().GetString(constants.FlagNativeAccountPrefix)
	if err != nil {
		panic(err)
	}
	orcConfig.NativeConfig.GasAdjustment, err = cmd.Flags().GetFloat64(constants.FlagNativeGasAdjustment)
	if err != nil {
		panic(err)
	}
	orcConfig.NativeConfig.GasPrices, err = cmd.Flags().GetString(constants.FlagNativeGasPrice)
	if err != nil {
		panic(err)
	}
	orcConfig.NativeConfig.CoinType, err = cmd.Flags().GetUint32(constants.FlagNativeCoinType)

	return orcConfig

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
	GRPCAddr      string  `json:"grpc_addr"`
	RPCAddr       string  `json:"rpc_addr"`
	AccountPrefix string  `json:"account_prefix"`
	GasAdjustment float64 `json:"gas_adjustment"`
	GasPrice      string  `json:"gas_price"`
	CoinType      uint32  `json:"coin_type"`
}

type NativeConfig struct {
	ChainID       string  `json:"chain_id"`
	Denom         string  `json:"denom"`
	RPCAddr       string  `json:"rpc_addr"`
	GRPCAddr      string  `json:"grpc_addr"`
	AccountPrefix string  `json:"account_prefix"`
	GasAdjustment float64 `json:"gas_adjustment"`
	GasPrices     string  `json:"gas_price"`
	CoinType      uint32  `json:"coin_type"`
}

func DefaultCosmosConfig() *CosmosConfig {
	return &CosmosConfig{
		ChainID:       constants.CosmosChainID,
		CustodialAddr: constants.CosmosCustodialAddr,
		Denom:         constants.CosmosDenom,
		GRPCAddr:      constants.CosmosGRPCAddr,
		RPCAddr:       constants.CosmosRPCAddr,
		AccountPrefix: constants.CosmosAccountPrefix,
		GasAdjustment: constants.CosmosGasAdjustment,
		GasPrice:      constants.CosmosGasPrice,
	}
}

func DefaultNativeConfig() *NativeConfig {
	return &NativeConfig{
		ChainID:       constants.NativeChainID,
		Denom:         constants.NativeDenom,
		RPCAddr:       constants.NativeRPCAddr,
		AccountPrefix: constants.NativeAccountPrefix,
		GasAdjustment: constants.NativeGasAdjustment,
		GasPrices:     constants.NativeGasPrice,
	}
}

func NewCosmosConfig() CosmosConfig {
	return CosmosConfig{
		ChainID:       constants.CosmosChainID,
		CustodialAddr: constants.CosmosCustodialAddr,
		Denom:         constants.CosmosDenom,
		RPCAddr:       constants.CosmosRPCAddr,
		AccountPrefix: constants.CosmosAccountPrefix,
		GasAdjustment: constants.CosmosGasAdjustment,
		GasPrice:      constants.CosmosGasPrice,
		CoinType:      constants.CosmosCoinType,
	}
}

func NewNativeConfig() NativeConfig {
	return NativeConfig{
		ChainID:       constants.NativeChainID,
		Denom:         constants.NativeDenom,
		RPCAddr:       constants.NativeRPCAddr,
		GRPCAddr:      constants.NativeGRPCAddr,
		AccountPrefix: constants.NativeAccountPrefix,
		GasAdjustment: constants.NativeGasAdjustment,
		GasPrices:     constants.NativeGasPrice,
		CoinType:      constants.NativeCoinType,
	}
}

func newConfig() Config {
	return Config{
		ValAddress:   constants.ValAddress,
		OrcSeeds:     constants.Seed,
		CosmosConfig: NewCosmosConfig(),
		NativeConfig: NewNativeConfig(),
	}
}

func InitializeConfigFromToml(homepath string) Config {
	var config = newConfig()
	_, _ = toml.DecodeFile(filepath.Join(homepath, "config.toml"), &config)
	//log.Fatalf("Error Decoding oracle config: %v\n", err.Error())
	fmt.Println(config)
	return config
}
