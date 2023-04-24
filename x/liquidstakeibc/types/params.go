package types

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"sigs.k8s.io/yaml"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultFeeAddress string = "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld" // TODO: Use correct address on launch
)

var (
	KeyFeeAddress = []byte("FeeAddress")
)

// ParamKeyTable for liquidstakeibc module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	feeAddress string,
) Params {

	return Params{
		FeeAddress: feeAddress,
	}
}

// DefaultParams returns the default set of parameters of the module
func DefaultParams() Params {
	return NewParams(
		DefaultFeeAddress,
	)
}

// Validate all liquidstakeibc module parameters
func (p Params) Validate() error {
	if err := isAddress(p.FeeAddress); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFeeAddress, &p.FeeAddress, isAddress),
	}
}

func isAddress(i interface{}) error {
	val, ok := i.(string)
	if !ok {
		return fmt.Errorf("parameter is not valid: %T", i)
	}

	_, err := sdktypes.AccAddressFromBech32(val)
	if err != nil {
		return fmt.Errorf("parameter %s must be a valid address", val)
	}

	return nil
}
