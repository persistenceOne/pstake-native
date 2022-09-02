package types

import (
	"errors"
	"strings"

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
		if _, err := ValAddressFromBech32(v.ValidatorAddress); err != nil {
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

func NewHostChainRewardAddress(address string) HostChainRewardAddress {
	return HostChainRewardAddress{
		Address: address,
	}
}

// ValAddressFromBech32 creates a ValAddress from a Bech32 string.
func ValAddressFromBech32(address string) (addr sdk.ValAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.ValAddress{}, errors.New("empty address string is not allowed")
	}

	bech32PrefixValAddr := CosmosValOperPrefix

	bz, err := sdk.GetFromBech32(address, bech32PrefixValAddr)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return bz, nil
}
