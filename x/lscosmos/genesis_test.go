package lscosmos_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/app/helpers"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # genesis/test/state
	}

	_, pStakeApp, ctx := helpers.CreateTestApp(t)
	k := pStakeApp.LSCosmosKeeper
	lscosmos.InitGenesis(ctx, k, genesisState)
	got := lscosmos.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	// this line is used by starport scaffolding # genesis/test/assert
}
