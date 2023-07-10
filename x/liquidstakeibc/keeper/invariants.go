package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// RegisterInvariants registers the bank module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "cvalue-limits", CValueLimits(k))
}

func CValueLimits(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		hostChains := k.GetAllHostChains(ctx)
		str := ""
		broken := false
		for _, hc := range hostChains {
			if !k.CValueWithinLimits(ctx, hc) {
				str = fmt.Sprintf("chainID: %s, cValue: %s \n", hc.ChainId, hc.CValue)
			}
		}
		if str != "" {
			broken = true
		}
		return sdk.FormatInvariant(
			types.ModuleName, "cvalue-limits",
			fmt.Sprintf("cvalue out of bounds as follows \n%s ", str),
		), broken
	}
}
