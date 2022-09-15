package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetUnbondingEpochCValue sets cvalue for unbonding epoch
func (k Keeper) SetUnbondingEpochEntry(ctx sdk.Context, unbondingEpochEntry types.UnbondingEpochEntry) {
	store := ctx.KVStore(k.storeKey)
	delAddr, err := sdk.AccAddressFromBech32(unbondingEpochEntry.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	bz := k.cdc.MustMarshal(&unbondingEpochEntry)
	store.Set(types.GetUnbondingEpochEntryKey(unbondingEpochEntry.EpochNumber, delAddr), bz)
}

// GetUnbondingEpochCValue sets cvalue for unbonding epoch
func (k Keeper) GetUnbondingEpochEntry(ctx sdk.Context, epochNumber int64, delegatorAddress sdk.AccAddress) types.UnbondingEpochEntry {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUnbondingEpochEntryKey(epochNumber, delegatorAddress))
	var unbondingEpochEntry types.UnbondingEpochEntry
	k.cdc.MustUnmarshal(bz, &unbondingEpochEntry)

	return unbondingEpochEntry
}
