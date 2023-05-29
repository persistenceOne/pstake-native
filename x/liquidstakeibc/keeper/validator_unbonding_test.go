package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) TestGetSetValidatorUnbonding() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	suite.app.LiquidStakeIBCKeeper.SetValidatorUnbonding(
		suite.ctx,
		&types.ValidatorUnbonding{
			ChainId:          suite.path.EndpointB.Chain.ChainID,
			ValidatorAddress: TestAddress,
			EpochNumber:      epoch,
		},
	)

	unbonding, found := suite.app.LiquidStakeIBCKeeper.GetValidatorUnbonding(
		suite.ctx,
		suite.path.EndpointB.Chain.ChainID,
		TestAddress,
		epoch,
	)

	suite.Require().Equal(true, found)
	suite.Require().Equal(TestAddress, unbonding.ValidatorAddress)
}

func (suite *IntegrationTestSuite) TestDeleteValidatorUnbonding() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbonding := &types.ValidatorUnbonding{
		ChainId:          suite.path.EndpointB.Chain.ChainID,
		ValidatorAddress: TestAddress,
		EpochNumber:      epoch,
	}

	suite.app.LiquidStakeIBCKeeper.SetValidatorUnbonding(suite.ctx, unbonding)
	suite.app.LiquidStakeIBCKeeper.DeleteValidatorUnbonding(suite.ctx, unbonding)

	unbonding, found := suite.app.LiquidStakeIBCKeeper.GetValidatorUnbonding(
		suite.ctx,
		TestAddress,
		suite.path.EndpointB.Chain.ChainID,
		epoch,
	)

	suite.Require().Equal(false, found)
}

func (suite *IntegrationTestSuite) TestDeleteValidatorUnbondingsForSequenceID() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbondings := []*types.ValidatorUnbonding{
		{
			ChainId:          suite.path.EndpointB.Chain.ChainID,
			ValidatorAddress: TestAddress,
			EpochNumber:      epoch,
			IbcSequenceId:    "sequence-1",
		},
		{
			ChainId:          suite.path.EndpointB.Chain.ChainID,
			ValidatorAddress: TestAddress,
			EpochNumber:      epoch + 1,
			IbcSequenceId:    "sequence-1",
		},
		{
			ChainId:          suite.path.EndpointB.Chain.ChainID,
			ValidatorAddress: TestAddress,
			EpochNumber:      epoch + 2,
			IbcSequenceId:    "sequence-2",
		},
	}

	for _, unbonding := range unbondings {
		suite.app.LiquidStakeIBCKeeper.SetValidatorUnbonding(suite.ctx, unbonding)
	}

	suite.app.LiquidStakeIBCKeeper.DeleteValidatorUnbondingsForSequenceID(suite.ctx, "sequence-1")

	updatedUnbondings := suite.app.LiquidStakeIBCKeeper.FilterValidatorUnbondings(
		suite.ctx,
		func(u types.ValidatorUnbonding) bool { return true },
	)

	suite.Require().Equal(1, len(updatedUnbondings))
	suite.Require().Equal("sequence-2", updatedUnbondings[0].IbcSequenceId)
}

func (suite *IntegrationTestSuite) TestGetAllValidatorUnbondedAmount() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbondings := []*types.ValidatorUnbonding{
		{
			ChainId:          suite.path.EndpointB.Chain.ChainID,
			ValidatorAddress: TestAddress,
			EpochNumber:      epoch,
			MatureTime:       time.Now(),
			Amount:           sdk.NewCoin(HostDenom, sdk.NewInt(100)),
		},
		{
			ChainId:          suite.path.EndpointB.Chain.ChainID,
			ValidatorAddress: TestAddress,
			EpochNumber:      epoch + 1,
			MatureTime:       time.Now(),
			Amount:           sdk.NewCoin(HostDenom, sdk.NewInt(100)),
		},
		{
			ChainId:          suite.path.EndpointB.Chain.ChainID,
			ValidatorAddress: TestAddress,
			EpochNumber:      epoch + 2,
			MatureTime:       time.Time{},
			Amount:           sdk.NewCoin(HostDenom, sdk.NewInt(100)),
		},
	}

	for _, unbonding := range unbondings {
		suite.app.LiquidStakeIBCKeeper.SetValidatorUnbonding(suite.ctx, unbonding)
	}

	hc, _ := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.ctx, suite.path.EndpointB.Chain.ChainID)
	amount := suite.app.LiquidStakeIBCKeeper.GetAllValidatorUnbondedAmount(suite.ctx, hc)

	suite.Require().Equal(int64(200), amount.Int64())
}

func (suite *IntegrationTestSuite) TestFilterValidatorUnbondings() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbonding := &types.ValidatorUnbonding{
		ChainId:          suite.path.EndpointB.Chain.ChainID,
		ValidatorAddress: TestAddress,
		EpochNumber:      epoch,
	}

	suite.app.LiquidStakeIBCKeeper.SetValidatorUnbonding(suite.ctx, unbonding)

	unbondings := suite.app.LiquidStakeIBCKeeper.FilterValidatorUnbondings(
		suite.ctx,
		func(u types.ValidatorUnbonding) bool {
			return u.ChainId == suite.path.EndpointB.Chain.ChainID &&
				u.ValidatorAddress == TestAddress &&
				u.EpochNumber == epoch
		},
	)

	suite.Require().Equal(1, len(unbondings))
	suite.Require().Equal(unbondings[0], unbonding)
}
