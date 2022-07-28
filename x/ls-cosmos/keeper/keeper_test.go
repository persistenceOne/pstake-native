package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	app        *app.PstakeApp
	ctx        sdk.Context
	govHandler govtypes.Handler
}

func (suite *IntegrationTestSuite) SetupTest() {
	_, app, ctx := helpers.CreateTestApp()

	keeper := app.LSCosmosKeeper

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	suite.app = &app
	suite.ctx = ctx
}

func testProposal(
	title,
	description,
	connection,
	channel,
	transfer,
	uatom,
	ustkatom string) *types.RegisterCosmosChainProposal {
	return types.NewRegisterCosmosChainProposal(
		title,
		description,
		connection,
		channel,
		transfer,
		uatom,
		ustkatom,
	)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
