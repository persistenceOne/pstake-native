package types

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ghodss/yaml"

	epochsTypes "github.com/persistenceOne/pstake-native/x/epochs/types"
)

// Default constants for period, mint and staking denom
const (
	DefaultPeriod       = time.Minute * 1 // 6 hours //TODO : Change back to 6 hours
	DefaultMintDenom    = "ustkatom"
	DefaultStakingDenom = "uatom"
	MinAmount           = 5000000
	MaxAmount           = 100000000000
	MinReward           = 1000
)

// DefaultBondDenom is a default bond denom param
var (
	DefaultBondDenom = []string{"uatom"}
)

// Parameter store key
var (
	KeyMinMintingAmount                  = []byte("MinMintingAmount")
	KeyMaxMintingAmount                  = []byte("MaxMintingAmount")
	KeyMinBurningAmount                  = []byte("MinBurningAmount")
	KeyMaxBurningAmount                  = []byte("MaxBurningAmount")
	KeyMinReward                         = []byte("MinReward")
	KeyMaxValidatorToDelegate            = []byte("MaxValidatorToDelegate")
	KeyWeightedDeveloperRewardsReceivers = []byte("WeightedDeveloperRewardsReceivers")
	KeyDistributionProportion            = []byte("DistributionProportion")
	KeyCosmosProposalParams              = []byte("CosmosProposalParams")
	KeyModuleEnabled                     = []byte("ModuleEnabled")
	KeyStakingEpochIdentifier            = []byte("StakeEpochIdentifier")
	KeyCustodialAddress                  = []byte("CustodialAddress")
	KeyUndelegateEpochIdentifier         = []byte("UndelegateEpochIdentifier")
	KeyChunkSize                         = []byte("ChunkSize")
	KeyBondDenom                         = []byte("BondDenom")
	KeyStakingDenom                      = []byte("StakingDenom")
	KeyMintDenom                         = []byte("MintDenom")
	KeyRetryLimit                        = []byte("RetryLimit")
	KeyRewardEpochIdentifier             = []byte("RewardEpochIdentifier")
)

// ParamKeyTable - Key declaration for parameters
func ParamKeyTable() paramsTypes.KeyTable {
	return paramsTypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams default parameters for deposits
func DefaultParams() Params {
	return Params{
		MinMintingAmount:       sdk.NewInt64Coin(DefaultStakingDenom, MinAmount),
		MaxMintingAmount:       sdk.NewInt64Coin(DefaultStakingDenom, MaxAmount),
		MinBurningAmount:       sdk.NewInt64Coin(DefaultStakingDenom, MinAmount),
		MaxBurningAmount:       sdk.NewInt64Coin(DefaultStakingDenom, MaxAmount),
		MinReward:              sdk.NewInt64Coin(DefaultStakingDenom, MinReward),
		MaxValidatorToDelegate: 3,
		WeightedDeveloperRewardsReceivers: []WeightedAddress{
			{
				Address: "persistence1g5lz0gq98y8tav477dltxgpdft0wr9rmqt7mvu",
				Weight:  sdk.NewDecWithPrec(10, 1),
			},
		},
		DistributionProportion: DistributionProportions{
			ValidatorRewards: sdk.NewDecWithPrec(5, 2),
			DeveloperRewards: sdk.NewDecWithPrec(5, 2),
		},
		CosmosProposalParams: CosmosChainProposalParams{
			ChainID:              "cosmoshub-4", //TODO use these as conditions for proposals
			ReduceVotingPeriodBy: DefaultPeriod,
		},
		ModuleEnabled:             false, //TODO : Make false before launch
		CustodialAddress:          "cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2",
		StakingEpochIdentifier:    "uatom",
		UndelegateEpochIdentifier: "undelegate",
		RewardEpochIdentifier:     "reward",
		ChunkSize:                 5,
		BondDenoms:                []string{DefaultStakingDenom},
		StakingDenom:              DefaultStakingDenom,
		MintDenom:                 DefaultMintDenom,
		RetryLimit:                10,
	}
}

// Validate runs basic stateless validity checks
func (p Params) Validate() error {
	if err := validateAmount(p.MinMintingAmount); err != nil {
		return err
	}
	if err := validateAmount(p.MaxMintingAmount); err != nil {
		return err
	}
	if err := validateAmount(p.MinBurningAmount); err != nil {
		return err
	}
	if err := validateAmount(p.MaxBurningAmount); err != nil {
		return err
	}
	if err := validateMaxValidatorToDelegate(p.MaxValidatorToDelegate); err != nil {
		return err
	}
	if err := validateWeightedDeveloperRewardsReceivers(p.WeightedDeveloperRewardsReceivers); err != nil {
		return err
	}
	if err := validateDistributionProportion(p.DistributionProportion); err != nil {
		return err
	}
	if err := validateCosmosProposalParams(p.CosmosProposalParams); err != nil {
		return err
	}
	if err := validateModuleEnabled(p.ModuleEnabled); err != nil {
		return err
	}
	if err := epochsTypes.ValidateEpochIdentifierInterface(p.StakingEpochIdentifier); err != nil {
		return err
	}
	if err := validateCustodialAddress(p.CustodialAddress); err != nil {
		return err
	}
	if err := epochsTypes.ValidateEpochIdentifierInterface(p.UndelegateEpochIdentifier); err != nil {
		return err
	}
	if err := validateWithdrawRewardsChunkSize(p.ChunkSize); err != nil {
		return err
	}
	if err := validateBondDenom(p.BondDenoms); err != nil {
		return err
	}
	if err := validateStakingDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateRetryLimit(p.RetryLimit); err != nil {
		return err
	}
	if err := epochsTypes.ValidateEpochIdentifierInterface(p.RewardEpochIdentifier); err != nil {
		return err
	}
	return nil
}

// String implements stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// of cosmos module's parameters.
func (p *Params) ParamSetPairs() paramsTypes.ParamSetPairs {
	return paramsTypes.ParamSetPairs{
		paramsTypes.NewParamSetPair(KeyMinMintingAmount, &p.MinMintingAmount, validateAmount),
		paramsTypes.NewParamSetPair(KeyMaxMintingAmount, &p.MaxMintingAmount, validateAmount),
		paramsTypes.NewParamSetPair(KeyMinBurningAmount, &p.MinBurningAmount, validateAmount),
		paramsTypes.NewParamSetPair(KeyMaxBurningAmount, &p.MaxBurningAmount, validateAmount),
		paramsTypes.NewParamSetPair(KeyMinReward, &p.MinReward, validateAmount),
		paramsTypes.NewParamSetPair(KeyMaxValidatorToDelegate, &p.MaxValidatorToDelegate, validateMaxValidatorToDelegate),
		paramsTypes.NewParamSetPair(KeyWeightedDeveloperRewardsReceivers, &p.WeightedDeveloperRewardsReceivers, validateWeightedDeveloperRewardsReceivers),
		paramsTypes.NewParamSetPair(KeyDistributionProportion, &p.DistributionProportion, validateDistributionProportion),
		paramsTypes.NewParamSetPair(KeyCosmosProposalParams, &p.CosmosProposalParams, validateCosmosProposalParams),
		paramsTypes.NewParamSetPair(KeyModuleEnabled, &p.ModuleEnabled, validateModuleEnabled),
		paramsTypes.NewParamSetPair(KeyStakingEpochIdentifier, &p.StakingEpochIdentifier, epochsTypes.ValidateEpochIdentifierInterface),
		paramsTypes.NewParamSetPair(KeyCustodialAddress, &p.CustodialAddress, validateCustodialAddress),
		paramsTypes.NewParamSetPair(KeyUndelegateEpochIdentifier, &p.UndelegateEpochIdentifier, epochsTypes.ValidateEpochIdentifierInterface),
		paramsTypes.NewParamSetPair(KeyChunkSize, &p.ChunkSize, validateWithdrawRewardsChunkSize),
		paramsTypes.NewParamSetPair(KeyBondDenom, &p.BondDenoms, validateBondDenom),
		paramsTypes.NewParamSetPair(KeyStakingDenom, &p.StakingDenom, validateStakingDenom),
		paramsTypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramsTypes.NewParamSetPair(KeyRetryLimit, &p.RetryLimit, validateRetryLimit),
		paramsTypes.NewParamSetPair(KeyRewardEpochIdentifier, &p.RewardEpochIdentifier, epochsTypes.ValidateEpochIdentifierInterface),
	}
}

// GetBondDenomOf returns the bond denom if present
func (p Params) GetBondDenomOf(s string) (string, error) {
	for _, element := range p.BondDenoms {
		if element == s {
			return element, nil
		}
	}
	return "", ErrInvalidBondDenom
}

// Equal returns if the other param set matches or not
func (p Params) Equal(other Params) bool {
	for i := range p.WeightedDeveloperRewardsReceivers {
		if p.WeightedDeveloperRewardsReceivers[i] != other.WeightedDeveloperRewardsReceivers[i] {
			return false
		}
	}

	for i := range p.BondDenoms {
		if p.BondDenoms[i] != other.BondDenoms[i] {
			return false
		}
	}
	return p.MintDenom == other.MintDenom &&
		p.CustodialAddress == other.CustodialAddress &&
		p.MinMintingAmount.IsEqual(other.MinMintingAmount) &&
		p.MaxMintingAmount.IsEqual(other.MaxMintingAmount) &&
		p.MinBurningAmount.IsEqual(other.MinBurningAmount) &&
		p.MaxBurningAmount.IsEqual(other.MaxBurningAmount) &&
		p.MaxValidatorToDelegate == other.MaxValidatorToDelegate &&
		p.DistributionProportion == other.DistributionProportion &&
		p.CosmosProposalParams == other.CosmosProposalParams &&
		p.CustodialAddress == other.CustodialAddress &&
		p.ModuleEnabled == other.ModuleEnabled &&
		p.StakingEpochIdentifier == other.StakingEpochIdentifier &&
		p.ChunkSize == other.ChunkSize &&
		p.UndelegateEpochIdentifier == other.UndelegateEpochIdentifier &&
		p.RetryLimit == other.RetryLimit &&
		p.RewardEpochIdentifier == other.RewardEpochIdentifier
}

func validateAmount(i interface{}) error {
	coin, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if coin.IsNegative() {
		return fmt.Errorf("amount cannot be negative")
	}
	return nil
}

func validateMaxValidatorToDelegate(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func ValidateValidatorSetCosmosChain(i interface{}) error {
	v, ok := i.([]WeightedAddressAmount)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// fund community pool when rewards address is empty
	if len(v) == 0 {
		return nil
	}

	weightSum := sdk.NewDec(0)
	for i, w := range v {
		// we allow address to be "" to go to community pool
		if w.Address != "" {
			_, err := ValAddressFromBech32(w.Address, Bech32PrefixValAddr)
			if err != nil {
				return fmt.Errorf("invalid address at %dth", i)
			}
		}
		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(sdk.NewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
		if w.Amount.IsNegative() {
			return fmt.Errorf("non-positive current delegation amount at %dth", i)
		}
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}

func ValidateValidatorSetNativeChain(i interface{}) error {
	v, ok := i.([]WeightedAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// fund community pool when rewards address is empty
	if len(v) == 0 {
		return nil
	}

	weightSum := sdk.NewDec(0)
	for i, w := range v {
		// we allow address to be "" to go to community pool
		if w.Address != "" {
			_, err := sdk.ValAddressFromBech32(w.Address)
			if err != nil {
				return fmt.Errorf("invalid address at %dth", i)
			}
		}
		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(sdk.NewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}

func validateWeightedDeveloperRewardsReceivers(i interface{}) error {
	v, ok := i.([]WeightedAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// fund community pool when rewards address is empty
	if len(v) == 0 {
		return nil
	}

	weightSum := sdk.NewDec(0)
	for i, w := range v {
		// we allow address to be "" to go to community pool
		if w.Address != "" {
			_, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return fmt.Errorf("invalid address at %dth", i)
			}
		}
		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(sdk.NewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}

func validateDistributionProportion(i interface{}) error {
	v, ok := i.(DistributionProportions)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.ValidatorRewards.IsNegative() {
		return errors.New("validator rewards distribution ratio should not be negative")
	}

	if v.DeveloperRewards.IsNegative() {
		return errors.New("developer rewards distribution ratio should not be negative")
	}

	totalProportions := v.ValidatorRewards.Add(v.DeveloperRewards)

	if !totalProportions.Equal(sdk.NewDecWithPrec(1, 1)) {
		return errors.New("total distributions ratio should be 0.1")
	}

	return nil
}

func validateCosmosProposalParams(i interface{}) error {
	v, ok := i.(CosmosChainProposalParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.ReduceVotingPeriodBy <= 0 {
		return fmt.Errorf("incorrect voting Period %T", i)
	}

	return nil
}

func validateModuleEnabled(i interface{}) error {
	_, ok := i.(bool)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateCustodialAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v != "" {
		_, err := AccAddressFromBech32(v, Bech32Prefix)
		if err != nil {
			return fmt.Errorf("invalid custodial address")
		}
	}
	return nil
}

func validateWithdrawRewardsChunkSize(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid ")
	}
	if v <= 0 {
		return fmt.Errorf("non-positive chunk size in invalid : %d", i)
	}
	return nil
}

func validateBondDenom(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("bond denom cannot be empty")
	}
	return nil
}

func validateStakingDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("staking denom cannot be empty")
	}
	return nil
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("mint denom cannot be empty")
	}
	return nil
}

func validateRetryLimit(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("retry limit must be non negative")
	}
	return nil
}
