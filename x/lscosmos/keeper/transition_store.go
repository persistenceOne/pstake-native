package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetIBCTransitionStore sets tokens that are in ibc transition
func (k Keeper) SetIBCTransitionStore(ctx sdk.Context, ibcAmountTransitionStore types.IbcAmountTransitionStore) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.IBCTransitionStore, k.cdc.MustMarshal(&ibcAmountTransitionStore))
}

// GetIBCTransitionStore gets tokens that are in ibc transition
func (k Keeper) GetIBCTransitionStore(ctx sdk.Context) types.IbcAmountTransitionStore {
	store := ctx.KVStore(k.storeKey)
	var ibcAmountTransitionStore types.IbcAmountTransitionStore
	k.cdc.MustUnmarshal(store.Get(types.IBCTransitionStore), &ibcAmountTransitionStore)

	return ibcAmountTransitionStore
}

// AddIBCTransferToTransitionStore adds ibctransfer tokens that are in ibc transition
// CONTRACT: to be used atomically with IBCTransfer of tokens from delegation account to it's host counterpart
func (k Keeper) AddIBCTransferToTransitionStore(ctx sdk.Context, amount sdk.Coin) {
	transitionStore := k.GetIBCTransitionStore(ctx)
	transitionStore.IbcTransfer = transitionStore.IbcTransfer.Add(amount)
	k.SetIBCTransitionStore(ctx, transitionStore)
}

// RemoveIBCTransferFromTransitionStore removes ibctransfer tokens that are in ibc transition
// CONTRACT: to be used atomically with AddBalanceToDelegationState
func (k Keeper) RemoveIBCTransferFromTransitionStore(ctx sdk.Context, amount sdk.Coin) {
	transitionStore := k.GetIBCTransitionStore(ctx)
	transitionStore.IbcTransfer = transitionStore.IbcTransfer.Sub(sdk.NewCoins(amount))
	k.SetIBCTransitionStore(ctx, transitionStore)
}

// AddICADelegateToTransitionStore adds ibctransfer tokens that are in ibc transition
// CONTRACT: to be used atomically with RemoveBalanceFromDelegationState
func (k Keeper) AddICADelegateToTransitionStore(ctx sdk.Context, amount sdk.Coin) {
	transitionStore := k.GetIBCTransitionStore(ctx)
	transitionStore.IcaDelegate = transitionStore.IcaDelegate.Add(amount)
	k.SetIBCTransitionStore(ctx, transitionStore)
}

// RemoveICADelegateFromTransitionStore removes ibctransfer tokens that are in ibc transition
// Contract: to be used atomically with AddHostAccountDelegation
func (k Keeper) RemoveICADelegateFromTransitionStore(ctx sdk.Context, amount sdk.Coin) {
	transitionStore := k.GetIBCTransitionStore(ctx)
	transitionStore.IcaDelegate = transitionStore.IcaDelegate.Sub(amount)
	k.SetIBCTransitionStore(ctx, transitionStore)
}
