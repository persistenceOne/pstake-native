package keeper_test

import liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"

func (suite *IntegrationTestSuite) TestBeginBlocker() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	_ = createNChain(keeper, ctx, 1)
	suite.Require().NotPanics(func() {
		keeper.BeginBlock(ctx)
	})
	hcs := keeper.GetAllHostChain(ctx)
	hc := hcs[0]
	hc.Features.LiquidStake.Enabled = true
	keeper.SetHostChain(ctx, hc)
	suite.Require().NotPanics(func() {
		keeper.BeginBlock(ctx)
	})
}

func (suite *IntegrationTestSuite) TestDoRecreateICA() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	hc := ValidHostChainInMsg(1)
	keeper.SetHostChain(ctx, hc)
	suite.Require().NotPanics(func() {
		keeper.DoRecreateICA(ctx, hc)
	})
	hc.ICAAccount.ChannelState = liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATED
	suite.Require().NotPanics(func() {
		keeper.DoRecreateICA(ctx, hc)
	})

}
