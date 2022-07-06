package keeper

import (
	"encoding/binary"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

// autoIncrementID
// a specialized function used for iterating store counters, handling
// returning, initializing and incrementing all at once. This is particularly
// used for the transaction pool and batch pool where each batch or transaction is
// assigned a unique ID.
func (k Keeper) autoIncrementID(ctx sdkTypes.Context, idKey []byte) uint64 {
	id := k.getID(ctx, idKey)
	id++
	k.setID(ctx, id, idKey)
	return id
}

// getID gets a generic uint64 counter from the store, initializing to 1 if no value exists
func (k Keeper) getID(ctx sdkTypes.Context, idKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(idKey)
	id := binary.BigEndian.Uint64(bz)
	return id
}

// setID sets a generic uint64 counter in the store
func (k Keeper) setID(ctx sdkTypes.Context, id uint64, idKey []byte) {
	store := ctx.KVStore(k.storeKey)
	bz := sdkTypes.Uint64ToBigEndian(id)
	store.Set(idKey, bz)
}
