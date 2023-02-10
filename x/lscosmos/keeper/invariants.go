package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// RegisterInvariants registers the lscosmos module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "c-value-range", CValueRangeInvariant(k))
}

// CValueRangeInvariant checks that if CValue is within module safety range
func CValueRangeInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		if !k.GetModuleState(ctx) {
			return "Module is disabled, cannot check invariant", false
		}

		cValue := k.GetCValue(ctx)
		if !cValue.IsPositive() || cValue.GT(types.MaxCValue) {
			return sdk.FormatInvariant(
				types.ModuleName, "C-value out of range",
				fmt.Sprintf("cValue is expected between %s\n%s, currently is %s", sdk.ZeroDec(), types.MaxCValue, cValue),
			), true
		}
		return sdk.FormatInvariant(
			types.ModuleName, "C-value is in range",
			fmt.Sprintf("cValue is expected between %s\n%s, currently is %s", sdk.ZeroDec(), types.MaxCValue, cValue),
		), false
	}
}
