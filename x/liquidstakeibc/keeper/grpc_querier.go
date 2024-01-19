package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (k *Keeper) LSMDeposits(
	goCtx context.Context,
	request *types.QueryLSMDepositsRequest,
) (*types.QueryLSMDepositsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	hc, found := k.GetHostChain(ctx, request.ChainId)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	deposits := k.FilterLSMDeposits(
		ctx,
		func(d types.LSMDeposit) bool {
			return d.ChainId == hc.ChainId
		},
	)

	return &types.QueryLSMDepositsResponse{Deposits: deposits}, nil
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

func (k *Keeper) Unbonding(
	goCtx context.Context,
	request *types.QueryUnbondingRequest,
) (*types.QueryUnbondingResponse, error) {
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
			return u.ChainId == request.ChainId && u.EpochNumber == request.Epoch
		},
	)

	if len(unbondings) == 0 {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryUnbondingResponse{Unbonding: unbondings[0]}, nil
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

func (k *Keeper) HostChainUserUnbondings(
	goCtx context.Context,
	request *types.QueryHostChainUserUnbondingsRequest,
) (*types.QueryHostChainUserUnbondingsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.ChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "chain id cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	userUnbondingStore := prefix.NewStore(store, types.UserUnbondingKey)

	var userUnbondings []*types.UserUnbonding
	pageRes, err := query.FilteredPaginate(
		userUnbondingStore,
		request.Pagination,
		func(key, value []byte, accumulate bool) (bool, error) {
			if accumulate {
				var uu types.UserUnbonding
				if err := k.cdc.Unmarshal(value, &uu); err != nil {
					return false, err
				}

				if uu.ChainId == request.ChainId {
					userUnbondings = append(userUnbondings, &uu)
					return true, nil
				}

				return false, nil
			}

			return true, nil
		})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryHostChainUserUnbondingsResponse{UserUnbondings: userUnbondings, Pagination: pageRes}, nil
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

func (k *Keeper) DepositAccountBalance(
	goCtx context.Context,
	request *types.QueryDepositAccountBalanceRequest,
) (*types.QueryDepositAccountBalanceResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	hc, found := k.GetHostChain(ctx, request.ChainId)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryDepositAccountBalanceResponse{
		Balance: k.bankKeeper.GetBalance(
			ctx,
			authtypes.NewModuleAddress(types.DepositModuleAccount), hc.IBCDenom()),
	}, nil
}

func (k *Keeper) ExchangeRate(
	goCtx context.Context,
	request *types.QueryExchangeRateRequest,
) (*types.QueryExchangeRateResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	hc, found := k.GetHostChain(ctx, request.ChainId)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryExchangeRateResponse{Rate: hc.CValue}, nil
}

func (k *Keeper) Redelegations(goCtx context.Context, request *types.QueryRedelegationsRequest) (*types.QueryRedelegationsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.ChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "chain_id cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	hc, found := k.GetHostChain(ctx, request.ChainId)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}
	redels, _ := k.GetRedelegations(ctx, hc.ChainId)

	return &types.QueryRedelegationsResponse{Redelegations: redels}, nil
}

func (k *Keeper) RedelegationTx(goCtx context.Context, request *types.QueryRedelegationTxRequest) (*types.QueryRedelegationTxResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.ChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "chain_id cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	hc, found := k.GetHostChain(ctx, request.ChainId)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}
	redelTxs := k.FilterRedelegationTx(ctx, func(d types.RedelegateTx) bool {
		return d.ChainId == hc.ChainId
	})
	return &types.QueryRedelegationTxResponse{RedelegationTx: redelTxs}, nil
}
