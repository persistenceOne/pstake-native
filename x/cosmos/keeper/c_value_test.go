package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"testing"
)

func TestKeeper_GetCValue(t *testing.T) {
	_, app, ctx := helpers.CreateTestApp()
	fmt.Println(ctx.BlockHeight())
	keeper := app.CosmosKeeper

	cValue := keeper.GetCValue(ctx)
	fmt.Println(cValue)

	keeper.AddToMinted(ctx, sdk.NewInt64Coin("uatom", 100))
	keeper.AddToVirtuallyStaked(ctx, sdk.NewInt64Coin("uatom", 60))
	keeper.AddToStaked(ctx, sdk.NewInt64Coin("uatom", 39))

	cValue = keeper.GetCValue(ctx)
	fmt.Println(cValue)

	cValue = keeper.GetCValue(ctx)
	fmt.Println(cValue)

	keeper.SlashingEvent(ctx, sdk.NewInt64Coin("uatom", 10))
	keeper.AddToMinted(ctx, sdk.NewInt64Coin("uatom", 100))
	keeper.AddToVirtuallyStaked(ctx, sdk.NewInt64Coin("uatom", 60))
	keeper.AddToStaked(ctx, sdk.NewInt64Coin("uatom", 39))

	fmt.Println(keeper.GetMintedAmount(ctx))
	fmt.Println(keeper.GetVirtuallyStakedAmount(ctx))
	fmt.Println(keeper.GetStakedAmount(ctx))

	cValue = keeper.GetCValue(ctx)
	fmt.Println(cValue)

	cValue = keeper.GetCValue(ctx)
	fmt.Println(cValue)
}
