package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/persistenceOne/pstake-native/v2/x/lspersistence/types"
)

// ParamChanges defines the parameters that can be modified by legacy param change proposals
// on the simulation.
func ParamChanges(r *rand.Rand) []simtypes.LegacyParamChange {
	return []simtypes.LegacyParamChange{

		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyWhitelistedValidators),
			func(r *rand.Rand) string {
				bz, err := json.Marshal(genWhitelistedValidator(r))
				if err != nil {
					panic(err)
				}
				return string(bz)
			},
		),

		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyLiquidBondDenom),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", genLiquidBondDenom(r))
			},
		),

		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyUnstakeFeeRate),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", genUnstakeFeeRate(r).String())
			},
		),

		simulation.NewSimLegacyParamChange(types.ModuleName, string(types.KeyMinLiquidStakingAmount),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", genMinLiquidStakingAmount(r))
			},
		),
	}
}
