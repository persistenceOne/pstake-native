package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestKeeper_GetCValue() {
	app, ctx := suite.app, suite.ctx
	keeper := app.CosmosKeeper

	cValue := keeper.GetCValue(ctx)
	suite.Equal(sdk.NewDec(1), cValue)

	keeper.AddToMinted(ctx, sdk.NewInt64Coin("uatom", 100))
	keeper.AddToVirtuallyStaked(ctx, sdk.NewInt64Coin("uatom", 60))
	keeper.AddToStaked(ctx, sdk.NewInt64Coin("uatom", 39))

	cValue = keeper.GetCValue(ctx)
	totalStaked := keeper.GetVirtuallyStakedAmount(ctx).Amount.Add(keeper.GetStakedAmount(ctx).Amount).Sub(keeper.GetVirtuallyUnbonded(ctx).Amount)
	calculatedCValue := sdk.NewDecFromInt(keeper.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	suite.Equal(calculatedCValue, cValue)

	cValue = keeper.GetCValue(ctx)
	totalStaked = keeper.GetVirtuallyStakedAmount(ctx).Amount.Add(keeper.GetStakedAmount(ctx).Amount).Sub(keeper.GetVirtuallyUnbonded(ctx).Amount)
	calculatedCValue = sdk.NewDecFromInt(keeper.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	suite.Equal(calculatedCValue, cValue)

	keeper.SlashingEvent(ctx, sdk.NewInt64Coin("uatom", 10))
	keeper.AddToMinted(ctx, sdk.NewInt64Coin("uatom", 100))
	keeper.AddToVirtuallyStaked(ctx, sdk.NewInt64Coin("uatom", 60))
	keeper.AddToStaked(ctx, sdk.NewInt64Coin("uatom", 39))

	cValue = keeper.GetCValue(ctx)
	totalStaked = keeper.GetVirtuallyStakedAmount(ctx).Amount.Add(keeper.GetStakedAmount(ctx).Amount).Sub(keeper.GetVirtuallyUnbonded(ctx).Amount)
	calculatedCValue = sdk.NewDecFromInt(keeper.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	suite.Equal(calculatedCValue, cValue)

	cValue = keeper.GetCValue(ctx)
	totalStaked = keeper.GetVirtuallyStakedAmount(ctx).Amount.Add(keeper.GetStakedAmount(ctx).Amount).Sub(keeper.GetVirtuallyUnbonded(ctx).Amount)
	calculatedCValue = sdk.NewDecFromInt(keeper.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	suite.Equal(calculatedCValue, cValue)
}
