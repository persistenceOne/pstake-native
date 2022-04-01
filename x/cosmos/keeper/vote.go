package keeper

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// AddVote adds a vote on a specific proposal
func (k Keeper) AddVote(ctx sdkTypes.Context, proposalID uint64, voterAddr sdkTypes.AccAddress, options cosmosTypes.WeightedVoteOptions) error {
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		return sdkErrors.Wrapf(cosmosTypes.ErrUnknownProposal, "%d", proposalID)
	}
	if proposal.Status != cosmosTypes.StatusVotingPeriod {
		return sdkErrors.Wrapf(cosmosTypes.ErrInactiveProposal, "%d", proposalID)
	}

	for _, option := range options {
		if !cosmosTypes.ValidWeightedVoteOption(option) {
			return sdkErrors.Wrap(cosmosTypes.ErrInvalidVote, option.String())
		}
	}

	vote := cosmosTypes.NewVote(proposalID, voterAddr, options)
	k.SetVote(ctx, vote)

	// called after a vote on a proposal is cast
	k.AfterProposalVote(ctx, proposalID, voterAddr)

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			cosmosTypes.EventTypeProposalVote,
			sdkTypes.NewAttribute(cosmosTypes.AttributeKeyOption, options.String()),
			sdkTypes.NewAttribute(cosmosTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

// SetVote sets a Vote to the gov store
func (keeper Keeper) SetVote(ctx sdkTypes.Context, vote cosmosTypes.Vote) {
	// vote.Option is a deprecated field, we don't set it in state
	if vote.Option != cosmosTypes.OptionEmpty { //nolint
		vote.Option = cosmosTypes.OptionEmpty //nolint
	}

	store := ctx.KVStore(keeper.storeKey)
	bz, err := vote.Marshal()
	if err != nil {
		panic(err)
	}
	//bz := keeper.cdc.MustMarshal(&vote)
	addr, err := sdkTypes.AccAddressFromBech32(vote.Voter)
	if err != nil {
		panic(err)
	}
	store.Set(cosmosTypes.VoteKey(vote.ProposalId, addr), bz)
}

func (keeper Keeper) GetVotes(ctx sdkTypes.Context, proposalID uint64) (votes cosmosTypes.Votes) {
	keeper.IterateVotes(ctx, proposalID, func(vote cosmosTypes.Vote) bool {
		populateLegacyOption(&vote)
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVote gets the vote from an address on a specific proposal
func (keeper Keeper) GetVote(ctx sdkTypes.Context, proposalID uint64, voterAddr sdkTypes.AccAddress) (vote cosmosTypes.Vote, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(cosmosTypes.VoteKey(proposalID, voterAddr))
	if bz == nil {
		return vote, false
	}

	err := vote.Unmarshal(bz)
	if err != nil {
		return vote, false
	}
	populateLegacyOption(&vote)

	return vote, true
}

// populateLegacyOption adds graceful fallback of deprecated `Option` field, in case
// there's only 1 VoteOption.
func populateLegacyOption(vote *cosmosTypes.Vote) {
	if len(vote.Options) == 1 && vote.Options[0].Weight.Equal(sdkTypes.MustNewDecFromStr("1.0")) {
		vote.Option = vote.Options[0].Option //nolint
	}
}

// deleteVote deletes a vote from a given proposalID and voter from the store
func (keeper Keeper) deleteVote(ctx sdkTypes.Context, proposalID uint64, voterAddr sdkTypes.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(cosmosTypes.VoteKey(proposalID, voterAddr))
}

func (keeper Keeper) IterateVotes(ctx sdkTypes.Context, proposalID uint64, cb func(vote cosmosTypes.Vote) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdkTypes.KVStorePrefixIterator(store, cosmosTypes.VotesKey(proposalID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote cosmosTypes.Vote
		err := vote.Unmarshal(iterator.Value())
		if err != nil {
			return
		}
		populateLegacyOption(&vote)

		if cb(vote) {
			break
		}
	}
}
