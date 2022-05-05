package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

//TODO : Use this once module is enabled to set multisig account
func (k Keeper) setAccountState(ctx sdk.Context, acc authTypes.AccountI) {
	addr := acc.GetAddress()
	store := ctx.KVStore(k.storeKey)

	bz, err := k.authKeeper.MarshalAccount(acc)
	if err != nil {
		panic(err)
	}

	store.Set(cosmosTypes.MultisigAccountStoreKey(addr), bz)
}

func (k Keeper) getAccountState(ctx sdk.Context, accAddress sdk.AccAddress) authTypes.AccountI {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(cosmosTypes.MultisigAccountStoreKey(accAddress))
	if bz == nil {
		return nil
	}

	acc, err := k.authKeeper.UnmarshalAccount(bz)
	if err != nil {
		panic(err)
	}

	return acc
}

func (k Keeper) getCurrentAddress(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get(cosmosTypes.CurrentMultisigAddressKey())

}

func (k Keeper) setCurrentAddress(ctx sdk.Context, accAddress sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(cosmosTypes.CurrentMultisigAddressKey(), accAddress)
}
