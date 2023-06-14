package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

var _ types.QueryServer = &Keeper{}

func (k *Keeper) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k *Keeper) HostChain(
	goCtx context.Context,
	request *types.QueryHostChainRequest,
) (*types.QueryHostChainResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	hc, found := k.GetHostChain(ctx, request.ChainId)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryHostChainResponse{HostChain: *hc}, nil
}

func (k *Keeper) HostChains(
	goCtx context.Context,
	request *types.QueryHostChainsRequest,
) (*types.QueryHostChainsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryHostChainsResponse{HostChains: k.GetAllHostChains(ctx)}, nil
}

func (k *Keeper) Deposits(
	goCtx context.Context,
	request *types.QueryDepositsRequest,
) (*types.QueryDepositsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	hc, found := k.GetHostChain(ctx, request.ChainId)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryDepositsResponse{Deposits: k.GetDepositsForHostChain(ctx, hc.ChainId)}, nil
}

func (k *Keeper) Unbondings(
	goCtx context.Context,
	request *types.QueryUnbondingsRequest,
) (*types.QueryUnbondingsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.ChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "chain_id cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	unbondings := k.FilterUnbondings(
		ctx,
		func(u types.Unbonding) bool {
			return u.ChainId == request.ChainId
		},
	)

	return &types.QueryUnbondingsResponse{Unbondings: unbondings}, nil
}

func (k *Keeper) UserUnbondings(
	goCtx context.Context,
	request *types.QueryUserUnbondingsRequest,
) (*types.QueryUserUnbondingsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, sdkerrors.ErrKeyNotFound
	}

	userUnbondings := k.FilterUserUnbondings(
		ctx,
		func(u types.UserUnbonding) bool {
			return u.Address == address.String()
		},
	)

	return &types.QueryUserUnbondingsResponse{UserUnbondings: userUnbondings}, nil
}

func (k *Keeper) ValidatorUnbondings(
	goCtx context.Context,
	request *types.QueryValidatorUnbondingRequest,
) (*types.QueryValidatorUnbondingResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.ChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "chain_id cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	validatorUnbondings := k.FilterValidatorUnbondings(
		ctx,
		func(u types.ValidatorUnbonding) bool { return u.ChainId == request.ChainId },
	)

	return &types.QueryValidatorUnbondingResponse{ValidatorUnbondings: validatorUnbondings}, nil
}
