package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetAllowListedValidators sets allowlisted validator set
func (k Keeper) SetAllowListedValidators(ctx sdk.Context, allowlistedValidators types.AllowListedValidators) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AllowListedValidatorsKey, k.cdc.MustMarshal(&allowlistedValidators))
}

// GetAllowListedValidators gets the allow listed validator set
func (k Keeper) GetAllowListedValidators(ctx sdk.Context) types.AllowListedValidators {
	store := ctx.KVStore(k.storeKey)
	var allowListedValidators types.AllowListedValidators
	k.cdc.MustUnmarshal(store.Get(types.AllowListedValidatorsKey), &allowListedValidators)

	return allowListedValidators
}

// GetAllValidatorsState returns the combined allowed listed validators set and combined
// delegation state. It is done to keep the old validators in the loop while calculating weighted amounts
// for delegation and undelegation
func (k Keeper) GetAllValidatorsState(ctx sdk.Context) (types.AllowListedValidators, types.DelegationState) {
	hostChainParams := k.GetHostChainParams(ctx)

	// Get current active val set and make a map of it
	currentAllowListedValSet := k.GetAllowListedValidators(ctx)
	currentAllowListedValSetMap := make(map[string]sdk.Dec)
	for _, val := range currentAllowListedValSet.AllowListedValidators {
		currentAllowListedValSetMap[val.ValidatorAddress] = val.TargetWeight
	}

	// get delegation state and make a map with address as
	currentDelegationState := k.GetDelegationState(ctx)
	currentDelegationStateMap := make(map[string]sdk.Coin)
	for _, delegation := range currentDelegationState.HostAccountDelegations {
		currentDelegationStateMap[delegation.ValidatorAddress] = delegation.Amount
	}

	// get validator list from allow listed validators
	delList := make([]string, len(currentAllowListedValSet.AllowListedValidators))
	for i, delegation := range currentAllowListedValSet.AllowListedValidators {
		delList[i] = delegation.ValidatorAddress
	}

	// get validator list from current delegation state
	valList := make([]string, len(currentDelegationState.HostAccountDelegations))
	for i, val := range currentDelegationState.HostAccountDelegations {
		valList[i] = val.ValidatorAddress
	}

	// Assign zero weights to all the validators not present in the current allow listed validator set map
	for _, val := range valList {
		_, ok := currentAllowListedValSetMap[val]
		if !ok {
			currentAllowListedValSetMap[val] = sdk.ZeroDec()
		}
	}

	// Convert the updated val set map to a slice of types.AllowListedValidator
	//nolint:prealloc,len_not_fixed
	var updatedValSet []types.AllowListedValidator
	for key, value := range currentAllowListedValSetMap {
		updatedValSet = append(updatedValSet, types.AllowListedValidator{ValidatorAddress: key, TargetWeight: value})
	}

	// Assign zero coin to all the validators not present in the current delegation state map
	for _, val := range delList {
		_, ok := currentDelegationStateMap[val]
		if !ok {
			currentDelegationStateMap[val] = sdk.NewCoin(hostChainParams.BaseDenom, sdk.ZeroInt())
		}
	}

	// Convert the updated delegation state map to slice of types.HostChainDelegation
	//nolint:prealloc,len_not_fixed
	var updatedDelegationState []types.HostAccountDelegation
	for key, value := range currentDelegationStateMap {
		updatedDelegationState = append(updatedDelegationState, types.HostAccountDelegation{ValidatorAddress: key, Amount: value})
	}

	// returns the two updated lists
	return types.AllowListedValidators{AllowListedValidators: updatedValSet}, types.DelegationState{HostAccountDelegations: updatedDelegationState}
}
