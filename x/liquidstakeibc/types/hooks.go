package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistence-sdk/v3/utils"
)

type LiquidStakeIBCHooks interface {
	PostCValueUpdate(ctx sdk.Context, mintDenom, hostDenom string, cValue sdk.Dec) error
}

var _ LiquidStakeIBCHooks = MultiLiquidStakeIBCHooks{}

// MultiLiquidStakeIBCHooks combine multiple liquidstake ibc hooks, all hook functions are run in array sequence
type MultiLiquidStakeIBCHooks []LiquidStakeIBCHooks

func NewMultiLiquidStakeIBCHooks(hooks ...LiquidStakeIBCHooks) MultiLiquidStakeIBCHooks {
	return hooks
}

func (h MultiLiquidStakeIBCHooks) PostCValueUpdate(ctx sdk.Context, mintDenom, hostDenom string, cValue sdk.Dec) error {
	for i := range h {
		wrappedHookFn := func(ctx sdk.Context) error {
			//nolint:scopelint // the variables will be same for each loop, ok to use global
			return h[i].PostCValueUpdate(ctx, mintDenom, hostDenom, cValue)
		}

		err := utils.ApplyFuncIfNoError(ctx, wrappedHookFn)
		if err != nil {
			ctx.Logger().Error("Error occurred in calling PostCValueUpdate hooks, ", "err: ", err, "module:", ModuleName, "index:", i)
		}
	}

	return nil
}
