package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// AddVote adds a vote on a specific proposal
func (k Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options cosmosTypes.WeightedVoteOptions) error {
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		return sdkerrors.Wrapf(cosmosTypes.ErrUnknownProposal, "%d", proposalID)
	}
	if proposal.Status != cosmosTypes.StatusVotingPeriod {
		return sdkerrors.Wrapf(cosmosTypes.ErrInactiveProposal, "%d", proposalID)
	}

	for _, option := range options {
		if !cosmosTypes.ValidWeightedVoteOption(option) {
			return sdkerrors.Wrap(cosmosTypes.ErrInvalidVote, option.String())
		}
	}

	vote := cosmosTypes.NewVote(proposalID, voterAddr, options)
	k.SetVote(ctx, vote)

	// called after a vote on a proposal is cast
	k.AfterProposalVote(ctx, proposalID, voterAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			cosmosTypes.EventTypeProposalVote,
			sdk.NewAttribute(cosmosTypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(cosmosTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

// SetVote sets a Vote to the gov store
func (keeper Keeper) SetVote(ctx sdk.Context, vote cosmosTypes.Vote) {
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
	addr, err := sdk.AccAddressFromBech32(vote.Voter)
	if err != nil {
		panic(err)
	}
	store.Set(cosmosTypes.VoteKey(vote.ProposalId, addr), bz)
}
