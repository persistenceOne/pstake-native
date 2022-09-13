package lscosmos_test

import (
	"testing"

	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/lscosmos"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # genesis/test/state
	}

	_, pStakeApp, ctx := helpers.CreateTestApp()
	k := pStakeApp.LSCosmosKeeper
	lscosmos.InitGenesis(ctx, k, genesisState)
	got := lscosmos.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	// this line is used by starport scaffolding # genesis/test/assert
}
