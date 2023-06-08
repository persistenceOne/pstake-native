package types

import (
	"fmt"
)

// Validate performs basic validation of supply genesis data returning an
// error for any failed validation criteria.
func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	hostChainMap := make(map[string]HostChain)
	for _, hc := range gs.HostChains {
		if _, ok := hostChainMap[hc.ChainId]; ok {
			return fmt.Errorf("duplicated host chain: %s", hc.ChainId)
		}
		hostChainMap[hc.ChainId] = *hc

		if err := hc.Validate(); err != nil {
			return err
		}
	}

	for _, deposit := range gs.Deposits {
		if err := deposit.Validate(); err != nil {
			return err
		}
		hc, ok := hostChainMap[deposit.ChainId]
		if !ok {
			return fmt.Errorf("deposit for chain %s doesnt have a valid chain id", deposit.ChainId)
		}
		if hc.HostDenom != deposit.Amount.Denom {
			return fmt.Errorf(
				"deposit for chain %s doesnt have the correct host chain denom: %s, should be %s",
				deposit.ChainId,
				deposit.Amount.Denom,
				hc.HostDenom,
			)
		}
	}
	for _, unbonding := range gs.Unbondings {
		hc, ok := hostChainMap[unbonding.ChainId]
		if !ok {
			return fmt.Errorf("unbonding for chain %s doesnt have a valid chain id", unbonding.ChainId)
		}
		if hc.MintDenom() == unbonding.BurnAmount.Denom {
			return fmt.Errorf(
				"unbonding for chain %s doesnt have the correct burn amount denom: %s, should be %s",
				hc.ChainId,
				unbonding.BurnAmount.Denom,
				hc.MintDenom(),
			)
		}
		if hc.HostDenom == unbonding.UnbondAmount.Denom {
			return fmt.Errorf(
				"unbonding for chain %s doesnt have the correct host chain denom: %s, should be %s",
				hc.ChainId,
				unbonding.UnbondAmount.Denom,
				hc.HostDenom,
			)
		}
		if err := unbonding.Validate(); err != nil {
			return err
		}
	}
	for _, userUnbonding := range gs.UserUnbondings {
		hc, ok := hostChainMap[userUnbonding.ChainId]
		if !ok {
			return fmt.Errorf("user unbonding for chain %s doesnt have a valid chain id", userUnbonding.ChainId)
		}
		if hc.MintDenom() == userUnbonding.StkAmount.Denom {
			return fmt.Errorf(
				"user unbonding for chain %s doesnt have the correct mint amount denom: %s, should be %s",
				hc.ChainId,
				userUnbonding.StkAmount.Denom,
				hc.MintDenom(),
			)
		}
		if hc.HostDenom == userUnbonding.UnbondAmount.Denom {
			return fmt.Errorf(
				"user unbonding for chain %s doesnt have the correct host chain denom: %s, should be %s",
				hc.ChainId,
				userUnbonding.UnbondAmount.Denom,
				hc.HostDenom,
			)
		}
		if err := userUnbonding.Validate(); err != nil {
			return err
		}
	}
	for _, valUnbonding := range gs.ValidatorUnbondings {
		hc, ok := hostChainMap[valUnbonding.ChainId]
		if !ok {
			return fmt.Errorf("validator unbonding for chain %s doesnt have a valid chain id", valUnbonding.ChainId)
		}
		if hc.HostDenom == valUnbonding.Amount.Denom {
			return fmt.Errorf(
				"validator unbonding for chain %s doesnt have the correct host chain denom: %s, should be %s",
				hc.ChainId,
				valUnbonding.Amount.Denom,
				hc.HostDenom,
			)
		}

		if err := valUnbonding.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// DefaultGenesisState returns a default liquidstakeibc module genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:              DefaultParams(),
		HostChains:          []*HostChain{},
		Deposits:            []*Deposit{},
		Unbondings:          []*Unbonding{},
		UserUnbondings:      []*UserUnbonding{},
		ValidatorUnbondings: []*ValidatorUnbonding{},
	}
}
