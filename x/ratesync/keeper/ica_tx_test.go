package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestExecuteLiquidStakeRateTx() {
	k, ctx := suite.app.RatesyncKeeper, suite.ctx
	_ = createNChain(k, ctx, 2)
	hc, _ := k.GetHostChain(ctx, 1)
	suite.Require().NoError(k.ExecuteLiquidStakeRateTx(ctx, hc.Features.LiquidStakeIBC,
		"stk/uatom", "uatom", sdk.OneDec(), hc.ID, suite.ratesyncPathAB.EndpointA.ConnectionID, hc.ICAAccount))
	suite.Require().NoError(k.InstantiateLiquidStakeContract(ctx, hc.ICAAccount,
		hc.Features.LiquidStake, hc.ID, suite.ratesyncPathAB.EndpointA.ConnectionID, hc.TransferChannelID, hc.TransferPortID))
}
