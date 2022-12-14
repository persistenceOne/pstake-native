package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// DelegateMsgs gives the list of Delegate Txs to be executed based on the current state and params.
// CONTRACT: allowlistedValList.len > 0, amount > 0
func (k Keeper) DelegateMsgs(ctx sdk.Context, amount sdk.Int, denom string, delegationState types.DelegationState) ([]sdk.Msg, error) {
	// fetch a combined updated val set list and delegation state
	updateValList, hostAccountDelegations := k.GetAllValidatorsState(ctx, denom)

	// assign the updated validator delegation state to the current delegation state
	delegationState.HostAccountDelegations = hostAccountDelegations

	updatedAllowListedValidators := types.AllowListedValidators{AllowListedValidators: updateValList}

	valAddressAmount, err := FetchValidatorsToDelegate(updatedAllowListedValidators, delegationState, sdk.NewCoin(denom, amount))
	if err != nil {
		return nil, err
	}

	var msgs []sdk.Msg
	for _, val := range valAddressAmount {
		if val.Amount.IsPositive() {
			msg := &stakingtypes.MsgDelegate{
				DelegatorAddress: delegationState.HostChainDelegationAddress,
				ValidatorAddress: val.ValidatorAddr,
				Amount:           val.Amount,
			}
			msgs = append(msgs, msg)
		}
	}

	if len(msgs) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsgs, "No msgs to delegate")
	}

	return msgs, nil
}

// UndelegateMsgs gives the list of Undelegate Txs to be executed based on the current state and params.
// CONTRACT: allowlistedValList.len > 0, amount > 0
func (k Keeper) UndelegateMsgs(ctx sdk.Context, amount sdk.Int, denom string, delegationState types.DelegationState) ([]sdk.Msg, []types.UndelegationEntry, error) {
	// fetch a combined updated val set list and delegation state
	updateValList, hostAccountDelegations := k.GetAllValidatorsState(ctx, denom)

	// assign the updated validator delegation state to the current delegation state
	delegationState.HostAccountDelegations = hostAccountDelegations

	updatedAllowListedValidators := types.AllowListedValidators{AllowListedValidators: updateValList}

	valAddressAmount, err := FetchValidatorsToUndelegate(updatedAllowListedValidators, delegationState, sdk.NewCoin(denom, amount))
	if err != nil {
		return nil, nil, err
	}

	//nolint:prealloc,len_not_fixed
	var msgs []sdk.Msg
	//nolint:prealloc,len_not_fixed
	var undelegationEntries []types.UndelegationEntry
	for _, val := range valAddressAmount {

		msg := &stakingtypes.MsgUndelegate{
			DelegatorAddress: delegationState.HostChainDelegationAddress,
			ValidatorAddress: val.ValidatorAddr,
			Amount:           val.Amount,
		}
		msgs = append(msgs, msg)

		undelegationEntry := types.UndelegationEntry{
			ValidatorAddress: val.ValidatorAddr,
			Amount:           val.Amount,
		}
		undelegationEntries = append(undelegationEntries, undelegationEntry)
	}

	// should never come ideally
	if len(msgs) == 0 || len(undelegationEntries) == 0 {
		return nil, nil, sdkerrors.Wrap(types.ErrInvalidMsgs, "No msgs to undelegate")
	}

	return msgs, undelegationEntries, nil
}

// FetchValidatorsToDelegate gives a list of all validators having weighted amount for few and 1uatom for rest in order to auto claim all rewards accumulated in current epoch
func FetchValidatorsToDelegate(valList types.AllowListedValidators, delegationState types.DelegationState, amount sdk.Coin) (types.ValAddressAmounts, error) {
	curDiffDistribution := GetIdealCurrentDelegations(valList, delegationState, amount, false)
	sort.Sort(sort.Reverse(curDiffDistribution))

	return DivideAmountIntoValidatorSet(curDiffDistribution, amount)
}

// FetchValidatorsToUndelegate gives a list of all validators having weighted amount for few and 1uatom for rest in order to auto claim all rewards accumulated in current epoch
func FetchValidatorsToUndelegate(valList types.AllowListedValidators, delegationState types.DelegationState, amount sdk.Coin) (types.ValAddressAmounts, error) {
	currDiffDistribution := GetIdealCurrentDelegations(valList, delegationState, amount, true)
	sort.Sort(sort.Reverse(currDiffDistribution))
	return DivideUndelegateAmountIntoValidatorSet(currDiffDistribution, amount)
}

// GetIdealCurrentDelegations returns ideal amount of delegations to validators on host chain
func GetIdealCurrentDelegations(valList types.AllowListedValidators, delegationState types.DelegationState, amt sdk.Coin, reverse bool) types.WeightedAddressAmounts {
	totalDelegations := delegationState.TotalDelegations(amt.Denom)

	curDiffDistribution := types.WeightedAddressAmounts{}
	delegationMap := types.GetHostAccountDelegationMap(delegationState.HostAccountDelegations)
	var idealTokens, curTokens sdk.Int

	for _, valT := range valList.AllowListedValidators {
		totalAmt := totalDelegations.Amount.Add(amt.Amount)
		if reverse {
			totalAmt = totalDelegations.Amount.Sub(amt.Amount)
		}
		idealTokens = valT.TargetWeight.Mul(sdk.NewDecFromInt(totalAmt)).TruncateInt()
		curCoins, ok := delegationMap[valT.ValidatorAddress]
		if !ok {
			curCoins = sdk.NewCoin(amt.Denom, sdk.ZeroInt())
		}
		curTokens = curCoins.Amount
		diffAmt := idealTokens.Sub(curTokens)
		if reverse {
			diffAmt = curTokens.Sub(idealTokens)
		}
		curDiffDistribution = append(curDiffDistribution, types.WeightedAddressAmount{
			Address: valT.ValidatorAddress,
			Weight:  valT.TargetWeight,
			Denom:   amt.Denom,
			Amount:  diffAmt,
		})
	}

	return curDiffDistribution
}

// divideAmountWeightedSet : divides amount to be delegated or undelegated w.r.t weights.
func divideAmountWeightedSet(valAmounts types.ValAddressAmounts, coin sdk.Coin, valAddressWeightMap map[string]sdk.Dec) types.ValAddressAmounts {
	var newValAmounts types.ValAddressAmounts

	totalWeight := sdk.ZeroDec()
	for _, weight := range valAddressWeightMap {
		totalWeight = totalWeight.Add(weight)
	}

	for _, valAmt := range valAmounts {
		weight := valAddressWeightMap[valAmt.ValidatorAddr].Quo(totalWeight)
		amt := weight.MulInt(coin.Amount).RoundInt()
		newValAmounts = append(newValAmounts, types.ValAddressAmount{
			ValidatorAddr: valAmt.ValidatorAddr,
			Amount:        sdk.NewCoin(valAmt.Amount.Denom, valAmt.Amount.Amount.Add(amt)),
		})
	}
	return newValAmounts
}

// distributeCoinsAmongstValSet takes the validator distribution and coins to distribute and returns the
// validator address amount to distribute and the remaining amount to make
func distributeCoinsAmongstValSet(ws types.WeightedAddressAmounts, coin sdk.Coin) (types.ValAddressAmounts, sdk.Coin) {
	var valAddrAmts types.ValAddressAmounts

	for _, w := range ws {
		if coin.Amount.LTE(w.Amount) {
			valAddrAmts = append(valAddrAmts, types.ValAddressAmount{ValidatorAddr: w.Address, Amount: coin})
			return valAddrAmts, sdk.NewInt64Coin(coin.Denom, 0)
		}
		valAddrAmts = append(valAddrAmts, types.ValAddressAmount{ValidatorAddr: w.Address, Amount: w.Coin()})
		coin = coin.SubAmount(w.Amount)
	}

	return valAddrAmts, coin
}

// DivideAmountIntoValidatorSet : divides amount into validator set
func DivideAmountIntoValidatorSet(sortedValDiff types.WeightedAddressAmounts, coin sdk.Coin) (types.ValAddressAmounts, error) {
	if coin.IsZero() {
		return nil, nil
	}

	// Delegate to non-zero weighted validator set only
	_, nonZeroWeighted := types.GetZeroNonZeroWightedAddrAmts(sortedValDiff)
	sort.Sort(sort.Reverse(nonZeroWeighted))

	valAmounts, remainderCoin := distributeCoinsAmongstValSet(nonZeroWeighted, coin)

	// If the remaining amount is not possitive, return early
	if !remainderCoin.IsPositive() {
		return valAmounts, nil
	}

	// Remaining token is the slippage from the multiplication with dec,
	// Ideally this amount is not going to be alot, hence assigning to
	// validator with index zero.
	valAmounts[0].Amount = valAmounts[0].Amount.Add(remainderCoin)

	sort.Sort(valAmounts)

	return valAmounts, nil
}

// DivideUndelegateAmountIntoValidatorSet : divides undelegation amount into validator set
//
//nolint:gocritic,len_not_fixed
func DivideUndelegateAmountIntoValidatorSet(sortedValDiff types.WeightedAddressAmounts, coin sdk.Coin) (types.ValAddressAmounts, error) {
	if coin.IsZero() {
		return nil, nil
	}

	// Undelegate first from zero weighted validators then nonzero weighted
	zeroWeighted, nonZeroWeighted := types.GetZeroNonZeroWightedAddrAmts(sortedValDiff)
	sort.Sort(sort.Reverse(zeroWeighted))
	sort.Sort(sort.Reverse(nonZeroWeighted))
	valWeighted := append(zeroWeighted, nonZeroWeighted...)

	valAmounts, remainderCoin := distributeCoinsAmongstValSet(valWeighted, coin)

	// If the remaining amount is not positive, return early
	if !remainderCoin.IsPositive() {
		return valAmounts, nil
	}

	// Divide the remaining amount amongst the validators a/c to weight
	zeroValued := sortedValDiff.GetZeroValued()
	valAddressMap := types.GetWeightedAddressMap(zeroValued)
	valAmounts = divideAmountWeightedSet(valAmounts, remainderCoin, valAddressMap)

	// sort the val address amount based on address to avoid generating different lists
	// by all validators
	sort.Sort(valAmounts)

	return valAmounts, nil
}

// GetAllValidatorsState returns the combined allowed listed validators set and combined
// delegation state. It is done to keep the old validators in the loop while calculating weighted amounts
// for delegation and undelegation
func (k Keeper) GetAllValidatorsState(ctx sdk.Context, denom string) (types.AllowListedVals, types.HostAccountDelegations) {
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
	var updatedValSet types.AllowListedVals
	for key, value := range currentAllowListedValSetMap {
		updatedValSet = append(updatedValSet, types.AllowListedValidator{ValidatorAddress: key, TargetWeight: value})
	}

	// Assign zero coin to all the validators not present in the current delegation state map
	for _, val := range delList {
		_, ok := currentDelegationStateMap[val]
		if !ok {
			currentDelegationStateMap[val] = sdk.NewCoin(denom, sdk.ZeroInt())
		}
	}

	// Convert the updated delegation state map to slice of types.HostChainDelegation
	var updatedDelegationState types.HostAccountDelegations
	for key, value := range currentDelegationStateMap {
		updatedDelegationState = append(updatedDelegationState, types.HostAccountDelegation{ValidatorAddress: key, Amount: value})
	}

	// sort both updatedValList and hostAccountDelegations
	sort.Sort(updatedValSet)
	sort.Sort(updatedDelegationState)

	// returns the two updated lists
	return updatedValSet, updatedDelegationState
}

// DelegateStrategy returns a list of types.ValAddressAmounts when provided with an input
// types.DelegationState, types.AllowListedValidators and amount to be delegated
func DelegateStrategy(valList types.AllowListedValidators, delegationState types.DelegationState, newAmount sdk.Coin) (types.ValAddressAmounts, error) {
	// return nil if the new amount to be delegated is already zero
	if newAmount.IsZero() {
		return nil, nil
	}

	// get validators weight map
	validatorWeightsMap := types.GetValidatorWeightsMap(valList)

	// get total delegations post adding more to it
	totalDelegationAfterAddingNewAmount := delegationState.TotalDelegations(newAmount.Denom).Add(newAmount).Amount

	// get the difference, current and ideally delegated amounts for each validator
	differenceCurrentAndIdealAmounts := GetDifferenceCurrentAndIdealAmounts(valList, delegationState, totalDelegationAfterAddingNewAmount)

	// sort the difference, current and ideally delegated amounts for each validator based
	// on the difference in descending order
	sort.SliceStable(
		differenceCurrentAndIdealAmounts,
		func(i, j int) bool {
			return differenceCurrentAndIdealAmounts[i].Diff.GT(differenceCurrentAndIdealAmounts[j].Diff)
		},
	)

	// this is the first round of distribution to all the validators who are not
	// yet matched to the ideal distribution
	finalDistribution := make(map[string]sdk.Coin)
	for _, i := range differenceCurrentAndIdealAmounts {

		if i.Diff.GT(sdk.ZeroDec()) {
			if newAmount.Amount.ToDec().LT(i.Diff.Abs()) && newAmount.Amount.Sub(newAmount.Amount).GTE(sdk.OneInt()) {
				finalDistribution[i.Address] = newAmount
				newAmount.Amount = newAmount.Amount.Sub(newAmount.Amount)
			} else if newAmount.Amount.ToDec().GTE(i.Diff.Abs()) && newAmount.Amount.Sub(newAmount.Amount).GTE(sdk.OneInt()) {
				finalDistribution[i.Address] = sdk.NewCoin(newAmount.Denom, i.Diff.Abs().TruncateInt())
				newAmount.Amount = newAmount.Amount.Sub(i.Diff.Abs().TruncateInt())
			}
		}
	}

	// this is the second round of distribution to all the validators based on their
	// current weights
	temporaryCoin := sdk.NewCoin(newAmount.Denom, sdk.ZeroInt())
	for _, i := range differenceCurrentAndIdealAmounts {
		if newAmount.IsZero() {
			break
		}
		amountForValidator := newAmount.Amount.ToDec().Mul(validatorWeightsMap[i.Address]).TruncateInt()
		if !amountForValidator.IsZero() {
			finalDistribution[i.Address] = sdk.NewCoin(newAmount.Denom, amountForValidator)
			temporaryCoin.Amount = temporaryCoin.Amount.Add(amountForValidator)
		}
	}

	// convert distribution map into types.ValAddressAmounts array
	var finalValAddressAmounts types.ValAddressAmounts
	for key, value := range finalDistribution {
		finalValAddressAmounts = append(finalValAddressAmounts, types.ValAddressAmount{
			ValidatorAddr: key,
			Amount:        value,
		})
	}

	// sort the types.ValAddressAmounts based on the validator addresses
	sort.SliceStable(finalValAddressAmounts, func(i, j int) bool {
		return finalValAddressAmounts[i].ValidatorAddr < finalValAddressAmounts[j].ValidatorAddr
	})

	// this is a boundary case when there is still some residue left after
	// distributing it to all the validators
	if !newAmount.Sub(temporaryCoin).IsZero() && len(finalValAddressAmounts) > 0 {
		// give all the leftover to the first validator in the above computed list
		finalValAddressAmounts[0].Amount = finalValAddressAmounts[0].Amount.Add(newAmount.Sub(temporaryCoin))
	} else if !newAmount.Sub(temporaryCoin).IsZero() && len(finalValAddressAmounts) == 0 {
		// give all the leftover to the very first validator to the sorted
		// differenceCurrentAndIdealAmounts slice
		return types.ValAddressAmounts{{ValidatorAddr: differenceCurrentAndIdealAmounts[0].Address, Amount: newAmount.Sub(temporaryCoin)}}, nil
	}

	return finalValAddressAmounts, nil
}

// GetDifferenceCurrentAndIdealAmounts return an array of types.DifferenceCurrentAndIdealAmount from the given
// types.AllowListedValidators, delegation state map and
func GetDifferenceCurrentAndIdealAmounts(valList types.AllowListedValidators, delegationState types.DelegationState, totalDelegationAfterAddingNewAmount sdk.Int) []types.DifferenceCurrentAndIdealAmount {
	// get the delegation state map
	delegationStateMap := types.GetDelegationStateMap(delegationState)

	// find out current ideal and difference in delegations for every validator
	var differenceCurrentAndIdealAmounts []types.DifferenceCurrentAndIdealAmount
	for _, i := range valList.AllowListedValidators {
		// if validator address not present in current delegation state then introduce it in the
		// if is present then continue as normal addition
		currentDelegation, ok := delegationStateMap[i.ValidatorAddress]
		if !ok {
			idealDelegation := totalDelegationAfterAddingNewAmount.ToDec().Mul(i.TargetWeight)
			diff := idealDelegation.Sub(sdk.ZeroDec())
			differenceCurrentAndIdealAmounts = append(differenceCurrentAndIdealAmounts, types.DifferenceCurrentAndIdealAmount{
				Current: sdk.ZeroDec(),
				Ideal:   idealDelegation,
				Diff:    diff,
				Address: i.ValidatorAddress,
			})
		} else {
			idealDelegation := totalDelegationAfterAddingNewAmount.ToDec().Mul(i.TargetWeight)
			diff := idealDelegation.Sub(currentDelegation.Amount.ToDec())
			differenceCurrentAndIdealAmounts = append(differenceCurrentAndIdealAmounts, types.DifferenceCurrentAndIdealAmount{
				Current: currentDelegation.Amount.ToDec(),
				Ideal:   idealDelegation,
				Diff:    diff,
				Address: i.ValidatorAddress,
			})
		}
	}

	return differenceCurrentAndIdealAmounts
}
