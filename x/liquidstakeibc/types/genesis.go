package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Validate performs basic validation of supply genesis data returning an
// error for any failed validation criteria.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	hostChainMap := make(map[string]HostChain)
	for _, hc := range gs.HostChains {
		if _, ok := hostChainMap[hc.ChainId]; ok {
			return fmt.Errorf("duplicated host chain: %s", hc.ChainId)
		}
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
	}

	for _, deposit := range gs.Deposits {
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
		if _, ok := hostChainMap[deposit.ChainId]; !ok {
			return fmt.Errorf("deposit for chain %s doesnt have a valid chain id", deposit.ChainId)
		}
		if hc, _ := hostChainMap[deposit.ChainId]; hc.HostDenom != deposit.Amount.Denom { //nolint:gosimple
			return fmt.Errorf(
				"deposit for chain %s doesnt have the correct host chain denom: %s, should be %s",
				deposit.ChainId,
				deposit.Amount.Denom,
				hc.HostDenom,
			)
		}
	}

	return nil
}

// DefaultGenesisState returns a default liquidstakeibc module genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		HostChains: []*HostChain{},
		Deposits:   []*Deposit{},
	}
}
