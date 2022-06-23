package keeper_test

import (
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
	// Set weighted developer address to empty to not conflict with sdk addr prefix
	params.WeightedDeveloperRewardsReceivers = []types.WeightedAddress{}
	// Set DelegationThreshold to 10 unit
	app.CosmosKeeper.SetParams(ctx, params)

	suite.app = &app
	suite.ctx = ctx
}

func (suite *IntegrationTestSuite) SetupValWeightedAmounts(ws types.WeightedAddressAmounts) {
	suite.app.CosmosKeeper.SetCosmosValidatorSet(suite.ctx, ws)
	for _, w := range ws {
		valAddr, _ := types.ValAddressFromBech32(w.Address, types.Bech32PrefixValAddr)
		suite.app.CosmosKeeper.UpdateDelegationCosmosValidator(suite.ctx, valAddr, w.Coin())
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
