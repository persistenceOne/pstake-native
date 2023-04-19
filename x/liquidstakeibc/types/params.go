package types

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"sigs.k8s.io/yaml"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultDepositFee    string = "0.00"
	DefaultRestakeFee    string = "0.05"
	DefaultUnstakeFee    string = "0.00"
	DefaultRedemptionFee string = "0.005"
	DefaultFeeAddress    string = "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld" // TODO: Use correct address on launch
)

var (
	KeyDepositFee    = []byte("DepositFee")
	KeyRestakeFee    = []byte("RestakeFee")
	KeyUnstakeFee    = []byte("UnstakeFee")
	KeyRedemptionFee = []byte("RedemptionFee")
	KeyFeeAddress    = []byte("FeeAddress")
)

// ParamKeyTable for liquidstakeibc module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	depositFee string,
	restakeFee string,
	unstakeFee string,
	redemptionFee string,
	feeAddress string,
) Params {

	depositFeeDec, _ := sdktypes.NewDecFromStr(depositFee)
	restakeFeeDec, _ := sdktypes.NewDecFromStr(restakeFee)
	unstakeFeeDec, _ := sdktypes.NewDecFromStr(unstakeFee)
	redemptionFeeDec, _ := sdktypes.NewDecFromStr(redemptionFee)

	return Params{
		DepositFee:    depositFeeDec,
		RestakeFee:    restakeFeeDec,
		UnstakeFee:    unstakeFeeDec,
		RedemptionFee: redemptionFeeDec,
		FeeAddress:    feeAddress,
	}
}

// DefaultParams returns the default set of parameters of the module
func DefaultParams() Params {
	return NewParams(
		DefaultDepositFee,
		DefaultRestakeFee,
		DefaultUnstakeFee,
		DefaultRedemptionFee,
		DefaultFeeAddress,
	)
}

// Validate all liquidstakeibc module parameters
func (p Params) Validate() error {
	if err := isPositive(p.DepositFee); err != nil {
		return err
	}
	if err := isPositive(p.RestakeFee); err != nil {
		return err
	}
	if err := isPositive(p.UnstakeFee); err != nil {
		return err
	}
	if err := isPositive(p.RedemptionFee); err != nil {
		return err
	}
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
		paramtypes.NewParamSetPair(KeyDepositFee, &p.DepositFee, isPositive),
		paramtypes.NewParamSetPair(KeyRestakeFee, &p.RestakeFee, isPositive),
		paramtypes.NewParamSetPair(KeyUnstakeFee, &p.UnstakeFee, isPositive),
		paramtypes.NewParamSetPair(KeyRedemptionFee, &p.RedemptionFee, isPositive),
		paramtypes.NewParamSetPair(KeyFeeAddress, &p.FeeAddress, isAddress),
	}
}

func isPositive(i interface{}) error {
	val, ok := i.(sdktypes.Dec)
	if !ok {
		return fmt.Errorf("parameter is not valid: %T", i)
	}

	if val.LT(sdktypes.NewDec(0)) {
		return fmt.Errorf("parameter %d must be positive", val)
	}
	return nil
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
