package utils

import (
	"io/ioutil"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// EnableModuleProposalReq defines a community pool spend proposal request body.
type EnableModuleProposalReq struct {
	BaseReq      rest.BaseReq             `json:"base_req" yaml:"base_req"`
	EnableModule EnableModuleProposalJSON `json:"enable_module" yaml:"enable_module"`
}

// ChangeMultisigPropsoalReq defines a community pool spend proposal request body.
type ChangeMultisigPropsoalReq struct {
	BaseReq        rest.BaseReq               `json:"base_req" yaml:"base_req"`
	ChangeMultisig ChangeMultisigPropsoalJSON `json:"change_multisig" yaml:"change_multisig"`
}

type ChangeCosmosValidatorWeightsProposalReq struct {
	BaseReq            rest.BaseReq                             `json:"base_req" yaml:"base_req"`
	CosmosValidatorSet ChangeCosmosValidatorWeightsProposalJSON `json:"cosmos_validator_set" yaml:"cosmos_validator_set"`
}

type ChangeOracleValidatorWeightsProposalReq struct {
	BaseReq            rest.BaseReq                             `json:"base_req" yaml:"base_req"`
	OracleValidatorSet ChangeOracleValidatorWeightsProposalJSON `json:"oracle_validator_set" yaml:"oracle_validator_set"`
}

type EnableModuleProposalJSON struct {
	Title                 string   `json:"title" yaml:"title"`
	Description           string   `json:"description" yaml:"description"`
	Threshold             uint64   `json:"threshold" yaml:"threshold"`
	AccountNumber         uint64   `json:"account_number" yaml:"account_number"`
	SequenceNumber        uint64   `json:"sequence_number" yaml:"sequence_number"`
	OrchestratorAddresses []string `json:"orchestrator_addresses" yaml:"orchestrator_addresses"`
	Depositor             string   `json:"depositor" yaml:"depositor"`
	Deposit               string   `json:"deposit" yaml:"deposit"`
}

type ChangeMultisigPropsoalJSON struct {
	Title                 string   `json:"title" yaml:"title"`
	Description           string   `json:"description" yaml:"description"`
	Threshold             uint64   `json:"threshold" yaml:"threshold"`
	OrchestratorAddresses []string `json:"orchestrator_addresses" yaml:"orchestrator_addresses"`
	AccountNumber         uint64   `json:"account_number" yaml:"account_number"`
	Depositor             string   `json:"depositor" yaml:"depositor"`
	Deposit               string   `json:"deposit" yaml:"deposit"`
}

type ChangeCosmosValidatorWeightsProposalJSON struct {
	Title             string              `json:"title" yaml:"title"`
	Description       string              `json:"description" yaml:"description"`
	WeightedAddresses []WeightedAddresses `json:"weighted_addresses" yaml:"weighted_addresses"`
	Depositor         string              `json:"depositor" yaml:"depositor"`
	Deposit           string              `json:"deposit" yaml:"deposit"`
}

type ChangeOracleValidatorWeightsProposalJSON struct {
	Title             string              `json:"title" yaml:"title"`
	Description       string              `json:"description" yaml:"description"`
	WeightedAddresses []WeightedAddresses `json:"weighted_addresses" yaml:"weighted_addresses"`
	Depositor         string              `json:"depositor" yaml:"depositor"`
	Deposit           string              `json:"deposit" yaml:"deposit"`
}

type WeightedAddresses struct {
	ValAddress string `json:"val_address" yaml:"val_address"`
	Weight     string `json:"weight" yaml:"weight"`
}

// NormalizeVoteOption - normalize user specified vote option
func NormalizeVoteOption(option string) string {
	switch option {
	case "Yes", "yes":
		return cosmosTypes.OptionYes.String()

	case "Abstain", "abstain":
		return cosmosTypes.OptionAbstain.String()

	case "No", "no":
		return cosmosTypes.OptionNo.String()

	case "NoWithVeto", "no_with_veto":
		return cosmosTypes.OptionNoWithVeto.String()

	default:
		return option
	}
}

// NormalizeWeightedVoteOptions - normalize vote options param string
func NormalizeWeightedVoteOptions(options string) string {
	newOptions := []string{}
	for _, option := range strings.Split(options, ",") {
		fields := strings.Split(option, "=")
		fields[0] = NormalizeVoteOption(fields[0])
		if len(fields) < 2 {
			fields = append(fields, "1")
		}
		newOptions = append(newOptions, strings.Join(fields, "="))
	}
	return strings.Join(newOptions, ",")
}

// NormalizeProposalStatus - normalize user specified proposal status.
func NormalizeProposalStatus(status string) string {
	switch status {
	case "VotingPeriod", "voting_period":
		return cosmosTypes.StatusVotingPeriod.String()
	case "Passed", "passed":
		return cosmosTypes.StatusPassed.String()
	case "Rejected", "rejected":
		return cosmosTypes.StatusRejected.String()
	default:
		return status
	}
}

// ParseEnableModuleProposalJSON reads and parses a ParseEnableModuleProposalJSON from
// file.
func ParseEnableModuleProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (EnableModuleProposalJSON, error) {
	proposal := EnableModuleProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

func ParseChangeMultisigProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (ChangeMultisigPropsoalJSON, error) {
	proposal := ChangeMultisigPropsoalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

func ParseChangeCosmosValidatorWeightsProposalJSON(
	cdc *codec.LegacyAmino, proposalFile string) (ChangeCosmosValidatorWeightsProposalJSON, error) {
	proposal := ChangeCosmosValidatorWeightsProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

func ParseChangeOracleValidatorWeightsProposalJSON(
	cdc *codec.LegacyAmino, proposalFile string) (ChangeOracleValidatorWeightsProposalJSON, error) {
	proposal := ChangeOracleValidatorWeightsProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
