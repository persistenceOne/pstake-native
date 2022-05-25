package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) addToRewardsClaimedPool(ctx sdk.Context, msg cosmosTypes.MsgRewardsClaimedOnCosmosChain) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	key := []byte(cosmosTypes.GetChainIDAndBlockHeightKey(msg.ChainID, msg.BlockHeight))
	totalValidatorCount := k.GetTotalValidatorOrchestratorCount(ctx)
	// store has key in it or not
	if !rewardsClaimedStore.Has(key) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewRewardsClaimedValue(msg, msg.OrchestratorAddress, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		rewardsClaimedStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	var rewardsClaimedValue cosmosTypes.RewardsClaimedValue
	k.cdc.MustUnmarshal(rewardsClaimedStore.Get(key), &rewardsClaimedValue)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotRewardsClaimed(rewardsClaimedValue, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewRewardsClaimedValue(msg, msg.OrchestratorAddress, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		rewardsClaimedStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	// if equal then check if orchestrator has already sent same details previously
	if !rewardsClaimedValue.Find(msg.OrchestratorAddress) {
		rewardsClaimedValue.UpdateValues(msg.OrchestratorAddress, k.GetTotalValidatorOrchestratorCount(ctx))
		rewardsClaimedStore.Set(key, k.cdc.MustMarshal(&rewardsClaimedValue))
		return
	}
}

func (k Keeper) getAllFromRewardsClaimedPool(ctx sdk.Context) (list []cosmosTypes.RewardsClaimedValue, keys [][]byte, err error) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	iterator := rewardsClaimedStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rewardsClaimedValue cosmosTypes.RewardsClaimedValue
		err = k.cdc.Unmarshal(iterator.Value(), &rewardsClaimedValue)
		if err != nil {
			return list, keys, err
		}
		list = append(list, rewardsClaimedValue)
		keys = append(keys, iterator.Key())
	}
	return list, keys, err
}

func (k Keeper) setAddedToCurrentEpochTrue(ctx sdk.Context, key []byte) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)

	var rewardsClaimedValue cosmosTypes.RewardsClaimedValue
	k.cdc.MustUnmarshal(rewardsClaimedStore.Get(key), &rewardsClaimedValue)
	rewardsClaimedValue.AddedToCurrentEpoch = true
	rewardsClaimedStore.Set(key, k.cdc.MustMarshal(&rewardsClaimedValue))
}

func (k Keeper) deleteFromRewardsClaimedPool(ctx sdk.Context, key []byte) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	rewardsClaimedStore.Delete(key)
}

//______________________________________________________________________________________________________________________

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

func (k Keeper) getFromRewardsInCurrentEpochAmount(ctx sdk.Context, epochNumber int64) (amount sdk.Coin) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	if rewardsInCurrentEpochStore.Has(cosmosTypes.Int64Bytes(epochNumber)) {
		return sdk.Coin{}
	}
	k.cdc.MustUnmarshal(rewardsInCurrentEpochStore.Get(cosmosTypes.Int64Bytes(epochNumber)), &amount)
	return amount
}

func (k Keeper) deleteFromRewardsInCurrentEpoch(ctx sdk.Context, epochNumber int64) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	rewardsInCurrentEpochStore.Delete(cosmosTypes.Int64Bytes(epochNumber))
}

//______________________________________________________________________________________________________________________

func (k Keeper) ProcessRewards(ctx sdk.Context) {
	rewardsList, keys, err := k.getAllFromRewardsClaimedPool(ctx)
	if err != nil {
		panic(err)
	}
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
