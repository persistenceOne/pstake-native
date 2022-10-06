package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (k Keeper) Unclaimed(c context.Context, request *types.QueryUnclaimedRequest) (*types.QueryUnclaimedResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// get delegator account address from request
	delegatorAddress, err := sdk.AccAddressFromBech32(request.DelegatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	var queryResponse types.QueryUnclaimedResponse

	delegatorUnbondingEpochEntries := k.IterateDelegatorUnbondingEpochEntry(ctx, delegatorAddress)
	for _, entry := range delegatorUnbondingEpochEntries {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, entry.EpochNumber)

		// sort for all the cases
		if unbondingEpochCValue.IsMatured && unbondingEpochCValue.EpochNumber > 0 {
			// append to ready to claim entries
			queryResponse.Unclaimed = append(queryResponse.Unclaimed, unbondingEpochCValue)
		}
	}

	return &queryResponse, nil
}

func (k Keeper) FailedUnbondings(c context.Context, request *types.QueryFailedUnbondingsRequest) (*types.QueryFailedUnbondingsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// get delegator account address from request
	delegatorAddress, err := sdk.AccAddressFromBech32(request.DelegatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	var queryResponse types.QueryFailedUnbondingsResponse

	delegatorUnbondingEpochEntries := k.IterateDelegatorUnbondingEpochEntry(ctx, delegatorAddress)
	for _, entry := range delegatorUnbondingEpochEntries {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, entry.EpochNumber)
		if unbondingEpochCValue.IsFailed && unbondingEpochCValue.EpochNumber > 0 {
			// append to failed entries for which stkAtom should be claimed again
			queryResponse.FailedUnbondings = append(queryResponse.FailedUnbondings, unbondingEpochCValue)
		}
	}

	return &queryResponse, nil
}

func (k Keeper) PendingUnbondings(c context.Context, request *types.QueryPendingUnbondingsRequest) (*types.QueryPendingUnbondingsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// get delegator account address from request
	delegatorAddress, err := sdk.AccAddressFromBech32(request.DelegatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	var queryResponse types.QueryPendingUnbondingsResponse

	delegatorUnbondingEpochEntries := k.IterateDelegatorUnbondingEpochEntry(ctx, delegatorAddress)
	for _, entry := range delegatorUnbondingEpochEntries {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, entry.EpochNumber)
		if !unbondingEpochCValue.IsFailed && !unbondingEpochCValue.IsMatured && unbondingEpochCValue.EpochNumber > 0 {
			// append to in progress entries
			queryResponse.PendingUnbondings = append(queryResponse.PendingUnbondings, unbondingEpochCValue)
		}
	}

	return &queryResponse, nil
}

func (k Keeper) UnbondingEpochCValue(c context.Context, request *types.QueryUnbondingEpochCValueRequest) (*types.QueryUnbondingEpochCValueResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.EpochNumber <= 0 {
		return nil, status.Error(codes.InvalidArgument, "epoch number less than equal to 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, request.EpochNumber)

	return &types.QueryUnbondingEpochCValueResponse{UnbondingEpochCValue: unbondingEpochCValue}, nil
}

func (k Keeper) HostAccountUndelegation(c context.Context, request *types.QueryHostAccountUndelegationRequest) (*types.QueryHostAccountUndelegationResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.EpochNumber <= 0 {
		return nil, status.Error(codes.InvalidArgument, "epoch number less than equal to 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	hostAccountUndelegation, err := k.GetHostAccountUndelegationForEpoch(ctx, request.EpochNumber)
	if err != nil {
		return nil, err
	}

	return &types.QueryHostAccountUndelegationResponse{HostAccountUndelegation: hostAccountUndelegation}, nil
}

func (k Keeper) DelegatorUnbondingEpochEntry(c context.Context, request *types.QueryDelegatorUnbondingEpochEntryRequest) (*types.QueryDelegatorUnbondingEpochEntryResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if request.EpochNumber <= 0 {
		return nil, status.Error(codes.InvalidArgument, "epoch number less than equal to 0")
	}
	if request.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// get delegator account address from request
	delegatorAddress, err := sdk.AccAddressFromBech32(request.DelegatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	unbondingEpochEntry := k.GetDelegatorUnbondingEpochEntry(ctx, delegatorAddress, request.EpochNumber)

	return &types.QueryDelegatorUnbondingEpochEntryResponse{DelegatorUnbodingEpochEntry: unbondingEpochEntry}, nil
}

func (k Keeper) RewardsBoosterAccount(c context.Context, request *types.QueryRewardBoosterAccountRequest) (*types.QueryRewardBoosterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	balance := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.RewardBoosterModuleAccount), k.GetIBCDenom(ctx))

	return &types.QueryRewardBoosterAccountResponse{Balance: balance}, nil
}

func (k Keeper) HostAccounts(c context.Context, request *types.QueryHostAccountsRequest) (*types.QueryHostAccountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	hostAccounts := k.GetHostAccounts(ctx)
	return &types.QueryHostAccountsResponse{HostAccounts: hostAccounts}, nil
}
