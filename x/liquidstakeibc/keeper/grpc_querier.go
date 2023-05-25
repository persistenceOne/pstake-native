package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
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

	store := ctx.KVStore(k.storeKey)
	hcStore := prefix.NewStore(store, types.HostChainKey)

	var hostChains []types.HostChain
	pagination, err := query.Paginate(hcStore, request.Pagination, func(key []byte, value []byte) error {
		var hc types.HostChain
		if err := k.cdc.Unmarshal(value, &hc); err != nil {
			return err
		}

		hostChains = append(hostChains, hc)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryHostChainsResponse{
		HostChains: hostChains,
		Pagination: pagination,
	}, nil
}

func (k *Keeper) Deposits(
	goCtx context.Context,
	request *types.QueryDepositsRequest,
) (*types.QueryDepositsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	udStore := prefix.NewStore(store, types.DepositKey)

	var deposits []types.Deposit
	pagination, err := query.Paginate(udStore, request.Pagination, func(key []byte, value []byte) error {
		var deposit types.Deposit
		if err := k.cdc.Unmarshal(value, &deposit); err != nil {
			return err
		}

		deposits = append(deposits, deposit)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDepositsResponse{
		Deposits:   deposits,
		Pagination: pagination,
	}, nil
}

func (k *Keeper) Unbondings(
	goCtx context.Context,
	request *types.QueryUnbondingsRequest,
) (*types.QueryUnbondingsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.HostDenom == "" {
		return nil, status.Error(codes.InvalidArgument, "host_denom cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	hc, found := k.GetHostChainFromHostDenom(ctx, request.HostDenom)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	unbondings := k.FilterUnbondings(
		ctx,
		func(u types.Unbonding) bool {
			return u.ChainId == hc.ChainId
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
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
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
	if request.HostDenom == "" {
		return nil, status.Error(codes.InvalidArgument, "host_denom cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	hc, found := k.GetHostChainFromHostDenom(ctx, request.HostDenom)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	validatorUnbondings := k.FilterValidatorUnbondings(
		ctx,
		func(u types.ValidatorUnbonding) bool { return u.ChainId == hc.ChainId },
	)

	return &types.QueryValidatorUnbondingResponse{ValidatorUnbondings: validatorUnbondings}, nil
}
