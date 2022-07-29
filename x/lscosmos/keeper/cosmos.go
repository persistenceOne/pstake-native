package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetCosmosIBCParams sets the cosmos IBC params in store
func (k Keeper) SetCosmosIBCParams(ctx sdk.Context, proposal types.RegisterCosmosChainProposal) {
	ibcParams := types.NewCosmosIBCParams(
		proposal.IBCConnection,
		proposal.TokenTransferChannel,
		proposal.TokenTransferPort,
		proposal.BaseDenom,
		proposal.MintDenom)
	store := ctx.KVStore(k.storeKey)

	store.Set(types.CosmosIBCParamsKey, k.cdc.MustMarshal(&ibcParams))
}

// GetCosmosIBCParams gets the cosmos IBC params in store
func (k Keeper) GetCosmosIBCParams(ctx sdk.Context) types.CosmosIBCParams {
	store := ctx.KVStore(k.storeKey)

	var cosmosIBCParams types.CosmosIBCParams
	k.cdc.MustUnmarshal(store.Get(types.CosmosIBCParamsKey), &cosmosIBCParams)

	return cosmosIBCParams
}
