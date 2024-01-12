package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	DefaultAdminAddress     = authtypes.NewModuleAddress("placeholder") // will be set manually upon module initialisation
	DefaultFeeAddress       = authtypes.NewModuleAddress("placeholder") // will be set manually upon module initialisation
	DefaultUpperCValueLimit = sdktypes.MustNewDecFromStr("1.1")
	DefaultLowerCValueLimit = sdktypes.MustNewDecFromStr("0.85")
)

// NewParams creates a new Params object
func NewParams(
	adminAddress string,
	feeAddress string,
	upperCValueLimit sdktypes.Dec,
	lowerCValueLimit sdktypes.Dec,
) Params {
	return Params{
		AdminAddress:     adminAddress,
		FeeAddress:       feeAddress,
		UpperCValueLimit: upperCValueLimit,
		LowerCValueLimit: lowerCValueLimit,
	}
}

// DefaultParams returns the default set of parameters of the module
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAddress.String(),
		DefaultFeeAddress.String(),
		DefaultUpperCValueLimit,
		DefaultLowerCValueLimit,
	)
}

// Validate all liquidstakeibc module parameters
func (p *Params) Validate() error {
	if _, err := sdktypes.AccAddressFromBech32(p.AdminAddress); err != nil {
		return err
	}
	if _, err := sdktypes.AccAddressFromBech32(p.FeeAddress); err != nil {
		return err
	}
	if p.LowerCValueLimit.GT(sdktypes.OneDec()) || p.LowerCValueLimit.GTE(p.UpperCValueLimit) {
		return ErrInvalidParams.Wrapf("LowerCValue limit should be less than both 1 and UpperCValue limit, lowerCValue: %s, UpperCValue: %s", p.LowerCValueLimit, p.UpperCValueLimit)
	}

	return nil
}
