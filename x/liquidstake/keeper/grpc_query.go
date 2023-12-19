package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries the parameters of the liquidstake module.
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// LiquidValidators queries all liquid validators.
func (k Querier) LiquidValidators(c context.Context, req *types.QueryLiquidValidatorsRequest) (*types.QueryLiquidValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryLiquidValidatorsResponse{LiquidValidators: k.GetAllLiquidValidatorStates(ctx)}, nil
}

// States queries states of liquid stake module.
func (k Querier) States(c context.Context, req *types.QueryStatesRequest) (*types.QueryStatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryStatesResponse{NetAmountState: k.GetNetAmountState(ctx)}, nil
}
