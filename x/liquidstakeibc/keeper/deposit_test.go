package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) TestGetSetDeposit() {
	pstakeApp, ctx := suite.app, suite.ctx

	pstakeApp.LiquidStakeIBCKeeper.SetDeposit(ctx, &types.Deposit{ChainId: "hc1"})
	deposits := pstakeApp.LiquidStakeIBCKeeper.GetAllDeposits(ctx)

	suite.Require().Equal(len(deposits), 1)
	suite.Require().Equal(deposits[0].ChainId, "hc1")
}

func (suite *IntegrationTestSuite) TestDeleteDeposit() {
	pstakeApp, ctx := suite.app, suite.ctx

	deposit := &types.Deposit{ChainId: "hc1"}

	pstakeApp.LiquidStakeIBCKeeper.SetDeposit(ctx, deposit)
	pstakeApp.LiquidStakeIBCKeeper.DeleteDeposit(ctx, deposit)
	deposits := pstakeApp.LiquidStakeIBCKeeper.GetAllDeposits(ctx)

	suite.Require().Equal(len(deposits), 0)
}

func (suite *IntegrationTestSuite) TestCreateDeposits() {
	pstakeApp, ctx := suite.app, suite.ctx

	pstakeApp.LiquidStakeIBCKeeper.CreateDeposits(ctx, 10)

	deposits := pstakeApp.LiquidStakeIBCKeeper.GetAllDeposits(ctx)

	suite.Require().Equal(len(deposits), 1)
}

func (suite *IntegrationTestSuite) TestRevertDepositState() {
	pstakeApp, ctx := suite.app, suite.ctx

	deposits := []*types.Deposit{
		{
			ChainId:       "chain-1",
			Amount:        sdk.Coin{},
			Epoch:         sdk.NewInt(1),
			State:         types.Deposit_DEPOSIT_PENDING,
			IbcSequenceId: "",
		},
		{
			ChainId:       "chain-1",
			Amount:        sdk.Coin{},
			Epoch:         sdk.NewInt(2),
			State:         types.Deposit_DEPOSIT_SENT,
			IbcSequenceId: "",
		},
		{
			ChainId:       "chain-1",
			Amount:        sdk.Coin{},
			Epoch:         sdk.NewInt(3),
			State:         types.Deposit_DEPOSIT_RECEIVED,
			IbcSequenceId: "",
		},
		{
			ChainId:       "chain-1",
			Amount:        sdk.Coin{},
			Epoch:         sdk.NewInt(4),
			State:         types.Deposit_DEPOSIT_DELEGATING,
			IbcSequenceId: "",
		},
	}

	pstakeApp.LiquidStakeIBCKeeper.RevertDepositsState(ctx, deposits)
	revertedDeposits := pstakeApp.LiquidStakeIBCKeeper.GetAllDeposits(ctx)

	suite.Require().Equal(len(revertedDeposits), 4)

	for _, deposit := range revertedDeposits {
		switch deposit.Epoch.Int64() {
		case 1:
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_PENDING)
		case 2:
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_PENDING)
		case 3:
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_SENT)
		case 4:
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_RECEIVED)
		}
	}
}

func (suite *IntegrationTestSuite) TestGetDepositSequenceID() {
	pstakeApp := suite.app

	sequenceID := pstakeApp.LiquidStakeIBCKeeper.GetDepositSequenceID("channel-0", 1)

	suite.Require().Equal(sequenceID, "channel-0-sequence-1")
}
