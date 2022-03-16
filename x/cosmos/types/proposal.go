package types

import (
	"github.com/ghodss/yaml"
	"time"
)

func NewProposal(id uint64, title string, description string, submitTime time.Time, votingPeriod time.Duration) (Proposal, error) {
	p := Proposal{
		ProposalId:       id,
		Title:            title,
		Description:      description,
		Status:           StatusVotingPeriod,
		FinalTallyResult: EmptyTallyResult(),
		SubmitTime:       submitTime,
		VotingEndTime:    submitTime.Add(votingPeriod),
		VotingStartTime:  submitTime,
	}

	return p, nil
}

// String implements stringer interface
func (p Proposal) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
