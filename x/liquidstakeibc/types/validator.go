package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (v *Validator) SharesToTokens(shares sdk.Dec) sdk.Int { //nolint:staticcheck
	if v.DelegatorShares.IsZero() {
		return sdk.ZeroInt()
	}

	return sdk.NewDecFromInt(v.TotalAmount).Quo(v.DelegatorShares).Mul(shares).TruncateInt()
}
