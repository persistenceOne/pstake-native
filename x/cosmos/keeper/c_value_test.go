package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper_GetCValue(t *testing.T) {
	_, app, ctx := helpers.CreateTestApp()
	fmt.Println(ctx.BlockHeight())
	keeper := app.CosmosKeeper

	cValue := keeper.GetCValue(ctx)
	require.Equal(t, sdk.NewDec(1), cValue)

	keeper.AddToMinted(ctx, sdk.NewInt64Coin("uatom", 100))
	keeper.AddToVirtuallyStaked(ctx, sdk.NewInt64Coin("uatom", 60))
	keeper.AddToStaked(ctx, sdk.NewInt64Coin("uatom", 39))

	cValue = keeper.GetCValue(ctx)
	totalStaked := keeper.GetVirtuallyStakedAmount(ctx).Amount.Add(keeper.GetStakedAmount(ctx).Amount).Sub(keeper.GetVirtuallyUnbonded(ctx).Amount)
	calculatedCValue := sdk.NewDecFromInt(keeper.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	require.Equal(t, calculatedCValue, cValue)

	cValue = keeper.GetCValue(ctx)
	totalStaked = keeper.GetVirtuallyStakedAmount(ctx).Amount.Add(keeper.GetStakedAmount(ctx).Amount).Sub(keeper.GetVirtuallyUnbonded(ctx).Amount)
	calculatedCValue = sdk.NewDecFromInt(keeper.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	require.Equal(t, calculatedCValue, cValue)

	keeper.SlashingEvent(ctx, sdk.NewInt64Coin("uatom", 10))
	keeper.AddToMinted(ctx, sdk.NewInt64Coin("uatom", 100))
	keeper.AddToVirtuallyStaked(ctx, sdk.NewInt64Coin("uatom", 60))
	keeper.AddToStaked(ctx, sdk.NewInt64Coin("uatom", 39))

	cValue = keeper.GetCValue(ctx)
	totalStaked = keeper.GetVirtuallyStakedAmount(ctx).Amount.Add(keeper.GetStakedAmount(ctx).Amount).Sub(keeper.GetVirtuallyUnbonded(ctx).Amount)
	calculatedCValue = sdk.NewDecFromInt(keeper.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	require.Equal(t, calculatedCValue, cValue)

	cValue = keeper.GetCValue(ctx)
	totalStaked = keeper.GetVirtuallyStakedAmount(ctx).Amount.Add(keeper.GetStakedAmount(ctx).Amount).Sub(keeper.GetVirtuallyUnbonded(ctx).Amount)
	calculatedCValue = sdk.NewDecFromInt(keeper.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	require.Equal(t, calculatedCValue, cValue)
}
