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
			ChainId:      chainID,
			EpochNumber:  epochNumber,
			MatureTime:   time.Time{},
			BurnAmount:   burnAmount,
			UnbondAmount: unbondAmount,
			Failed:       false,
		}
	} else {
		unbonding.UnbondAmount = unbonding.UnbondAmount.Add(unbondAmount)
		unbonding.BurnAmount = unbonding.BurnAmount.Add(burnAmount)
	}

	k.SetUnbonding(ctx, unbonding)
}
