package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	MultipleTestSize int = 10

	TestAddress string = "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld"
)

func (suite *IntegrationTestSuite) TestQueryParams() {
	suite.app.LiquidStakeIBCKeeper.SetParams(suite.ctx, types.DefaultParams())

	tc := []struct {
		name string
		req  *types.QueryParamsRequest
		resp *types.QueryParamsResponse
		err  error
	}{
		{
			name: "Success",
			req:  &types.QueryParamsRequest{},
			resp: &types.QueryParamsResponse{
				Params: types.Params{
					AdminAddress: "persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr",
					FeeAddress:   "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld"},
			},
		},
		{
			name: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {

			resp, err := suite.app.LiquidStakeIBCKeeper.Params(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryHostChain() {

	hostChains := suite.app.LiquidStakeIBCKeeper.GetAllHostChains(suite.ctx)

	tc := []struct {
		name string
		req  *types.QueryHostChainRequest
		resp *types.QueryHostChainResponse
		err  error
	}{
		{
			name: "Success",
			req:  &types.QueryHostChainRequest{ChainId: hostChains[0].ChainId},
			resp: &types.QueryHostChainResponse{HostChain: *hostChains[0]},
		},
		{
			name: "NotFound",
			req:  &types.QueryHostChainRequest{ChainId: "not-registered-chain"},
			err:  sdkerrors.ErrKeyNotFound,
		},
		{
			name: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {

			resp, err := suite.app.LiquidStakeIBCKeeper.HostChain(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryHostChains() {
	hcs := suite.app.LiquidStakeIBCKeeper.GetAllHostChains(suite.ctx)

	tc := []struct {
		name string
		req  *types.QueryHostChainsRequest
		resp *types.QueryHostChainsResponse
		err  error
	}{
		{
			name: "Success",
			req:  &types.QueryHostChainsRequest{},
			resp: &types.QueryHostChainsResponse{HostChains: hcs},
		},
		{
			name: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {

			resp, err := suite.app.LiquidStakeIBCKeeper.HostChains(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryDeposits() {

	deposits := make([]*types.Deposit, 0)
	for i := 0; i < MultipleTestSize; i += 1 {
		deposit := &types.Deposit{
			ChainId: suite.path.EndpointB.Chain.ChainID,
			Epoch:   sdk.NewInt(int64(i)),
		}
		suite.app.LiquidStakeIBCKeeper.SetDeposit(suite.ctx, deposit)
		deposits = append(deposits, deposit)
	}

	tc := []struct {
		name string
		req  *types.QueryDepositsRequest
		resp *types.QueryDepositsResponse
		err  error
	}{
		{
			name: "Success",
			req:  &types.QueryDepositsRequest{ChainId: suite.path.EndpointB.Chain.ChainID},
			resp: &types.QueryDepositsResponse{Deposits: deposits},
		},
		{
			name: "NotFound",
			req:  &types.QueryDepositsRequest{ChainId: "chain-1"},
			err:  sdkerrors.ErrKeyNotFound,
		},
		{
			name: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {

			resp, err := suite.app.LiquidStakeIBCKeeper.Deposits(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryUnbondings() {
	unbondings := make([]*types.Unbonding, 0)
	for i := 0; i < MultipleTestSize; i += 1 {
		unbonding := &types.Unbonding{ChainId: suite.path.EndpointB.Chain.ChainID, EpochNumber: int64(i)}
		suite.app.LiquidStakeIBCKeeper.SetUnbonding(suite.ctx, unbonding)
		unbondings = append(unbondings, unbonding)
	}

	tc := []struct {
		name string
		req  *types.QueryUnbondingsRequest
		resp *types.QueryUnbondingsResponse
		err  error
	}{
		{
			name: "Success",
			req:  &types.QueryUnbondingsRequest{ChainId: suite.path.EndpointB.Chain.ChainID},
			resp: &types.QueryUnbondingsResponse{Unbondings: unbondings},
		},
		{
			name: "NotFound",
			req:  &types.QueryUnbondingsRequest{ChainId: "chain-1"},
			err:  sdkerrors.ErrKeyNotFound,
		},
		{
			name: "InvalidRequest",
			req:  &types.QueryUnbondingsRequest{ChainId: ""},
			err:  status.Error(codes.InvalidArgument, "chain_id cannot be empty"),
		},
		{
			name: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {

			resp, err := suite.app.LiquidStakeIBCKeeper.Unbondings(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryUserUnbondings() {
	userUnbondings := make([]*types.UserUnbonding, 0)
	for i := 0; i < MultipleTestSize; i += 1 {
		userUnbonding := &types.UserUnbonding{
			ChainId:     suite.path.EndpointB.Chain.ChainID,
			Address:     TestAddress,
			EpochNumber: int64(i),
		}
		suite.app.LiquidStakeIBCKeeper.SetUserUnbonding(suite.ctx, userUnbonding)
		userUnbondings = append(userUnbondings, userUnbonding)
	}

	tc := []struct {
		name string
		req  *types.QueryUserUnbondingsRequest
		resp *types.QueryUserUnbondingsResponse
		err  error
	}{
		{
			name: "Success",
			req:  &types.QueryUserUnbondingsRequest{Address: TestAddress},
			resp: &types.QueryUserUnbondingsResponse{UserUnbondings: userUnbondings},
		},
		{
			name: "NotFound",
			req:  &types.QueryUserUnbondingsRequest{Address: "persistence1234"},
			err:  sdkerrors.ErrKeyNotFound,
		},
		{
			name: "InvalidRequest",
			req:  &types.QueryUserUnbondingsRequest{Address: ""},
			err:  status.Error(codes.InvalidArgument, "address cannot be empty"),
		},
		{
			name: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {

			resp, err := suite.app.LiquidStakeIBCKeeper.UserUnbondings(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryValidatorUnbondings() {
	validatorUnbondings := make([]*types.ValidatorUnbonding, 0)
	for i := 0; i < MultipleTestSize; i += 1 {
		validatorUnbonding := &types.ValidatorUnbonding{
			ChainId:          suite.path.EndpointB.Chain.ChainID,
			ValidatorAddress: TestAddress,
			EpochNumber:      int64(i),
		}
		suite.app.LiquidStakeIBCKeeper.SetValidatorUnbonding(suite.ctx, validatorUnbonding)
		validatorUnbondings = append(validatorUnbondings, validatorUnbonding)
	}

	tc := []struct {
		name string
		req  *types.QueryValidatorUnbondingRequest
		resp *types.QueryValidatorUnbondingResponse
		err  error
	}{
		{
			name: "Success",
			req:  &types.QueryValidatorUnbondingRequest{ChainId: suite.path.EndpointB.Chain.ChainID},
			resp: &types.QueryValidatorUnbondingResponse{ValidatorUnbondings: validatorUnbondings},
		},
		{
			name: "NotFound",
			req:  &types.QueryValidatorUnbondingRequest{ChainId: "chain-1"},
			err:  sdkerrors.ErrKeyNotFound,
		},
		{
			name: "InvalidRequest",
			req:  &types.QueryValidatorUnbondingRequest{ChainId: ""},
			err:  status.Error(codes.InvalidArgument, "chain_id cannot be empty"),
		},
		{
			name: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {

			resp, err := suite.app.LiquidStakeIBCKeeper.ValidatorUnbondings(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}
