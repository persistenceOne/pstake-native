package orc

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	orc "github.com/persistenceOne/pStake-native/oracle/command"
	"github.com/persistenceOne/pStake-native/oracle/constants"
	"github.com/persistenceOne/pStake-native/oracle/oracle"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

//func StartCommand() *cobra.Command {
//	startCommand := &cobra.Command{
//		Use:   "start",
//		Short: "Start the orc server",
//		Long:  `Start the orc server`,
//		Run: func(cmd *cobra.Command, args []string) {
//			cmd.Help()
//		},
//	}
//}

func StartCommand() *cobra.Command {
	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start the orc server",
		Long:  `Start the orc server`,
		Run: func(cmd *cobra.Command, args []string) {
			homepath, err := cmd.Flags().GetString(constants.FlagOrcHomeDir)
			if err != nil {
				fmt.Println(err)
				return
			}

			cosmosChain, err := orc.InitCosmosChain(homepath)
			if err != nil {
				fmt.Println(err)
			}

			nativeChain, err := orc.InitNativeChain(homepath)
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

			//_ := codec.NewProtoCodec(clientContextNative.InterfaceRegistry)

			fmt.Println("start rpc server")

			fmt.Println("start to listen for txs cosmos side")
			go oracle.StartGettingDepositTx(clientContextNative, clientContextCosmos, cosmosChain, nativeChain, cosmosProtoCodec)
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			for sig := range signalChan {
				fmt.Sprintf("Stopping the oracle %v", sig.String())

			}

		},
	}
	startCommand.Flags().String(constants.FlagOrcHomeDir, "", "home directory")
	return startCommand
}
