package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Parameter store keys
var (
	DefaultLiquidBondDenom = "stk/uxprt"

	// DefaultUnstakeFeeRate is the default Unstake Fee Rate.
	DefaultUnstakeFeeRate = sdk.ZeroDec()

	// DefaultAutocompoundFeeRate is the default fee rate for auto redelegating the stake rewards.
	DefaultAutocompoundFeeRate = sdk.MustNewDecFromStr("0.05")

	// DefaultMinLiquidStakeAmount is the default minimum liquid stake amount.
	DefaultMinLiquidStakeAmount = math.NewInt(1000)

	// Const variables

	// RebalancingTrigger if the maximum difference and needed each redelegation amount exceeds it, asset rebalacing will be executed.
	RebalancingTrigger = math.LegacyNewDecWithPrec(1, 3) // "0.001000000000000000"

	// AutocompoundTrigger If the sum of balance and the upcoming rewards of LiquidStakeProxyAcc exceeds it,
	// the reward is automatically autocompounded, according to the weights.
	AutocompoundTrigger = math.LegacyNewDecWithPrec(1, 3) // "0.001000000000000000"

	// LiquidStakeProxyAcc is a proxy reserve account for delegation and undelegation.
	LiquidStakeProxyAcc = authtypes.NewModuleAddress(ModuleName + "-LiquidStakeProxyAcc")

	// DummyFeeAccountAcc is a dummy fee collection account that should be replaced via params.
	DummyFeeAccountAcc = authtypes.NewModuleAddress(ModuleName + "-FeeAcc")
)

// DefaultParams returns the default liquidstake module parameters.
func DefaultParams() Params {
	return Params{
		WhitelistedValidators: []WhitelistedValidator{},
		LiquidBondDenom:       DefaultLiquidBondDenom,
		UnstakeFeeRate:        DefaultUnstakeFeeRate,
		MinLiquidStakeAmount:  DefaultMinLiquidStakeAmount,
		FeeAccountAddress:     DummyFeeAccountAcc.String(),
		AutocompoundFeeRate:   DefaultAutocompoundFeeRate,
		CwLockedPoolAddress:   "",
		WhitelistAdminAddress: "",
	}
}

// String returns a human-readable string representation of the parameters.
func (p Params) String() string {
	out, _ := json.MarshalIndent(p, "", "")
	return string(out)
}

func (p Params) WhitelistedValsMap() WhitelistedValsMap {
	return GetWhitelistedValsMap(p.WhitelistedValidators)
}

// Validate validates parameters.
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.LiquidBondDenom, validateLiquidBondDenom},
		{p.WhitelistedValidators, validateWhitelistedValidators},
		{p.UnstakeFeeRate, validateUnstakeFeeRate},
		{p.MinLiquidStakeAmount, validateMinLiquidStakeAmount},
		{p.AutocompoundFeeRate, validateAutocompoundFeeRate},
		{p.FeeAccountAddress, validateFeeAccountAddress},
		{p.CwLockedPoolAddress, validateCwLockedPoolAddress},
		{p.WhitelistAdminAddress, validateWhitelistAdminAddress},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

func validateLiquidBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("liquid bond denom cannot be blank")
	}

	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}
	return nil
}

// validateWhitelistedValidators validates liquidstake validator and total weight.
func validateWhitelistedValidators(i interface{}) error {
	wvs, ok := i.([]WhitelistedValidator)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	valsMap := map[string]struct{}{}
	for _, wv := range wvs {
		_, valErr := sdk.ValAddressFromBech32(wv.ValidatorAddress)
		if valErr != nil {
			return valErr
		}

		if wv.TargetWeight.IsNil() {
			return fmt.Errorf("liquidstake validator target weight must not be nil")
		}

		if !wv.TargetWeight.IsPositive() {
			return fmt.Errorf("liquidstake validator target weight must be positive: %s", wv.TargetWeight)
		}

		if _, ok := valsMap[wv.ValidatorAddress]; ok {
			return fmt.Errorf("liquidstake validator cannot be duplicated: %s", wv.ValidatorAddress)
		}
		valsMap[wv.ValidatorAddress] = struct{}{}
	}
	return nil
}

func validateUnstakeFeeRate(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("unstake fee rate must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("unstake fee rate must not be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("unstake fee rate too large: %s", v)
	}

	return nil
}

func validateMinLiquidStakeAmount(i interface{}) error {
	v, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("min liquid stake amount must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("min liquid stake amount must not be negative: %s", v)
	}

	return nil
}

func validateAutocompoundFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("autocompound fee rate must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("autocompound fee rate must not be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("autocompound fee rate too large: %s", v)
	}

	return nil
}

func validateFeeAccountAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("cannot convert fee account address to bech32, invalid address: %s, err: %v", v, err)
	}
	return nil
}

func validateCwLockedPoolAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// allow empty address
	if len(v) == 0 {
		return nil
	}

	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("cannot convert cw contract address to bech32, invalid address: %s, err: %v", v, err)
	}
	return nil
}

func validateWhitelistAdminAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// allow empty address
	if len(v) == 0 {
		return nil
	}

	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("cannot convert whitelist admin address to bech32, invalid address: %s, err: %v", v, err)
	}
	return nil
}
