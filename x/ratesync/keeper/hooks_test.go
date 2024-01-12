package keeper_test

import sdk "github.com/cosmos/cosmos-sdk/types"

func (suite *IntegrationTestSuite) TestPostCValueUpdate() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	_ = createNChain(keeper, ctx, 10)
	suite.Require().NoError(keeper.PostCValueUpdate(ctx, "uatom", "stk/uatom", sdk.OneDec()))
	hc, _ := keeper.GetHostChain(ctx, 1)
	hc.Features.LiquidStakeIBC.Enabled = true
	hc.Features.LiquidStakeIBC.Denoms = []string{"*"}
	keeper.SetHostChain(ctx, hc)
	suite.Require().NoError(keeper.PostCValueUpdate(ctx, "uatom", "stk/uatom", sdk.OneDec()))

	hc.ICAAccount.Address = "InvalidAddr" // outer functions do not return errors
	keeper.SetHostChain(ctx, hc)
	suite.Require().NoError(keeper.PostCValueUpdate(ctx, "uatom", "stk/uatom", sdk.OneDec()))
}

func (suite *IntegrationTestSuite) TestAfterEpochEnd() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	_ = createNChain(keeper, ctx, 10)
	suite.Require().NoError(keeper.AfterEpochEnd(ctx, "hour", 1))
	hc, _ := keeper.GetHostChain(ctx, 1)
	hc.Features.LiquidStake.Enabled = true
	hc.Features.LiquidStake.Denoms = []string{"*"}
	keeper.SetHostChain(ctx, hc)
	suite.Require().NoError(keeper.AfterEpochEnd(ctx, "hour", 1))

	hc.ICAAccount.Address = "InvalidAddr" // outer functions do not return errors
	keeper.SetHostChain(ctx, hc)
	suite.Require().NoError(keeper.AfterEpochEnd(ctx, "hour", 1))
}
