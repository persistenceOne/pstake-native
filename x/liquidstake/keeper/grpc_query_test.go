package keeper_test

import (
	"cosmossdk.io/math"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

func (s *KeeperTestSuite) TestGRPCParams() {
	resp, err := s.querier.Params(sdk.WrapSDKContext(s.ctx), &types.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(s.keeper.GetParams(s.ctx), resp.Params)
}

func (s *KeeperTestSuite) TestGRPCQueries() {
	_, valOpers, _ := s.CreateValidators([]int64{1000000, 2000000, 3000000})
	params := s.keeper.GetParams(s.ctx)
	params.MinLiquidStakeAmount = math.NewInt(50000)
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx)

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(1)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx)

	// Test LiquidValidators grpc query
	res := s.keeper.GetAllLiquidValidatorStates(s.ctx)
	resp, err := s.querier.LiquidValidators(sdk.WrapSDKContext(s.ctx), &types.QueryLiquidValidatorsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(resp.LiquidValidators, res)

	resp, err = s.querier.LiquidValidators(sdk.WrapSDKContext(s.ctx), nil)
	s.Require().Nil(resp)
	s.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "invalid request"))

	// Test States grpc query
	respStates, err := s.querier.States(sdk.WrapSDKContext(s.ctx), &types.QueryStatesRequest{})
	resNetAmountState := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(respStates.NetAmountState, resNetAmountState)

	respStates, err = s.querier.States(sdk.WrapSDKContext(s.ctx), nil)
	s.Require().Nil(respStates)
	s.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "invalid request"))

	// Test Params grpc query
	respParams, err := s.querier.Params(sdk.WrapSDKContext(s.ctx), &types.QueryParamsRequest{})
	resParams := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(respParams.Params.LiquidBondDenom, resParams.LiquidBondDenom)
	s.Require().Equal(respParams.Params.WhitelistedValidators[0].ValidatorAddress, valOpers[0].String())
	s.Require().Equal(respParams.Params.WhitelistedValidators[1].ValidatorAddress, valOpers[1].String())
	s.Require().Equal(respParams.Params.WhitelistedValidators[2].ValidatorAddress, valOpers[2].String())
}
