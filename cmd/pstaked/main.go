package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v4/app"
	"github.com/persistenceOne/pstake-native/v4/cmd/pstaked/cmd"
)

func main() {
	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	configuration.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	configuration.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	configuration.SetCoinType(app.CoinType)
	configuration.SetPurpose(app.Purpose)
	configuration.Seal()

	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
