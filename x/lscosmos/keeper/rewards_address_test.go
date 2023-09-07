package keeper_test

import "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"

func (suite *IntegrationTestSuite) TestHostChainRewardsAddress() {
	app, ctx := suite.app, suite.ctx

	lscosmosKeeper := app.LSCosmosKeeper

	err := lscosmosKeeper.SetHostChainRewardAddressIfEmpty(ctx, types.HostChainRewardAddress{Address: "address________________"})
	suite.NoError(err)

	err = lscosmosKeeper.SetHostChainRewardAddressIfEmpty(ctx, types.HostChainRewardAddress{Address: "address________________"})
	suite.Error(err)

	hostChainRewardAddress := lscosmosKeeper.GetHostChainRewardAddress(ctx)
	suite.Equal(types.HostChainRewardAddress{Address: "address________________"}, hostChainRewardAddress)
}
