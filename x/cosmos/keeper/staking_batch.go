package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) addToStakingEpoch(ctx sdk.Context, keyAndValueForMinting cosmosTypes.KeyAndValueForMinting) {
	stakingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyStakingEpochStore)
	currentEpoch := k.epochsKeeper.GetEpochInfo(ctx, k.GetParams(ctx).StakingEpochIdentifier).CurrentEpoch
	storeKey := cosmosTypes.Int64Bytes(currentEpoch)
	if stakingEpochStore.Has(storeKey) {
		var stakingEpochValue cosmosTypes.StakingEpochValue
		err := k.cdc.Unmarshal(stakingEpochStore.Get(storeKey), &stakingEpochValue)
		if err != nil {
			panic(err)
		}
		stakingEpochValue.EpochMintingTxns = append(stakingEpochValue.EpochMintingTxns, keyAndValueForMinting)
		bz, err := k.cdc.Marshal(&stakingEpochValue)
		if err != nil {
			panic(err)
		}
		stakingEpochStore.Set(storeKey, bz)
		return
	}
	stakingEpochValue := cosmosTypes.NewStakingEpochValue(keyAndValueForMinting)
	bz, err := k.cdc.Marshal(&stakingEpochValue)
	if err != nil {
		panic(err)
	}
	stakingEpochStore.Set(storeKey, bz)
}

func (k Keeper) getFromStakingEpoch(ctx sdk.Context, epochNumber int64) (stakingEpochValue cosmosTypes.StakingEpochValue, err error) {
	stakingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyStakingEpochStore)
	key := cosmosTypes.Int64Bytes(epochNumber)
	err = k.cdc.Unmarshal(stakingEpochStore.Get(key), &stakingEpochValue)
	if err != nil {
		return stakingEpochValue, err
	}
	return stakingEpochValue, nil
}

func (k Keeper) deleteFromStakingEpoch(ctx sdk.Context, epochNumber int64) {
	stakingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyStakingEpochStore)
	key := cosmosTypes.Int64Bytes(epochNumber)
	stakingEpochStore.Delete(key)
}

func getTotalStakingAmount(stakingEpochValue cosmosTypes.StakingEpochValue, denom string) sdk.Coin {
	amount := sdk.NewInt64Coin(denom, 0)
	for _, st := range stakingEpochValue.EpochMintingTxns {
		amount = amount.Add(sdk.NewCoin(denom, st.Value.Amount.Amount))
	}
	return amount
}
