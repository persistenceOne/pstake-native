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

func GetAddressMap(validators AllowListedValidators) map[string]sdk.Dec {
	addressMap := map[string]sdk.Dec{}

	for _, val := range validators.AllowListedValidators {
		addressMap[val.ValidatorAddress] = val.TargetWeight
	}
	return addressMap
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

// TotalDelegations gives the amount of total delegations on Host Chain.
func (ds DelegationState) TotalDelegations(denom string) sdk.Coin {
	total := sdk.NewCoin(denom, sdk.ZeroInt())

	for _, val := range ds.HostAccountDelegations {
		if val.Amount.Denom == denom {
			total.Amount = total.Amount.Add(val.Amount.Amount)
		}
	}
	return total
}

func NewDelegatorUnbondingEpochEntry(delegatorAddress string, epochNumber int64, amount sdk.Coin) DelegatorUnbondingEpochEntry {
	return DelegatorUnbondingEpochEntry{
		DelegatorAddress: delegatorAddress,
		EpochNumber:      epochNumber,
		Amount:           amount,
	}
}

// GetUnbondingEpochCValue returns the calculated c value from the UnbondingEpochCValue struct entries.
func (uec *UnbondingEpochCValue) GetUnbondingEpochCValue() sdk.Dec {
	return uec.STKBurn.Amount.ToDec().Quo(uec.AmountUnbonded.Amount.ToDec())
}

func CurrentUnbondingEpoch(epochNumber int64) int64 {
	if epochNumber%UndelegationEpochNumberFactor == 0 {
		return epochNumber
	}
	return epochNumber + UndelegationEpochNumberFactor - epochNumber%UndelegationEpochNumberFactor
}
func PreviousUnbondingEpoch(epochNumber int64) int64 {
	if epochNumber%UndelegationEpochNumberFactor == 0 {
		return epochNumber - UndelegationEpochNumberFactor
	}
	return epochNumber - epochNumber%UndelegationEpochNumberFactor
}
