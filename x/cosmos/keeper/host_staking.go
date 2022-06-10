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

func GetIdealCurrentDelegations(validatorState types.WeightedAddressAmounts, amt sdk.Coin, reverse bool) types.WeightedAddressAmounts {
	totalDelegations := validatorState.TotalAmount(amt.Denom)
	curDiffDistribution := types.WeightedAddressAmounts{}
	var idealTokens, curTokens sdk.Int
	for _, valState := range validatorState {
		// Note this can lead to some leaks
		// Considering additional amount in ideal distribution
		totalAmt := totalDelegations.Amount.Add(amt.Amount)
		if reverse {
			totalAmt = totalDelegations.Amount.Sub(amt.Amount)
		}
		idealTokens = valState.Weight.Mul(sdk.NewDecFromInt(totalAmt)).TruncateInt()
		curTokens = valState.Amount
		diffAmt := idealTokens.Sub(curTokens)
		if reverse {
			diffAmt = curTokens.Sub(idealTokens)
		}
		curDiffDistribution = append(curDiffDistribution, types.WeightedAddressAmount{
			Address: valState.Address,
			Weight:  valState.Weight,
			Denom:   valState.Denom,
			Amount:  diffAmt,
		})
	}
	return curDiffDistribution
}

func divideAmountWeightedSet(valAmounts []ValAddressAmount, coin sdk.Coin, valAddressWeightMap map[string]sdk.Dec) []ValAddressAmount {
	newValAmounts := []ValAddressAmount{}

	totalWeight := sdk.ZeroDec()
	for _, weight := range valAddressWeightMap {
		totalWeight = totalWeight.Add(weight)
	}

	for _, valAmt := range valAmounts {
		weight := valAddressWeightMap[valAmt.Validator.String()].Quo(totalWeight)
		amt := weight.MulInt(coin.Amount).RoundInt()
		newValAmounts = append(newValAmounts, ValAddressAmount{
			Validator: valAmt.Validator,
			Amount:    sdk.NewCoin(valAmt.Amount.Denom, valAmt.Amount.Amount.Add(amt)),
		})
	}
	return newValAmounts
}

// distributeCoinsAmongstValSet takes the validator distribution and coins to distribute and returns the
// validator address amount to distribute and the remaining amount to make
func distributeCoinsAmongstValSet(ws types.WeightedAddressAmounts, coin sdk.Coin) ([]ValAddressAmount, sdk.Coin, error) {
	valAddrAmts := []ValAddressAmount{}

	for _, w := range ws {
		// Create val address
		valAddr, err := types.ValAddressFromBech32(w.Address, types.Bech32PrefixValAddr)
		if err != nil {
			return nil, coin, err
		}
		if coin.Amount.LTE(w.Amount) {
			valAddrAmts = append(valAddrAmts, ValAddressAmount{Validator: valAddr, Amount: coin})
			return valAddrAmts, sdk.NewInt64Coin(coin.Denom, 0), nil
		}
		valAddrAmts = append(valAddrAmts, ValAddressAmount{Validator: valAddr, Amount: w.Coin()})
		coin = coin.SubAmount(w.Amount)
	}

	return valAddrAmts, coin, nil
}

func DivideAmountIntoValidatorSet(sortedValDiff types.WeightedAddressAmounts, coin sdk.Coin) ([]ValAddressAmount, error) {
	if coin.IsZero() {
		return nil, nil
	}

	// Delegate to non zero weighted validator set only
	_, nonZeroWeighted := types.GetZeroNonZeroWightedAddrAmts(sortedValDiff)
	sort.Sort(sort.Reverse(nonZeroWeighted))

	valAmounts, remainderCoin, err := distributeCoinsAmongstValSet(nonZeroWeighted, coin)
	if err != nil {
		return nil, err
	}

	// If the remaining amount is not possitive, return early
	if !remainderCoin.IsPositive() {
		return valAmounts, nil
	}

	// Divide the remaining amount amongst the validators a/c to weight
	// Get zero valued val address to divide the remaing value a/c to weight
	zeroValued := sortedValDiff.GetZeroValued()
	valAddressMap := types.GetWeightedAddressMap(zeroValued)
	valAmounts = divideAmountWeightedSet(valAmounts, remainderCoin, valAddressMap)

	return valAmounts, nil
}

func DivideUndelegateAmountIntoValidatorSet(sortedValDiff types.WeightedAddressAmounts, coin sdk.Coin) ([]ValAddressAmount, error) {
	if coin.IsZero() {
		return nil, nil
	}

	// Undelegate first from zero weighted validators then nonzero weighted
	zeroWeighted, nonZeroWeighted := types.GetZeroNonZeroWightedAddrAmts(sortedValDiff)
	sort.Sort(sort.Reverse(zeroWeighted))
	sort.Sort(sort.Reverse(nonZeroWeighted))
	valWeighted := append(zeroWeighted, nonZeroWeighted...)

	valAmounts, remainderCoin, err := distributeCoinsAmongstValSet(valWeighted, coin)
	if err != nil {
		return nil, err
	}

	// If the remaining amount is not possitive, return early
	if !remainderCoin.IsPositive() {
		return valAmounts, nil
	}

	// Divide the remaining amount amongst the validators a/c to weight
	zeroValued := sortedValDiff.GetZeroValued()
	valAddressMap := types.GetWeightedAddressMap(zeroValued)
	valAmounts = divideAmountWeightedSet(valAmounts, remainderCoin, valAddressMap)

	return valAmounts, nil
}

// gives a list of all validators having weighted amount for few and 1uatom for rest in order to auto claim all rewards accumulated in current epoch
func (k Keeper) FetchValidatorsToDelegate(ctx sdk.Context, amount sdk.Coin) ([]ValAddressAmount, error) {
	params := k.GetParams(ctx)

	// Return nil list if amount is less than delegation threshold
	if amount.IsLT(params.DelegationThreshold) {
		return nil, nil
	}

	valWeightedAmt := k.GetAllCosmosValidatorSet(ctx)
	curDiffDistribution := GetIdealCurrentDelegations(valWeightedAmt, amount, false)

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

	valWeightedAmt := k.GetAllCosmosValidatorSet(ctx)

	// Check if amount asked to undelegate is more than total delegations
	totalStaked := valWeightedAmt.TotalAmount(params.StakingDenom)
	if totalStaked.Amount.LT(amount.Amount) {
		return nil, fmt.Errorf("undelegate amount %d more than total staked %d", amount.Amount, totalStaked.Amount)
	}

	curDiffDistribution := GetIdealCurrentDelegations(valWeightedAmt, amount, true)

	sort.Sort(sort.Reverse(curDiffDistribution))

	return DivideUndelegateAmountIntoValidatorSet(curDiffDistribution, amount)
}
