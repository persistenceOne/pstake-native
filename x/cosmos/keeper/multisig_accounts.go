package keeper

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// sets account state from the given Account Interface
func (k Keeper) setAccountState(ctx sdk.Context, acc authTypes.AccountI) {
	addr, _ := sdk.AccAddressFromHex(acc.GetPubKey().Address().String()) // todo : remove this
	store := ctx.KVStore(k.storeKey)

	bz, err := k.authKeeper.MarshalAccount(acc)
	if err != nil {
		panic(err)
	}

	store.Set(cosmosTypes.MultisigAccountStoreKey(addr), bz)
}

// gets account state of the given account address
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

// Gets the current multisig address
func (k Keeper) getCurrentAddress(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get(cosmosTypes.CurrentMultisigAddressKey())

}

// Sets a new given multsig address
func (k Keeper) setCurrentAddress(ctx sdk.Context, accAddress sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(cosmosTypes.CurrentMultisigAddressKey(), accAddress)
}

// Checks if the orchestrator address is present in the current multisig address or not
func (k Keeper) checkOrchestratorAddressPresentInMultisig(ctx sdk.Context, orch sdk.AccAddress) bool {
	// fetch orch address pub key on chain
	orchPubKey := k.authKeeper.GetAccount(ctx, orch).GetPubKey()
	if orchPubKey == nil {
		panic("pub key for orch address not found")
	}

	// fetch multisig pub key
	multsigPubKey := k.getAccountState(ctx, k.getCurrentAddress(ctx)).GetPubKey().(*multisig.LegacyAminoPubKey).GetPubKeys()

	for _, pb := range multsigPubKey {
		if pb.Equals(orchPubKey) {
			return true
		}
	}
	return false
}
