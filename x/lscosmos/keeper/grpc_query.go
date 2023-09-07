package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

var _ types.QueryServer = Keeper{} //nolint:staticcheck

// AllState returns genesis state
func (k Keeper) AllState(c context.Context, request *types.QueryAllStateRequest) (*types.QueryAllStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	state := k.GetGenesisState(ctx)

	return &types.QueryAllStateResponse{
		Genesis: *state,
	}, nil
}

// HostChainParams queries the host chain params
func (k Keeper) HostChainParams(c context.Context, in *types.QueryHostChainParamsRequest) (*types.QueryHostChainParamsResponse, error) {
	return nil, types.ErrDeprecated
}

// DelegationState queries the current delegation state
func (k Keeper) DelegationState(c context.Context, request *types.QueryDelegationStateRequest) (*types.QueryDelegationStateResponse, error) {
	return nil, types.ErrDeprecated

}

// AllowListedValidators queries the current allow listed validators set
func (k Keeper) AllowListedValidators(c context.Context, request *types.QueryAllowListedValidatorsRequest) (*types.QueryAllowListedValidatorsResponse, error) {
	return nil, types.ErrDeprecated
}

// CValue computes and returns the c value
func (k Keeper) CValue(c context.Context, request *types.QueryCValueRequest) (*types.QueryCValueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	hc, found := k.liquidStakeIBCKeeper.GetHostChain(ctx, "cosmoshub-4")
	if !found {
		return nil, types.ErrDeprecated
	}
	return &types.QueryCValueResponse{
		CValue: hc.CValue,
	}, nil
}

// ModuleState queries the current module state
func (k Keeper) ModuleState(c context.Context, request *types.QueryModuleStateRequest) (*types.QueryModuleStateResponse, error) {
	return nil, types.ErrDeprecated
}

// IBCTransientStore queries the current IBC transient store
func (k Keeper) IBCTransientStore(c context.Context, request *types.QueryIBCTransientStoreRequest) (*types.QueryIBCTransientStoreResponse, error) {
	return nil, types.ErrDeprecated
}

// Unclaimed queries the unclaimed entries corresponding to the input delegator address in types.QueryUnclaimedRequest
func (k Keeper) Unclaimed(c context.Context, request *types.QueryUnclaimedRequest) (*types.QueryUnclaimedResponse, error) {
	return nil, types.ErrDeprecated
}

// FailedUnbondings queries the failed unbonding entries corresponding to the input delegator address in
// types.QueryUnclaimedRequest
func (k Keeper) FailedUnbondings(c context.Context, request *types.QueryFailedUnbondingsRequest) (*types.QueryFailedUnbondingsResponse, error) {
	return nil, types.ErrDeprecated
}

// PendingUnbondings queries the pending unbonding entries corresponding to the input delegator address in
// types.QueryUnclaimedRequest
func (k Keeper) PendingUnbondings(c context.Context, request *types.QueryPendingUnbondingsRequest) (*types.QueryPendingUnbondingsResponse, error) {
	return nil, types.ErrDeprecated
}

// UnbondingEpochCValue queries the unbonding epoch c value details corresponding to the input epoch number
// in types.QueryUnbondingEpochCValueRequest
func (k Keeper) UnbondingEpochCValue(c context.Context, request *types.QueryUnbondingEpochCValueRequest) (*types.QueryUnbondingEpochCValueResponse, error) {
	return nil, types.ErrDeprecated
}

// HostAccountUndelegation queries the host account undelegation details corresponding to the input epoch number
// in types.QueryHostAccountUndelegationRequest
func (k Keeper) HostAccountUndelegation(c context.Context, request *types.QueryHostAccountUndelegationRequest) (*types.QueryHostAccountUndelegationResponse, error) {
	return nil, types.ErrDeprecated
}

// DelegatorUnbondingEpochEntry queries the delegator unbonding epoch entry details corresponding to the
// input epoch number and delegator address in types.QueryDelegatorUnbondingEpochEntryRequest
func (k Keeper) DelegatorUnbondingEpochEntry(c context.Context, request *types.QueryDelegatorUnbondingEpochEntryRequest) (*types.QueryDelegatorUnbondingEpochEntryResponse, error) {
	return nil, types.ErrDeprecated
}

// HostAccounts queries the host accounts
func (k Keeper) HostAccounts(c context.Context, request *types.QueryHostAccountsRequest) (*types.QueryHostAccountsResponse, error) {
	return nil, types.ErrDeprecated
}

// DepositModuleAccount queries the deposit module account balance
func (k Keeper) DepositModuleAccount(c context.Context, request *types.QueryDepositModuleAccountRequest) (*types.QueryDepositModuleAccountResponse, error) {
	return nil, types.ErrDeprecated
}

// DelegatorUnbondingEpochEntries queries all the delegator unbonding epoch entries corresponding to
// the input delegator address in types.QueryAllDelegatorUnbondingEpochEntriesRequest
func (k Keeper) DelegatorUnbondingEpochEntries(c context.Context, request *types.QueryAllDelegatorUnbondingEpochEntriesRequest) (*types.QueryAllDelegatorUnbondingEpochEntriesResponse, error) {
	return nil, types.ErrDeprecated
}
