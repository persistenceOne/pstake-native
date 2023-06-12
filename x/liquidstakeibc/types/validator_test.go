package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func TestValidator_SharesToTokens(t *testing.T) {
	for _, tc := range []struct {
		name      string
		validator *types.Validator
		shares    sdk.Dec
		result    sdk.Int
	}{
		{
			name: "Success",
			validator: &types.Validator{
				OperatorAddress: "valoper1",
				Status:          stakingtypes.BondStatusBonded,
				TotalAmount:     sdk.NewInt(500),
				DelegatorShares: sdk.NewDec(10000),
			},
			shares: sdk.NewDec(1000),
			result: sdk.NewInt(50),
		},
		{
			name: "ZeroShares",
			validator: &types.Validator{
				OperatorAddress: "valoper1",
				Status:          stakingtypes.BondStatusBonded,
				TotalAmount:     sdk.NewInt(0),
				DelegatorShares: sdk.NewDec(0),
			},
			shares: sdk.NewDec(0),
			result: sdk.NewInt(0),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.result, tc.validator.SharesToTokens(tc.shares))
		})
	}
}
