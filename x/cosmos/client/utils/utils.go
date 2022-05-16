package utils

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"io/ioutil"
	"strings"
)

type (
	// EnableModuleProposalReq defines a community pool spend proposal request body.
	EnableModuleProposalReq struct {
		BaseReq      rest.BaseReq             `json:"base_req" yaml:"base_req"`
		EnableModule EnableModuleProposalJSON `json:"enable_module" yaml:"enable_module"`
	}

	// ChangeMultisigPropsoalReq defines a community pool spend proposal request body.
	ChangeMultisigPropsoalReq struct {
		BaseReq        rest.BaseReq               `json:"base_req" yaml:"base_req"`
		ChangeMultisig ChangeMultisigPropsoalJSON `json:"change_multisig" yaml:"change_multisig"`
	}

	EnableModuleProposalJSON struct {
		Title         string `json:"title" yaml:"title"`
		Description   string `json:"description" yaml:"description"`
		Threshold     uint64 `json:"threshold" yaml:"threshold"`
		AccountNumber uint64 `json:"account_number" yaml:"account_number"`
		Depositor     string `json:"depositor" yaml:"depositor"`
		Deposit       string `json:"deposit" yaml:"deposit"`
	}

	ChangeMultisigPropsoalJSON struct {
		Title                 string   `json:"title" yaml:"title"`
		Description           string   `json:"description" yaml:"description"`
		Threshold             uint64   `json:"threshold" yaml:"threshold"`
		OrchestratorAddresses []string `json:"orchestrator_addresses" yaml:"orchestrator_addresses"`
		AccountNumber         uint64   `json:"account_number" yaml:"account_number"`
		Depositor             string   `json:"depositor" yaml:"depositor"`
		Deposit               string   `json:"deposit" yaml:"deposit"`
	}
)

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
