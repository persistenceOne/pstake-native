package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"gopkg.in/yaml.v2"
)

var DefaultAdmin = authtypes.NewModuleAddress(govtypes.ModuleName)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		Admin: DefaultAdmin.String(),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params
func (p Params) Validate() error {
	_, err := sdk.AccAddressFromBech32(p.Admin)
	return err
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
