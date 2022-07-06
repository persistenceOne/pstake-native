package keeper

import (
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	sdkCodec "github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier returns query handler for the module
func NewQuerier(k Keeper, legacyQuerierCdc *sdkCodec.LegacyAmino) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abciTypes.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case cosmosTypes.QueryParameters:
			return queryParams(ctx, k, legacyQuerierCdc)
		case cosmosTypes.QueryTxByID:
			return queryTxByID(ctx, req, k, legacyQuerierCdc)
		case cosmosTypes.QueryProposal:
			return queryProposal(ctx, req, k, legacyQuerierCdc)
		case cosmosTypes.QueryVote:
			return queryVote(ctx, req, k, legacyQuerierCdc)
		case cosmosTypes.QueryVotes:
			return queryVotes(ctx, req, k, legacyQuerierCdc)
		case cosmosTypes.QueryProposals:
			return queryProposals(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkErrors.Wrapf(sdkErrors.ErrUnknownRequest, "unknown %s query endpoint", cosmosTypes.ModuleName)
		}
	}
}

func queryParams(ctx sdkTypes.Context, k Keeper, legacyQuerierCdc *sdkCodec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := sdkCodec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTxByID(ctx sdkTypes.Context, req abciTypes.RequestQuery, k Keeper, legacyQuerierCdc *sdkCodec.LegacyAmino) ([]byte, error) {
	var txByIDRequest cosmosTypes.QueryOutgoingTxByIDRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &txByIDRequest)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
	}

	tx, err := k.GetTxnFromOutgoingPoolByID(ctx, txByIDRequest.TxID)
	if err != nil {
		return nil, err
	}

	res, err := sdkCodec.MarshalJSONIndent(legacyQuerierCdc, tx)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryProposal(ctx sdkTypes.Context, req abciTypes.RequestQuery, keeper Keeper, legacyQuerierCdc *sdkCodec.LegacyAmino) ([]byte, error) {
	var params cosmosTypes.QueryProposalRequest
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
	}

	proposal, ok := keeper.GetProposal(ctx, params.ProposalId)
	if !ok {
		return nil, sdkErrors.Wrapf(cosmosTypes.ErrUnknownProposal, "%d", params.ProposalId)
	}

	bz, err := sdkCodec.MarshalJSONIndent(legacyQuerierCdc, proposal)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryProposals(ctx sdkTypes.Context, req abciTypes.RequestQuery, keeper Keeper, legacyQuerierCdc *sdkCodec.LegacyAmino) ([]byte, error) {
	var params cosmosTypes.QueryProposalsRequest
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
	}

	proposals := keeper.GetProposalsFiltered(ctx, params)
	if proposals == nil {
		proposals = cosmosTypes.Proposals{}
	}

	bz, err := sdkCodec.MarshalJSONIndent(legacyQuerierCdc, proposals)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryVote(ctx sdkTypes.Context, req abciTypes.RequestQuery, keeper Keeper, legacyQuerierCdc *sdkCodec.LegacyAmino) ([]byte, error) {
	var params cosmosTypes.QueryVoteRequest
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
	}

	voterAddr, err := sdkTypes.AccAddressFromBech32(params.Voter)
	if err != nil {
		return nil, err
	}

	vote, _ := keeper.GetVote(ctx, params.ProposalId, voterAddr)
	bz, err := sdkCodec.MarshalJSONIndent(legacyQuerierCdc, vote)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryVotes(ctx sdkTypes.Context, req abciTypes.RequestQuery, keeper Keeper, legacyQuerierCdc *sdkCodec.LegacyAmino) ([]byte, error) {
	var params cosmosTypes.QueryVotesRequest
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
	}

	votes := keeper.GetVotes(ctx, params.ProposalId)
	if votes == nil {
		votes = cosmosTypes.Votes{}
	} else {
		start, end := sdkClient.Paginate(len(votes), 10, 10, 100)
		if start < 0 || end < 0 {
			votes = cosmosTypes.Votes{}
		} else {
			votes = votes[start:end]
		}
	}

	bz, err := sdkCodec.MarshalJSONIndent(legacyQuerierCdc, votes)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
