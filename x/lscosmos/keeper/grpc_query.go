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

func (k Keeper) ReadyToClaim(c context.Context, request *types.QueryReadyToClaimRequest) (*types.QueryReadyToClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// get delegator account address from request
	delegatorAddress, err := sdk.AccAddressFromBech32(request.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	var queryResponse types.QueryReadyToClaimResponse

	delegatorUnbondingEpochEntries := k.IterateDelegatorUnbondingEpochEntry(ctx, delegatorAddress)
	for _, entry := range delegatorUnbondingEpochEntries {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, entry.EpochNumber)

		// sort for all the cases
		if unbondingEpochCValue.IsMatured {
			var responseEntry types.Entry

			// get c value from the UnbondingEpochCValue struct
			// calculate claimable amount from un inverse c value
			claimableAmount := entry.Amount.Amount.ToDec().Quo(unbondingEpochCValue.GetUnbondingEpochCValue())

			// calculate claimable coin and community coin to be sent to delegator account and community pool respectively
			claimableCoin, _ := sdk.NewDecCoinFromDec(k.GetIBCDenom(ctx), claimableAmount).TruncateDecimal()

			// fill in the details of responseEntry
			// amount : claimable coin
			responseEntry.DelegatorAddress = entry.DelegatorAddress
			responseEntry.Amount = claimableCoin
			responseEntry.EpochNumber = entry.EpochNumber
			responseEntry.BatchCValue = unbondingEpochCValue.GetUnbondingEpochCValue()

			// append to ready to claim entries
			queryResponse.ReadyToClaim = append(queryResponse.ReadyToClaim, responseEntry)
		}
	}

	return &queryResponse, nil
}

func (k Keeper) UnbondingFailed(c context.Context, request *types.QueryUnbondFailRequest) (*types.QueryUnbondFailResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// get delegator account address from request
	delegatorAddress, err := sdk.AccAddressFromBech32(request.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	var queryResponse types.QueryUnbondFailResponse

	delegatorUnbondingEpochEntries := k.IterateDelegatorUnbondingEpochEntry(ctx, delegatorAddress)
	for _, entry := range delegatorUnbondingEpochEntries {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, entry.EpochNumber)
		if unbondingEpochCValue.IsTimedOut {
			var responseEntry types.Entry

			// fill in the details of responseEntry
			// amount : entry amount (failed unbonding)
			responseEntry.DelegatorAddress = entry.DelegatorAddress
			responseEntry.Amount = entry.Amount

			// append to failed entries for which stkAtom should be claimed again
			queryResponse.UnbondFail = append(queryResponse.UnbondFail, responseEntry)
		}
	}

	return &queryResponse, nil
}

func (k Keeper) UnbondInProgress(c context.Context, request *types.QueryUnbondInProgressRequest) (*types.QueryUnbondInProgressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// get delegator account address from request
	delegatorAddress, err := sdk.AccAddressFromBech32(request.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	var queryResponse types.QueryUnbondInProgressResponse

	delegatorUnbondingEpochEntries := k.IterateDelegatorUnbondingEpochEntry(ctx, delegatorAddress)
	for _, entry := range delegatorUnbondingEpochEntries {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, entry.EpochNumber)

		if !unbondingEpochCValue.IsTimedOut && !unbondingEpochCValue.IsMatured {
			var responseEntry types.Entry

			// get c value from the UnbondingEpochCValue struct
			// calculate claimable amount from un inverse c value
			claimableAmount := entry.Amount.Amount.ToDec().Quo(unbondingEpochCValue.GetUnbondingEpochCValue())

			// calculate claimable coin and community coin to be sent to delegator account and community pool respectively
			claimableCoin, _ := sdk.NewDecCoinFromDec(k.GetIBCDenom(ctx), claimableAmount).TruncateDecimal()

			hostAccountUndelegationForEpoch, err := k.GetHostAccountUndelegationForEpoch(ctx, entry.EpochNumber)
			if err != nil {
				return nil, err
			}

			// fill in the details of responseEntry
			// amount : claimable amount
			responseEntry.DelegatorAddress = entry.DelegatorAddress
			responseEntry.Amount = claimableCoin
			responseEntry.BatchCValue = unbondingEpochCValue.GetUnbondingEpochCValue()
			responseEntry.EpochNumber = entry.EpochNumber
			responseEntry.CompletionTime = hostAccountUndelegationForEpoch.CompletionTime

			// append to in progress entries
			queryResponse.InProgress = append(queryResponse.InProgress, responseEntry)
		}
	}

	return &queryResponse, nil
}
