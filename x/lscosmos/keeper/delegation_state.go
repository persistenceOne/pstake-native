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

func (k Keeper) AddBalanceToDelegationState(ctx sdk.Context, coin sdk.Coin) {
	delegationState := k.GetDelegationState(ctx)
	delegationState.HostDelegationAccountBalance = delegationState.HostDelegationAccountBalance.Add(coin)
	k.SetDelegationState(ctx, delegationState)
}

func (k Keeper) RemoveBalanceFromDelegationState(ctx sdk.Context, coins sdk.Coins) {
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

func (k Keeper) AddHostAccountDelegation(ctx sdk.Context, delegation types.HostAccountDelegation) {
	delegationState := k.GetDelegationState(ctx)
	delegationState = appendHostAccountDelegation(delegationState, delegation)
	k.SetDelegationState(ctx, delegationState)
}

func (k Keeper) SubtractHostAccountDelegation(ctx sdk.Context, delegation types.HostAccountDelegation) error {
	delegationState := k.GetDelegationState(ctx)
	delegationState, err := removeHostAccountDelegation(delegationState, delegation)
	if err != nil {
		return err
	}
	k.SetDelegationState(ctx, delegationState)
	return nil
}

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

func (k Keeper) AddHostAccountUndelegation(ctx sdk.Context, undelegationEntry types.HostAccountUndelegation) {
	delegationState := k.GetDelegationState(ctx)
	delegationState.HostAccountUndelegations = append(delegationState.HostAccountUndelegations, undelegationEntry)
	k.SetDelegationState(ctx, delegationState)
}

func (k Keeper) AddTotalUndelegationForEpoch(ctx sdk.Context, epochNumber int64, amount sdk.Coin) {
	delegationState := k.GetDelegationState(ctx)
	found := false
	for i, undelegation := range delegationState.HostAccountUndelegations {
		if undelegation.EpochNumber == epochNumber && !found {
			delegationState.HostAccountUndelegations[i].TotalUndelegationAmount = delegationState.HostAccountUndelegations[i].TotalUndelegationAmount.Add(amount)
			found = true
		}
	}
	if !found {
		delegationState.HostAccountUndelegations = append(delegationState.HostAccountUndelegations, types.HostAccountUndelegation{
			EpochNumber:             epochNumber,
			TotalUndelegationAmount: amount,
			CompletionTime:          time.Time{},
			UndelegationEntries:     []types.UndelegationEntry{},
		})
	}
	k.SetDelegationState(ctx, delegationState)
}

func (k Keeper) AddEntriesForUndelegationEpoch(ctx sdk.Context, epochNumber int64, entries []types.UndelegationEntry) {
	delegationState := k.GetDelegationState(ctx)
	found := false
	for i, undelegation := range delegationState.HostAccountUndelegations {
		if undelegation.EpochNumber == epochNumber && !found {
			delegationState.HostAccountUndelegations[i].UndelegationEntries = append(delegationState.HostAccountUndelegations[i].UndelegationEntries, entries...)
			found = true
		}
	}
	if !found {
		panic("Adding Unbonding entries for non existing epoch")
	}
	k.SetDelegationState(ctx, delegationState)
}

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

func (k Keeper) GetHostAccountUndelegationForEpoch(ctx sdk.Context, epochNumber int64) (types.HostAccountUndelegation, error) {
	delegationState := k.GetDelegationState(ctx)
	for _, undelegation := range delegationState.HostAccountUndelegations {
		if undelegation.EpochNumber == epochNumber {
			return undelegation, nil
		}
	}
	return types.HostAccountUndelegation{}, types.ErrUndelegationEpochNotFound
}

func (k Keeper) GetHostAccountMaturedUndelegations(ctx sdk.Context) []types.HostAccountUndelegation {
	undelegations := k.GetDelegationState(ctx).HostAccountUndelegations
	var maturedUndelegations []types.HostAccountUndelegation
	for _, undelegation := range undelegations {
		if ctx.BlockTime().After(undelegation.CompletionTime) && !undelegation.CompletionTime.Equal(time.Time{}) {
			maturedUndelegations = append(maturedUndelegations, undelegation)
		}
	}
	return maturedUndelegations
}
