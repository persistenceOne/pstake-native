package types

import (
	"sigs.k8s.io/yaml"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// ParamKeyTable for liquidstakeibc module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the liquidstakeibc module
func NewParams() Params {
	return Params{}
}

// DefaultParams is the default parameter configuration for the liquidstakeibc module
func DefaultParams() Params {
	return NewParams()
}

// Validate all liquidstakeibc module parameters
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}
