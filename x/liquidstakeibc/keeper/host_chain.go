package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// SetHostChain sets a host chain in the store
func (k *Keeper) SetHostChain(ctx sdk.Context, hostZone *types.HostChain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	bytes := k.cdc.MustMarshal(hostZone)
	store.Set([]byte(hostZone.ChainId), bytes)
}

// GetHostChain returns a host chain given its id
func (k *Keeper) GetHostChain(ctx sdk.Context, chainID string) (types.HostChain, bool) {
	hc := types.HostChain{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	bytes := store.Get([]byte(chainID))
	if len(bytes) == 0 {
		return hc, false
	}

	k.cdc.MustUnmarshal(bytes, &hc)
	return hc, true
}

// GetHostChainFromLocalDenom returns a host chain given its ibc denomination on Persistence
func (k *Keeper) GetHostChainFromLocalDenom(ctx sdk.Context, localDenom string) (types.HostChain, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	found := false
	hc := types.HostChain{}
	for ; iterator.Valid(); iterator.Next() {
		chain := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &chain)

		if chain.LocalDenom == localDenom {
			hc = chain
			found = true
			break
		}
	}

	return hc, found
}
