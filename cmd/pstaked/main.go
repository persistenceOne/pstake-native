package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/cmd/pstaked/cmd"
)

func main() {
	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	configuration.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	configuration.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	configuration.SetCoinType(app.CoinType)
	configuration.SetFullFundraiserPath(app.FullFundraiserPath)
	configuration.Seal()

	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
