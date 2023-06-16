package types

import (
	"fmt"
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"sigs.k8s.io/yaml"
)

const (
	DefaultAdminAddress     string = "persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr" // TODO: Use correct address on launch
	DefaultFeeAddress       string = "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld" // TODO: Use correct address on launch
	DefaultUpperCValueLimit string = "1.1"
	DefaultLowerCValueLimit string = "0.85"
)

// NewParams creates a new Params object
func NewParams(
	adminAddress string,
	feeAddress string,
	upperCValueLimit string,
	lowerCValueLimit string,
) Params {

	upperLimit, _ := sdktypes.NewDecFromStr(upperCValueLimit)
	lowerLimit, _ := sdktypes.NewDecFromStr(lowerCValueLimit)

	return Params{
		AdminAddress:     adminAddress,
		FeeAddress:       feeAddress,
		UpperCValueLimit: upperLimit,
		LowerCValueLimit: lowerLimit,
	}
}

// DefaultParams returns the default set of parameters of the module
func DefaultParams() Params {
	return NewParams(
		DefaultAdminAddress,
		DefaultFeeAddress,
		DefaultUpperCValueLimit,
		DefaultLowerCValueLimit,
	)
}

// Validate all liquidstakeibc module parameters
func (p *Params) Validate() error {
	if err := isAddress(p.AdminAddress); err != nil {
		return err
	}
	if err := isAddress(p.FeeAddress); err != nil {
		return err
	}
	if err := isGTOne(p.UpperCValueLimit); err != nil {
		return err
	}
	if err := isLTOne(p.LowerCValueLimit); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p *Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// checks

func isAddress(i interface{}) error {
	val, ok := i.(string)
	if !ok {
		return fmt.Errorf("parameter is not valid: %T", i)
	}

	if len(strings.TrimSpace(val)) == 0 {
		return fmt.Errorf("empty address string is not allowed")
	}

	_, err := sdktypes.GetFromBech32(val, "persistence")
	if err != nil {
		return err
	}

	return nil
}

func isGTOne(i interface{}) error {
	val, ok := i.(sdktypes.Dec)
	if !ok {
		return fmt.Errorf("parameter is not valid: %T", i)
	}

	if !val.GT(sdktypes.OneDec()) {
		return fmt.Errorf("upper limit must be higher than 1")
	}

	return nil
}

func isLTOne(i interface{}) error {
	val, ok := i.(sdktypes.Dec)
	if !ok {
		return fmt.Errorf("parameter is not valid: %T", i)
	}

	if !val.LT(sdktypes.OneDec()) {
		return fmt.Errorf("lower limit must be lower than 1")
	}

	return nil
}
