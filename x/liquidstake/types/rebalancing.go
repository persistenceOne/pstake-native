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

// DivideByCurrentWeight divide the input value by the ratio of the weight of the liquid validator's liquid token and return it with crumb
// which is may occur while dividing according to the weight of liquid validators by decimal error, outputs is truncated decimal.
func DivideByCurrentWeight(lvs LiquidValidators, input math.LegacyDec, totalLiquidTokens math.Int, liquidTokenMap map[string]math.Int) (outputs []math.LegacyDec, crumb math.LegacyDec) {
	if !totalLiquidTokens.IsPositive() {
		return []math.LegacyDec{}, sdk.ZeroDec()
	}

	totalOutput := sdk.ZeroDec()
	unitInput := input.QuoTruncate(math.LegacyNewDecFromInt(totalLiquidTokens))
	for _, val := range lvs {
		output := unitInput.MulTruncate(math.LegacyNewDecFromInt(liquidTokenMap[val.OperatorAddress])).TruncateDec()
		totalOutput = totalOutput.Add(output)
		outputs = append(outputs, output)
	}

	return outputs, input.Sub(totalOutput)
}
