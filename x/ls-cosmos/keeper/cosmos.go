package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

// SetCosmosIBCParams sets the cosmos IBC params in store
func (k Keeper) SetCosmosIBCParams(ctx sdk.Context, proposal types.RegisterCosmosChainProposal) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.CosmosIBCParams, k.cdc.MustMarshal(&proposal))
}

// GetCosmosIBCParams gets the cosmos IBC params in store
func (k Keeper) GetCosmosIBCParams(ctx sdk.Context) types.RegisterCosmosChainProposal {
	store := ctx.KVStore(k.storeKey)

	var cosmosIBCParams types.RegisterCosmosChainProposal
	k.cdc.MustUnmarshal(store.Get(types.CosmosIBCParams), &cosmosIBCParams)

	return cosmosIBCParams
}
