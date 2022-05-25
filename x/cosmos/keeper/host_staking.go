package keeper

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type ValAddressAmount struct {
	Validator sdk.ValAddress
	Amount    sdk.Coin
}

// normalizedWeightedAddressAmounts function takes input as the weighted address amounts
// finds the smallest amount or zero from the array and returns a new array with normalized amounts
func NormalizedWeightedAddressAmounts(weightedAddrAmt types.WeightedAddressAmounts) types.WeightedAddressAmounts {
	// Find smallest diff less than zero
	smallestVal := sdk.ZeroInt()
	normalizedDistribution := types.WeightedAddressAmounts{}

	for _, w := range weightedAddrAmt {
		if w.Amount.LT(smallestVal) {
			smallestVal = w.Amount
		}
	}
	// Return early incase the smallest value is zero 
	if smallestVal.Equal(sdk.ZeroInt()) {
		return weightedAddrAmt
	}
	// Normalize based on smallest diff
	for _, w := range weightedAddrAmt {
		normCoin := sdk.NewCoin(w.Denom, w.Amount.Sub(smallestVal))
		normalizedDistribution = append(
			normalizedDistribution,
			types.NewWeightedAddressAmount(w.Address, w.Weight, normCoin),
		)
	}
	return normalizedDistribution
}

func GetIdealCurrentDelegations(validatorState types.WeightedAddressAmounts, stakingDenom string, reverse bool) types.WeightedAddressAmounts {
	totalDelegations := validatorState.TotalAmount(stakingDenom)
	curDiffDistribution := types.WeightedAddressAmounts{}
	var idealTokens, curTokens sdk.Int
	for _, valState := range validatorState {
		// Note this can lead to some leaks
		idealTokens = valState.Weight.Mul(totalDelegations.Amount.ToDec()).RoundInt()
		curTokens = valState.Amount
		amt := idealTokens.Sub(curTokens)
		if reverse {
			amt = curTokens.Sub(idealTokens)
		}
		curDiffDistribution = append(curDiffDistribution, types.WeightedAddressAmount{
			Address: valState.Address,
			Weight: valState.Weight,
			Denom: valState.Denom,
			Amount: amt,
		})
	}
	return curDiffDistribution
}

func divideAmountWeightedSet(valAmounts []ValAddressAmount, coin sdk.Coin, valAddressWeightMap map[string]sdk.Dec) []ValAddressAmount {
	newValAmounts := []ValAddressAmount{}
	for _, valAmt := range valAmounts {
		weight := valAddressWeightMap[valAmt.Validator.String()]
		amt := weight.MulInt(coin.Amount).RoundInt()
		newValAmounts = append(newValAmounts, ValAddressAmount{
			Validator: valAmt.Validator,
			Amount: sdk.NewCoin(valAmt.Amount.Denom, valAmt.Amount.Amount.Add(amt)),
		})
	}
	return newValAmounts
}

func DivideAmountIntoValidatorSet(sortedValDiff types.WeightedAddressAmounts, coin sdk.Coin) ([]ValAddressAmount, error) {
	if coin.IsZero() {
		return nil, nil
	}

	valAmounts := []ValAddressAmount{}
	
	for _, w := range sortedValDiff {
		// Skip validators with zero weights
		if w.Weight.IsZero() {
			continue
		}
		// Create val address
		valAddr, err := sdk.ValAddressFromBech32(w.Address)
		if err != nil {
			return nil, err
		}
		if w.Amount.GTE(coin.Amount) {
			valAmounts = append(valAmounts, ValAddressAmount{Validator: valAddr, Amount: coin})
			return valAmounts, nil
		}
		valAmounts = append(valAmounts, ValAddressAmount{Validator: valAddr, Amount: w.Coin()})
		coin = coin.SubAmount(w.Amount)
	}

	// If the remaining amount is not possitive, return early
	if !coin.IsPositive() {
		return valAmounts, nil
	}

	// Divide the remaining amount amongst the validators a/c to weight
	// Note: Maybe there is some slippage due to multiplication
	valAddressMap := types.GetWeightedAddressMap(sortedValDiff)
	valAmounts = divideAmountWeightedSet(valAmounts, coin, valAddressMap)

	return valAmounts, nil
}

func DivideUndelegateAmountIntoValidatorSet(sortedValDiff types.WeightedAddressAmounts, coin sdk.Coin) ([]ValAddressAmount, error) {
	if coin.IsZero() {
		return nil, nil
	}

	valAmounts := []ValAddressAmount{}

	zeroVals := sortedValDiff.GetZeroWeighted()
	sort.Sort(zeroVals)
	for _, w := range zeroVals {
		valAddr, err := sdk.ValAddressFromBech32(w.Address)
		if err != nil {
			return nil, err
		}
		if w.Amount.LTE(coin.Amount) {
			valAmounts = append(valAmounts, ValAddressAmount{Validator: valAddr, Amount: coin})
			return valAmounts, nil
		}
		fmt.Printf("Coins: %s and val amount: %s", coin.Amount, w.Amount)
		valAmounts = append(valAmounts, ValAddressAmount{Validator: valAddr, Amount: w.Coin()})
		coin = coin.SubAmount(w.Amount) // gets negative
	}
	
	for _, w := range sortedValDiff {
		// Skip validators with zero weights
		if w.Weight.IsZero() {
			continue
		}
		// Create val address
		valAddr, err := sdk.ValAddressFromBech32(w.Address)
		if err != nil {
			return nil, err
		}
		// ideal - current < coin
		if w.Amount.LTE(coin.Amount) {
			valAmounts = append(valAmounts, ValAddressAmount{Validator: valAddr, Amount: coin})
			return valAmounts, nil
		}
		fmt.Printf("Coins: %s and val amount: %s", coin.Amount, w.Amount)
		// ideal - current > coin
		valAmounts = append(valAmounts, ValAddressAmount{Validator: valAddr, Amount: w.Coin()})
		coin = coin.SubAmount(w.Amount) // gets negative
	}

	// If the remaining amount is not possitive, return early
	if !coin.IsPositive() {
		return valAmounts, nil
	}

	// Divide the remaining amount amongst the validators a/c to weight
	// Note: Maybe there is some slippage due to multiplication
	valAddressMap := types.GetWeightedAddressMap(sortedValDiff)
	valAmounts = divideAmountWeightedSet(valAmounts, coin, valAddressMap)

	return valAmounts, nil
}

// gives a list of all validators having weighted amount for few and 1uatom for rest in order to auto claim all rewards accumulated in current epoch
func (k Keeper) FetchValidatorsToDelegate(ctx sdk.Context, amount sdk.Coin) ([]ValAddressAmount, error) {
	params := k.GetParams(ctx)

	// Return nil list if amount is less than delegation threshold
	if amount.IsLT(params.DelegationThreshold) {
		return nil, nil
	}

	valWeightedAmt := k.getAllCosmosValidatorSet(ctx)
	
	curDiffDistribution := GetIdealCurrentDelegations(valWeightedAmt, params.StakingDenom, false)
	curDiffDistribution = NormalizedWeightedAddressAmounts(curDiffDistribution)
	
	sort.Sort(sort.Reverse(curDiffDistribution))

	return DivideAmountIntoValidatorSet(curDiffDistribution, amount)
}

// gives a list of validators having weighted amount for few validators
func (k Keeper) FetchValidatorsToUndelegate(ctx sdk.Context, amount sdk.Coin) ([]ValAddressAmount, error) {
	params := k.GetParams(ctx)

	// Return nil list if amount is less than delegation threshold
	if amount.IsLT(params.DelegationThreshold) {
		return nil, nil
	}

	valWeightedAmt := k.getAllCosmosValidatorSet(ctx)

	// Check if amount asked to undelegate is more than total delegations
	totalStaked := valWeightedAmt.TotalAmount(params.StakingDenom)
	if totalStaked.Amount.LT(amount.Amount) {
		return nil, fmt.Errorf("undelegate amount %d more than total staked %d", amount.Amount, totalStaked.Amount)
	}
	
	curDiffDistribution := GetIdealCurrentDelegations(valWeightedAmt, params.StakingDenom, true)
	curDiffDistribution = NormalizedWeightedAddressAmounts(curDiffDistribution)
	
	sort.Sort(sort.Reverse(curDiffDistribution))

	return DivideUndelegateAmountIntoValidatorSet(curDiffDistribution, amount)
}
