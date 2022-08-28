package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetDelegationState sets the delegation state in store
func (k Keeper) SetDelegationState(ctx sdk.Context, delegationState types.DelegationState) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DelegationStateKey, k.cdc.MustMarshal(&delegationState))
}

// GetDelegationState gets the delegation state in store
func (k Keeper) GetDelegationState(ctx sdk.Context) types.DelegationState {
	store := ctx.KVStore(k.storeKey)

	var delegationState types.DelegationState
	k.cdc.MustUnmarshal(store.Get(types.DelegationStateKey), &delegationState)

	return delegationState
}

func (k Keeper) AddBalanceToDelegationState(ctx sdk.Context, coin sdk.Coin) {
	delegationState := k.GetDelegationState(ctx)
	delegationState.HostDelegationAccountBalance = delegationState.HostDelegationAccountBalance.Add(coin)
	k.SetDelegationState(ctx, delegationState)
}

func (k Keeper) RemoveBalanceToDelegationState(ctx sdk.Context, coins sdk.Coins) {
	delegationState := k.GetDelegationState(ctx)
	delegationState.HostDelegationAccountBalance = delegationState.HostDelegationAccountBalance.Sub(coins)
	k.SetDelegationState(ctx, delegationState)
}

func (k Keeper) SetHostChainDelegationAddress(ctx sdk.Context, addr string) error {
	delegationState := k.GetDelegationState(ctx)
	if delegationState.HostChainDelegationAddress != "" {
		return icatypes.ErrInterchainAccountAlreadySet
	}
	delegationState.HostChainDelegationAddress = addr
	k.SetDelegationState(ctx, delegationState)
	return nil
}

func (k Keeper) AddHostChainDelegation(ctx sdk.Context, delegation types.HostAccountDelegation) {
	delegationState := k.GetDelegationState(ctx)
	delegationState = appendHostChainDelegation(delegationState, delegation)
	k.SetDelegationState(ctx, delegationState)
}
func (k Keeper) SubtractHostChainDelegation(ctx sdk.Context, delegation types.HostAccountDelegation) {
	delegationState := k.GetDelegationState(ctx)
	delegationState = removeHostChainDelegation(delegationState, delegation)
	k.SetDelegationState(ctx, delegationState)
}
func appendHostChainDelegation(delegationState types.DelegationState, delegation types.HostAccountDelegation) types.DelegationState {
	// optimise this // do we want to have it sorted?
	for _, existingDelegation := range delegationState.HostAccountDelegations {
		if existingDelegation.ValidatorAddress == delegation.ValidatorAddress {
			existingDelegation.Amount = existingDelegation.Amount.Add(delegation.Amount)
			return delegationState
		}
	}

	delegationState.HostAccountDelegations = append(delegationState.HostAccountDelegations, delegation)
	return delegationState
}
func removeHostChainDelegation(delegationState types.DelegationState, delegation types.HostAccountDelegation) types.DelegationState {
	// optimise this // do we want to have it sorted?
	for _, existingDelegation := range delegationState.HostAccountDelegations {
		if existingDelegation.ValidatorAddress == delegation.ValidatorAddress {
			existingDelegation.Amount = existingDelegation.Amount.Sub(delegation.Amount)
			if existingDelegation.Amount.IsNegative() {
				panic("WTF?")
			}
			return delegationState
		}
	}
	panic("WTF?")
}
