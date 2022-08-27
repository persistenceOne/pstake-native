package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetCosmosParams sets the cosmos IBC params in store
func (k Keeper) SetCosmosParams(ctx sdk.Context, cosmosParams types.CosmosParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CosmosParamsKey, k.cdc.MustMarshal(&cosmosParams))
}

// GetCosmosParams gets the cosmos IBC params in store
func (k Keeper) GetCosmosParams(ctx sdk.Context) types.CosmosParams {
	store := ctx.KVStore(k.storeKey)

	var cosmosParams types.CosmosParams
	k.cdc.MustUnmarshal(store.Get(types.CosmosParamsKey), &cosmosParams)

	return cosmosParams
}
