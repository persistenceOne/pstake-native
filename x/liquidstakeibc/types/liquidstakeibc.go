package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func IsUnbondingEpoch(factor, epochNumber int64) bool {
	return epochNumber%factor == 0
}

// CurrentUnbondingEpoch computes and returns the current unbonding epoch to the next nearest
// multiple of the host chain Undelegation Factor
func CurrentUnbondingEpoch(factor, epochNumber int64) int64 {
	if epochNumber%factor == 0 {
		return epochNumber
	}
	return epochNumber + factor - epochNumber%factor
}

// DefaultDelegateAccountPortOwner generates a delegate ICA port owner given the chain id
// Only Use this function while registering a new chain
func DefaultDelegateAccountPortOwner(chainID string) string {
	return chainID + "." + DelegateICAType
}

// DefaultRewardsAccountPortOwner generates a rewards ICA port owner given the chain id
// Only Use this function while registering a new chain
func DefaultRewardsAccountPortOwner(chainID string) string {
	return chainID + "." + RewardsICAType
}

func (deposit *Deposit) Validate() error {
	if deposit.State != Deposit_DEPOSIT_PENDING &&
		deposit.State != Deposit_DEPOSIT_RECEIVED &&
		deposit.State != Deposit_DEPOSIT_DELEGATING {
		return fmt.Errorf(
			"host chain %s deposit has an invalid state: %s",
			deposit.ChainId,
			deposit.State,
		)
	}
	if deposit.Amount.Amount.LT(sdk.ZeroInt()) {
		return fmt.Errorf("deposit for chain %s has negative amount", deposit.ChainId)
	}
	return nil
}

func (hc *HostChain) Validate() error {
	if hc.Params.DepositFee.LT(sdk.ZeroDec()) {
		return fmt.Errorf("host chain %s has negative deposit fee", hc.ChainId)
	}
	if hc.Params.RestakeFee.LT(sdk.ZeroDec()) {
		return fmt.Errorf("host chain %s has negative restake fee", hc.ChainId)
	}
	if hc.Params.RedemptionFee.LT(sdk.ZeroDec()) {
		return fmt.Errorf("host chain %s has negative redemption fee", hc.ChainId)
	}
	if hc.Params.UnstakeFee.LT(sdk.ZeroDec()) {
		return fmt.Errorf("host chain %s has negative unstake fee", hc.ChainId)
	}

	if hc.MinimumDeposit.LT(sdk.ZeroInt()) {
		return fmt.Errorf("host chain %s has negative minimum deposit", hc.ChainId)
	}
	if hc.CValue.LT(sdk.ZeroDec()) || hc.CValue.GT(sdk.OneDec()) {
		return fmt.Errorf("host chain %s has c value out of bounds: %d", hc.ChainId, hc.CValue)
	}

	for _, validator := range hc.Validators {
		if validator.Status != stakingtypes.Unspecified.String() &&
			validator.Status != stakingtypes.Unbonded.String() &&
			validator.Status != stakingtypes.Unbonding.String() &&
			validator.Status != stakingtypes.Bonded.String() {
			return fmt.Errorf(
				"host chain %s validator %s has an invalid status: %s",
				hc.ChainId,
				validator.OperatorAddress,
				validator.Status,
			)
		}

		if validator.Weight.LT(sdk.ZeroDec()) || validator.Weight.GT(sdk.OneDec()) {
			return fmt.Errorf(
				"host chain %s validator %s has weight out of bounds: %d",
				hc.ChainId,
				validator.OperatorAddress,
				validator.Weight)
		}

		if validator.DelegatedAmount.LT(sdk.ZeroInt()) {
			return fmt.Errorf(
				"host chain %s validator %s has negative delegated amount: %s",
				hc.ChainId,
				validator.OperatorAddress,
				validator.DelegatedAmount.String(),
			)
		}
	}
	return nil
}

func (u *Unbonding) Validate() error {
	if u.BurnAmount.IsNegative() {
		return fmt.Errorf("unbonding entry %s has negative burn amount: %s", u.String(), u.BurnAmount)
	}
	if u.UnbondAmount.IsNegative() {
		return fmt.Errorf("unbonding entry %s has negative unbond amount: %s", u.String(), u.UnbondAmount)
	}
	return nil
}

func (ub *UserUnbonding) Validate() error {
	if _, err := sdk.AccAddressFromBech32(ub.Address); err != nil {
		return sdkerrors.ErrInvalidAddress
	}
	if ub.UnbondAmount.IsNegative() {
		return fmt.Errorf("user unbonding %s has negative unbonding amount, amount: %s", ub.String(), ub.UnbondAmount)
	}
	return nil
}

func (vb *ValidatorUnbonding) Validate() error {
	if _, _, err := bech32.DecodeAndConvert(vb.ValidatorAddress); err != nil {
		return err
	}
	if vb.Amount.IsNegative() {
		return fmt.Errorf("validator unbonding %s has negative amount, amount: %s", vb.String(), vb.Amount)
	}
	return nil
}
