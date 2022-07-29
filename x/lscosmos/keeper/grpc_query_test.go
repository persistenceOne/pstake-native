package keeper_test

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
