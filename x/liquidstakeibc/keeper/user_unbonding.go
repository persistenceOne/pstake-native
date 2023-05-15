package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
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

func (k *Keeper) IncreaseUserUndelegatingAmountForEpoch(
	ctx sdk.Context,
	chainID string,
	delegatorAddress string,
	epochNumber int64,
	amount sdk.Coin,
) {
	userUnbonding, found := k.GetUserUnbonding(ctx, chainID, delegatorAddress, epochNumber)
	if !found {
		userUnbonding = &types.UserUnbonding{
			ChainId:     chainID,
			EpochNumber: epochNumber,
			Address:     delegatorAddress,
			Amount:      amount,
		}
	} else {
		userUnbonding.Amount = userUnbonding.Amount.Add(amount)
	}

	k.SetUserUnbonding(ctx, userUnbonding)
}
