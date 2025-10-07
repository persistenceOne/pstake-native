package liquidstake

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v5/x/liquidstake/keeper"
	"github.com/persistenceOne/pstake-native/v5/x/liquidstake/types"
)

func BeginBlock(ctx context.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if !k.GetParams(sdkCtx).ModulePaused {
		// return value of UpdateLiquidValidatorSet is useful only in testing
		_ = k.UpdateLiquidValidatorSet(sdkCtx, false)
	}
}
