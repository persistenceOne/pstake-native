package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// SetIBCTransientStore sets tokens that are in ibc transition
func (k Keeper) SetIBCTransientStore(ctx sdk.Context, ibcAmountTransientStore types.IBCAmountTransientStore) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.IBCTransientStoreKey, k.cdc.MustMarshal(&ibcAmountTransientStore))
}

// GetIBCTransientStore gets tokens that are in ibc transition
func (k Keeper) GetIBCTransientStore(ctx sdk.Context) types.IBCAmountTransientStore {
	store := ctx.KVStore(k.storeKey)
	var ibcAmountTransientStore types.IBCAmountTransientStore
	k.cdc.MustUnmarshal(store.Get(types.IBCTransientStoreKey), &ibcAmountTransientStore)

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
	transientStore.IBCTransfer = transientStore.IBCTransfer.Sub(sdk.NewCoins(amount)...)
	k.SetIBCTransientStore(ctx, transientStore)
}

// AddICADelegateToTransientStore adds ibctransfer tokens that are in ibc transition
// CONTRACT: to be used atomically with RemoveBalanceFromDelegationState
func (k Keeper) AddICADelegateToTransientStore(ctx sdk.Context, amount sdk.Coin) {
	transientStore := k.GetIBCTransientStore(ctx)
	if !transientStore.ICADelegate.IsValid() || transientStore.ICADelegate.IsZero() { // only because initialised coin is invalid
		transientStore.ICADelegate = amount
	} else {
		transientStore.ICADelegate = transientStore.ICADelegate.Add(amount)
	}
	k.SetIBCTransientStore(ctx, transientStore)
}

// RemoveICADelegateFromTransientStore removes ibctransfer tokens that are in ibc transition
// Contract: to be used atomically with AddHostAccountDelegation and AddBalanceToDelegationState(incase of failed txns)
func (k Keeper) RemoveICADelegateFromTransientStore(ctx sdk.Context, amount sdk.Coin) {
	transientStore := k.GetIBCTransientStore(ctx)
	transientStore.ICADelegate = transientStore.ICADelegate.Sub(amount)
	k.SetIBCTransientStore(ctx, transientStore)
}

// AddUndelegationTransferToTransientStore adds ibctransfer tokens that are in ibc transition from host chain to controller chain
// CONTRACT: to be used atomically with RemoveHostAccountUndelegation (after successful undelegations)
func (k Keeper) AddUndelegationTransferToTransientStore(ctx sdk.Context, undelegationTransfer types.TransientUndelegationTransfer) {
	transientStore := k.GetIBCTransientStore(ctx)
	transientStore.UndelegatonCompleteIBCTransfer = append(transientStore.UndelegatonCompleteIBCTransfer, undelegationTransfer)
	k.SetIBCTransientStore(ctx, transientStore)
}

// RemoveUndelegationTransferFromTransientStore removes ibctransfer tokens that are in ibc transition from host chain to controller chain
// Contract: to be used atomically with MatureUnbondingEpochCValue (after successful undelegations) and AddHostAccountUndelegation ( after failed ICA+IBC txn - matured undelegations)
func (k Keeper) RemoveUndelegationTransferFromTransientStore(ctx sdk.Context, amount sdk.Coin) (types.TransientUndelegationTransfer, error) {
	transientStore := k.GetIBCTransientStore(ctx)
	for i, undelegationTransfer := range transientStore.UndelegatonCompleteIBCTransfer {
		if undelegationTransfer.AmountUnbonded.IsEqual(amount) {
			transientStore.UndelegatonCompleteIBCTransfer = append(transientStore.UndelegatonCompleteIBCTransfer[:i], transientStore.UndelegatonCompleteIBCTransfer[i+1:]...)
			k.SetIBCTransientStore(ctx, transientStore)
			return undelegationTransfer, nil
		}
	}
	return types.TransientUndelegationTransfer{}, types.ErrTransientUndelegationTransferNotFound
}
