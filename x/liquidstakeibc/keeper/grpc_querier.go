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

func (k *Keeper) Unbonding(
	goCtx context.Context,
	request *types.QueryUnbondingRequest,
) (*types.QueryUnbondingResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.EpochNumber <= 0 {
		return nil, status.Error(codes.InvalidArgument, "epoch number less than equal to 0")
	}
	if request.HostDenom == "" {
		return nil, status.Error(codes.InvalidArgument, "host_denom cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	hc, found := k.GetHostChainFromHostDenom(ctx, request.HostDenom)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	unbonding, found := k.GetUnbonding(ctx, hc.ChainId, request.EpochNumber)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryUnbondingResponse{Unbonding: *unbonding}, nil
}

func (k *Keeper) UserUnbonding(
	goCtx context.Context,
	request *types.QueryUserUnbondingsRequest,
) (*types.QueryUserUnbondingsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.EpochNumber <= 0 {
		return nil, status.Error(codes.InvalidArgument, "epoch number less than equal to 0")
	}
	if request.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}
	if request.HostDenom == "" {
		return nil, status.Error(codes.InvalidArgument, "host_denom cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	hc, found := k.GetHostChainFromHostDenom(ctx, request.HostDenom)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	userUnbonding, found := k.GetUserUnbonding(ctx, hc.ChainId, address.String(), request.EpochNumber)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryUserUnbondingsResponse{UserUnbonding: *userUnbonding}, nil
}
