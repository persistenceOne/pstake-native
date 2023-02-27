package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestAfterEpochEnd() {
	app, ctx := suite.app, suite.ctx

	lscosmosKeeper := app.LSCosmosKeeper

	// enable the module
	lscosmosKeeper.SetModuleState(ctx, true)

	// calling the rewards epoch identifier without setting delegation state
	// to go into len check of Rewards workflow
	suite.Require().NoError(app.LSCosmosKeeper.AfterEpochEnd(ctx, types.RewardEpochIdentifier, 1))

	// get host chain params
	hostChainParams := lscosmosKeeper.GetHostChainParams(ctx)

	// create delegations state for reward epoch
	delegationState := types.DelegationState{
		HostAccountDelegations: []types.HostAccountDelegation{
			{
				ValidatorAddress: "",
				Amount:           sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 600000),
			},
			{
				ValidatorAddress: "",
				Amount:           sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 200000),
			},
			{
				ValidatorAddress: "",
				Amount:           sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 100000),
			},
			{
				ValidatorAddress: "",
				Amount:           sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 100000),
			},
		},
		HostDelegationAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 1000)),
	}

	// set the above created delegation state in module
	lscosmosKeeper.SetDelegationState(ctx, delegationState)

	// create ibcDenom from host chain params
	ibcDenom := ibctransfertypes.ParseDenomTrace(
		ibctransfertypes.GetPrefixedDenom(
			hostChainParams.TransferPort, hostChainParams.TransferChannel, hostChainParams.BaseDenom,
		),
	).IBCDenom()

	// mint coins in module with ibcDenom
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ibcDenom, sdk.NewInt(20000))))
	suite.NoError(err)

	// send the minted coins to types.DepositModuleAccount
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx,
		types.ModuleName,
		types.DepositModuleAccount,
		sdk.NewCoins(sdk.NewCoin(ibcDenom, sdk.NewInt(10000))),
	)
	suite.NoError(err)

	// call the after epoch end of LSCosmosKeeper to perform the actions
	suite.Require().NoError(app.LSCosmosKeeper.AfterEpochEnd(ctx, types.DelegationEpochIdentifier, 1))
	suite.Require().NoError(app.LSCosmosKeeper.AfterEpochEnd(ctx, types.UndelegationEpochIdentifier, 1))
}
