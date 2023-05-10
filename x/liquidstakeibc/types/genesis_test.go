package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesisState(),
			valid:    true,
		},
		{
			desc:     "invalid genesis state, params not set",
			genState: &types.GenesisState{},
			valid:    false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
