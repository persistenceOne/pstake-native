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

	sequenceID := pstakeApp.LiquidStakeIBCKeeper.GetTransactionSequenceID("channel-0", 1)

	suite.Require().Equal(sequenceID, "channel-0-sequence-1")
}

func (suite *IntegrationTestSuite) TestAdjustDepositsForRedemption() {
	tc := []struct {
		name             string
		deposits         []*types.Deposit
		expected         map[int64]sdk.Coin
		redemptionAmount sdk.Coin
		success          bool
	}{
		{
			name: "one deposit that can fill the request",
			deposits: []*types.Deposit{
				{
					ChainId: "chain-1",
					Epoch:   sdk.NewInt(1),
					Amount:  sdk.Coin{Denom: "stake", Amount: sdk.NewInt(10000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected: map[int64]sdk.Coin{
				1: {Denom: "stake", Amount: sdk.NewInt(5000)},
			},
			redemptionAmount: sdk.Coin{Denom: "stake", Amount: sdk.NewInt(5000)},
			success:          true,
		},
		{
			name: "one deposit that can't fill the request",
			deposits: []*types.Deposit{
				{
					ChainId: "chain-1",
					Epoch:   sdk.NewInt(1),
					Amount:  sdk.Coin{Denom: "stake", Amount: sdk.NewInt(3500)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected: map[int64]sdk.Coin{
				1: {Denom: "stake", Amount: sdk.NewInt(3500)},
			},
			redemptionAmount: sdk.Coin{Denom: "stake", Amount: sdk.NewInt(5000)},
			success:          true,
		},
		{
			name: "one deposit that can exactly fill the request",
			deposits: []*types.Deposit{
				{
					ChainId: "chain-1",
					Epoch:   sdk.NewInt(1),
					Amount:  sdk.Coin{Denom: "stake", Amount: sdk.NewInt(5000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected:         map[int64]sdk.Coin{},
			redemptionAmount: sdk.Coin{Denom: "stake", Amount: sdk.NewInt(5000)},
			success:          true,
		},
		{
			name: "redemption filled with first deposit",
			deposits: []*types.Deposit{
				{
					ChainId: "chain-1",
					Epoch:   sdk.NewInt(1),
					Amount:  sdk.Coin{Denom: "stake", Amount: sdk.NewInt(10000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
				{
					ChainId: "chain-1",
					Epoch:   sdk.NewInt(2),
					Amount:  sdk.Coin{Denom: "stake", Amount: sdk.NewInt(5000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected: map[int64]sdk.Coin{
				1: {Denom: "stake", Amount: sdk.NewInt(5000)},
				2: {Denom: "stake", Amount: sdk.NewInt(5000)},
			},
			redemptionAmount: sdk.Coin{Denom: "stake", Amount: sdk.NewInt(5000)},
			success:          true,
		},
		{
			name: "redemption filled with second deposit",
			deposits: []*types.Deposit{
				{
					ChainId: "chain-1",
					Epoch:   sdk.NewInt(1),
					Amount:  sdk.Coin{Denom: "stake", Amount: sdk.NewInt(5000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
				{
					ChainId: "chain-1",
					Epoch:   sdk.NewInt(2),
					Amount:  sdk.Coin{Denom: "stake", Amount: sdk.NewInt(10000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected: map[int64]sdk.Coin{
				2: {Denom: "stake", Amount: sdk.NewInt(5000)},
			},
			redemptionAmount: sdk.Coin{Denom: "stake", Amount: sdk.NewInt(10000)},
			success:          true,
		},
		{
			name: "redemption exactly filled with two deposits",
			deposits: []*types.Deposit{
				{
					ChainId: "chain-1",
					Epoch:   sdk.NewInt(1),
					Amount:  sdk.Coin{Denom: "stake", Amount: sdk.NewInt(10000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
				{
					ChainId: "chain-1",
					Epoch:   sdk.NewInt(2),
					Amount:  sdk.Coin{Denom: "stake", Amount: sdk.NewInt(5000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected:         map[int64]sdk.Coin{},
			redemptionAmount: sdk.Coin{Denom: "stake", Amount: sdk.NewInt(15000)},
			success:          true,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			pstakeApp, ctx := suite.app, suite.ctx

			for _, deposit := range t.deposits {
				pstakeApp.LiquidStakeIBCKeeper.SetDeposit(ctx, deposit)
			}

			err := pstakeApp.LiquidStakeIBCKeeper.AdjustDepositsForRedemption(
				ctx,
				&types.HostChain{ChainId: "chain-1"},
				t.redemptionAmount,
			)

			if t.success {
				suite.Require().Equal(err, nil)

				deposits := pstakeApp.LiquidStakeIBCKeeper.GetAllDeposits(ctx)
				for _, deposit := range deposits {
					suite.Require().Equal(deposit.Amount, t.expected[deposit.Epoch.Int64()])
				}

				suite.Require().Equal(len(deposits), len(t.expected))
			} else {
				suite.Require().NotEqual(err, nil)
			}
		})
	}
}
