package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// SetModuleState allows all module transactions
func (k Keeper) SetModuleState(ctx sdk.Context, enable bool) {
	store := ctx.KVStore(k.storeKey)
	storeBool := "false"
	if enable {
		storeBool = "true"
	}
	store.Set(types.ModuleEnableKey, []byte(storeBool))
}

// GetModuleState blocks all module transactions except for register proposal or valset update
func (k Keeper) GetModuleState(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	storeBool := store.Get(types.ModuleEnableKey)
	moduleState := false
	if string(storeBool) == "true" {
		moduleState = true
	}
	return moduleState
}
