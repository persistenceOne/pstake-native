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

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) HostChain(
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

	return &types.QueryHostChainResponse{HostChain: hc}, nil
}

func (k Keeper) HostChains(
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
