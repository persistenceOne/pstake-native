package keeper_test

import (
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestGetParams() {
	app, ctx := suite.app, suite.ctx

	params := types.DefaultParams()
	suite.Equal(params, app.LSCosmosKeeper.GetParams(ctx))
}
