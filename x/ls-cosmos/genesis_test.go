package ls_cosmos_test

import (
	"testing"

	keepertest "github.com/persistenceOne/pstake-native/testutil/keeper"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LscosmosKeeper(t)
	ls_cosmos.InitGenesis(ctx, *k, genesisState)
	got := ls_cosmos.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
