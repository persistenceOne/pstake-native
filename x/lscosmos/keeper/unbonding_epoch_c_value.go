package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetUnbondingEpochCValue sets cvalue for unbonding epoch
func (k Keeper) SetUnbondingEpochCValue(ctx sdk.Context, unbondingEpochCValue types.UnbondingEpochCValue) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&unbondingEpochCValue)
	store.Set(types.GetUnbondingEpochCValueKey(unbondingEpochCValue.EpochNumber), bz)
}

// GetUnbondingEpochCValue sets cvalue for unbonding epoch
func (k Keeper) GetUnbondingEpochCValue(ctx sdk.Context, epochNumber int64) types.UnbondingEpochCValue {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUnbondingEpochCValueKey(epochNumber))
	var unbondingEpochCValue types.UnbondingEpochCValue
	k.cdc.MustUnmarshal(bz, &unbondingEpochCValue)
	return unbondingEpochCValue
}

// MatureUnbondingEpochCValue sets unbonding epochCValue as matured
func (k Keeper) MatureUnbondingEpochCValue(ctx sdk.Context, epochNumber int64) {
	unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, epochNumber)
	unbondingEpochCValue.IsMatured = true
	k.SetUnbondingEpochCValue(ctx, unbondingEpochCValue)
}

// FailUnbondingEpochCValue sets unbonding epochCValue as timeout for undelegation
func (k Keeper) FailUnbondingEpochCValue(ctx sdk.Context, epochNumber int64) {
	unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, epochNumber)
	unbondingEpochCValue.IsFailed = true
	k.SetUnbondingEpochCValue(ctx, unbondingEpochCValue)
}
