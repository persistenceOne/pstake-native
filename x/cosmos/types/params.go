package types

import (
	"errors"
	"fmt"
	epochsTypes "github.com/persistenceOne/pstake-native/x/epochs/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ghodss/yaml"
)

const (
	DefaultPeriod    time.Duration = time.Minute * 1 // 6 hours //TODO : Change back to 6 hours
	DefaultMintDenom string        = "ustkxprt"
)

var (
	DefaultBondDenom = []string{"uatom"}
)

var (
	KeyMinMintingAmount                  = []byte("MinMintingAmount")
	KeyMaxMintingAmount                  = []byte("MaxMintingAmount")
	KeyMinBurningAmount                  = []byte("MinBurningAmount")
	KeyMaxBurningAmount                  = []byte("MaxBurningAmount")
	KeyMaxValidatorToDelegate            = []byte("MaxValidatorToDelegate")
	KeyValidatorSetCosmosChain           = []byte("ValidatorSetCosmosChain")
	KeyValidatorSetNativeChain           = []byte("ValidatorSetNativeChain")
	KeyWeightedDeveloperRewardsReceivers = []byte("WeightedDeveloperRewardsReceivers")
	KeyDistributionProportion            = []byte("DistributionProportion")
	KeyEpochs                            = []byte("Epochs")
	KeyMaxIncomingAndOutgoingTxns        = []byte("MaxIncomingAndOutgoingTxns")
	KeyCosmosProposalParams              = []byte("CosmosProposalParams")
	KeyDelegationThreshold               = []byte("DelegationThreshold")
	KeyModuleEnabled                     = []byte("ModuleEnabled")
	KeyStakingEpochIdentifier            = []byte("StakingEpochIdentifier")
	KeyCustodialAddress                  = []byte("CustodialAddress")
	KeyUndelegateEpochIdentifier         = []byte("UndelegateEpochIdentifier")
	KeyChunkSize                         = []byte("ChunkSize")
	KeyBondDenom                         = []byte("BondDenom")
	KeyMintDenom                         = []byte("MintDenom")
	KeyMultisigThreshold                 = []byte("MultisigThreshold")
	KeyRetryLimit                        = []byte("RetryLimit")
)

func ParamKeyTable() paramsTypes.KeyTable {
	return paramsTypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(minMintingAmount sdk.Coin, maxMintingAmount sdk.Coin, minBurningAmount sdk.Coin, maxBurningAmount sdk.Coin,
	maxValidatorToDelegate uint64, weightedDeveloperRewardsReceivers []WeightedAddress,
	distributionProportion DistributionProportions, epochs int64, maxIncomingAndOutgoingTxns int64,
	cosmosProposalParams CosmosChainProposalParams, stakingEpochIdentifier string, custodialAddress string,
	undelegateEpochIdentifier string, ChunkSize int64, bondDenom []string, mintDenom string,
	multiSigThreshold uint64, retryLimit uint64) Params {
	return Params{
		MinMintingAmount:                  minMintingAmount,
		MaxMintingAmount:                  maxMintingAmount,
		MinBurningAmount:                  minBurningAmount,
		MaxBurningAmount:                  maxBurningAmount,
		MaxValidatorToDelegate:            maxValidatorToDelegate,
		WeightedDeveloperRewardsReceivers: weightedDeveloperRewardsReceivers,
		DistributionProportion:            distributionProportion,
		Epochs:                            epochs,
		MaxIncomingAndOutgoingTxns:        maxIncomingAndOutgoingTxns,
		CosmosProposalParams:              cosmosProposalParams,
		ModuleEnabled:                     false,
		StakingEpochIdentifier:            stakingEpochIdentifier,
		CustodialAddress:                  custodialAddress,
		UndelegateEpochIdentifier:         undelegateEpochIdentifier,
		ChunkSize:                         ChunkSize,
		BondDenoms:                        bondDenom,
		MintDenom:                         mintDenom,
		MultisigThreshold:                 multiSigThreshold,
		RetryLimit:                        retryLimit,
	}
}

func DefaultParams() Params {
	return Params{
		MinMintingAmount:       sdk.NewInt64Coin("uatom", 5000000),
		MaxMintingAmount:       sdk.NewInt64Coin("uatom", 100000000000),
		MinBurningAmount:       sdk.NewInt64Coin("uatom", 5000000),
		MaxBurningAmount:       sdk.NewInt64Coin("uatom", 100000000000),
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
		Epochs:                     0,
		MaxIncomingAndOutgoingTxns: 10000,
		CosmosProposalParams: CosmosChainProposalParams{
			ChainID:              "cosmoshub-4", //TODO use these as conditions for proposals
			ReduceVotingPeriodBy: DefaultPeriod,
		},
		DelegationThreshold:       sdk.NewInt64Coin("uatom", 2000000000),
		ModuleEnabled:             false, //TODO : Make false before launch
		StakingEpochIdentifier:    "3hour",
		CustodialAddress:          "cosmos15vm0p2x990762txvsrpr26ya54p5qlz9xqlw5z",
		UndelegateEpochIdentifier: "3.5day",
		ChunkSize:                 5,
		BondDenoms:                DefaultBondDenom,
		MintDenom:                 DefaultMintDenom,
		MultisigThreshold:         3,
		RetryLimit:                10,
	}
}

func (p Params) Validate() error {
	if err := validateMinMintingAmount(p.MinMintingAmount); err != nil {
		return err
	}
	if err := validateMaxMintingAmount(p.MaxMintingAmount); err != nil {
		return err
	}
	if err := validateMinBurningAmount(p.MinBurningAmount); err != nil {
		return err
	}
	if err := validateMaxBurningAmount(p.MaxBurningAmount); err != nil {
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
	if err := validateEpochs(p.Epochs); err != nil {
		return err
	}
	if err := validateMaxIncomingAndOutgoingTxns(p.MaxIncomingAndOutgoingTxns); err != nil {
		return err
	}
	if err := validateCosmosProposalParams(p.CosmosProposalParams); err != nil {
		return err
	}
	if err := validateDelegationThreshold(p.DelegationThreshold); err != nil {
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
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateMultisigThreshold(p.MultisigThreshold); err != nil {
		return err
	}
	if err := validateRetryLimit(p.RetryLimit); err != nil {
		return err
	}
	return nil
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func (p *Params) ParamSetPairs() paramsTypes.ParamSetPairs {
	return paramsTypes.ParamSetPairs{
		paramsTypes.NewParamSetPair(KeyMinMintingAmount, &p.MinMintingAmount, validateMinMintingAmount),
		paramsTypes.NewParamSetPair(KeyMaxMintingAmount, &p.MaxMintingAmount, validateMaxMintingAmount),
		paramsTypes.NewParamSetPair(KeyMinBurningAmount, &p.MinBurningAmount, validateMinBurningAmount),
		paramsTypes.NewParamSetPair(KeyMaxBurningAmount, &p.MaxBurningAmount, validateMaxBurningAmount),
		paramsTypes.NewParamSetPair(KeyMaxValidatorToDelegate, &p.MaxValidatorToDelegate, validateMaxValidatorToDelegate),
		paramsTypes.NewParamSetPair(KeyWeightedDeveloperRewardsReceivers, &p.WeightedDeveloperRewardsReceivers, validateWeightedDeveloperRewardsReceivers),
		paramsTypes.NewParamSetPair(KeyDistributionProportion, &p.DistributionProportion, validateDistributionProportion),
		paramsTypes.NewParamSetPair(KeyEpochs, &p.Epochs, validateEpochs),
		paramsTypes.NewParamSetPair(KeyMaxIncomingAndOutgoingTxns, &p.MaxIncomingAndOutgoingTxns, validateMaxIncomingAndOutgoingTxns),
		paramsTypes.NewParamSetPair(KeyCosmosProposalParams, &p.CosmosProposalParams, validateCosmosProposalParams),
		paramsTypes.NewParamSetPair(KeyDelegationThreshold, &p.DelegationThreshold, validateDelegationThreshold),
		paramsTypes.NewParamSetPair(KeyModuleEnabled, &p.ModuleEnabled, validateModuleEnabled),
		paramsTypes.NewParamSetPair(KeyStakingEpochIdentifier, &p.StakingEpochIdentifier, epochsTypes.ValidateEpochIdentifierInterface),
		paramsTypes.NewParamSetPair(KeyCustodialAddress, &p.CustodialAddress, validateCustodialAddress),
		paramsTypes.NewParamSetPair(KeyUndelegateEpochIdentifier, &p.UndelegateEpochIdentifier, epochsTypes.ValidateEpochIdentifierInterface),
		paramsTypes.NewParamSetPair(KeyChunkSize, &p.ChunkSize, validateWithdrawRewardsChunkSize),
		paramsTypes.NewParamSetPair(KeyBondDenom, &p.BondDenoms, validateBondDenom),
		paramsTypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramsTypes.NewParamSetPair(KeyMultisigThreshold, &p.MultisigThreshold, validateMultisigThreshold),
		paramsTypes.NewParamSetPair(KeyRetryLimit, &p.RetryLimit, validateRetryLimit),
	}
}

func (p Params) GetBondDenomOf(s string) (string, error) {
	for _, element := range p.BondDenoms {
		if element == s {
			return element, nil
		}
	}
	return "", ErrInvalidBondDenom
}

func validateMinMintingAmount(i interface{}) error {
	coin, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if coin.IsNegative() {
		return errors.New("min minting amount cannot be negative")
	}
	return nil
}

func validateMaxMintingAmount(i interface{}) error {
	coin, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if coin.IsNegative() {
		return errors.New("max minting amount cannot be negative")
	}
	return nil
}

func validateMinBurningAmount(i interface{}) error {
	coin, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if coin.IsNegative() {
		return errors.New("min burning amount cannot be negative")
	}
	return nil
}

func validateMaxBurningAmount(i interface{}) error {
	coin, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if coin.IsNegative() {
		return errors.New("max burning amount cannot be negative")
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

func validateEpochs(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("epoch must be non-negative")
	}

	return nil
}

func validateMaxIncomingAndOutgoingTxns(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("total incoming or outgoing transaction must be non-negative")
	}

	return nil
}

func validateCosmosProposalParams(i interface{}) error {
	v, ok := i.(CosmosChainProposalParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.ChainID != "cosmoshub-4" {
		return fmt.Errorf("invalid chain-id for cosmos %T", i)
	}

	if v.ReduceVotingPeriodBy <= 0 {
		return fmt.Errorf("incorrect voting Period %T", i)
	}

	return nil
}

func validateDelegationThreshold(i interface{}) error {
	v, ok := i.(sdk.Coin)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("delegation threshold cannot be negative")
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

	if len(v) <= 0 {
		return fmt.Errorf("bond denom cannot be empty")
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

func validateMultisigThreshold(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("multisig threshold must be non negative")
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
