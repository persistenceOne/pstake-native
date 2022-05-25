package cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	// TODO implement EndBlocker
	k.ProcessAllMintingStoreValue(ctx)
	k.ProcessProposals(ctx)
	k.ProcessRewards(ctx)
	k.ProcessAllTxAndDetails(ctx)
	k.ProcessAllUndelegateSuccess(ctx)
	k.ProcessAllSignature(ctx)
	return
}
