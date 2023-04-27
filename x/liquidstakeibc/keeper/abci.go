package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistence-sdk/v2/utils"
)

func (k *Keeper) BeginBlock(ctx sdk.Context) {
	err := utils.ApplyFuncIfNoError(ctx, k.DoDelegate)
	if err != nil {
		k.Logger(ctx).Error("Unable to Delegate tokens with ", "err: ", err)
	}

	// TODO: Submit validator set queries
}

func (k *Keeper) DoDelegate(ctx sdk.Context) error {
	return nil
}
