package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/keeper"
)

func (suite *IntegrationTestSuite) TestCValueLimits() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)
	str, broken := keeper.CValueLimits(k)(ctx)
	suite.False(broken)
	suite.Equal("liquidstakeibc: cvalue-limits invariant\ncvalue out of bounds: false, values as follows \n  \n", str)

	hc.CValue = sdk.MustNewDecFromStr("2")
	k.SetHostChain(ctx, hc)
	str, broken = keeper.CValueLimits(k)(ctx)
	suite.True(broken)
	suite.Equal("liquidstakeibc: cvalue-limits invariant\ncvalue out of bounds: true, values as follows \n chainID: testchain2-1, cValue: 2.000000000000000000 \n \n", str)
}
