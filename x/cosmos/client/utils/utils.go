package utils

import (
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"strings"
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
