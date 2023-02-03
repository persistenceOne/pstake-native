package keeper

import (
	"time"

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

// AddBalanceToDelegationState adds balance in the HostDelegationAccountBalance of types.DelegationState
func (k Keeper) AddBalanceToDelegationState(ctx sdk.Context, coin sdk.Coin) {
	delegationState := k.GetDelegationState(ctx)
	delegationState.HostDelegationAccountBalance = delegationState.HostDelegationAccountBalance.Add(coin)
	k.SetDelegationState(ctx, delegationState)
}

// RemoveBalanceFromDelegationState subtracts balance in the HostDelegationAccountBalance
// of types.DelegationState
func (k Keeper) RemoveBalanceFromDelegationState(ctx sdk.Context, coins sdk.Coins) {
	delegationState := k.GetDelegationState(ctx)
	delegationState.HostDelegationAccountBalance = delegationState.HostDelegationAccountBalance.Sub(coins)
	k.SetDelegationState(ctx, delegationState)
}

// SetHostChainDelegationAddress sets the host chain delegator address in types.DelegationState
func (k Keeper) SetHostChainDelegationAddress(ctx sdk.Context, addr string) error {
	delegationState := k.GetDelegationState(ctx)
	if delegationState.HostChainDelegationAddress != "" {
		return icatypes.ErrInterchainAccountAlreadySet
	}
	delegationState.HostChainDelegationAddress = addr
	k.SetDelegationState(ctx, delegationState)
	return nil
}

// AddHostAccountDelegation append the host account delegations in types.DelegationState provided
// in the input
func (k Keeper) AddHostAccountDelegation(ctx sdk.Context, delegation types.HostAccountDelegation) {
	delegationState := k.GetDelegationState(ctx)
	delegationState = appendHostAccountDelegation(delegationState, delegation)
	k.SetDelegationState(ctx, delegationState)
}

// SubtractHostAccountDelegation calls the removeHostAccountDelegation function to
// subtract host account delegations in types.DelegationState
func (k Keeper) SubtractHostAccountDelegation(ctx sdk.Context, delegation types.HostAccountDelegation) error {
	delegationState := k.GetDelegationState(ctx)
	delegationState, err := removeHostAccountDelegation(delegationState, delegation)
	if err != nil {
		return err
	}
	k.SetDelegationState(ctx, delegationState)
	return nil
}

// appendHostAccountDelegation is a helper function to append the input delegation to the
// input delegationState
func appendHostAccountDelegation(delegationState types.DelegationState, delegation types.HostAccountDelegation) types.DelegationState {
	// optimise this // do we want to have it sorted?
	for i, existingDelegation := range delegationState.HostAccountDelegations {
		if existingDelegation.ValidatorAddress == delegation.ValidatorAddress {
			delegationState.HostAccountDelegations[i].Amount = existingDelegation.Amount.Add(delegation.Amount)
			return delegationState
		}
	}

	delegationState.HostAccountDelegations = append(delegationState.HostAccountDelegations, delegation)
	return delegationState
}

// removeHostAccountDelegation is a helper function to remove the input delegation from the input
// delegationState
func removeHostAccountDelegation(delegationState types.DelegationState, delegation types.HostAccountDelegation) (types.DelegationState, error) {
	// optimise this // do we want to have it sorted?
	for i, existingDelegation := range delegationState.HostAccountDelegations {
		if existingDelegation.ValidatorAddress == delegation.ValidatorAddress {
			delegationState.HostAccountDelegations[i].Amount = existingDelegation.Amount.Sub(delegation.Amount) //This will panic if coin goes negative
			return delegationState, nil
		}
	}
	return types.DelegationState{}, types.ErrCannotRemoveNonExistentDelegation
}

// AddHostAccountUndelegation appends the input undelegationEntry in types.DelegationState
func (k Keeper) AddHostAccountUndelegation(ctx sdk.Context, undelegationEntry types.HostAccountUndelegation) {
	delegationState := k.GetDelegationState(ctx)
	delegationState.HostAccountUndelegations = append(delegationState.HostAccountUndelegations, undelegationEntry)
	k.SetDelegationState(ctx, delegationState)
}

// AddTotalUndelegationForEpoch adds the total undelegations corresponding to the input epoch number in
// types.DelegationState
func (k Keeper) AddTotalUndelegationForEpoch(ctx sdk.Context, epochNumber int64, amount sdk.Coin) {
	delegationState := k.GetDelegationState(ctx)
	found := false
	for i, undelegation := range delegationState.HostAccountUndelegations {
		if undelegation.EpochNumber == epochNumber && !found {
			delegationState.HostAccountUndelegations[i].TotalUndelegationAmount =
				delegationState.HostAccountUndelegations[i].TotalUndelegationAmount.Add(amount)
			found = true
		}
	}
	if !found {
		delegationState.HostAccountUndelegations = append(
			delegationState.HostAccountUndelegations,
			types.HostAccountUndelegation{
				EpochNumber:             epochNumber,
				TotalUndelegationAmount: amount,
				CompletionTime:          time.Time{},
				UndelegationEntries:     []types.UndelegationEntry{},
			})
	}
	k.SetDelegationState(ctx, delegationState)
}

// AddEntriesForUndelegationEpoch adds the input entries corresponding to the input epochNumber
// in types.DelegationState
func (k Keeper) AddEntriesForUndelegationEpoch(ctx sdk.Context, epochNumber int64, entries []types.UndelegationEntry) {
	delegationState := k.GetDelegationState(ctx)
	found := false
	for i, undelegation := range delegationState.HostAccountUndelegations {
		if undelegation.EpochNumber == epochNumber && !found {
			delegationState.HostAccountUndelegations[i].UndelegationEntries =
				append(delegationState.HostAccountUndelegations[i].UndelegationEntries, entries...)
			found = true
		}
	}
	if !found {
		panic("Adding Unbonding entries for non existing epoch")
	}
	k.SetDelegationState(ctx, delegationState)
}

// UpdateCompletionTimeForUndelegationEpoch updates the completion time for undelegation epoch
// corresponding to the input epoch number in types.DelegationState
func (k Keeper) UpdateCompletionTimeForUndelegationEpoch(ctx sdk.Context, epochNumber int64, completionTime time.Time) {
	delegationState := k.GetDelegationState(ctx)
	found := false
	for i, undelegation := range delegationState.HostAccountUndelegations {
		if undelegation.EpochNumber == epochNumber && !found {
			delegationState.HostAccountUndelegations[i].CompletionTime = completionTime
			found = true
		}
	}

	k.SetDelegationState(ctx, delegationState)
}

// RemoveHostAccountUndelegation removes the completion time for undelegation epoch
// corresponding to the input epoch number in types.DelegationState
func (k Keeper) RemoveHostAccountUndelegation(ctx sdk.Context, epochNumber int64) error {
	delegationState := k.GetDelegationState(ctx)
	for i, undelegation := range delegationState.HostAccountUndelegations {
		if undelegation.EpochNumber == epochNumber {
			delegationState.HostAccountUndelegations = append(delegationState.HostAccountUndelegations[:i], delegationState.HostAccountUndelegations[i+1:]...)
			k.SetDelegationState(ctx, delegationState)
			return nil
		}
	}
	return types.ErrCannotRemoveNonExistentUndelegation
}

// GetHostAccountUndelegationForEpoch returns the host account undelegation the input epoch number
func (k Keeper) GetHostAccountUndelegationForEpoch(ctx sdk.Context, epochNumber int64) (types.HostAccountUndelegation, error) {
	delegationState := k.GetDelegationState(ctx)
	for _, undelegation := range delegationState.HostAccountUndelegations {
		if undelegation.EpochNumber == epochNumber {
			return undelegation, nil
		}
	}
	return types.HostAccountUndelegation{}, types.ErrUndelegationEpochNotFound
}

// GetHostAccountMaturedUndelegations returns the host account matured undelegations
func (k Keeper) GetHostAccountMaturedUndelegations(ctx sdk.Context) []types.HostAccountUndelegation {
	undelegations := k.GetDelegationState(ctx).HostAccountUndelegations
	var maturedUndelegations []types.HostAccountUndelegation
	for _, undelegation := range undelegations {
		if !ctx.BlockTime().Before(undelegation.CompletionTime) && !undelegation.CompletionTime.Equal(time.Time{}) {
			maturedUndelegations = append(maturedUndelegations, undelegation)
		}
	}
	return maturedUndelegations
}
