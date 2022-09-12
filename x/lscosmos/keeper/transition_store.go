package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetIBCTransitionStore sets allowlisted validator set
func (k Keeper) SetIBCTransitionStore(ctx sdk.Context, ibcAmountTransitionStore types.IbcAmountTransitionStore) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.IBCTransitionStore, k.cdc.MustMarshal(&ibcAmountTransitionStore))
}

// GetIBCTransitionStore gets the allow listed validator set
func (k Keeper) GetIBCTransitionStore(ctx sdk.Context) types.IbcAmountTransitionStore {
	store := ctx.KVStore(k.storeKey)
	var ibcAmountTransitionStore types.IbcAmountTransitionStore
	k.cdc.MustUnmarshal(store.Get(types.IBCTransitionStore), &ibcAmountTransitionStore)

	return ibcAmountTransitionStore
}
