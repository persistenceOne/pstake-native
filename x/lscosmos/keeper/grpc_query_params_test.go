package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestParamsQuery() {
	app, ctx := suite.app, suite.ctx

	c := sdk.WrapSDKContext(ctx)
	response, err := app.LSCosmosKeeper.Params(c, &types.QueryParamsRequest{})
	suite.NoError(err)
	suite.Equal(&types.QueryParamsResponse{Params: types.DefaultParams()}, response)

	_, err = app.LSCosmosKeeper.Params(c, nil)
	suite.Error(err)
}
