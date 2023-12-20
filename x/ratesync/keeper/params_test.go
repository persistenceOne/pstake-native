package keeper_test

import (
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

func (suite *IntegrationTestSuite) TestGetParams() {
	k, ctx := suite.app.RatesyncKeeper, suite.ctx
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	suite.Require().EqualValues(params, k.GetParams(ctx))
}
