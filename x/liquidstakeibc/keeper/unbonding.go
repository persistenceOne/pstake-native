package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) SetUnbonding(ctx sdk.Context, ub *types.Unbonding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingKey)
	bytes := k.cdc.MustMarshal(ub)
	store.Set(types.GetUnbondingStoreKey(ub.ChainId, ub.EpochNumber), bytes)
}

func (k *Keeper) GetUnbonding(ctx sdk.Context, chainID string, epochNumber int64) (*types.Unbonding, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingKey)
	bz := store.Get(types.GetUnbondingStoreKey(chainID, epochNumber))
	if bz == nil {
		return nil, false
	}

	var unbonding types.Unbonding
	k.cdc.MustUnmarshal(bz, &unbonding)
	return &unbonding, true
}

func (k *Keeper) DeleteUnbonding(ctx sdk.Context, ub *types.Unbonding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingKey)
	store.Delete(types.GetUnbondingStoreKey(ub.ChainId, ub.EpochNumber))
}

func (k *Keeper) FilterUnbondings(ctx sdk.Context, filter func(u types.Unbonding) bool) []*types.Unbonding {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	unbondings := make([]*types.Unbonding, 0)
	for ; iterator.Valid(); iterator.Next() {
		unbonding := types.Unbonding{}
		k.cdc.MustUnmarshal(iterator.Value(), &unbonding)
		if filter(unbonding) {
			unbondings = append(unbondings, &unbonding)
		}
	}

	return unbondings
}

func (k *Keeper) IncreaseUndelegatingAmountForEpoch(
	ctx sdk.Context,
	chainID string,
	epochNumber int64,
	burnAmount sdk.Coin,
	unbondAmount sdk.Coin,
) {
	unbonding, found := k.GetUnbonding(ctx, chainID, epochNumber)
	if !found {
		unbonding = &types.Unbonding{
			ChainId:       chainID,
			EpochNumber:   epochNumber,
			MatureTime:    time.Time{},
			BurnAmount:    burnAmount,
			UnbondAmount:  unbondAmount,
			IbcSequenceId: "",
			State:         types.Unbonding_UNBONDING_PENDING,
		}
	} else {
		unbonding.UnbondAmount = unbonding.UnbondAmount.Add(unbondAmount)
		unbonding.BurnAmount = unbonding.BurnAmount.Add(burnAmount)
	}

	k.SetUnbonding(ctx, unbonding)
}

func (k *Keeper) FailAllUnbondingsForSequenceID(ctx sdk.Context, sequenceID string) {
	unbondings := k.FilterUnbondings(ctx, func(u types.Unbonding) bool { return u.IbcSequenceId == sequenceID })

	for _, unbonding := range unbondings {
		unbonding.IbcSequenceId = ""
		unbonding.State = types.Unbonding_UNBONDING_FAILED
		k.SetUnbonding(ctx, unbonding)
	}
}

func (k *Keeper) RevertUnbondingsState(ctx sdk.Context, unbondings []*types.Unbonding) {
	for _, unbonding := range unbondings {
		unbonding.IbcSequenceId = ""

		if unbonding.State != types.Unbonding_UNBONDING_PENDING &&
			unbonding.State != types.Unbonding_UNBONDING_FAILED {
			unbonding.State--
		}

		k.SetUnbonding(ctx, unbonding)
	}
}
