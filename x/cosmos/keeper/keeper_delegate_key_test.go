package keeper_test

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test(t *testing.T) {
	_, app, ctx := createTestInput()
	keeper := app.CosmosKeeper

	valAddr1 := sdkTypes.ValAddress("Val1")
	valAddr2 := sdkTypes.ValAddress("Val2")
	valAddrNotExists := sdkTypes.ValAddress("ValNotExists")
	orchAddr1 := sdkTypes.AccAddress("orch1")
	orchAddr2 := sdkTypes.AccAddress("orch2")
	app.StakingKeeper.SetValidator(ctx, stakingTypes.Validator{
		OperatorAddress: valAddr1.String(),
	})
	app.StakingKeeper.SetValidator(ctx, stakingTypes.Validator{
		OperatorAddress: valAddr2.String(),
	})
	err := keeper.SetValidatorOrchestrator(ctx, valAddr1, orchAddr1)
	require.Nil(t, err, "Could not set valAddr1")

	err = keeper.SetValidatorOrchestrator(ctx, valAddr2, orchAddr2)
	require.Nil(t, err, "Could not set valAddr2")

	val1, ok := keeper.GetValidatorOrchestrator(ctx, valAddr1)
	require.Equal(t, valAddr1, val1)
	require.Equal(t, true, ok)

	valNotExists, ok := keeper.GetValidatorOrchestrator(ctx, valAddrNotExists)
	require.Nil(t, valNotExists)
	require.Equal(t, false, ok)

}
