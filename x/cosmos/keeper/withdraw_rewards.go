package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// addToRewardsInCurrentEpoch Add the rewards claimed amount to the current epoch
func (k Keeper) addToRewardsInCurrentEpoch(ctx sdk.Context, amount sdk.Coin) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	currentEpoch := k.epochsKeeper.GetEpochInfo(ctx, k.GetParams(ctx).StakingEpochIdentifier).CurrentEpoch
	key := cosmosTypes.Int64Bytes(currentEpoch)

	// if store does not have key in it then create a new one
	if !rewardsInCurrentEpochStore.Has(key) {
		rewardsInCurrentEpochStore.Set(key, k.cdc.MustMarshal(&amount))
		return
	}

	// if store already has the key then add the amount to the previous value
	var newAmount sdk.Coin
	k.cdc.MustUnmarshal(rewardsInCurrentEpochStore.Get(key), &newAmount)
	newAmount = newAmount.Add(amount)
	rewardsInCurrentEpochStore.Set(key, k.cdc.MustMarshal(&newAmount))
}

// getFromRewardsInCurrentEpochAmount Get the amount of rewards claimed mapped to the given epoch number
func (k Keeper) getFromRewardsInCurrentEpochAmount(ctx sdk.Context, epochNumber int64) (amount sdk.Coin) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	if !rewardsInCurrentEpochStore.Has(cosmosTypes.Int64Bytes(epochNumber)) {
		return sdk.NewInt64Coin("uatom", 0)
	}
	k.cdc.MustUnmarshal(rewardsInCurrentEpochStore.Get(cosmosTypes.Int64Bytes(epochNumber)), &amount)
	return amount
}

// shiftRewardsToNextEpoch shifts the rewards in the given epoch number to the next epoch number for rewards delegation
func (k Keeper) shiftRewardsToNextEpoch(ctx sdk.Context, epochNumber int64) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)

	// get given epoch number amount
	var amount sdk.Coin
	k.cdc.MustUnmarshal(rewardsInCurrentEpochStore.Get(cosmosTypes.Int64Bytes(epochNumber)), &amount)

	// check if next epoch number is present in the store, if not then set this amount to the next epoch number
	newEpochKey := cosmosTypes.Int64Bytes(epochNumber + 1)
	if !rewardsInCurrentEpochStore.Has(newEpochKey) {
		rewardsInCurrentEpochStore.Set(newEpochKey, k.cdc.MustMarshal(&amount))
		return
	}

	// if next epoch number is present then add this amount to the already existing amount for new epoch number
	var nextEpochAmount sdk.Coin
	k.cdc.MustUnmarshal(rewardsInCurrentEpochStore.Get(newEpochKey), &nextEpochAmount)
	nextEpochAmount = nextEpochAmount.Add(amount)
	rewardsInCurrentEpochStore.Set(newEpochKey, k.cdc.MustMarshal(&nextEpochAmount))
}

// deleteFromRewardsInCurrentEpoch Remove the given key from the rewards in current epoch store
func (k Keeper) deleteFromRewardsInCurrentEpoch(ctx sdk.Context, epochNumber int64) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	rewardsInCurrentEpochStore.Delete(cosmosTypes.Int64Bytes(epochNumber))
}
