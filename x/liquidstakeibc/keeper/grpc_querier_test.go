package keeper_test

import (
	"strconv"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	MultipleTestSize  int = 10
	PaginatedTestSize int = 98
	PaginatedStep     int = 5

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
					AdminAddress:     "persistence1gztc3y3k52hjds5nqvl7h9jvfnc33spz47zcjy",
					FeeAddress:       "persistence1gztc3y3k52hjds5nqvl7h9jvfnc33spz47zcjy",
					UpperCValueLimit: decFromStr("1.1"),
					LowerCValueLimit: decFromStr("0.85"),
				},
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
	deposits := suite.app.LiquidStakeIBCKeeper.GetAllDeposits(suite.ctx)
	for _, deposit := range deposits {
		suite.app.LiquidStakeIBCKeeper.DeleteDeposit(suite.ctx, deposit)
	}
	deposits = make([]*types.Deposit, 0)
	for i := 0; i < MultipleTestSize; i += 1 {
		deposit := &types.Deposit{
			ChainId: suite.chainB.ChainID,
			Epoch:   int64(i),
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
			req:  &types.QueryDepositsRequest{ChainId: suite.chainB.ChainID},
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

func (suite *IntegrationTestSuite) TestQueryLSMDeposits() {
	deposits := make([]*types.LSMDeposit, 0)
	for i := 0; i < MultipleTestSize; i += 1 {
		deposit := &types.LSMDeposit{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: "cosmos" + strconv.Itoa(i),
			Denom:            "u" + strconv.Itoa(i),
			Shares:           sdktypes.ZeroDec(),
		}
		suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(suite.ctx, deposit)
		deposits = append(deposits, deposit)
	}

	tc := []struct {
		name string
		req  *types.QueryLSMDepositsRequest
		resp *types.QueryLSMDepositsResponse
		err  error
	}{
		{
			name: "Success",
			req:  &types.QueryLSMDepositsRequest{ChainId: suite.chainB.ChainID},
			resp: &types.QueryLSMDepositsResponse{Deposits: deposits},
		},
		{
			name: "NotFound",
			req:  &types.QueryLSMDepositsRequest{ChainId: "chain-1"},
			err:  sdkerrors.ErrKeyNotFound,
		},
		{
			name: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			resp, err := suite.app.LiquidStakeIBCKeeper.LSMDeposits(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryUnbondings() {
	unbondings := make([]*types.Unbonding, 0)
	for i := 0; i < MultipleTestSize; i += 1 {
		unbonding := &types.Unbonding{ChainId: suite.chainB.ChainID, EpochNumber: int64(i)}
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
			req:  &types.QueryUnbondingsRequest{ChainId: suite.chainB.ChainID},
			resp: &types.QueryUnbondingsResponse{Unbondings: unbondings},
		},
		{
			name: "NotFound",
			req:  &types.QueryUnbondingsRequest{ChainId: "chain-1"},
			resp: &types.QueryUnbondingsResponse{Unbondings: make([]*types.Unbonding, 0)},
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
			ChainId:     suite.chainB.ChainID,
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

func (suite *IntegrationTestSuite) TestQueryHostChainUserUnbondings() {
	chainAUserUnbondings := make([]*types.UserUnbonding, 0)
	for i := 0; i < PaginatedTestSize; i += 1 {
		userUnbonding := &types.UserUnbonding{
			ChainId:     suite.chainA.ChainID,
			Address:     TestAddress,
			EpochNumber: int64(i),
		}
		suite.app.LiquidStakeIBCKeeper.SetUserUnbonding(suite.ctx, userUnbonding)
		chainAUserUnbondings = append(chainAUserUnbondings, userUnbonding)
	}

	chainBUserUnbondings := make([]*types.UserUnbonding, 0)
	for i := 0; i < PaginatedTestSize; i += 1 {
		userUnbonding := &types.UserUnbonding{
			ChainId:     suite.chainB.ChainID,
			Address:     TestAddress,
			EpochNumber: int64(i),
		}
		suite.app.LiquidStakeIBCKeeper.SetUserUnbonding(suite.ctx, userUnbonding)
		chainBUserUnbondings = append(chainBUserUnbondings, userUnbonding)
	}

	request := func(
		chainID string,
		next []byte,
		offset, limit uint64,
		total bool,
	) *types.QueryHostChainUserUnbondingsRequest {
		return &types.QueryHostChainUserUnbondingsRequest{
			ChainId:    chainID,
			Pagination: &query.PageRequest{Key: next, Offset: offset, Limit: limit, CountTotal: total},
		}
	}

	suite.T().Run("ByOffset", func(t *testing.T) {
		for i := 0; i < PaginatedTestSize; i += PaginatedStep {
			resp, err := suite.app.LiquidStakeIBCKeeper.HostChainUserUnbondings(
				suite.ctx,
				request(suite.chainB.ChainID, nil, uint64(i), uint64(PaginatedStep), false),
			)
			suite.Require().NoError(err)
			suite.Require().LessOrEqual(len(resp.UserUnbondings), PaginatedStep)
			suite.Require().Subset(chainBUserUnbondings, resp.UserUnbondings)
		}
	})

	suite.T().Run("ByKey", func(t *testing.T) {
		var next []byte
		for i := 0; i < PaginatedTestSize; i += PaginatedStep {
			resp, err := suite.app.LiquidStakeIBCKeeper.HostChainUserUnbondings(
				suite.ctx,
				request(suite.chainA.ChainID, next, 0, uint64(PaginatedStep), false),
			)
			suite.Require().NoError(err)
			suite.Require().LessOrEqual(len(resp.UserUnbondings), PaginatedStep)
			suite.Require().Subset(chainAUserUnbondings, resp.UserUnbondings)
			next = resp.Pagination.NextKey
		}
	})

	suite.T().Run("Total", func(t *testing.T) {
		resp, err := suite.app.LiquidStakeIBCKeeper.HostChainUserUnbondings(
			suite.ctx,
			request(suite.chainB.ChainID, nil, 0, 0, true),
		)
		suite.Require().NoError(err)
		suite.Require().Equal(len(chainBUserUnbondings), int(resp.Pagination.Total))
		suite.Require().ElementsMatch(chainBUserUnbondings, resp.UserUnbondings)
	})

	suite.T().Run("Total Empty", func(t *testing.T) {
		resp, err := suite.app.LiquidStakeIBCKeeper.HostChainUserUnbondings(
			suite.ctx,
			request("non-existing-chain", nil, 0, 0, true),
		)
		suite.Require().NoError(err)
		suite.Require().Equal(0, int(resp.Pagination.Total))
	})

	suite.T().Run("Invalid Request", func(t *testing.T) {
		_, err := suite.app.LiquidStakeIBCKeeper.HostChainUserUnbondings(
			suite.ctx,
			request("", nil, 0, 0, true),
		)
		suite.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "chain id cannot be empty"))
	})

	suite.T().Run("Invalid Request", func(t *testing.T) {
		_, err := suite.app.LiquidStakeIBCKeeper.HostChainUserUnbondings(suite.ctx, nil)
		suite.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "empty request"))
	})
}

func (suite *IntegrationTestSuite) TestQueryValidatorUnbondings() {
	validatorUnbondings := make([]*types.ValidatorUnbonding, 0)
	for i := 0; i < MultipleTestSize; i += 1 {
		validatorUnbonding := &types.ValidatorUnbonding{
			ChainId:          suite.chainB.ChainID,
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
			req:  &types.QueryValidatorUnbondingRequest{ChainId: suite.chainB.ChainID},
			resp: &types.QueryValidatorUnbondingResponse{ValidatorUnbondings: validatorUnbondings},
		},
		{
			name: "NotFound",
			req:  &types.QueryValidatorUnbondingRequest{ChainId: "chain-1"},
			resp: &types.QueryValidatorUnbondingResponse{ValidatorUnbondings: make([]*types.ValidatorUnbonding, 0)},
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

func (suite *IntegrationTestSuite) TestQueryUnbonding() {
	unbonding := &types.Unbonding{ChainId: suite.chainB.ChainID, EpochNumber: int64(1)}
	suite.app.LiquidStakeIBCKeeper.SetUnbonding(suite.ctx, unbonding)

	tc := []struct {
		name string
		req  *types.QueryUnbondingRequest
		resp *types.QueryUnbondingResponse
		err  error
	}{
		{
			name: "Valid",
			req: &types.QueryUnbondingRequest{
				ChainId: suite.chainB.ChainID,
				Epoch:   1,
			},
			resp: &types.QueryUnbondingResponse{Unbonding: unbonding},
			err:  nil,
		}, {
			name: "NotFound",
			req:  &types.QueryUnbondingRequest{ChainId: "chain-1"},
			resp: nil,
			err:  sdkerrors.ErrKeyNotFound,
		}, {
			name: "InvalidRequest",
			req:  &types.QueryUnbondingRequest{ChainId: ""},
			err:  status.Error(codes.InvalidArgument, "chain_id cannot be empty"),
		}, {
			name: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "empty request"),
		},
	}
	for _, t := range tc {
		suite.Run(t.name, func() {
			resp, err := suite.app.LiquidStakeIBCKeeper.Unbonding(suite.ctx, t.req)

			suite.Require().Equal(t.err, err)
			suite.Require().Equal(t.resp, resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryDepositAccountBalance() {
	err := testutil.FundAccount(suite.app.BankKeeper, suite.ctx,
		authtypes.NewModuleAddress(types.DepositModuleAccount),
		sdktypes.NewCoins(sdktypes.NewInt64Coin("ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", 1000)))
	suite.Require().NoError(err)

	tc := []struct {
		name string
		req  *types.QueryDepositAccountBalanceRequest
		resp *types.QueryDepositAccountBalanceResponse
		err  error
	}{{
		name: "Valid",
		req:  &types.QueryDepositAccountBalanceRequest{ChainId: suite.chainB.ChainID},
		resp: &types.QueryDepositAccountBalanceResponse{Balance: sdktypes.NewInt64Coin("ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", 1000)},
		err:  nil,
	}, {
		name: "NotFound",
		req:  &types.QueryDepositAccountBalanceRequest{ChainId: "chain-1"},
		err:  sdkerrors.ErrKeyNotFound,
	}, {
		name: "InvalidRequest",
		err:  status.Error(codes.InvalidArgument, "empty request"),
	}}

	for _, t := range tc {
		suite.Run(t.name, func() {
			resp, err := suite.app.LiquidStakeIBCKeeper.DepositAccountBalance(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryExchangeRate() {
	tc := []struct {
		name string
		req  *types.QueryExchangeRateRequest
		resp *types.QueryExchangeRateResponse
		err  error
	}{{
		name: "Valid",
		req:  &types.QueryExchangeRateRequest{ChainId: suite.chainB.ChainID},
		resp: &types.QueryExchangeRateResponse{Rate: sdktypes.OneDec()},
		err:  nil,
	}, {
		name: "NotFound",
		req:  &types.QueryExchangeRateRequest{ChainId: "chain-1"},
		err:  sdkerrors.ErrKeyNotFound,
	}, {
		name: "InvalidRequest",
		err:  status.Error(codes.InvalidArgument, "empty request"),
	}}

	for _, t := range tc {
		suite.Run(t.name, func() {
			resp, err := suite.app.LiquidStakeIBCKeeper.ExchangeRate(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryRedelegations() {
	tc := []struct {
		name string
		req  *types.QueryRedelegationsRequest
		resp *types.QueryRedelegationsResponse
		err  error
	}{{
		name: "Valid",
		req:  &types.QueryRedelegationsRequest{ChainId: suite.chainB.ChainID},
		resp: &types.QueryRedelegationsResponse{Redelegations: nil},
		err:  nil,
	}, {
		name: "NotFound",
		req:  &types.QueryRedelegationsRequest{ChainId: "chain-1"},
		err:  sdkerrors.ErrKeyNotFound,
	}, {
		name: "InvalidRequest",
		err:  status.Error(codes.InvalidArgument, "empty request"),
	}}

	for _, t := range tc {
		suite.Run(t.name, func() {
			resp, err := suite.app.LiquidStakeIBCKeeper.Redelegations(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryRedelegationTx() {
	tc := []struct {
		name string
		req  *types.QueryRedelegationTxRequest
		resp *types.QueryRedelegationTxResponse
		err  error
	}{{
		name: "Valid",
		req:  &types.QueryRedelegationTxRequest{ChainId: suite.chainB.ChainID},
		resp: &types.QueryRedelegationTxResponse{RedelegationTx: []*types.RedelegateTx{}},
		err:  nil,
	}, {
		name: "NotFound",
		req:  &types.QueryRedelegationTxRequest{ChainId: "chain-1"},
		err:  sdkerrors.ErrKeyNotFound,
	}, {
		name: "InvalidRequest",
		err:  status.Error(codes.InvalidArgument, "empty request"),
	}}

	for _, t := range tc {
		suite.Run(t.name, func() {
			resp, err := suite.app.LiquidStakeIBCKeeper.RedelegationTx(suite.ctx, t.req)

			suite.Require().Equal(err, t.err)
			suite.Require().Equal(resp, t.resp)
		})
	}
}
