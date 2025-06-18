package liquidstakeibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/keeper"
	"github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
)

// InitGenesis initializes the liquidstakeibc module's state from a given genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState *types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	for _, hc := range genState.HostChains {
		k.SetHostChain(ctx, hc)
	}

	for _, deposit := range genState.Deposits {
		k.SetDeposit(ctx, deposit)
	}

	for _, unbonding := range genState.Unbondings {
		k.SetUnbonding(ctx, unbonding)
	}
	for _, userUnbonding := range genState.UserUnbondings {
		k.SetUserUnbonding(ctx, userUnbonding)
	}
	for _, valUnbonding := range genState.ValidatorUnbondings {
		k.SetValidatorUnbonding(ctx, valUnbonding)
	}

	k.GetDepositModuleAccount(ctx)
	k.GetUndelegationModuleAccount(ctx)
}

// ExportGenesis returns the liquidstakeibc module's genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:              k.GetParams(ctx),
		HostChains:          k.GetAllHostChains(ctx),
		Deposits:            k.GetAllDeposits(ctx),
		Unbondings:          k.FilterUnbondings(ctx, func(u types.Unbonding) bool { return true }),         // GetAll
		UserUnbondings:      k.FilterUserUnbondings(ctx, func(u types.UserUnbonding) bool { return true }), // GetAll
		ValidatorUnbondings: k.FilterValidatorUnbondings(ctx, func(u types.ValidatorUnbonding) bool { return true }),
	}
}
