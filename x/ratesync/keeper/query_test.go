package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

func (suite *IntegrationTestSuite) TestParamsQuery() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryParamsResponse{Params: params}, response)
}

func (suite *IntegrationTestSuite) TestChainQuerySingle() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNChain(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetHostChainRequest
		response *types.QueryGetHostChainResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetHostChainRequest{
				ID: msgs[0].ID,
			},
			response: &types.QueryGetHostChainResponse{HostChain: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetHostChainRequest{
				ID: msgs[1].ID,
			},
			response: &types.QueryGetHostChainResponse{HostChain: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetHostChainRequest{
				ID: uint64(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		suite.T().Run(tc.desc, func(t *testing.T) {
			response, err := keeper.HostChain(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response, response)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestAllHostChainsQueryPaginated() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNChain(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllHostChainsRequest {
		return &types.QueryAllHostChainsRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	suite.T().Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.AllHostChains(wctx, request(nil, uint64(i), uint64(step), false))
			suite.Require().NoError(err)
			suite.Require().LessOrEqual(len(resp.HostChains), step)
			suite.Require().Subset(msgs, resp.HostChains)
		}
	})
	suite.T().Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.AllHostChains(wctx, request(next, 0, uint64(step), false))
			suite.Require().NoError(err)
			suite.Require().LessOrEqual(len(resp.HostChains), step)
			suite.Require().Subset(msgs, resp.HostChains)
			next = resp.Pagination.NextKey
		}
	})
	suite.T().Run("Total", func(t *testing.T) {
		resp, err := keeper.AllHostChains(wctx, request(nil, 0, 0, true))
		suite.Require().NoError(err)
		suite.Require().Equal(len(msgs), int(resp.Pagination.Total))
		suite.Require().ElementsMatch(msgs, resp.HostChains)
	})
	suite.T().Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AllHostChains(wctx, nil)
		suite.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
