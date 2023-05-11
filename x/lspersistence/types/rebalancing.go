package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Redelegation struct {
	Delegator    sdk.AccAddress
	SrcValidator LiquidValidator
	DstValidator LiquidValidator
	Amount       math.Int
	Last         bool
	Error        error
}

// DivideByWeight divide the input value by the ratio of the param weight of the liquid validator and return it with crumb
// which is may occur while dividing according to the weight of active liquid validators by decimal error.
func DivideByWeight(avs ActiveLiquidValidators, input math.Int, whitelistedValsMap WhitelistedValsMap) (outputs []math.Int, crumb math.Int) {
	totalWeight := avs.TotalWeight(whitelistedValsMap)
	if !totalWeight.IsPositive() {
		return []math.Int{}, sdk.ZeroInt()
	}
	totalOutput := sdk.ZeroInt()
	unitInput := sdk.NewDecFromInt(input).QuoTruncate(sdk.NewDecFromInt(totalWeight))
	for _, val := range avs {
		output := unitInput.MulInt(val.GetWeight(whitelistedValsMap, true)).TruncateInt()
		totalOutput = totalOutput.Add(output)
		outputs = append(outputs, output)
	}
	return outputs, input.Sub(totalOutput)
}

// DivideByCurrentWeight divide the input value by the ratio of the weight of the liquid validator's liquid token and return it with crumb
// which is may occur while dividing according to the weight of liquid validators by decimal error, outputs is truncated decimal.
func DivideByCurrentWeight(lvs LiquidValidators, input sdk.Dec, totalLiquidTokens math.Int, liquidTokenMap map[string]math.Int) (outputs []sdk.Dec, crumb sdk.Dec) {
	if !totalLiquidTokens.IsPositive() {
		return []sdk.Dec{}, sdk.ZeroDec()
	}
	totalOutput := sdk.ZeroDec()
	unitInput := input.QuoTruncate(sdk.NewDecFromInt(totalLiquidTokens))
	for _, val := range lvs {
		output := unitInput.MulTruncate(sdk.NewDecFromInt(liquidTokenMap[val.OperatorAddress])).TruncateDec()
		totalOutput = totalOutput.Add(output)
		outputs = append(outputs, output)
	}
	return outputs, input.Sub(totalOutput)
}
