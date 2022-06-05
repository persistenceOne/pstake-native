package keeper_test

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	valAddr1         = sdkTypes.ValAddress("Val1")
	valAddr2         = sdkTypes.ValAddress("Val2")
	valAddrNotExists = sdkTypes.ValAddress("ValNotExists")
	valAddrInvalid   = sdkTypes.ValAddress(make([]byte, address.MaxAddrLen+1))

	orchAddr1       = sdkTypes.AccAddress("orch1")
	orchAddr2       = sdkTypes.AccAddress("orch2,0")
	orchAddr21      = sdkTypes.AccAddress("orch2,1")
	orchAddrInvalid = sdkTypes.AccAddress(make([]byte, address.MaxAddrLen+1))
)

func SetupValidators(t *testing.T, app app.PstakeApp, ctx sdkTypes.Context) {

	app.StakingKeeper.SetValidator(ctx, stakingTypes.Validator{
		OperatorAddress: valAddr1.String(),
	})
	app.StakingKeeper.SetValidator(ctx, stakingTypes.Validator{
		OperatorAddress: valAddr2.String(),
	})
	err := app.CosmosKeeper.SetValidatorOrchestrator(ctx, valAddr1, orchAddr1)
	require.Nil(t, err, "Could not set valAddr1")

	err = app.CosmosKeeper.SetValidatorOrchestrator(ctx, valAddr2, orchAddr2)
	require.Nil(t, err, "Could not set valAddr2")

}

func TestCheckValidator(t *testing.T) {
	_, pstakeApp, ctx := helpers.CreateTestApp()
	keeper := pstakeApp.CosmosKeeper

	SetupValidators(t, pstakeApp, ctx)

	val1, ok := keeper.CheckValidator(ctx, valAddr1)
	require.Equal(t, valAddr1, val1)
	require.Equal(t, true, ok)

	valNotExists, ok := keeper.CheckValidator(ctx, valAddrNotExists)
	require.Nil(t, valNotExists)
	require.Equal(t, false, ok)

	valInvalid, ok := keeper.CheckValidator(ctx, valAddrInvalid)
	require.Nil(t, valInvalid)
	require.Equal(t, false, ok)

}

func TestSetValidatorOrchestrator(t *testing.T) {
	_, pstakeApp, ctx := helpers.CreateTestApp()
	keeper := pstakeApp.CosmosKeeper
	require.Panics(t, func() { _ = keeper.SetValidatorOrchestrator(ctx, valAddr1, orchAddrInvalid) })
	require.Panics(t, func() { _ = keeper.SetValidatorOrchestrator(ctx, valAddrInvalid, orchAddr1) })
	require.Error(t, keeper.SetValidatorOrchestrator(ctx, valAddrNotExists, orchAddr1))
	SetupValidators(t, pstakeApp, ctx)
	require.Error(t, keeper.SetValidatorOrchestrator(ctx, valAddr1, orchAddr1))
	require.Nil(t, keeper.SetValidatorOrchestrator(ctx, valAddr2, orchAddr21))
}

func TestGetTotalValidatorOrchestratorCount(t *testing.T) {
	_, pstakeApp, ctx := helpers.CreateTestApp()
	keeper := pstakeApp.CosmosKeeper
	count := keeper.GetTotalValidatorOrchestratorCount(ctx)
	require.Equal(t, int64(0), count)

	SetupValidators(t, pstakeApp, ctx)
	count = keeper.GetTotalValidatorOrchestratorCount(ctx)
	require.Equal(t, int64(2), count)
}
