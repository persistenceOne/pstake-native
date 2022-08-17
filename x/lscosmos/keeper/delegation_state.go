package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetCosmosIBCParams sets the cosmos IBC params in store
func (k Keeper) SetDelegationState(ctx sdk.Context, delegationState types.DelegationState) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DelegationStateKey, k.cdc.MustMarshal(&delegationState))
}

// GetCosmosIBCParams gets the cosmos IBC params in store
func (k Keeper) GetDelegationState(ctx sdk.Context) types.DelegationState {
	store := ctx.KVStore(k.storeKey)

	var delegationState types.DelegationState
	k.cdc.MustUnmarshal(store.Get(types.DelegationStateKey), &delegationState)

	return delegationState
}
