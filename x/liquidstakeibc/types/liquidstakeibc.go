package types

import (
	"fmt"
	"github.com/cometbft/cometbft/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func IsLiquidStakingDenom(denom string) bool {
	return strings.HasPrefix(denom, fmt.Sprintf("%s/", LiquidStakeDenomPrefix))
}

func MintDenomToHostDenom(mintDenom string) (string, bool) {
	return strings.CutPrefix(mintDenom, fmt.Sprintf("%s/", LiquidStakeDenomPrefix))
}

func HostDenomToMintDenom(hostDenom string) string {
	return fmt.Sprintf("%s/%s", LiquidStakeDenomPrefix, hostDenom)
}

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
	return fmt.Sprintf("%s.%s", chainID, DelegateICAType)
}

// DefaultRewardsAccountPortOwner generates a rewards ICA port owner given the chain id
// Only Use this function while registering a new chain
func DefaultRewardsAccountPortOwner(chainID string) string {
	return fmt.Sprintf("%s.%s", chainID, RewardsICAType)
}

func (deposit *Deposit) Validate() error {
	if deposit.State != Deposit_DEPOSIT_PENDING &&
		deposit.State != Deposit_DEPOSIT_SENT &&
		deposit.State != Deposit_DEPOSIT_RECEIVED &&
		deposit.State != Deposit_DEPOSIT_DELEGATING {
		return fmt.Errorf(
			"host chain %s deposit has an invalid state: %s",
			deposit.ChainId,
			deposit.State,
		)
	}
	if err := deposit.Amount.Validate(); err != nil {
		return fmt.Errorf("deposit amount is invalid, err: %v", err)
	}

	return nil
}

func (hc *HostChain) Validate() error {
	err := hc.Params.Validate()
	if err != nil {
		return fmt.Errorf("host chain %s validation failed with err, err: %s", hc.ChainId, err)
	}
	if hc.MinimumDeposit.LT(sdk.ZeroInt()) {
		return fmt.Errorf("host chain %s has negative minimum deposit", hc.ChainId)
	}
	if hc.CValue.LT(sdk.ZeroDec()) { // GT limits should be checked by module level params, invariants.
		return fmt.Errorf("host chain %s has c value out of bounds: %d", hc.ChainId, hc.CValue)
	}
	if strings.TrimSpace(hc.ChainId) == "" {
		return fmt.Errorf("chain_id must be non-empty")
	}
	if len(hc.ChainId) > types.MaxChainIDLen {
		return fmt.Errorf("chain_id is too long (max: %d)", types.MaxChainIDLen)
	}
	err = host.ConnectionIdentifierValidator(hc.ConnectionId)
	if err != nil {
		return fmt.Errorf("hostchain connectionID invalid err: %v", err)
	}
	err = host.PortIdentifierValidator(hc.PortId)
	if err != nil {
		return err
	}
	err = host.ChannelIdentifierValidator(hc.ChannelId)
	if err != nil {
		return err
	}
	if hc.DelegationAccount != nil {
		err = hc.DelegationAccount.Validate()
		if err != nil {
			return err
		}
	}
	if hc.RewardsAccount != nil {
		err = hc.RewardsAccount.Validate()
		if err != nil {
			return err
		}
	}
	for _, validator := range hc.Validators {
		err := validator.Validate()
		if err != nil {
			return fmt.Errorf("host chain %s validator is invalid, err: %s", hc.ChainId, err)
		}
	}
	if hc.RewardParams != nil {
		err = hc.RewardParams.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (icaAccount *ICAAccount) Validate() error {
	if icaAccount.ChannelState != ICAAccount_ICA_CHANNEL_CREATING &&
		icaAccount.ChannelState != ICAAccount_ICA_CHANNEL_CREATED {
		return fmt.Errorf("invalid channel state")
	}
	portID, err := icatypes.NewControllerPortID(icaAccount.Owner)
	if err != nil {
		return err
	}
	err = host.PortIdentifierValidator(portID)
	if err != nil {
		return err
	}
	err = icaAccount.Balance.Validate()
	if err != nil {
		return err
	}
	if icaAccount.Address != "" {
		_, _, err = bech32.DecodeAndConvert(icaAccount.Address)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rewardParams *RewardParams) Validate() error {
	_, _, err := bech32.DecodeAndConvert(rewardParams.Destination)
	if err != nil {
		return err
	}
	return sdk.ValidateDenom(rewardParams.Denom)

}
func (params *HostChainLSParams) Validate() error {
	if params.DepositFee.LT(sdk.ZeroDec()) || params.DepositFee.GT(sdk.OneDec()) {
		return fmt.Errorf("host chain lsparams has invalid deposit fee, should be 0<=fee<=1")
	}
	if params.RestakeFee.LT(sdk.ZeroDec()) || params.RestakeFee.GT(sdk.OneDec()) {
		return fmt.Errorf("host chain lsparams has invalid restake fee, should be 0<=fee<=1\"")
	}
	if params.RedemptionFee.LT(sdk.ZeroDec()) || params.RedemptionFee.GT(sdk.OneDec()) {
		return fmt.Errorf("host chain lsparams has invalid redemption fee, should be 0<=fee<=1\"")
	}
	if params.UnstakeFee.LT(sdk.ZeroDec()) || params.UnstakeFee.GT(sdk.OneDec()) {
		return fmt.Errorf("host chain lsparams has invalid unstake fee, should be 0<=fee<=1\"")
	}
	return nil
}

func (validator *Validator) Validate() error {
	if validator.Status != stakingtypes.Unspecified.String() &&
		validator.Status != stakingtypes.Unbonded.String() &&
		validator.Status != stakingtypes.Unbonding.String() &&
		validator.Status != stakingtypes.Bonded.String() {
		return fmt.Errorf(
			"host chain validator %s has an invalid status: %s",
			validator.OperatorAddress,
			validator.Status,
		)
	}

	if validator.Weight.LT(sdk.ZeroDec()) || validator.Weight.GT(sdk.OneDec()) {
		return fmt.Errorf(
			"host chain validator %s has weight out of bounds: %d",
			validator.OperatorAddress,
			validator.Weight)
	}

	if validator.DelegatedAmount.LT(sdk.ZeroInt()) {
		return fmt.Errorf(
			"host chain validator %s has negative delegated amount: %s",
			validator.OperatorAddress,
			validator.DelegatedAmount.String(),
		)
	}

	_, _, err := bech32.DecodeAndConvert(validator.OperatorAddress)
	if err != nil {
		return fmt.Errorf(
			"host chain validator %s is invalid bech32 addr, err: %s",
			validator.OperatorAddress,
			err,
		)
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
	if u.State != Unbonding_UNBONDING_PENDING &&
		u.State != Unbonding_UNBONDING_INITIATED &&
		u.State != Unbonding_UNBONDING_MATURING &&
		u.State != Unbonding_UNBONDING_MATURED &&
		u.State != Unbonding_UNBONDING_CLAIMABLE &&
		u.State != Unbonding_UNBONDING_FAILED {
		return fmt.Errorf(
			"host chain %s unbonding has an invalid state: %s",
			u.ChainId,
			u.State,
		)
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
