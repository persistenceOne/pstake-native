package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	epochsTypes "github.com/persistenceOne/pstake-native/x/epochs/types"
)

// Implements GovHooks interface
var _ cosmosTypes.GovHooks = Keeper{}

// AfterProposalSubmission - call hook if registered
func (k Keeper) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {
	if k.hooks != nil {
		k.hooks.AfterProposalSubmission(ctx, proposalID)
	}
}

// AfterProposalVote - call hook if registered
func (k Keeper) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.AfterProposalVote(ctx, proposalID, voterAddr)
	}
}

// AfterProposalVotingPeriodEnded - call hook if registered
func (k Keeper) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {
	if k.hooks != nil {
		k.hooks.AfterProposalVotingPeriodEnded(ctx, proposalID)
	}
}

// BeforeEpochStart - call hook if registered
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

/*
AfterEpochEnd handle the "stake", "reward" and "undelegate" epoch and their respective actions
1. "stake" generates delegate transaction for delegating the amount of stake accumulated over the "stake" epoch
2. "reward" generates delegate transaction for delegating the amount of stake accumulated over the "reward" epochs
and shift the amount to next epoch if the min amount is not reached
3. "undelegate" generated the undelegate transaction for undelegating the amount accumulated over the "undelegate" epoch
*/
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

// ___________________________________________________________________________________________________

// RewardsHooks wrapper struct for incentives keeper
type RewardsHooks struct {
	k Keeper
}

var _ epochsTypes.EpochHooks = RewardsHooks{}

// Return the wrapper struct
func (k Keeper) Hooks() RewardsHooks {
	return RewardsHooks{k}
}

// epochs hooks
func (h RewardsHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h RewardsHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
