package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) TestGetSetUserUnbonding() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	suite.app.LiquidStakeIBCKeeper.SetUserUnbonding(
		suite.ctx,
		&types.UserUnbonding{
			ChainId:     suite.chainB.ChainID,
			Address:     TestAddress,
			EpochNumber: epoch,
		},
	)

	unbonding, found := suite.app.LiquidStakeIBCKeeper.GetUserUnbonding(
		suite.ctx,
		suite.chainB.ChainID,
		TestAddress,
		epoch,
	)

	suite.Require().Equal(true, found)
	suite.Require().Equal(TestAddress, unbonding.Address)
}

func (suite *IntegrationTestSuite) TestDeleteUserUnbonding() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbonding := &types.UserUnbonding{
		ChainId:     suite.chainB.ChainID,
		Address:     TestAddress,
		EpochNumber: epoch,
	}

	suite.app.LiquidStakeIBCKeeper.SetUserUnbonding(suite.ctx, unbonding)
	suite.app.LiquidStakeIBCKeeper.DeleteUserUnbonding(suite.ctx, unbonding)

	unbonding, found := suite.app.LiquidStakeIBCKeeper.GetUserUnbonding(
		suite.ctx,
		TestAddress,
		suite.chainB.ChainID,
		epoch,
	)

	suite.Require().Equal(false, found)
}

func (suite *IntegrationTestSuite) TestFilterUserUnbondings() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch

	unbonding := &types.UserUnbonding{
		ChainId:     suite.chainB.ChainID,
		Address:     TestAddress,
		EpochNumber: epoch,
	}

	suite.app.LiquidStakeIBCKeeper.SetUserUnbonding(suite.ctx, unbonding)

	unbondings := suite.app.LiquidStakeIBCKeeper.FilterUserUnbondings(
		suite.ctx,
		func(u types.UserUnbonding) bool {
			return u.ChainId == suite.chainB.ChainID &&
				u.Address == TestAddress &&
				u.EpochNumber == epoch
		},
	)

	suite.Require().Equal(1, len(unbondings))
	suite.Require().Equal(unbondings[0], unbonding)
}

func (suite *IntegrationTestSuite) TestIncreaseUserUnbondingAmountForEpoch() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch
	ubd1 := &types.UserUnbonding{
		ChainId:      suite.chainB.ChainID,
		Address:      TestAddress,
		EpochNumber:  epoch,
		StkAmount:    sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
		UnbondAmount: sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
	}
	suite.app.LiquidStakeIBCKeeper.SetUserUnbonding(suite.ctx, ubd1)
	tc := []struct {
		name      string
		burn      sdk.Coin
		unbond    sdk.Coin
		unbonding *types.UserUnbonding
	}{
		{
			name:      "Success",
			burn:      sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
			unbond:    sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
			unbonding: ubd1,
		},
		{
			name:   "NotFound",
			burn:   sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
			unbond: sdk.NewCoin(HostDenom, sdk.NewInt(1000)),
			unbonding: &types.UserUnbonding{
				ChainId:      suite.chainB.ChainID,
				Address:      TestAddress,
				EpochNumber:  epoch + 1,
				StkAmount:    sdk.NewCoin(HostDenom, sdk.NewInt(0)),
				UnbondAmount: sdk.NewCoin(HostDenom, sdk.NewInt(0)),
			},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			suite.app.LiquidStakeIBCKeeper.IncreaseUserUnbondingAmountForEpoch(
				suite.ctx,
				t.unbonding.ChainId,
				t.unbonding.Address,
				t.unbonding.EpochNumber,
				t.burn,
				t.unbond,
			)

			unbonding, _ := suite.app.LiquidStakeIBCKeeper.GetUserUnbonding(
				suite.ctx,
				t.unbonding.ChainId,
				t.unbonding.Address,
				t.unbonding.EpochNumber,
			)

			suite.Require().Equal(t.unbonding.StkAmount.Add(t.unbond), unbonding.StkAmount)
			suite.Require().Equal(t.unbonding.UnbondAmount.Add(t.unbond), unbonding.UnbondAmount)
		})
	}
}
