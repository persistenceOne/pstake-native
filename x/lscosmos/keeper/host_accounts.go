package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetHostAccounts sets host account port ids in store
func (k Keeper) SetHostAccounts(ctx sdk.Context, hostAccounts types.HostAccounts) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.HostAccountsKey, k.cdc.MustMarshal(&hostAccounts))
}

// GetHostAccounts gets host account port ids from store
func (k Keeper) GetHostAccounts(ctx sdk.Context) types.HostAccounts {
	store := ctx.KVStore(k.storeKey)

	var hostAccounts types.HostAccounts
	k.cdc.MustUnmarshal(store.Get(types.HostAccountsKey), &hostAccounts)

	return hostAccounts
}
