package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) TestGetSetDeposit() {
	suite.app.LiquidStakeIBCKeeper.SetDeposit(suite.ctx, &types.Deposit{ChainId: suite.path.EndpointB.Chain.ChainID})
	deposits := suite.app.LiquidStakeIBCKeeper.GetAllDeposits(suite.ctx)

	suite.Require().Equal(1, len(deposits))
	suite.Require().Equal(suite.path.EndpointB.Chain.ChainID, deposits[0].ChainId)
}

func (suite *IntegrationTestSuite) TestDeleteDeposit() {
	deposit := &types.Deposit{ChainId: suite.path.EndpointB.Chain.ChainID}

	suite.app.LiquidStakeIBCKeeper.SetDeposit(suite.ctx, deposit)
	suite.app.LiquidStakeIBCKeeper.DeleteDeposit(suite.ctx, deposit)
	deposits := suite.app.LiquidStakeIBCKeeper.GetAllDeposits(suite.ctx)

	suite.Require().Equal(0, len(deposits))
}

func (suite *IntegrationTestSuite) TestCreateDeposits() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	suite.app.LiquidStakeIBCKeeper.CreateDeposits(suite.ctx, epoch)

	deposits := suite.app.LiquidStakeIBCKeeper.GetAllDeposits(suite.ctx)

	suite.Require().Equal(1, len(deposits))
	suite.Require().Equal(epoch, deposits[0].Epoch.Int64())
}

func (suite *IntegrationTestSuite) TestRevertDepositState() {
	// ibc sequence id is used as index
	deposits := []*types.Deposit{
		{
			State:         types.Deposit_DEPOSIT_PENDING,
			IbcSequenceId: "1",
		},
		{
			State:         types.Deposit_DEPOSIT_SENT,
			IbcSequenceId: "2",
		},
		{
			State:         types.Deposit_DEPOSIT_RECEIVED,
			IbcSequenceId: "3",
		},
		{
			State:         types.Deposit_DEPOSIT_DELEGATING,
			IbcSequenceId: "4",
		},
	}

	suite.app.LiquidStakeIBCKeeper.RevertDepositsState(suite.ctx, deposits)

	for _, deposit := range suite.app.LiquidStakeIBCKeeper.GetAllDeposits(suite.ctx) {
		switch deposit.IbcSequenceId {
		case "1":
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_PENDING)
		case "2":
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_PENDING)
		case "3":
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_SENT)
		case "4":
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_RECEIVED)
		}
	}
}

func (suite *IntegrationTestSuite) TestTransactionSequenceID() {
	sequenceID := suite.app.LiquidStakeIBCKeeper.GetTransactionSequenceID("channel-0", 1)

	suite.Require().Equal("channel-0-sequence-1", sequenceID)
}

func (suite *IntegrationTestSuite) TestAdjustDepositsForRedemption() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	tc := []struct {
		name             string
		deposits         []*types.Deposit
		expected         map[int64]sdk.Coin
		redemptionAmount sdk.Coin
		err              error
	}{
		{
			name: "Case 1",
			deposits: []*types.Deposit{
				{
					ChainId: suite.path.EndpointB.Chain.ChainID,
					Epoch:   sdk.NewInt(epoch),
					Amount:  sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(10000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected: map[int64]sdk.Coin{
				1: {Denom: HostDenom, Amount: sdk.NewInt(5000)},
			},
			redemptionAmount: sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(5000)},
		},
		{
			name: "Case 2",
			deposits: []*types.Deposit{
				{
					ChainId: suite.path.EndpointB.Chain.ChainID,
					Epoch:   sdk.NewInt(epoch),
					Amount:  sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(3500)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected: map[int64]sdk.Coin{
				1: {Denom: HostDenom, Amount: sdk.NewInt(3500)},
			},
			redemptionAmount: sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(5000)},
		},
		{
			name: "Case 3",
			deposits: []*types.Deposit{
				{
					ChainId: suite.path.EndpointB.Chain.ChainID,
					Epoch:   sdk.NewInt(epoch),
					Amount:  sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(5000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected:         map[int64]sdk.Coin{},
			redemptionAmount: sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(5000)},
		},
		{
			name: "Case 4",
			deposits: []*types.Deposit{
				{
					ChainId: suite.path.EndpointB.Chain.ChainID,
					Epoch:   sdk.NewInt(epoch),
					Amount:  sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(10000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
				{
					ChainId: suite.path.EndpointB.Chain.ChainID,
					Epoch:   sdk.NewInt(epoch + 1),
					Amount:  sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(5000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected: map[int64]sdk.Coin{
				1: {Denom: HostDenom, Amount: sdk.NewInt(5000)},
				2: {Denom: HostDenom, Amount: sdk.NewInt(5000)},
			},
			redemptionAmount: sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(5000)},
		},
		{
			name: "Case 5",
			deposits: []*types.Deposit{
				{
					ChainId: suite.path.EndpointB.Chain.ChainID,
					Epoch:   sdk.NewInt(epoch),
					Amount:  sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(5000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
				{
					ChainId: suite.path.EndpointB.Chain.ChainID,
					Epoch:   sdk.NewInt(epoch + 1),
					Amount:  sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(10000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected: map[int64]sdk.Coin{
				2: {Denom: HostDenom, Amount: sdk.NewInt(5000)},
			},
			redemptionAmount: sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(10000)},
		},
		{
			name: "Case 6",
			deposits: []*types.Deposit{
				{
					ChainId: suite.path.EndpointB.Chain.ChainID,
					Epoch:   sdk.NewInt(epoch),
					Amount:  sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(10000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
				{
					ChainId: suite.path.EndpointB.Chain.ChainID,
					Epoch:   sdk.NewInt(epoch + 1),
					Amount:  sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(5000)},
					State:   types.Deposit_DEPOSIT_PENDING,
				},
			},
			expected:         map[int64]sdk.Coin{},
			redemptionAmount: sdk.Coin{Denom: HostDenom, Amount: sdk.NewInt(15000)},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {

			for _, deposit := range t.deposits {
				suite.app.LiquidStakeIBCKeeper.SetDeposit(suite.ctx, deposit)
			}

			err := suite.app.LiquidStakeIBCKeeper.AdjustDepositsForRedemption(
				suite.ctx,
				&types.HostChain{ChainId: suite.path.EndpointB.Chain.ChainID},
				t.redemptionAmount,
			)

			suite.Require().Equal(t.err, err)

			deposits := suite.app.LiquidStakeIBCKeeper.GetAllDeposits(suite.ctx)
			for _, deposit := range deposits {
				suite.Require().Equal(t.expected[deposit.Epoch.Int64()], deposit.Amount)
			}

			suite.Require().Equal(len(t.expected), len(deposits))
		})
	}
}

func (suite *IntegrationTestSuite) TestGetDepositForChainAndEpoch() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	tc := []struct {
		name     string
		deposits []types.Deposit
		chainID  string
		epoch    int64
		expected types.Deposit
		found    bool
	}{
		{
			name: "Success",
			deposits: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch)},
				{ChainId: suite.path.EndpointA.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1)},
				{ChainId: suite.path.EndpointA.Chain.ChainID, Epoch: sdk.NewInt(epoch)},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1)},
			},
			chainID:  suite.path.EndpointB.Chain.ChainID,
			epoch:    epoch,
			expected: types.Deposit{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch)},
			found:    true,
		},
		{
			name: "unsuccessful test",
			deposits: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch)},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1)},
				{ChainId: suite.path.EndpointA.Chain.ChainID, Epoch: sdk.NewInt(epoch)},
				{ChainId: suite.path.EndpointA.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1)},
			},
			chainID:  suite.path.EndpointA.Chain.ChainID,
			epoch:    epoch + 2,
			expected: types.Deposit{},
			found:    false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			for _, deposit := range t.deposits {
				suite.app.LiquidStakeIBCKeeper.SetDeposit(suite.ctx, &deposit)
			}

			hc, found := suite.app.LiquidStakeIBCKeeper.GetDepositForChainAndEpoch(suite.ctx, t.chainID, t.epoch)

			if found {
				suite.Require().Equal(t.chainID, hc.ChainId)
				suite.Require().Equal(sdk.NewInt(t.epoch), hc.Epoch)
			}

			suite.Require().Equal(t.found, found)
		})
	}
}

func (suite *IntegrationTestSuite) TestGetDepositsWithSequenceID() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	tc := []struct {
		name       string
		deposits   []types.Deposit
		sequenceID string
		expected   []types.Deposit
	}{
		{
			name: "Success",
			deposits: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch), IbcSequenceId: "seq-1"},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1), IbcSequenceId: "seq-2"},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 2), IbcSequenceId: "seq-3"},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 3), IbcSequenceId: "seq-4"},
			},
			sequenceID: "seq-1",
			expected: []types.Deposit{
				{
					ChainId:       suite.path.EndpointB.Chain.ChainID,
					Epoch:         sdk.NewInt(1),
					IbcSequenceId: "seq-1",
				},
			},
		},
		{
			name: "NotFound",
			deposits: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch), IbcSequenceId: "seq-1"},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1), IbcSequenceId: "seq-2"},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 2), IbcSequenceId: "seq-3"},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 3), IbcSequenceId: "seq-4"},
			},
			sequenceID: "seq-8",
			expected:   []types.Deposit{},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			for _, deposit := range t.deposits {
				suite.app.LiquidStakeIBCKeeper.SetDeposit(suite.ctx, &deposit)
			}

			hcs := suite.app.LiquidStakeIBCKeeper.GetDepositsWithSequenceID(suite.ctx, t.sequenceID)
			suite.Require().Equal(len(t.expected), len(hcs))

			for _, hc := range hcs {
				suite.Require().Equal(t.sequenceID, hc.IbcSequenceId)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetPendingDepositsBeforeEpoch() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	tc := []struct {
		name     string
		deposits []types.Deposit
		epoch    int64
		expected []types.Deposit
	}{
		{
			name: "Success",
			deposits: []types.Deposit{
				{Epoch: sdk.NewInt(epoch), State: types.Deposit_DEPOSIT_PENDING},
				{Epoch: sdk.NewInt(epoch + 1), State: types.Deposit_DEPOSIT_PENDING},
				{Epoch: sdk.NewInt(epoch + 2), State: types.Deposit_DEPOSIT_RECEIVED},
				{Epoch: sdk.NewInt(epoch + 3), State: types.Deposit_DEPOSIT_DELEGATING},
			},
			epoch: epoch + 1,
			expected: []types.Deposit{
				{Epoch: sdk.NewInt(epoch), State: types.Deposit_DEPOSIT_PENDING},
				{Epoch: sdk.NewInt(epoch + 1), State: types.Deposit_DEPOSIT_PENDING},
			},
		},
		{
			name: "NotFound",
			deposits: []types.Deposit{
				{Epoch: sdk.NewInt(epoch), State: types.Deposit_DEPOSIT_RECEIVED},
				{Epoch: sdk.NewInt(epoch + 1), State: types.Deposit_DEPOSIT_DELEGATING},
				{Epoch: sdk.NewInt(epoch + 2), State: types.Deposit_DEPOSIT_PENDING},
				{Epoch: sdk.NewInt(epoch + 3), State: types.Deposit_DEPOSIT_PENDING},
			},
			epoch:    epoch + 1,
			expected: []types.Deposit{},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			for _, deposit := range t.deposits {
				suite.app.LiquidStakeIBCKeeper.SetDeposit(suite.ctx, &deposit)
			}

			hcs := suite.app.LiquidStakeIBCKeeper.GetPendingDepositsBeforeEpoch(suite.ctx, t.epoch)
			suite.Require().Equal(len(t.expected), len(hcs))

			for _, hc := range hcs {
				suite.Require().LessOrEqual(hc.Epoch.Int64(), t.epoch)
				suite.Require().Equal(types.Deposit_DEPOSIT_PENDING, hc.State)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetDelegableDepositsForChain() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	tc := []struct {
		name     string
		deposits []types.Deposit
		chainID  string
		expected []types.Deposit
	}{
		{
			name: "Success",
			deposits: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch), State: types.Deposit_DEPOSIT_RECEIVED},
				{ChainId: suite.path.EndpointA.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1), State: types.Deposit_DEPOSIT_RECEIVED},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 2), State: types.Deposit_DEPOSIT_RECEIVED},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 3), State: types.Deposit_DEPOSIT_PENDING},
			},
			chainID: suite.path.EndpointB.Chain.ChainID,
			expected: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch), State: types.Deposit_DEPOSIT_RECEIVED},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 2), State: types.Deposit_DEPOSIT_RECEIVED},
			},
		},
		{
			name: "NotFound",
			deposits: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch), State: types.Deposit_DEPOSIT_RECEIVED},
				{ChainId: suite.path.EndpointA.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1), State: types.Deposit_DEPOSIT_RECEIVED},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 2), State: types.Deposit_DEPOSIT_RECEIVED},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 3), State: types.Deposit_DEPOSIT_PENDING},
			},
			chainID:  "test-chain-id",
			expected: []types.Deposit{},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			for _, deposit := range t.deposits {
				suite.app.LiquidStakeIBCKeeper.SetDeposit(suite.ctx, &deposit)
			}

			hcs := suite.app.LiquidStakeIBCKeeper.GetDelegableDepositsForChain(suite.ctx, t.chainID)
			suite.Require().Equal(len(t.expected), len(hcs))

			for _, hc := range hcs {
				suite.Require().Equal(t.chainID, hc.ChainId)
				suite.Require().Equal(types.Deposit_DEPOSIT_RECEIVED, hc.State)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetDelegatingDepositsForChain() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	tc := []struct {
		name     string
		deposits []types.Deposit
		chainID  string
		expected []types.Deposit
	}{
		{
			name: "found test",
			deposits: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch), State: types.Deposit_DEPOSIT_DELEGATING},
				{ChainId: suite.path.EndpointA.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1), State: types.Deposit_DEPOSIT_DELEGATING},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 2), State: types.Deposit_DEPOSIT_DELEGATING},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 3), State: types.Deposit_DEPOSIT_PENDING},
			},
			chainID: suite.path.EndpointB.Chain.ChainID,
			expected: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch), State: types.Deposit_DEPOSIT_DELEGATING},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 2), State: types.Deposit_DEPOSIT_DELEGATING},
			},
		},
		{
			name: "not found test",
			deposits: []types.Deposit{
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch), State: types.Deposit_DEPOSIT_DELEGATING},
				{ChainId: suite.path.EndpointA.Chain.ChainID, Epoch: sdk.NewInt(epoch + 1), State: types.Deposit_DEPOSIT_DELEGATING},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 2), State: types.Deposit_DEPOSIT_DELEGATING},
				{ChainId: suite.path.EndpointB.Chain.ChainID, Epoch: sdk.NewInt(epoch + 3), State: types.Deposit_DEPOSIT_PENDING},
			},
			chainID:  "test-host-chain",
			expected: []types.Deposit{},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			for _, deposit := range t.deposits {
				suite.app.LiquidStakeIBCKeeper.SetDeposit(suite.ctx, &deposit)
			}

			hcs := suite.app.LiquidStakeIBCKeeper.GetDelegatingDepositsForChain(suite.ctx, t.chainID)
			suite.Require().Equal(len(t.expected), len(hcs))

			for _, hc := range hcs {
				suite.Require().Equal(t.chainID, hc.ChainId)
				suite.Require().Equal(types.Deposit_DEPOSIT_DELEGATING, hc.State)
			}
		})
	}
}
