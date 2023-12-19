package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(genState *types.GenesisState)
		expectedErr string
	}{
		{
			"default is valid",
			func(genState *types.GenesisState) {},
			"",
		},
		{
			"invalid liquid validator address",
			func(genState *types.GenesisState) {
				genState.LiquidValidators = []types.LiquidValidator{
					{
						OperatorAddress: "invalidAddr",
					},
				}
			},
			"invalid liquid validator {invalidAddr}: decoding bech32 failed: string not all lowercase or all uppercase: invalid address",
		},
		{
			"empty liquid validator address",
			func(genState *types.GenesisState) {
				genState.LiquidValidators = []types.LiquidValidator{
					{
						OperatorAddress: "",
					},
				}
			},
			"invalid liquid validator {}: empty address string is not allowed: invalid address",
		},
		{
			"invalid params(UnstakeFeeRate)",
			func(genState *types.GenesisState) {
				genState.Params.UnstakeFeeRate = math.LegacyDec{}
			},
			"unstake fee rate must not be nil",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			genState := types.DefaultGenesisState()
			tc.malleate(genState)
			err := types.ValidateGenesis(*genState)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
