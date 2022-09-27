package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// DelegateMsgs gives the list of Delegate Txs to be executed based on the current state and params.
func (k Keeper) DelegateMsgs(ctx sdk.Context, delegatorAddr string, amount sdk.Int, denom string) ([]sdk.Msg, error) {
	valList := k.GetAllowListedValidators(ctx)
	delegationState := k.GetDelegationState(ctx)

	valAddressAmount, err := FetchValidatorsToDelegate(valList, delegationState, sdk.NewCoin(denom, amount))
	if err != nil {
		return nil, err
	}

	msgs := make([]sdk.Msg, len(valAddressAmount))

	for i, val := range valAddressAmount {

		msg := &stakingtypes.MsgDelegate{
			DelegatorAddress: delegatorAddr,
			ValidatorAddress: val.ValidatorAddr,
			Amount:           val.Amount,
		}
		msgs[i] = msg
	}

	return msgs, nil
}

// UndelegateMsgs gives the list of Undelegate Txs to be executed based on the current state and params.
func (k Keeper) UndelegateMsgs(ctx sdk.Context, delegatorAddr string, amount sdk.Int, denom string) ([]sdk.Msg, []types.UndelegationEntry, error) {
	valList := k.GetAllowListedValidators(ctx)
	delegationState := k.GetDelegationState(ctx)

	valAddressAmount, err := FetchValidatorsToUndelegate(valList, delegationState, sdk.NewCoin(denom, amount))
	if err != nil {
		return nil, nil, err
	}

	msgs := make([]sdk.Msg, len(valAddressAmount))
	undelegationEntries := make([]types.UndelegationEntry, len(valAddressAmount))

	for i, val := range valAddressAmount {

		msg := &stakingtypes.MsgUndelegate{
			DelegatorAddress: delegatorAddr,
			ValidatorAddress: val.ValidatorAddr,
			Amount:           val.Amount,
		}
		msgs[i] = msg

		undelegationEntry := types.UndelegationEntry{
			ValidatorAddress: val.ValidatorAddr,
			Amount:           val.Amount,
		}
		undelegationEntries[i] = undelegationEntry
	}

	return msgs, undelegationEntries, nil
}

// FetchValidatorsToDelegate gives a list of all validators having weighted amount for few and 1uatom for rest in order to auto claim all rewards accumulated in current epoch
func FetchValidatorsToDelegate(valList types.AllowListedValidators, delegationState types.DelegationState, amount sdk.Coin) ([]types.ValAddressAmount, error) {
	curDiffDistribution := GetIdealCurrentDelegations(valList, delegationState, amount, false)
	sort.Sort(sort.Reverse(curDiffDistribution))

	return DivideAmountIntoValidatorSet(curDiffDistribution, amount)
}

// FetchValidatorsToUndelegate gives a list of all validators having weighted amount for few and 1uatom for rest in order to auto claim all rewards accumulated in current epoch
func FetchValidatorsToUndelegate(valList types.AllowListedValidators, delegationState types.DelegationState, amount sdk.Coin) ([]types.ValAddressAmount, error) {
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
//
//nolint:prealloc,len_not_fixed
func divideAmountWeightedSet(valAmounts []types.ValAddressAmount, coin sdk.Coin, valAddressWeightMap map[string]sdk.Dec) []types.ValAddressAmount {
	var newValAmounts []types.ValAddressAmount

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
//
//nolint:prealloc,len_not_fixed
func distributeCoinsAmongstValSet(ws types.WeightedAddressAmounts, coin sdk.Coin) ([]types.ValAddressAmount, sdk.Coin) {
	var valAddrAmts []types.ValAddressAmount

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
func DivideAmountIntoValidatorSet(sortedValDiff types.WeightedAddressAmounts, coin sdk.Coin) ([]types.ValAddressAmount, error) {
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

	return valAmounts, nil
}

// DivideUndelegateAmountIntoValidatorSet : divides undelegation amount into validator set
//
//nolint:gocritic,len_not_fixed
func DivideUndelegateAmountIntoValidatorSet(sortedValDiff types.WeightedAddressAmounts, coin sdk.Coin) ([]types.ValAddressAmount, error) {
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

	return valAmounts, nil
}
