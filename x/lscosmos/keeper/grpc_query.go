package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

var _ types.QueryServer = Keeper{}

// HostChainParams returns the stored host chain params set through proposal.
func (k Keeper) HostChainParams(c context.Context, in *types.QueryHostChainParamsRequest) (*types.QueryHostChainParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	hostChainParams := k.GetHostChainParams(ctx)

	return &types.QueryHostChainParamsResponse{
		HostChainParams: hostChainParams,
	}, nil
}

func (k Keeper) DelegationState(c context.Context, request *types.QueryDelegationStateRequest) (*types.QueryDelegationStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	delegationState := k.GetDelegationState(ctx)

	return &types.QueryDelegationStateResponse{
		DelegationState: delegationState,
	}, nil

}

func (k Keeper) AllowListedValidators(c context.Context, request *types.QueryAllowListedValidatorsRequest) (*types.QueryAllowListedValidatorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	allowListedValidators := k.GetAllowListedValidators(ctx)

	return &types.QueryAllowListedValidatorsResponse{
		AllowListedValidators: allowListedValidators,
	}, nil
}

func (k Keeper) CValue(c context.Context, request *types.QueryCValueRequest) (*types.QueryCValueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	cValue := k.GetCValue(ctx)

	return &types.QueryCValueResponse{
		CValue: cValue,
	}, nil
}

func (k Keeper) ModuleState(c context.Context, request *types.QueryModuleStateRequest) (*types.QueryModuleStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	moduleState := k.GetModuleState(ctx)

	return &types.QueryModuleStateResponse{
		ModuleState: moduleState,
	}, nil
}

func (k Keeper) IBCTransientStore(c context.Context, request *types.QueryIBCTransientStoreRequest) (*types.QueryIBCTransientStoreResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ibcTransientStore := k.GetIBCTransientStore(ctx)

	return &types.QueryIBCTransientStoreResponse{
		IBCTransientStore: ibcTransientStore,
	}, nil
}
