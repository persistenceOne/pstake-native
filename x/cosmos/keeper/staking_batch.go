package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) addToStakingEpoch(ctx sdk.Context, amount sdk.Coin) {
	stakingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyStakingEpochStore)
	currentEpoch := k.epochsKeeper.GetEpochInfo(ctx, k.GetParams(ctx).StakingEpochIdentifier).CurrentEpoch
	key := cosmosTypes.Int64Bytes(currentEpoch)

	// if store does not have key in it then create one
	if !stakingEpochStore.Has(key) {
		stakingEpochStore.Set(key, k.cdc.MustMarshal(&amount))
	}

	// if store has key in it then add the amount to the previous value and put it back in store
	var newAmount sdk.Coin
	k.cdc.MustUnmarshal(stakingEpochStore.Get(key), &newAmount)
	newAmount.Add(amount)
	stakingEpochStore.Set(key, k.cdc.MustMarshal(&newAmount))
}

func (k Keeper) getAmountFromStakingEpoch(ctx sdk.Context, epochNumber int64) (amount sdk.Coin) {
	stakingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyStakingEpochStore)
	key := cosmosTypes.Int64Bytes(epochNumber)
	k.cdc.MustUnmarshal(stakingEpochStore.Get(key), &amount)
	return amount
}

func (k Keeper) deleteFromStakingEpoch(ctx sdk.Context, epochNumber int64) {
	stakingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyStakingEpochStore)
	key := cosmosTypes.Int64Bytes(epochNumber)
	stakingEpochStore.Delete(key)
}
