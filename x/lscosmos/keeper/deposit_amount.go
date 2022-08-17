package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetDepositAmount sets the deposit amount in store
func (k Keeper) SetDepositAmount(ctx sdk.Context, amount types.DepositAmount) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DepositAmountKey, k.cdc.MustMarshal(&amount))
}

// GetDepositAmount gets the deposit amount in store
func (k Keeper) GetDepositAmount(ctx sdk.Context) types.DepositAmount {
	store := ctx.KVStore(k.storeKey)
	var depositAmount types.DepositAmount
	k.cdc.MustUnmarshal(store.Get(types.DepositAmountKey), &depositAmount)

	return depositAmount
}
