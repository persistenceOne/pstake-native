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
	DefaultBondDenom string        = "uatom"
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
	KeyRewardsEpochIdentifier            = []byte("EpochIdentifier")
	KeyCustodialAddress                  = []byte("CustodialAddress")
	KeyUnbondingEpochIdentifier          = []byte("UnbondingKeyIdentifier")
	KeyChunkSize                         = []byte("ChunkSize")
	KeyBondDenom                         = []byte("BondDenom")
	KeyMintDenom                         = []byte("MintDenom")
)

func ParamKeyTable() paramsTypes.KeyTable {
	return paramsTypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(minMintingAmount sdk.Coin, maxMintingAmount sdk.Coin, minBurningAmount sdk.Coin, maxBurningAmount sdk.Coin,
	maxValidatorToDelegate uint64, validatorSetCosmosChain []WeightedAddressCosmos, validatorSetNativeChain []WeightedAddress,
	weightedDeveloperRewardsReceivers []WeightedAddress, distributionProportion DistributionProportions, epochs int64,
	maxIncomingAndOutgoingTxns int64, cosmosProposalParams CosmosChainProposalParams, epochIdentifier string,
	custodialAddress string, unbondingEpochIdentifier string, ChunkSize int64, bondDenom string, mintDenom string) Params {
	return Params{
		MinMintingAmount:                  minMintingAmount,
		MaxMintingAmount:                  maxMintingAmount,
		MinBurningAmount:                  minBurningAmount,
		MaxBurningAmount:                  maxBurningAmount,
		MaxValidatorToDelegate:            maxValidatorToDelegate,
		ValidatorSetCosmosChain:           validatorSetCosmosChain,
		ValidatorSetNativeChain:           validatorSetNativeChain,
		WeightedDeveloperRewardsReceivers: weightedDeveloperRewardsReceivers,
		DistributionProportion:            distributionProportion,
		Epochs:                            epochs,
		MaxIncomingAndOutgoingTxns:        maxIncomingAndOutgoingTxns,
		CosmosProposalParams:              cosmosProposalParams,
		ModuleEnabled:                     false,
		RewardsEpochIdentifier:            epochIdentifier,
		CustodialAddress:                  custodialAddress,
		UnbondingEpochIdentifier:          unbondingEpochIdentifier,
		ChunkSize:                         ChunkSize,
		BondDenom:                         bondDenom,
		MintDenom:                         mintDenom,
	}
}

func DefaultParams() Params {
	return Params{
		MinMintingAmount:       sdk.NewInt64Coin("uatom", 5000000),
		MaxMintingAmount:       sdk.NewInt64Coin("uatom", 100000000000),
		MinBurningAmount:       sdk.NewInt64Coin("uatom", 5000000),
		MaxBurningAmount:       sdk.NewInt64Coin("uatom", 100000000000),
		MaxValidatorToDelegate: 3,
		ValidatorSetCosmosChain: []WeightedAddressCosmos{
			{Address: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
				Weight:                 sdk.NewDecWithPrec(5, 1),
				CurrentDelegatedAmount: sdk.NewInt64Coin("uatom", 0),
				IdealDelegatedAmount:   sdk.NewInt64Coin("uatom", 0),
				Difference:             sdk.NewInt64Coin("uatom", 0),
			},
			{Address: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
				Weight:                 sdk.NewDecWithPrec(2, 1),
				CurrentDelegatedAmount: sdk.NewInt64Coin("uatom", 0),
				IdealDelegatedAmount:   sdk.NewInt64Coin("uatom", 0),
				Difference:             sdk.NewInt64Coin("uatom", 0),
			},
			{Address: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5",
				Weight:                 sdk.NewDecWithPrec(1, 1),
				CurrentDelegatedAmount: sdk.NewInt64Coin("uatom", 0),
				IdealDelegatedAmount:   sdk.NewInt64Coin("uatom", 0),
				Difference:             sdk.NewInt64Coin("uatom", 0),
			},
			{Address: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2",
				Weight:                 sdk.NewDecWithPrec(2, 1),
				CurrentDelegatedAmount: sdk.NewInt64Coin("uatom", 0),
				IdealDelegatedAmount:   sdk.NewInt64Coin("uatom", 0),
				Difference:             sdk.NewInt64Coin("uatom", 0),
			},
		},
		ValidatorSetNativeChain:           []WeightedAddress{},
		WeightedDeveloperRewardsReceivers: []WeightedAddress{},
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
		DelegationThreshold:      sdk.NewInt64Coin(DefaultBondDenom, 2000000000),
		ModuleEnabled:            true, //TODO : Make false before launch
		RewardsEpochIdentifier:   "3hour",
		CustodialAddress:         "cosmos15vm0p2x990762txvsrpr26ya54p5qlz9xqlw5z",
		UnbondingEpochIdentifier: "3.5day",
		ChunkSize:                5,
		BondDenom:                DefaultBondDenom,
		MintDenom:                DefaultMintDenom,
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
	if err := validateValidatorSetCosmosChain(p.ValidatorSetCosmosChain); err != nil {
		return err
	}
	if err := validateValidatorSetNativeChain(p.ValidatorSetNativeChain); err != nil {
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
	if err := epochsTypes.ValidateEpochIdentifierInterface(p.RewardsEpochIdentifier); err != nil {
		return err
	}
	if err := validateCustodialAddress(p.CustodialAddress); err != nil {
		return err
	}
	if err := epochsTypes.ValidateEpochIdentifierInterface(p.UnbondingEpochIdentifier); err != nil {
		return err
	}
	if err := validateWithdrawRewardsChunkSize(p.ChunkSize); err != nil {
		return err
	}
	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}
	if err := validateMintDenom(p.MintDenom); err != nil {
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
		paramsTypes.NewParamSetPair(KeyValidatorSetCosmosChain, &p.ValidatorSetCosmosChain, validateValidatorSetCosmosChain),
		paramsTypes.NewParamSetPair(KeyValidatorSetNativeChain, &p.ValidatorSetNativeChain, validateValidatorSetNativeChain),
		paramsTypes.NewParamSetPair(KeyWeightedDeveloperRewardsReceivers, &p.WeightedDeveloperRewardsReceivers, validateWeightedDeveloperRewardsReceivers),
		paramsTypes.NewParamSetPair(KeyDistributionProportion, &p.DistributionProportion, validateDistributionProportion),
		paramsTypes.NewParamSetPair(KeyEpochs, &p.Epochs, validateEpochs),
		paramsTypes.NewParamSetPair(KeyMaxIncomingAndOutgoingTxns, &p.MaxIncomingAndOutgoingTxns, validateMaxIncomingAndOutgoingTxns),
		paramsTypes.NewParamSetPair(KeyCosmosProposalParams, &p.CosmosProposalParams, validateCosmosProposalParams),
		paramsTypes.NewParamSetPair(KeyDelegationThreshold, &p.DelegationThreshold, validateDelegationThreshold),
		paramsTypes.NewParamSetPair(KeyModuleEnabled, &p.ModuleEnabled, validateModuleEnabled),
		paramsTypes.NewParamSetPair(KeyRewardsEpochIdentifier, &p.RewardsEpochIdentifier, epochsTypes.ValidateEpochIdentifierInterface),
		paramsTypes.NewParamSetPair(KeyCustodialAddress, &p.CustodialAddress, validateCustodialAddress),
		paramsTypes.NewParamSetPair(KeyUnbondingEpochIdentifier, &p.UnbondingEpochIdentifier, epochsTypes.ValidateEpochIdentifierInterface),
		paramsTypes.NewParamSetPair(KeyChunkSize, &p.ChunkSize, validateWithdrawRewardsChunkSize),
		paramsTypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		paramsTypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
	}
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

func validateValidatorSetCosmosChain(i interface{}) error {
	v, ok := i.([]WeightedAddressCosmos)
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
		if w.CurrentDelegatedAmount.IsNegative() {
			return fmt.Errorf("non-positive current delegation amount at %dth", i)
		}
		if w.IdealDelegatedAmount.IsNegative() {
			return fmt.Errorf("non-positive ideal delegated amount at %dth", i)
		}
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}

func validateValidatorSetNativeChain(i interface{}) error {
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
		if w.Weight.GT(sdk.NewDecWithPrec(1, 1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(sdk.NewDecWithPrec(1, 1)) {
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
		_, err := sdk.AccAddressFromBech32(v)
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
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
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
