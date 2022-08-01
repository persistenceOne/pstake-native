package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (suite *IntegrationTestSuite) TestKeeper_MintRewardsClaimed() {
	app, ctx := suite.app, suite.ctx
	keeper := app.CosmosKeeper

	err := keeper.MintRewardsClaimed(ctx, sdk.NewCoin(cosmosTypes.DefaultStakingDenom, sdk.NewInt(0)))
	suite.Error(err)

	err = keeper.MintRewardsClaimed(ctx, sdk.NewCoin(cosmosTypes.DefaultStakingDenom, sdk.NewInt(100000)))
	suite.NoError(err)

	rewardAccount, err := sdk.AccAddressFromBech32(keeper.GetParams(ctx).WeightedDeveloperRewardsReceivers[0].Address)
	suite.NoError(err)

	amount := app.BankKeeper.GetBalance(ctx, rewardAccount, cosmosTypes.DefaultMintDenom)
	suite.Equal(sdk.NewInt(5319), amount.Amount)
}
