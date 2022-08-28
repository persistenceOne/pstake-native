package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetHostChainParams sets the host chain params in store
func (k Keeper) SetHostChainParams(ctx sdk.Context, hostChainParams types.HostChainParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.HostChainParamsKey, k.cdc.MustMarshal(&hostChainParams))
}

// GetHostChainParams gets the host chain params in store
func (k Keeper) GetHostChainParams(ctx sdk.Context) types.HostChainParams {
	store := ctx.KVStore(k.storeKey)

	var hostChainParams types.HostChainParams
	k.cdc.MustUnmarshal(store.Get(types.HostChainParamsKey), &hostChainParams)

	return hostChainParams
}
