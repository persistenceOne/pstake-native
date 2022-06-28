package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

/*
Adds the rewards claimed message entry to the rewards claimed store with the given validator address.
Performs the following actions :
  1. Checks if store has the key or not. If not then create new entry
  2. Checks if store has it and matches all the details present in the message. If not then create a new entry.
  3. Finally, if all the details match then append the validator address to keep track.
*/
func (k Keeper) addToRewardsClaimedPool(ctx sdk.Context, msg cosmosTypes.MsgRewardsClaimedOnCosmosChain, validatorAddress sdk.ValAddress) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	key := []byte(cosmosTypes.GetChainIDAndBlockHeightKey(msg.ChainID, msg.BlockHeight))
	totalValidatorCount := k.GetTotalValidatorOrchestratorCount(ctx)
	// store has key in it or not
	if !rewardsClaimedStore.Has(key) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewRewardsClaimedValue(msg, validatorAddress, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		rewardsClaimedStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	var rewardsClaimedValue cosmosTypes.RewardsClaimedValue
	k.cdc.MustUnmarshal(rewardsClaimedStore.Get(key), &rewardsClaimedValue)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotRewardsClaimed(rewardsClaimedValue, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewRewardsClaimedValue(msg, validatorAddress, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		rewardsClaimedStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	// if equal then check if orchestrator has already sent same details previously
	if !rewardsClaimedValue.Find(validatorAddress.String()) {
		rewardsClaimedValue.UpdateValues(validatorAddress.String(), k.GetTotalValidatorOrchestratorCount(ctx))
		rewardsClaimedStore.Set(key, k.cdc.MustMarshal(&rewardsClaimedValue))
		return
	}
}

// Gets the list of all rewards claimed requests from rewards claimed store
func (k Keeper) getAllFromRewardsClaimedPool(ctx sdk.Context) (list []cosmosTypes.RewardsClaimedValue, keys [][]byte) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	iterator := rewardsClaimedStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rewardsClaimedValue cosmosTypes.RewardsClaimedValue
		k.cdc.MustUnmarshal(iterator.Value(), &rewardsClaimedValue)
		list = append(list, rewardsClaimedValue)
		keys = append(keys, iterator.Key())
	}
	return list, keys
}

// Set added to current epoch true for the given key in rewards claimed store
func (k Keeper) setAddedToCurrentEpochTrue(ctx sdk.Context, key []byte) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)

	var rewardsClaimedValue cosmosTypes.RewardsClaimedValue
	k.cdc.MustUnmarshal(rewardsClaimedStore.Get(key), &rewardsClaimedValue)
	rewardsClaimedValue.AddedToCurrentEpoch = true
	rewardsClaimedStore.Set(key, k.cdc.MustMarshal(&rewardsClaimedValue))
}

// Remove the given key from the rewards claimed store
func (k Keeper) deleteFromRewardsClaimedPool(ctx sdk.Context, key []byte) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	rewardsClaimedStore.Delete(key)
}

//______________________________________________________________________________________________________________________

// Add the rewards claimed amount to the current epoch
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

// Get the amount of rewards claimed mapped to the given epoch number
func (k Keeper) getFromRewardsInCurrentEpochAmount(ctx sdk.Context, epochNumber int64) (amount sdk.Coin) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	if !rewardsInCurrentEpochStore.Has(cosmosTypes.Int64Bytes(epochNumber)) {
		return sdk.NewInt64Coin("stake", 0)
	}
	k.cdc.MustUnmarshal(rewardsInCurrentEpochStore.Get(cosmosTypes.Int64Bytes(epochNumber)), &amount)
	return amount
}

// shifts the rewards in the given epoch number to the next epoch number for rewards delegation
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
	return
}

// Remove the given key from the rewards in current epoch store
func (k Keeper) deleteFromRewardsInCurrentEpoch(ctx sdk.Context, epochNumber int64) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	rewardsInCurrentEpochStore.Delete(cosmosTypes.Int64Bytes(epochNumber))
}

//______________________________________________________________________________________________________________________

/*
ProcessRewards processes all the rewards requests
This function is called every EndBlocker to perform the defined set of actions as mentioned below :
   1. Get the list of all rewards requests
   2. Checks if the majority of the validator oracle have sent the minting request. Also checks the
      addedToCurrentEpoch flag.
   3. If majority is reached and other conditions match then rewards are added to current epoch and
      addedToCurrentEpoch flag is marked true.
   4. Another condition of ActiveBlockHeight is also checked whether to delete the entry or not.
*/
func (k Keeper) ProcessRewards(ctx sdk.Context) {
	rewardsList, keys := k.getAllFromRewardsClaimedPool(ctx)
	if len(rewardsList) != len(keys) {
		panic(fmt.Errorf("rewards list and keys do not have equal number of elements"))
	}
	for i, r := range rewardsList {
		if r.Ratio.GT(cosmosTypes.MinimumRatioForMajority) && !r.AddedToCurrentEpoch {
			r.AddedToCurrentEpoch = true

			k.addToRewardsInCurrentEpoch(ctx, r.RewardsClaimed.AmountClaimed)

			k.setAddedToCurrentEpochTrue(ctx, keys[i])
		}
		if r.ActiveBlockHeight <= ctx.BlockHeight() && r.AddedToCurrentEpoch {
			k.deleteFromRewardsClaimedPool(ctx, keys[i])
		}
	}
}

// StoreValueEqualOrNotRewardsClaimed Helper function for rewards claimed store to check if the relevant details in the message matches or not.
func StoreValueEqualOrNotRewardsClaimed(storeValue cosmosTypes.RewardsClaimedValue,
	msgValue cosmosTypes.MsgRewardsClaimedOnCosmosChain) bool {
	if !storeValue.RewardsClaimed.AmountClaimed.IsEqual(msgValue.AmountClaimed) {
		return false
	}
	if storeValue.RewardsClaimed.ChainID != msgValue.ChainID {
		return false
	}
	if storeValue.RewardsClaimed.BlockHeight != msgValue.BlockHeight {
		return false
	}
	return true
}
