package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/ghodss/yaml"
)

// NewProposal creates a new Proposal instance
func NewProposal(id uint64, title string, description string, submitTime time.Time, votingPeriod time.Duration, cosmosProposalID uint64) (Proposal, error) {
	p := Proposal{
		ProposalId:       id,
		Title:            title,
		Description:      description,
		Status:           StatusVotingPeriod,
		FinalTallyResult: EmptyTallyResult(),
		SubmitTime:       submitTime,
		VotingEndTime:    submitTime.Add(votingPeriod),
		VotingStartTime:  submitTime,
		CosmosProposalId: cosmosProposalID,
	}

	return p, nil
}

// String implements stringer interface
func (p Proposal) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Proposals is an array of proposal
type Proposals []Proposal

// Equal returns true if two slices (order-dependant) of proposals are equal.
func (p Proposals) Equal(other Proposals) bool {
	if len(p) != len(other) {
		return false
	}

	for i, proposal := range p {
		if !proposal.Equal(other[i]) {
			return false
		}
	}

	return true
}

// String implements stringer interface
func (p Proposals) String() string {
	out := "ID - (Status) [Type] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] \n",
			prop.ProposalId, prop.Status,
			prop.Title)
	}
	return strings.TrimSpace(out)
}

// ProposalStatusFromString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	num, ok := ProposalStatus_value[str]
	if !ok {
		return StatusNil, fmt.Errorf("'%s' is not a valid proposal status", str)
	}
	return ProposalStatus(num), nil
}

// ValidProposalStatus returns true if the proposal status is valid and false
// otherwise.
func ValidProposalStatus(status ProposalStatus) bool {
	if status == StatusVotingPeriod ||
		status == StatusPassed {
		return true
	}
	return false
}
