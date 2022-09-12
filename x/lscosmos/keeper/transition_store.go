package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetIBCTransientStore sets tokens that are in ibc transition
func (k Keeper) SetIBCTransientStore(ctx sdk.Context, ibcAmountTransientStore types.IBCAmountTransientStore) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.IBCTransitionStore, k.cdc.MustMarshal(&ibcAmountTransientStore))
}

// GetIBCTransientStore gets tokens that are in ibc transition
func (k Keeper) GetIBCTransientStore(ctx sdk.Context) types.IBCAmountTransientStore {
	store := ctx.KVStore(k.storeKey)
	var ibcAmountTransientStore types.IBCAmountTransientStore
	k.cdc.MustUnmarshal(store.Get(types.IBCTransitionStore), &ibcAmountTransientStore)

	return ibcAmountTransientStore
}

// AddIBCTransferToTransientStore adds ibctransfer tokens that are in ibc transition
// CONTRACT: to be used atomically with IBCTransfer of tokens from delegation account to it's host counterpart
func (k Keeper) AddIBCTransferToTransientStore(ctx sdk.Context, amount sdk.Coin) {
	transientStore := k.GetIBCTransientStore(ctx)
	transientStore.IBCTransfer = transientStore.IBCTransfer.Add(amount)
	k.SetIBCTransientStore(ctx, transientStore)
}

// RemoveIBCTransferFromTransientStore removes ibctransfer tokens that are in ibc transition
// CONTRACT: to be used atomically with AddBalanceToDelegationState
func (k Keeper) RemoveIBCTransferFromTransientStore(ctx sdk.Context, amount sdk.Coin) {
	transientStore := k.GetIBCTransientStore(ctx)
	transientStore.IBCTransfer = transientStore.IBCTransfer.Sub(sdk.NewCoins(amount))
	k.SetIBCTransientStore(ctx, transientStore)
}

// AddICADelegateToTransientStore adds ibctransfer tokens that are in ibc transition
// CONTRACT: to be used atomically with RemoveBalanceFromDelegationState
func (k Keeper) AddICADelegateToTransientStore(ctx sdk.Context, amount sdk.Coin) {
	transientStore := k.GetIBCTransientStore(ctx)
	transientStore.ICADelegate = transientStore.ICADelegate.Add(amount)
	k.SetIBCTransientStore(ctx, transientStore)
}

// RemoveICADelegateFromTransientStore removes ibctransfer tokens that are in ibc transition
// Contract: to be used atomically with AddHostAccountDelegation
func (k Keeper) RemoveICADelegateFromTransientStore(ctx sdk.Context, amount sdk.Coin) {
	transientStore := k.GetIBCTransientStore(ctx)
	transientStore.ICADelegate = transientStore.ICADelegate.Sub(amount)
	k.SetIBCTransientStore(ctx, transientStore)
}
