package types_test

import (
	"fmt"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProposalStatus_Format(t *testing.T) {
	statusVotingPeriod, _ := types.ProposalStatusFromString("PROPOSAL_STATUS_VOTING_PERIOD")
	statusProposalPassed, _ := types.ProposalStatusFromString("PROPOSAL_STATUS_PASSED")
	tests := []struct {
		pt                   types.ProposalStatus
		sprintFArgs          string
		expectedStringOutput string
	}{
		{statusVotingPeriod, "%s", "PROPOSAL_STATUS_VOTING_PERIOD"},
		{statusProposalPassed, "%v", "PROPOSAL_STATUS_PASSED"},
	}
	for _, tt := range tests {
		got := fmt.Sprintf(tt.sprintFArgs, tt.pt)
		require.Equal(t, tt.expectedStringOutput, got)
	}
}
