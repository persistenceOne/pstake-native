package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

//TODO : Use this once module is enabled to set multisig account
func (k Keeper) setAccountState(ctx sdk.Context, acc authTypes.AccountI) {
	addr := acc.GetAddress()
	store := ctx.KVStore(k.storeKey)

	bz, err := k.authKeeper.MarshalAccount(acc)
	if err != nil {
		panic(err)
	}

	store.Set(authTypes.AddressStoreKey(addr), bz)
}

func (k Keeper) getAccountState(ctx sdk.Context, accAddress sdk.AccAddress) authTypes.AccountI {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(authTypes.AddressStoreKey(accAddress))
	if bz == nil {
		return nil
	}

	acc, err := k.authKeeper.UnmarshalAccount(bz)
	if err != nil {
		panic(err)
	}

	return acc
}
