package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

// IncrementHostChainID increments and returns a unique ID for an unbonding operation
func (k Keeper) IncrementHostChainID(ctx sdk.Context) (hostChainID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HostChainIDKeyPrefix)
	if bz != nil {
		hostChainID = binary.BigEndian.Uint64(bz)
	}

	hostChainID++

	// Convert back into bytes for storage
	bz = make([]byte, 8)
	binary.BigEndian.PutUint64(bz, hostChainID)

	store.Set(types.HostChainIDKeyPrefix, bz)

	return hostChainID
}

// SetHostChain set a specific chain in the store from its index
func (k Keeper) SetHostChain(ctx sdk.Context, chain types.HostChain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKeyPrefix)
	b := k.cdc.MustMarshal(&chain)
	store.Set(types.HostChainKey(
		chain.ID,
	), b)
}

// GetHostChain returns a chain from its index
func (k Keeper) GetHostChain(
	ctx sdk.Context,
	id uint64,
) (val types.HostChain, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKeyPrefix)

	b := store.Get(types.HostChainKey(
		id,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetHostChain returns a chain from its index
func (k Keeper) GetHostChainsByChainID(
	ctx sdk.Context,
	chainID string,
) []types.HostChain {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)

	defer iterator.Close()

	var vals []types.HostChain
	for ; iterator.Valid(); iterator.Next() {
		var val types.HostChain
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.ChainID == chainID {
			vals = append(vals, val)
		}
	}
	return vals
}

// RemoveHostChain removes a chain from the store
func (k Keeper) RemoveHostChain(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKeyPrefix)
	store.Delete(types.HostChainKey(id))
}

// GetAllHostChain returns all chain
func (k Keeper) GetAllHostChain(ctx sdk.Context) []types.HostChain {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)

	defer iterator.Close()

	var list []types.HostChain
	for ; iterator.Valid(); iterator.Next() {
		var val types.HostChain
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}
