package cmd

import (
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/persistenceOne/pstake-native/orchestrator/config"
	"github.com/persistenceOne/pstake-native/orchestrator/constants"
	"github.com/persistenceOne/pstake-native/orchestrator/orchestrator"
	"github.com/spf13/cobra"
)

func StartCommand() *cobra.Command {
	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start the orc server",
		Long:  `Start the orc server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			homepath, err := cmd.Flags().GetString(constants.FlagOrcHomeDir)
			if err != nil {
				stdlog.Println(err)
				stdlog.Fatalln(err)
			}

			orcConfig := InitConfig(homepath)

			orcSeeds := orcConfig.OrcSeeds
			valAddr := orcConfig.ValAddress

			cosmosChain, err := InitCosmosChain(homepath, orcConfig.CosmosConfig)
			if err != nil {
				panic(any(err))
			}

			nativeChain, err := InitNativeChain(homepath, orcConfig.NativeConfig)
			if err != nil {
				panic(any(err))
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

			nativeProtoCodec := codec.NewProtoCodec(clientContextNative.InterfaceRegistry)

			stdlog.Println("start rpc server")

			stdlog.Println("start to listen for txs cosmos side")
			//TODO : use goroutines  implementation properly here
			go orchestrator.StartListeningCosmosEvent(valAddr, orcSeeds, clientContextNative, clientContextCosmos, cosmosChain, nativeChain, cosmosProtoCodec)
			stdlog.Println("started listening for deposits")
			go orchestrator.StartListeningCosmosDeposit(valAddr, orcSeeds, clientContextNative, clientContextCosmos, cosmosChain, nativeChain, cosmosProtoCodec)

			stdlog.Println("start to listen for txs native side")
			go orchestrator.StartListeningNativeSideActions(valAddr, orcSeeds, clientContextNative, clientContextCosmos, cosmosChain, nativeChain, nativeProtoCodec)

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			for sig := range signalChan {
				_ = fmt.Sprintf("Stopping the orchestrator %v", sig.String())

			}
			return nil
		},
	}
	startCommand.Flags().String(constants.FlagOrcHomeDir, "", "home directory")
	return startCommand
}

func InitConfig(homepath string) config.Config {
	cfg := config.InitializeConfigFromToml(homepath)
	return cfg
}
