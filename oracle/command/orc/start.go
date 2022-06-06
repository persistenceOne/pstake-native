package orc

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	orc "github.com/persistenceOne/pstake-native/oracle/command"
	"github.com/persistenceOne/pstake-native/oracle/configuration"
	"github.com/persistenceOne/pstake-native/oracle/constants"
	"github.com/persistenceOne/pstake-native/oracle/oracle"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func StartCommand() *cobra.Command {
	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start the orc server",
		Long:  `Start the orc server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			homepath, err := cmd.Flags().GetString(constants.FlagOrcHomeDir)
			if err != nil {
				fmt.Println(err)
				log.Fatalln(err)
			}

			orcConfig := InitConfig(homepath)

			orcSeeds := orcConfig.OrcSeeds
			valAddr := orcConfig.ValAddress

			cosmosChain, err := orc.InitCosmosChain(homepath, orcConfig.CosmosConfig)
			if err != nil {
				fmt.Println(err)
			}

			nativeChain, err := orc.InitNativeChain(homepath, orcConfig.NativeConfig)
			if err != nil {
				fmt.Println(err)
			}

			cosmosEncodingConfig := cosmosChain.MakeEncodingConfig()
			nativeEncodingConfig := nativeChain.MakeEncodingConfig()
			clientContextCosmos := client.Context{}.
				WithCodec(cosmosEncodingConfig.Marshaler).
				WithInterfaceRegistry(cosmosEncodingConfig.InterfaceRegistry).
				WithTxConfig(cosmosEncodingConfig.TxConfig).
				WithLegacyAmino(cosmosEncodingConfig.Amino).
				WithInput(os.Stdin).
				WithAccountRetriever(authTypes.AccountRetriever{}).
				WithHomeDir(homepath).
				WithNodeURI(cosmosChain.RPCAddr).
				WithClient(cosmosChain.Client).
				WithViper("")

			cosmosProtoCodec := codec.NewProtoCodec(clientContextCosmos.InterfaceRegistry)

			clientContextNative := client.Context{}.
				WithCodec(nativeEncodingConfig.Marshaler).
				WithInterfaceRegistry(nativeEncodingConfig.InterfaceRegistry).
				WithTxConfig(nativeEncodingConfig.TxConfig).
				WithLegacyAmino(nativeEncodingConfig.Amino).
				WithInput(os.Stdin).
				WithAccountRetriever(authTypes.AccountRetriever{}).
				WithNodeURI(nativeChain.RPCAddr).
				WithClient(nativeChain.Client).
				WithHomeDir(homepath).
				WithViper("")

			//nativeProtoCodec := codec.NewProtoCodec(clientContextNative.InterfaceRegistry)

			fmt.Println("start rpc server")

			fmt.Println("start to listen for txs cosmos side")
			go oracle.StartListeningCosmosSideActions(valAddr, orcSeeds, clientContextNative, clientContextCosmos, cosmosChain, nativeChain, cosmosProtoCodec)

			fmt.Println("start to listen for txs cosmos side")
			go oracle.StartListeningNativeSideActions(valAddr, orcSeeds, clientContextNative, clientContextCosmos, cosmosChain, nativeChain, cosmosProtoCodec)

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			for sig := range signalChan {
				_ = fmt.Sprintf("Stopping the oracle %v", sig.String())

			}
			return nil
		},
	}
	startCommand.Flags().String(constants.FlagOrcHomeDir, "", "home directory")
	return startCommand
}

func InitConfig(homepath string) configuration.Config {
	config := configuration.InitializeConfigFromToml(homepath)

	return config

}
