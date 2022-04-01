package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ghodss/yaml"
)

const (
	DefaultPeriod time.Duration = time.Minute * 1 // 6 hours //TODO : Change back to 6 hours
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
)

func ParamKeyTable() paramsTypes.KeyTable {
	return paramsTypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(minMintingAmount uint64, maxMintingAmount uint64, minBurningAmount uint64, maxBurningAmount uint64,
	maxValidatorToDelegate uint64, validatorSetCosmosChain []WeightedAddress, validatorSetNativeChain []WeightedAddress,
	weightedDeveloperRewardsReceivers []WeightedAddress, distributionProportion DistributionProportions, epochs int64,
	maxIncomingAndOutgoingTxns int64, cosmosProposalParams CosmosChainProposalParams) Params {
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
	}
}

func DefaultParams() Params {
	return Params{
		MinMintingAmount:       5000000,
		MaxMintingAmount:       100000000000,
		MinBurningAmount:       5000000,
		MaxBurningAmount:       100000000000,
		MaxValidatorToDelegate: 3,
		ValidatorSetCosmosChain: []WeightedAddress{
			{Address: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Weight: sdk.NewDecWithPrec(5, 1)},
			{Address: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2", Weight: sdk.NewDecWithPrec(2, 1)},
			{Address: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Weight: sdk.NewDecWithPrec(1, 1)},
			{Address: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Weight: sdk.NewDecWithPrec(2, 1)},
		},
		ValidatorSetNativeChain:           []WeightedAddress{},
		WeightedDeveloperRewardsReceivers: []WeightedAddress{},
		DistributionProportion: DistributionProportions{
			ValidatorRewards: sdk.NewDecWithPrec(7, 2),
			DeveloperRewards: sdk.NewDecWithPrec(3, 2),
		},
		Epochs:                     5000,
		MaxIncomingAndOutgoingTxns: 10000,
		CosmosProposalParams: CosmosChainProposalParams{
			ChainID:              "cosmoshub-4", //TODO use these as conditions for proposals
			ReduceVotingPeriodBy: DefaultPeriod,
		},
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
	//if err := validateDistributionProportion(p.DistributionProportion); err != nil {
	//	return err
	//}
	if err := validateEpochs(p.Epochs); err != nil {
		return err
	}
	if err := validateMaxIncomingAndOutgoingTxns(p.MaxIncomingAndOutgoingTxns); err != nil {
		return err
	}
	if err := validateCosmosProposalParams(p.CosmosProposalParams); err != nil {
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
		//paramsTypes.NewParamSetPair(KeyDistributionProportion, &p.DistributionProportion, validateDistributionProportion),
		paramsTypes.NewParamSetPair(KeyEpochs, &p.Epochs, validateEpochs),
		paramsTypes.NewParamSetPair(KeyMaxIncomingAndOutgoingTxns, &p.MaxIncomingAndOutgoingTxns, validateMaxIncomingAndOutgoingTxns),
		paramsTypes.NewParamSetPair(KeyCosmosProposalParams, &p.CosmosProposalParams, validateCosmosProposalParams),
	}
}

func validateMinMintingAmount(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMaxMintingAmount(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMinBurningAmount(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMaxBurningAmount(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
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
