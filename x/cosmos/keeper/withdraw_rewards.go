package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) addToRewardsClaimedPool(ctx sdk.Context, orchAddress sdk.AccAddress, amount sdk.Coin, chainID string, blockHeight int64) error {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	key := []byte(cosmosTypes.GetChainIDAndBlockHeightKey(chainID, blockHeight))
	if rewardsClaimedStore.Has(key) {
		bz := rewardsClaimedStore.Get(key)
		var rewardsClaimedValue cosmosTypes.RewardsClaimedValue
		err := rewardsClaimedValue.Unmarshal(bz)
		if err != nil {
			return err
		}
		if rewardsClaimedValue.Find(orchAddress.String()) {
			return fmt.Errorf("already sent the confirmation for rewards claimed")
		}

		rewardsClaimedValue.AddAndIncrement(orchAddress.String())
		rewardsClaimedValue.Amount = append(rewardsClaimedValue.Amount, amount)
		rewardsClaimedValue.Ratio = float32(rewardsClaimedValue.Counter) / float32(k.getTotalValidatorOrchestratorCount(ctx))

		bz1, err := rewardsClaimedValue.Marshal()
		if err != nil {
			return err
		}

		rewardsClaimedStore.Set(key, bz1)
		return nil
	}
	ratio := float32(1) / float32(k.getTotalValidatorOrchestratorCount(ctx))
	newRewardsClaimedValue := cosmosTypes.NewRewardsClaimedValue(orchAddress, amount, ratio, ctx.BlockHeight(), ctx.BlockHeight()+cosmosTypes.StorageWindow)
	bz, err := newRewardsClaimedValue.Marshal()
	if err != nil {
		return err
	}
	rewardsClaimedStore.Set(key, bz)
	return nil
}

func (k Keeper) getAllFromRewardsClaimedPool(ctx sdk.Context) (list []cosmosTypes.RewardsClaimedValue, keys [][]byte, err error) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	iterator := rewardsClaimedStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rewardsClaimedValue cosmosTypes.RewardsClaimedValue
		err = rewardsClaimedValue.Unmarshal(iterator.Value())
		if err != nil {
			return list, keys, err
		}
		list = append(list, rewardsClaimedValue)
		keys = append(keys, iterator.Key())
	}
	return list, keys, err
}

func (k Keeper) setAddedToCurrentEpochTrue(ctx sdk.Context, key []byte, val cosmosTypes.RewardsClaimedValue) error {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	bz, err := val.Marshal()
	if err != nil {
		return err
	}
	rewardsClaimedStore.Set(key, bz)
	return nil
}

func (k Keeper) deleteFromRewardsClaimedPool(ctx sdk.Context, key []byte) {
	rewardsClaimedStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyRewardsStore)
	rewardsClaimedStore.Delete(key)
}

//______________________________________________________________________________________________________________________

func (k Keeper) addToRewardsInCurrentEpoch(ctx sdk.Context, value cosmosTypes.RewardsClaimedValue) error {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	currentEpoch := k.epochsKeeper.GetEpochInfo(ctx, k.GetParams(ctx).StakingEpochIdentifier).CurrentEpoch
	key := cosmosTypes.Int64Bytes(currentEpoch)
	bz, err := value.Marshal()
	if err != nil {
		return err
	}
	rewardsInCurrentEpochStore.Set(key, bz)
	return nil
}

func (k Keeper) getFromRewardsInCurrentEpochAmount(ctx sdk.Context, epochNumber int64) (sdk.Coin, error) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	bz := rewardsInCurrentEpochStore.Get(cosmosTypes.Int64Bytes(epochNumber))
	var rewardsClaimedValue cosmosTypes.RewardsClaimedValue
	err := rewardsClaimedValue.Unmarshal(bz)
	if err != nil {
		return sdk.Coin{}, err
	}

	//Take the first element as all of them are same
	//TODO : check consistency of amount claimed, if not consistent then an average amount can be taken
	return rewardsClaimedValue.Amount[0], nil
}

func (k Keeper) deleteFromRewardsInCurrentEpoch(ctx sdk.Context, epochNumber int64) {
	rewardsInCurrentEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCurrentEpochRewardsStore)
	rewardsInCurrentEpochStore.Delete(cosmosTypes.Int64Bytes(epochNumber))
}

//______________________________________________________________________________________________________________________

func (k Keeper) ProcessRewards(ctx sdk.Context) error {
	rewardsList, keys, err := k.getAllFromRewardsClaimedPool(ctx)
	if err != nil {
		return err
	}
	if len(rewardsList) != len(keys) {
		return fmt.Errorf("rewards list and keys do not have equal number of elements")
	}
	for i, r := range rewardsList {
		if r.Ratio > cosmosTypes.MinimumRatioForMajority && !r.AddedToCurrentEpoch {
			r.AddedToCurrentEpoch = true

			err = k.addToRewardsInCurrentEpoch(ctx, r)
			if err != nil {
				return err
			}

			err = k.setAddedToCurrentEpochTrue(ctx, keys[i], r)
			if err != nil {
				return err
			}
		}
		if r.ActiveBlockHeight <= ctx.BlockHeight() {
			k.deleteFromRewardsClaimedPool(ctx, keys[i])
		}
	}
	return nil
}
