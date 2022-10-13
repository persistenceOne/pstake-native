package keeper_test

import (
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestAllowListedValidators() {
	app, ctx := suite.app, suite.ctx

	// set empty allow listed validators
	app.LSCosmosKeeper.SetAllowListedValidators(ctx, types.AllowListedValidators{})

	resAllowListedValidators := app.LSCosmosKeeper.GetAllowListedValidators(ctx)
	suite.Nil(resAllowListedValidators.AllowListedValidators)

	// set the filled allow listed validators
	app.LSCosmosKeeper.SetAllowListedValidators(ctx, allowListedValidators)

	resAllowListedValidators = app.LSCosmosKeeper.GetAllowListedValidators(ctx)
	suite.Equal(allowListedValidators, resAllowListedValidators)
}
