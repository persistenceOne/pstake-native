package cosmos

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	// TODO implement EndBlocker
	if !k.GetParams(ctx).ModuleEnabled {
		return
	}
	fmt.Println(k.GetCValue(ctx))
	fmt.Println("Minted Amount : ", k.GetMintedAmount(ctx))
	fmt.Println("Staked Amount : ", k.GetStakedAmount(ctx))
	fmt.Println("vStaked Amount : ", k.GetVirtuallyStakedAmount(ctx))
	k.ProcessAllMintingStoreValue(ctx)
	k.ProcessProposals(ctx)
	k.ProcessRewards(ctx)
	k.ProcessAllTxAndDetails(ctx)
	k.ProcessAllUndelegateSuccess(ctx)
	k.ProcessAllSignature(ctx)
}
