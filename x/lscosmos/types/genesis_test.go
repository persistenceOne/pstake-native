package types_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: &types.GenesisState{
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
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

func Test1GenesisState_Validate(t *testing.T) {
	fmt.Println("asdasd")
	fmt.Println(types.DepositAmount{}.Amount == nil)
	fmt.Println(types.DepositAmount{}.Amount.Add(sdk.NewInt64Coin("asd", 1)))
	fmt.Println(types.DepositAmount{}.Amount.IsEqual(sdk.Coins{}))
}
