package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetDelegatorUnbondingEpochEntry sets delegator entry for unbondign stkatom for an unbonding epoch
func (k Keeper) SetDelegatorUnbondingEpochEntry(ctx sdk.Context, unbondingEpochEntry types.DelegatorUnbondingEpochEntry) {
	store := ctx.KVStore(k.storeKey)
	delAddr, err := sdk.AccAddressFromBech32(unbondingEpochEntry.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	bz := k.cdc.MustMarshal(&unbondingEpochEntry)
	store.Set(types.GetDelegatorUnbondingEpochEntryKey(delAddr, unbondingEpochEntry.EpochNumber), bz)
}

// GetDelegatorUnbondingEpochEntry gets delegator entry for unbondign stkatom for an unbonding epoch
func (k Keeper) GetDelegatorUnbondingEpochEntry(ctx sdk.Context, delegatorAddress sdk.AccAddress, epochNumber int64) types.DelegatorUnbondingEpochEntry {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDelegatorUnbondingEpochEntryKey(delegatorAddress, epochNumber))
	var unbondingEpochEntry types.DelegatorUnbondingEpochEntry
	k.cdc.MustUnmarshal(bz, &unbondingEpochEntry)

	return unbondingEpochEntry
}
