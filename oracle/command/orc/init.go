package orc

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/pstake-native/oracle/configuration"
	"github.com/persistenceOne/pstake-native/oracle/constants"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func InitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialise oracle configuration in config.toml",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := configuration.SetConfig(cmd)
			var buf bytes.Buffer
			log.Println("init ")
			encoder := toml.NewEncoder(&buf)
			if err := encoder.Encode(config); err != nil {
				return err
			}
			homeDir, err := cmd.Flags().GetString(constants.FlagOrcHomeDir)
			if err != nil {
				return err
			}
			if err = os.MkdirAll(homeDir, os.ModePerm); err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(filepath.Join(homeDir, "config.toml"), buf.Bytes(), 0600); err != nil {
				panic(err)
			}
			log.Println("generated config.toml ", filepath.Join(homeDir, "config.toml"))

			return nil

		},
	}
	cmd.Flags().String(constants.FlagOrcHomeDir, "", "home directory")

	cmd.Flags().String(constants.FlagValAddress, "", "validator address")
	//	Cosmos Flag
	cmd.Flags().String(constants.FlagCosmosChainID, constants.CosmosChainID, "cosmos chain id")
	cmd.Flags().String(constants.FlagCosmosCustodialAddr, constants.CosmosCustodialAddr, "cosmos custodial address")
	cmd.Flags().String(constants.FlagCosmosDenom, constants.CosmosDenom, "cosmos denom")
	cmd.Flags().String(constants.FlagCosmosGRPCAddr, constants.CosmosGRPCAddr, "cosmos grpc address")
	cmd.Flags().String(constants.FlagCosmosRPCAddr, constants.CosmosRPCAddr, "cosmos rpc address")
	cmd.Flags().String(constants.FlagCosmosAccountPrefix, constants.CosmosAccountPrefix, "cosmos account prefix")
	cmd.Flags().Float64(constants.FlagCosmosGasAdjustment, constants.CosmosGasAdjustment, "cosmos fee")
	cmd.Flags().String(constants.FlagCosmosGasPrice, constants.CosmosGasPrice, "cosmos gas price")
	cmd.Flags().Uint32(constants.FlagCosmosCoinType, constants.CosmosCoinType, "cosmos coin type")

	// Native Flag
	cmd.Flags().String(constants.FlagNativeChainID, constants.NativeChainID, "")
	cmd.Flags().String(constants.FlagNativeDenom, constants.NativeDenom, "")
	cmd.Flags().String(constants.FlagNativeGRPCAddr, constants.NativeGRPCAddr, "")
	cmd.Flags().String(constants.FlagNativeRPCAddr, constants.NativeRPCAddr, "")
	cmd.Flags().String(constants.FlagNativeAccountPrefix, constants.NativeAccountPrefix, "")
	cmd.Flags().Float64(constants.FlagNativeGasAdjustment, constants.NativeGasAdjustment, "")
	cmd.Flags().String(constants.FlagNativeGasPrice, constants.NativeGasPrice, "")
	cmd.Flags().String(constants.FlagNativeModuleName, constants.ModuleName, "")
	cmd.Flags().Uint32(constants.FlagNativeCoinType, constants.NativeCoinType, "native coin type")

	return cmd
}
