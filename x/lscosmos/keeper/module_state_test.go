package keeper_test

func (suite *IntegrationTestSuite) TestModuleEnable() {
	app, ctx := suite.app, suite.ctx

	lscosmosKeeper := app.LSCosmosKeeper

	lscosmosKeeper.SetModuleState(ctx, true)
	suite.True(lscosmosKeeper.GetModuleState(ctx))

	lscosmosKeeper.SetModuleState(ctx, false)
	suite.False(lscosmosKeeper.GetModuleState(ctx))
}
