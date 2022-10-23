package types

import (
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
)

func (av *AllowListedValidators) Valid() bool {
	if av == nil {
		return false
	}
	if av.AllowListedValidators == nil || len(av.AllowListedValidators) == 0 {
		return false
	}

	noDuplicate := make(map[string]bool)
	sum := sdk.ZeroDec()
	for _, v := range av.AllowListedValidators {
		if _, err := ValAddressFromBech32(v.ValidatorAddress); err != nil {
			return false
		}
		_, ok := noDuplicate[v.ValidatorAddress]
		if ok {
			return false
		}
		noDuplicate[v.ValidatorAddress] = true
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

func (hostAccounts *HostAccounts) DelegatorAccountPortID() string {
	delegatorAccountPortID, err := icatypes.NewControllerPortID(hostAccounts.DelegatorAccountOwnerID)
	if err != nil {
		panic(err)
	}
	return delegatorAccountPortID
}

func (hostAccounts *HostAccounts) RewardsAccountPortID() string {
	rewardsAccountPortID, err := icatypes.NewControllerPortID(hostAccounts.RewardsAccountOwnerID)
	if err != nil {
		panic(err)
	}
	return rewardsAccountPortID
}

func (hostAccounts *HostAccounts) Validate() error {
	if hostAccounts.RewardsAccountOwnerID == "" || hostAccounts.DelegatorAccountOwnerID == "" {
		return ErrInvalidHostAccountOwnerIDs
	}
	return nil
}

func (pstakeParams *PstakeParams) Validate() error {
	_, err := sdk.AccAddressFromBech32(pstakeParams.PstakeFeeAddress)
	if err != nil {
		return err
	}

	if pstakeParams.PstakeDepositFee.IsNegative() || pstakeParams.PstakeDepositFee.GTE(MaxPstakeDepositFee) {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake deposit fee must be between %s and %s", sdk.ZeroDec(), MaxPstakeDepositFee)
	}

	if pstakeParams.PstakeRestakeFee.IsNegative() || pstakeParams.PstakeRestakeFee.GTE(MaxPstakeRestakeFee) {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake restake fee must be between %s and %s", sdk.ZeroDec(), MaxPstakeRestakeFee)
	}

	if pstakeParams.PstakeUnstakeFee.IsNegative() || pstakeParams.PstakeUnstakeFee.GTE(MaxPstakeUnstakeFee) {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake unstake fee must be between %s and %s", sdk.ZeroDec(), MaxPstakeUnstakeFee)
	}

	if pstakeParams.PstakeRedemptionFee.IsNegative() || pstakeParams.PstakeRedemptionFee.GTE(MaxPstakeRedemptionFee) {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake redemption fee must be between %s and %s", sdk.ZeroDec(), MaxPstakeRedemptionFee)
	}
	return nil
}
