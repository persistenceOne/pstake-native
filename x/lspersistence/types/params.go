package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Parameter store keys
var (
	KeyLiquidBondDenom        = []byte("LiquidBondDenom")
	KeyWhitelistedValidators  = []byte("WhitelistedValidators")
	KeyUnstakeFeeRate         = []byte("UnstakeFeeRate")
	KeyStakeFeeRate           = []byte("StakeFeeRate")
	KeyRestakeFeeRate         = []byte("RestakeFeeRate")
	KeyRedemptionFeeRate      = []byte("RedemptionFeeRate")
	KeyMinLiquidStakingAmount = []byte("MinLiquidStakingAmount")
	KeyAdminAddress           = []byte("AdminAddress")
	KeyFeeAddress             = []byte("FeeAddress")

	DefaultLiquidBondDenom = "stk/uxprt"

	// DefaultMinLiquidStakingAmount is the default minimum liquid staking amount.
	DefaultMinLiquidStakingAmount = sdk.NewInt(10000)

	// Const variables

	// RebalancingTrigger if the maximum difference and needed each redelegation amount exceeds it, asset rebalacing will be executed.
	RebalancingTrigger = sdk.NewDecWithPrec(1, 3) // "0.001000000000000000"

	// RewardTrigger If the sum of balance and the upcoming rewards of LiquidStakingProxyAcc exceeds it, the reward is automatically withdrawn and re-stake according to the weights.
	RewardTrigger = sdk.NewDecWithPrec(1, 3) // "0.001000000000000000"

	// LiquidStakingProxyAcc is a proxy reserve account for delegation and undelegation.
	LiquidStakingProxyAcc = authtypes.NewModuleAddress(ModuleName + "-LiquidStakingProxyAcc")
)

var _ paramstypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns the default liquidstaking module parameters.
func DefaultParams() Params {
	return Params{
		WhitelistedValidators:  []WhitelistedValidator{},
		LiquidBondDenom:        DefaultLiquidBondDenom,
		UnstakeFeeRate:         sdk.MustNewDecFromStr("0"),
		StakeFeeRate:           sdk.MustNewDecFromStr("0"),
		RestakeFeeRate:         sdk.MustNewDecFromStr("0.05"),
		RedemptionFeeRate:      sdk.MustNewDecFromStr("0.025"),
		MinLiquidStakingAmount: DefaultMinLiquidStakingAmount,
		AdminAddress:           authtypes.NewModuleAddress("dummy").String(),
		FeeAddress:             authtypes.NewModuleAddress("dummy").String(),
	}
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyLiquidBondDenom, &p.LiquidBondDenom, validateLiquidBondDenom),
		paramstypes.NewParamSetPair(KeyWhitelistedValidators, &p.WhitelistedValidators, validateWhitelistedValidators),
		paramstypes.NewParamSetPair(KeyUnstakeFeeRate, &p.UnstakeFeeRate, validateUnstakeFeeRate),
		paramstypes.NewParamSetPair(KeyStakeFeeRate, &p.StakeFeeRate, validateStakeFeeRate),
		paramstypes.NewParamSetPair(KeyRestakeFeeRate, &p.RestakeFeeRate, validateRestakeFeeRate),
		paramstypes.NewParamSetPair(KeyRedemptionFeeRate, &p.RedemptionFeeRate, validateRedemptionFeeRate),
		paramstypes.NewParamSetPair(KeyMinLiquidStakingAmount, &p.MinLiquidStakingAmount, validateMinLiquidStakingAmount),
		paramstypes.NewParamSetPair(KeyAdminAddress, &p.AdminAddress, validateAdminAddress),
		paramstypes.NewParamSetPair(KeyFeeAddress, &p.FeeAddress, validateFeeAddress),
	}
}

// String returns a human-readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
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
		{p.StakeFeeRate, validateStakeFeeRate},
		{p.RestakeFeeRate, validateRestakeFeeRate},
		{p.RedemptionFeeRate, validateRedemptionFeeRate},
		{p.MinLiquidStakingAmount, validateMinLiquidStakingAmount},
		{p.AdminAddress, validateAdminAddress},
		{p.FeeAddress, validateFeeAddress},
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

// validateWhitelistedValidators validates liquidstaking validator and total weight.
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
			return fmt.Errorf("liquidstaking validator target weight must not be nil")
		}

		if !wv.TargetWeight.IsPositive() {
			return fmt.Errorf("liquidstaking validator target weight must be positive: %s", wv.TargetWeight)
		}

		if _, ok := valsMap[wv.ValidatorAddress]; ok {
			return fmt.Errorf("liquidstaking validator cannot be duplicated: %s", wv.ValidatorAddress)
		}
		valsMap[wv.ValidatorAddress] = struct{}{}
	}
	return nil
}

func validateUnstakeFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
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

func validateStakeFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("stake fee rate must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("stake fee rate must not be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("stake fee rate too large: %s", v)
	}

	return nil
}

func validateRestakeFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("restake fee rate must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("restake fee rate must not be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("restake fee rate too large: %s", v)
	}

	return nil
}

func validateRedemptionFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("redemption fee rate must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("redemption fee rate must not be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("redemption fee rate too large: %s", v)
	}

	return nil
}

func validateMinLiquidStakingAmount(i interface{}) error {
	v, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("min liquid staking amount must not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("min liquid staking amount must not be negative: %s", v)
	}

	return nil
}

func validateAdminAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("cannot convert admin address to bech32, invalid address: %s, err: %v", v, err)
	}
	return nil
}

func validateFeeAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("cannot convert fee address to bech32, invalid address: %s, err: %v", v, err)
	}
	return nil
}
