package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ cosmosTypes.QueryServer = Keeper{}

// QueryParams queries all the params in genesis
func (k Keeper) QueryParams(context context.Context, _ *cosmosTypes.QueryParamsRequest) (*cosmosTypes.QueryParamsResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(context)
	params := k.GetParams(ctx)
	return &cosmosTypes.QueryParamsResponse{Params: params}, nil
}

// QueryTxByID Query txns by ID for orchestrators to sign
func (k Keeper) QueryTxByID(context context.Context, req *cosmosTypes.QueryOutgoingTxByIDRequest) (*cosmosTypes.QueryOutgoingTxByIDResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(context)
	cosmosTxDetails, err := k.GetTxnFromOutgoingPoolByID(ctx, req.TxID)
	if err != nil {
		return nil, err
	}
	return &cosmosTxDetails, nil
}

// Proposal Query proposal by ID which came in from cosmos side
func (k Keeper) Proposal(context context.Context, req *cosmosTypes.QueryProposalRequest) (*cosmosTypes.QueryProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	ctx := sdkTypes.UnwrapSDKContext(context)

	proposal, found := k.GetProposal(ctx, req.ProposalId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "proposal %d doesn't exist", req.ProposalId)
	}

	return &cosmosTypes.QueryProposalResponse{Proposal: proposal}, nil
}

// Proposals implements the Query/Proposals gRPC method
func (k Keeper) Proposals(c context.Context, req *cosmosTypes.QueryProposalsRequest) (*cosmosTypes.QueryProposalsResponse, error) {
	var filteredProposals cosmosTypes.Proposals
	ctx := sdkTypes.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	proposalStore := prefix.NewStore(store, cosmosTypes.ProposalsKeyPrefix)

	pageRes, err := query.FilteredPaginate(proposalStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var p cosmosTypes.Proposal
		if err := k.cdc.Unmarshal(value, &p); err != nil {
			return false, status.Error(codes.Internal, err.Error())
		}

		matchStatus := true

		// match status (if supplied/valid)
		if cosmosTypes.ValidProposalStatus(req.ProposalStatus) {
			matchStatus = p.Status == req.ProposalStatus
		}

		if matchStatus {
			if accumulate {
				filteredProposals = append(filteredProposals, p)
			}

			return true, nil
		}

		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cosmosTypes.QueryProposalsResponse{Proposals: filteredProposals, Pagination: pageRes}, nil
}

// Vote returns Voted information based on proposalID, voterAddr
func (k Keeper) Vote(c context.Context, req *cosmosTypes.QueryVoteRequest) (*cosmosTypes.QueryVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	if req.Voter == "" {
		return nil, status.Error(codes.InvalidArgument, "empty voter address")
	}

	ctx := sdkTypes.UnwrapSDKContext(c)

	voter, err := sdkTypes.AccAddressFromBech32(req.Voter)
	if err != nil {
		return nil, err
	}
	vote, found := k.GetVote(ctx, req.ProposalId, voter)
	if !found {
		return nil, status.Errorf(codes.InvalidArgument,
			"voter: %v not found for proposal: %v", req.Voter, req.ProposalId)
	}

	return &cosmosTypes.QueryVoteResponse{Vote: vote}, nil
}

// Votes returns single proposal's votes
func (k Keeper) Votes(c context.Context, req *cosmosTypes.QueryVotesRequest) (*cosmosTypes.QueryVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	var votes cosmosTypes.Votes
	ctx := sdkTypes.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	votesStore := prefix.NewStore(store, cosmosTypes.VotesKey(req.ProposalId))

	pageRes, err := query.Paginate(votesStore, req.Pagination, func(key []byte, value []byte) error {
		var vote cosmosTypes.Vote
		if err := k.cdc.Unmarshal(value, &vote); err != nil {
			return err
		}
		populateLegacyOption(&vote)

		votes = append(votes, vote)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cosmosTypes.QueryVotesResponse{Votes: votes, Pagination: pageRes}, nil
}

// CosmosValidatorSet returns the cosmos validator set and their respective weights
func (k Keeper) CosmosValidatorSet(c context.Context, _ *cosmosTypes.QueryCosmosValidatorSetRequest) (*cosmosTypes.QueryCosmosValidatorSetResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)
	weihtedAddresses := k.GetAllCosmosValidatorSet(ctx)
	return &cosmosTypes.QueryCosmosValidatorSetResponse{WeightedAddresses: weihtedAddresses}, nil
}

// OracleValidatorSet returns the oracle validator set and their respective weights
func (k Keeper) OracleValidatorSet(c context.Context, _ *cosmosTypes.QueryOracleValidatorSetRequest) (*cosmosTypes.QueryOracleValidatorSetResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)
	weightedAddresses := k.getAllOracleValidatorSet(ctx)
	return &cosmosTypes.QueryOracleValidatorSetResponse{WeightedAddresses: weightedAddresses}, nil
}

// ValidatorMapping returns the oracle address mapped to the given validator address
func (k Keeper) ValidatorMapping(c context.Context, query *cosmosTypes.QueryValidatorMappingRequest) (*cosmosTypes.QueryValidatorMappingResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	validatorAddress, err := sdkTypes.ValAddressFromBech32(query.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	validatorMapping := k.getValidatorMapping(ctx, validatorAddress)

	return &cosmosTypes.QueryValidatorMappingResponse{
		ValidatorAddress:    validatorAddress.String(),
		OrchestratorAddress: validatorMapping.OrchestratorAddresses,
	}, nil
}

// OracleHeight returns the last updated height of given oracle
func (k Keeper) OracleHeight(c context.Context, query *cosmosTypes.QueryOracleLastUpdateHeightRequest) (*cosmosTypes.QueryOracleLastUpdateHeightResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	oracleAddress, err := sdkTypes.AccAddressFromBech32(query.OracleAddress)
	if err != nil {
		return nil, err
	}

	oracleLastUpadteHeightNative := k.getOracleLastUpdateHeightNative(ctx, oracleAddress)
	oracleLastUpadteHeightCosmos := k.getOracleLastUpdateHeightCosmos(ctx, oracleAddress)

	return &cosmosTypes.QueryOracleLastUpdateHeightResponse{
		BlockHeightCosmos: oracleLastUpadteHeightCosmos,
		BlockHeightNative: oracleLastUpadteHeightNative,
	}, nil
}

// CosmosBalance returns the cosmos account balances
func (k Keeper) CosmosBalance(c context.Context, query *cosmosTypes.QueryCosmosBalanceRequest) (*cosmosTypes.QueryCosmosBalanceResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)
	balance := k.getCosmosBalances(ctx)
	return &cosmosTypes.QueryCosmosBalanceResponse{Balance: balance}, nil
}
