package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
)

func (k *Keeper) SetUserUnbonding(ctx sdk.Context, ub *types.UserUnbonding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UserUnbondingKey)
	bytes := k.cdc.MustMarshal(ub)
	store.Set(types.GetUserUnbondingStoreKey(ub.ChainId, ub.Address, ub.EpochNumber), bytes)
}

func (k *Keeper) GetUserUnbonding(
	ctx sdk.Context,
	chainID string,
	delegatorAddress string,
	epochNumber int64,
) (*types.UserUnbonding, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UserUnbondingKey)
	bz := store.Get(types.GetUserUnbondingStoreKey(chainID, delegatorAddress, epochNumber))
	if bz == nil {
		return &types.UserUnbonding{}, false
	}

	var userUnbonding types.UserUnbonding
	k.cdc.MustUnmarshal(bz, &userUnbonding)
	return &userUnbonding, true
}

func (k *Keeper) DeleteUserUnbonding(ctx sdk.Context, ub *types.UserUnbonding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UserUnbondingKey)
	store.Delete(types.GetUserUnbondingStoreKey(ub.ChainId, ub.Address, ub.EpochNumber))
}

func (k *Keeper) FilterUserUnbondings(ctx sdk.Context, filter func(u types.UserUnbonding) bool) []*types.UserUnbonding {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UserUnbondingKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	userUnbondings := make([]*types.UserUnbonding, 0)
	for ; iterator.Valid(); iterator.Next() {
		userUnbonding := types.UserUnbonding{}
		k.cdc.MustUnmarshal(iterator.Value(), &userUnbonding)
		if filter(userUnbonding) {
			userUnbondings = append(userUnbondings, &userUnbonding)
		}
	}

	return userUnbondings
}

func (k *Keeper) IncreaseUserUnbondingAmountForEpoch(
	ctx sdk.Context,
	chainID string,
	delegatorAddress string,
	epochNumber int64,
	stkAmount sdk.Coin,
	unbondAmount sdk.Coin,
) {
	userUnbonding, found := k.GetUserUnbonding(ctx, chainID, delegatorAddress, epochNumber)
	if !found {
		userUnbonding = &types.UserUnbonding{
			ChainId:      chainID,
			EpochNumber:  epochNumber,
			Address:      delegatorAddress,
			StkAmount:    stkAmount,
			UnbondAmount: unbondAmount,
		}
	} else {
		userUnbonding.StkAmount = userUnbonding.StkAmount.Add(stkAmount)
		userUnbonding.UnbondAmount = userUnbonding.UnbondAmount.Add(unbondAmount)
	}

	k.SetUserUnbonding(ctx, userUnbonding)
}
