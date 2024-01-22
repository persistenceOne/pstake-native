package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	DefaultAdminAddress = authtypes.NewModuleAddress("placeholder") // will be set manually upon module initialisation
	DefaultFeeAddress   = authtypes.NewModuleAddress("placeholder") // will be set manually upon module initialisation
)

// NewParams creates a new Params object
func NewParams(adminAddress string, feeAddress string) Params {
	return Params{
		AdminAddress: adminAddress,
		FeeAddress:   feeAddress,
	}
}

// DefaultParams returns the default set of parameters of the module
func DefaultParams() Params {
	return NewParams(DefaultAdminAddress.String(), DefaultFeeAddress.String())
}

// Validate all liquidstakeibc module parameters
func (p *Params) Validate() error {
	if _, err := sdktypes.AccAddressFromBech32(p.AdminAddress); err != nil {
		return err
	}
	if _, err := sdktypes.AccAddressFromBech32(p.FeeAddress); err != nil {
		return err
	}

	return nil
}
