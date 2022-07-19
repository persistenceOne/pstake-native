package keeper_test

import (
	"github.com/persistenceOne/pstake-native/x/cosmos/keeper"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	app *app.PstakeApp
	ctx sdk.Context
}

func (suite *IntegrationTestSuite) SetupTest() {
	_, app, ctx := helpers.CreateTestApp()

	params := types.DefaultParams()

	// Set some params required for testing
	params.BondDenoms = []string{"stake"}
	params.CosmosProposalParams.ChainID = "gaiad"
	params.CustodialAddress = "cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2"
	params.ModuleEnabled = true

	keeper.InitGenesis(ctx, app.CosmosKeeper, &types.GenesisState{Params: params})

	suite.app = &app
	suite.ctx = ctx
}

func (suite *IntegrationTestSuite) SetupValWeightedAmounts(ws types.WeightedAddressAmounts) {
	suite.app.CosmosKeeper.SetCosmosValidatorSet(suite.ctx, ws)
	for _, w := range ws {
		valAddr, _ := types.ValAddressFromBech32(w.Address, types.Bech32PrefixValAddr)
		suite.app.CosmosKeeper.UpdateDelegationCosmosValidator(suite.ctx, valAddr, w.Coin(), sdk.NewCoin("test", sdk.NewInt(0)))
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
