package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetAllowListedValidators sets allowlisted validator set
func (k Keeper) SetAllowListedValidators(ctx sdk.Context, allowlistedValidators types.AllowListedValidators) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AllowListedValidatorsKey, k.cdc.MustMarshal(&allowlistedValidators))
}

// GetAllowListedValidators gets the allow listed validator set
func (k Keeper) GetAllowListedValidators(ctx sdk.Context) types.AllowListedValidators {
	store := ctx.KVStore(k.storeKey)
	var allowListedValidators types.AllowListedValidators
	k.cdc.MustUnmarshal(store.Get(types.AllowListedValidatorsKey), &allowListedValidators)

	return allowListedValidators
}
