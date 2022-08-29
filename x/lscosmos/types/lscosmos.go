package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (av *AllowListedValidators) Valid() bool {
	if av == nil {
		return false
	}
	if av.AllowListedValidators == nil || len(av.AllowListedValidators) == 0 {
		return false
	}

	sum := sdk.ZeroDec()
	for _, v := range av.AllowListedValidators {
		if _, err := sdk.ValAddressFromBech32(v.ValidatorAddress); err != nil {
			return false
		}
		sum = sum.Add(v.TargetWeight)
	}
	return sum.Equal(sdk.OneDec())
}

func NewHostAccountDelegation(validatorAddress string, amount sdk.Coin) HostAccountDelegation {
	return HostAccountDelegation{
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}
