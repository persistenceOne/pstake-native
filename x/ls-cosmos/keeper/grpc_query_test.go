package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	ls_cosmos "github.com/persistenceOne/pstake-native/x/ls-cosmos"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

func (suite *IntegrationTestSuite) TestCosmosIBCParamsQuery() {
	app, ctx := suite.app, suite.ctx

	suite.govHandler = ls_cosmos.NewLSCosmosProposalHandler(suite.app.LSCosmosKeeper)
	propsal := testProposal("title", "description", "connection", "channel-1", "transfer", "uatom", "ustkatom")
	err := suite.govHandler(ctx, propsal)
	suite.NoError(err)

	c := sdk.WrapSDKContext(ctx)
	response, err := app.LSCosmosKeeper.CosmosIBCParams(c, &types.QueryCosmosIBCParamsRequest{})
	suite.NoError(err)
	suite.Equal(&types.QueryCosmosIBCParamsResponse{CosmosIBCParams: *propsal}, response)
}
