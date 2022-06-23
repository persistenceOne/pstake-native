package cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	if !k.GetParams(ctx).ModuleEnabled {
		return
	}
	logger := k.Logger(ctx)
	logger.Info("C value : ", k.GetCValue(ctx))
	logger.Info("Minted Amount : ", k.GetMintedAmount(ctx))
	logger.Info("Staked Amount : ", k.GetStakedAmount(ctx))
	logger.Info("vStaked Amount : ", k.GetVirtuallyStakedAmount(ctx))
	k.ProcessAllMintingStoreValue(ctx)
	k.ProcessProposals(ctx)
	k.ProcessRewards(ctx)
	k.ProcessAllTxAndDetails(ctx)
	k.ProcessAllUndelegateSuccess(ctx)
	k.ProcessAllSignature(ctx)
	k.ProcessAllSlashingEvents(ctx)
}
