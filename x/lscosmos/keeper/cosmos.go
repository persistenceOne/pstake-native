package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetCosmosIBCParams sets the cosmos IBC params in store
func (k Keeper) SetCosmosIBCParams(ctx sdk.Context, ibcParams types.CosmosIBCParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CosmosIBCParamsKey, k.cdc.MustMarshal(&ibcParams))
}

func (k Keeper) RegisterICAAccounts(ctx sdk.Context, ibcParams types.CosmosIBCParams) error {

	chainId, err := k.GetChainID(ctx, ibcParams.IBCConnection)
	if err != nil {
		return err
	}

	depositICAPort := chainId + ".deposit"
	err = k.icaControllerKeeper.RegisterInterchainAccount(ctx, ibcParams.IBCConnection, depositICAPort)
	if err != nil {
		return err
	}

	delegateICAPort := chainId + ".delegate"
	err = k.icaControllerKeeper.RegisterInterchainAccount(ctx, ibcParams.IBCConnection, delegateICAPort)
	if err != nil {
		return err
	}

	rewardsICAPort := chainId + ".rewards"
	err = k.icaControllerKeeper.RegisterInterchainAccount(ctx, ibcParams.IBCConnection, rewardsICAPort)
	if err != nil {
		return err
	}

	return nil
}

// GetCosmosIBCParams gets the cosmos IBC params in store
func (k Keeper) GetCosmosIBCParams(ctx sdk.Context) types.CosmosIBCParams {
	store := ctx.KVStore(k.storeKey)

	var cosmosIBCParams types.CosmosIBCParams
	k.cdc.MustUnmarshal(store.Get(types.CosmosIBCParamsKey), &cosmosIBCParams)

	return cosmosIBCParams
}
