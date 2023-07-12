package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) TestGetSetUnbonding() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	suite.app.LiquidStakeIBCKeeper.SetUnbonding(
		suite.ctx,
		&types.Unbonding{
			ChainId:     suite.chainB.ChainID,
			EpochNumber: epoch,
		},
	)

	unbonding, found := suite.app.LiquidStakeIBCKeeper.GetUnbonding(
		suite.ctx,
		suite.chainB.ChainID,
		epoch,
	)

	suite.Require().Equal(true, found)
	suite.Require().Equal(suite.chainB.ChainID, unbonding.ChainId)
}

func (suite *IntegrationTestSuite) TestDeleteUnbonding() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbonding := &types.Unbonding{
		ChainId:     suite.chainB.ChainID,
		EpochNumber: epoch,
	}

	suite.app.LiquidStakeIBCKeeper.SetUnbonding(suite.ctx, unbonding)
	suite.app.LiquidStakeIBCKeeper.DeleteUnbonding(suite.ctx, unbonding)

	unbonding, found := suite.app.LiquidStakeIBCKeeper.GetUnbonding(
		suite.ctx,
		suite.chainB.ChainID,
		epoch,
	)

	suite.Require().Equal(false, found)
}

func (suite *IntegrationTestSuite) TestFilterUnbondings() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbonding := &types.Unbonding{
		ChainId:     suite.chainB.ChainID,
		EpochNumber: epoch,
	}

	suite.app.LiquidStakeIBCKeeper.SetUnbonding(suite.ctx, unbonding)

	unbondings := suite.app.LiquidStakeIBCKeeper.FilterUnbondings(
		suite.ctx,
		func(u types.Unbonding) bool {
			return u.ChainId == suite.chainB.ChainID &&
				u.EpochNumber == epoch
		},
	)

	suite.Require().Equal(1, len(unbondings))
	suite.Require().Equal(suite.chainB.ChainID, unbondings[0].ChainId)
	suite.Require().Equal(epoch, unbondings[0].EpochNumber)
}

func (suite *IntegrationTestSuite) TestIncreaseUndelegatingAmountForEpoch() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	tc := []struct {
		name      string
		burn      sdk.Coin
		unbond    sdk.Coin
		unbonding *types.Unbonding
	}{
		{
			name:   "Success",
			burn:   sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
			unbond: sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
			unbonding: &types.Unbonding{
				ChainId:      suite.chainB.ChainID,
				EpochNumber:  epoch,
				BurnAmount:   sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
				UnbondAmount: sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
			},
		},
		{
			name:   "NotFound",
			burn:   sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
			unbond: sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
			unbonding: &types.Unbonding{
				ChainId:      suite.chainB.ChainID,
				EpochNumber:  epoch + 1,
				BurnAmount:   sdk.NewCoin(HostDenom, sdk.NewInt(0)),
				UnbondAmount: sdk.NewCoin(HostDenom, sdk.NewInt(0)),
			},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			suite.app.LiquidStakeIBCKeeper.SetUnbonding(suite.ctx, t.unbonding)

			suite.app.LiquidStakeIBCKeeper.IncreaseUndelegatingAmountForEpoch(
				suite.ctx,
				t.unbonding.ChainId,
				t.unbonding.EpochNumber,
				t.burn,
				t.unbond,
			)

			unbonding, _ := suite.app.LiquidStakeIBCKeeper.GetUnbonding(
				suite.ctx,
				t.unbonding.ChainId,
				t.unbonding.EpochNumber,
			)

			suite.Require().Equal(t.unbonding.BurnAmount.Add(t.unbond), unbonding.BurnAmount)
			suite.Require().Equal(t.unbonding.UnbondAmount.Add(t.unbond), unbonding.UnbondAmount)
		})
	}
}

func (suite *IntegrationTestSuite) TestFailAllUnbondingsForSequenceID() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbondings := []*types.Unbonding{
		{
			ChainId:       suite.chainB.ChainID,
			EpochNumber:   epoch,
			IbcSequenceId: "sequence-1",
			State:         types.Unbonding_UNBONDING_PENDING,
		},
		{
			ChainId:       suite.chainB.ChainID,
			EpochNumber:   epoch + 1,
			IbcSequenceId: "sequence-1",
			State:         types.Unbonding_UNBONDING_MATURING,
		},
		{
			ChainId:       suite.chainB.ChainID,
			EpochNumber:   epoch + 2,
			IbcSequenceId: "sequence-2",
			State:         types.Unbonding_UNBONDING_MATURED,
		},
	}

	for _, ub := range unbondings {
		suite.app.LiquidStakeIBCKeeper.SetUnbonding(suite.ctx, ub)
	}

	suite.app.LiquidStakeIBCKeeper.FailAllUnbondingsForSequenceID(suite.ctx, "sequence-1")

	updatedUnbondings := suite.app.LiquidStakeIBCKeeper.FilterUnbondings(
		suite.ctx,
		func(u types.Unbonding) bool { return true },
	)

	suite.Require().Equal(3, len(updatedUnbondings))

	for _, unbonding := range updatedUnbondings {
		if unbonding.IbcSequenceId == "sequence-1" {
			suite.Require().Equal(types.Unbonding_UNBONDING_FAILED, unbonding.IbcSequenceId)
		}
	}
}

func (suite *IntegrationTestSuite) TestRevertUnbondingsState() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbondings := []*types.Unbonding{
		{
			ChainId:     suite.chainB.ChainID,
			EpochNumber: epoch,
			State:       types.Unbonding_UNBONDING_PENDING,
		},
		{
			ChainId:     suite.chainB.ChainID,
			EpochNumber: epoch + 1,
			State:       types.Unbonding_UNBONDING_INITIATED,
		},
		{
			ChainId:     suite.chainB.ChainID,
			EpochNumber: epoch + 2,
			State:       types.Unbonding_UNBONDING_MATURING,
		},
		{
			ChainId:     suite.chainB.ChainID,
			EpochNumber: epoch + 3,
			State:       types.Unbonding_UNBONDING_MATURED,
		},
		{
			ChainId:     suite.chainB.ChainID,
			EpochNumber: epoch + 4,
			State:       types.Unbonding_UNBONDING_CLAIMABLE,
		},
		{
			ChainId:     suite.chainB.ChainID,
			EpochNumber: epoch + 5,
			State:       types.Unbonding_UNBONDING_FAILED,
		},
	}

	suite.app.LiquidStakeIBCKeeper.RevertUnbondingsState(suite.ctx, unbondings)

	updatedUnbondings := suite.app.LiquidStakeIBCKeeper.FilterUnbondings(
		suite.ctx,
		func(u types.Unbonding) bool { return true },
	)

	for _, unbonding := range updatedUnbondings {
		switch unbonding.EpochNumber {
		case epoch:
			suite.Assert().Equal(types.Unbonding_UNBONDING_PENDING, unbonding.State)
		case epoch + 1:
			suite.Assert().Equal(types.Unbonding_UNBONDING_PENDING, unbonding.State)
		case epoch + 2:
			suite.Assert().Equal(types.Unbonding_UNBONDING_INITIATED, unbonding.State)
		case epoch + 3:
			suite.Assert().Equal(types.Unbonding_UNBONDING_MATURING, unbonding.State)
		case epoch + 4:
			suite.Assert().Equal(types.Unbonding_UNBONDING_MATURED, unbonding.State)
		case epoch + 5:
			suite.Assert().Equal(types.Unbonding_UNBONDING_FAILED, unbonding.State)
		}
	}
}
