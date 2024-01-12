package keeper_test

import (
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) TestSetGetDeleteRedelegationTx() {
	suite.app.LiquidStakeIBCKeeper.SetRedelegationTx(
		suite.ctx,
		&types.RedelegateTx{
			ChainId:       suite.chainB.ChainID,
			IbcSequenceId: "channel-100-sequence-1",
			State:         0,
		},
	)

	redelegationTx, found := suite.app.LiquidStakeIBCKeeper.GetRedelegationTx(
		suite.ctx,
		suite.chainB.ChainID, "channel-100-sequence-1",
	)

	suite.Require().Equal(true, found)
	suite.Require().Equal(redelegationTx.ChainId, suite.chainB.ChainID)
	suite.Require().Equal(redelegationTx.State, types.RedelegateTx_REDELEGATE_SENT)
	suite.Require().Equal(redelegationTx.IbcSequenceId, "channel-100-sequence-1")

	txs := suite.app.LiquidStakeIBCKeeper.GetAllRedelegationTx(suite.ctx)
	suite.Require().Equal(len(txs), 1)

	filteredTxs := suite.app.LiquidStakeIBCKeeper.FilterRedelegationTx(suite.ctx, func(redel types.RedelegateTx) bool { return true })
	suite.Require().Equal(len(filteredTxs), 1)

	suite.app.LiquidStakeIBCKeeper.DeleteRedelegationTx(suite.ctx, suite.chainB.ChainID, "channel-100-sequence-1")
	txs2 := suite.app.LiquidStakeIBCKeeper.GetAllRedelegationTx(suite.ctx)
	suite.Require().Equal(len(txs2), 0)

	_, found = suite.app.LiquidStakeIBCKeeper.GetRedelegationTx(
		suite.ctx,
		suite.chainB.ChainID, "channel-100-sequence-1",
	)
	suite.Require().Equal(false, found)
}
