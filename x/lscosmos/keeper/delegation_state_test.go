package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestDelegationState() {
	app, ctx := suite.app, suite.ctx

	err := app.LSCosmosKeeper.SetHostChainDelegationAddress(ctx, "address_________________")
	suite.NoError(err)

	baseDenom := app.LSCosmosKeeper.GetHostChainParams(ctx).BaseDenom
	delegationState := types.DelegationState{
		HostDelegationAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 100)),
		HostChainDelegationAddress:   "address_________________",
		HostAccountDelegations: []types.HostAccountDelegation{
			{
				ValidatorAddress: "address_______________1",
				Amount:           sdk.NewInt64Coin(baseDenom, 25),
			},
			{
				ValidatorAddress: "address_______________2",
				Amount:           sdk.NewInt64Coin(baseDenom, 75),
			},
		},
	}
	app.LSCosmosKeeper.SetDelegationState(ctx, delegationState)

	err = app.LSCosmosKeeper.SetHostChainDelegationAddress(ctx, "address_________________")
	suite.Error(err)

	app.LSCosmosKeeper.AddBalanceToDelegationState(ctx, sdk.NewInt64Coin(baseDenom, 100))

	delegationState = app.LSCosmosKeeper.GetDelegationState(ctx)
	suite.Equal(sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 200)), delegationState.HostDelegationAccountBalance)

	app.LSCosmosKeeper.RemoveBalanceFromDelegationState(ctx, sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 100)))

	delegationState = app.LSCosmosKeeper.GetDelegationState(ctx)
	suite.Equal(sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 100)), delegationState.HostDelegationAccountBalance)

	app.LSCosmosKeeper.AddHostAccountDelegation(ctx, types.NewHostAccountDelegation("address_______________1", sdk.NewInt64Coin(baseDenom, 25)))
	app.LSCosmosKeeper.AddHostAccountDelegation(ctx, types.NewHostAccountDelegation("address_______________2", sdk.NewInt64Coin(baseDenom, 25)))

	delegationState = app.LSCosmosKeeper.GetDelegationState(ctx)
	suite.Equal(sdk.NewInt64Coin(baseDenom, 150), delegationState.HostAccountDelegations[0].Amount.Add(delegationState.HostAccountDelegations[1].Amount))

	err = app.LSCosmosKeeper.SubtractHostAccountDelegation(ctx, types.NewHostAccountDelegation("address_______________1", sdk.NewInt64Coin(baseDenom, 25)))
	suite.NoError(err)
	err = app.LSCosmosKeeper.SubtractHostAccountDelegation(ctx, types.NewHostAccountDelegation("address_______________2", sdk.NewInt64Coin(baseDenom, 25)))
	suite.NoError(err)

	delegationState = app.LSCosmosKeeper.GetDelegationState(ctx)
	suite.Equal(sdk.NewInt64Coin(baseDenom, 100), delegationState.HostAccountDelegations[0].Amount.Add(delegationState.HostAccountDelegations[1].Amount))

	err = app.LSCosmosKeeper.SubtractHostAccountDelegation(ctx, types.NewHostAccountDelegation("address_______________", sdk.NewInt64Coin(baseDenom, 25)))
	suite.Error(err)
	err = app.LSCosmosKeeper.SubtractHostAccountDelegation(ctx, types.NewHostAccountDelegation("address_______________", sdk.NewInt64Coin(baseDenom, 25)))
	suite.Error(err)
}
