package cosmos

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	if !k.GetParams(ctx).ModuleEnabled {
		return
	}
	logger := k.Logger(ctx)
	logger.Info(k.GetCValue(ctx).String())
	logger.Info(k.GetMintedAmount(ctx).String())
	logger.Info(k.GetStakedAmount(ctx).String())
	logger.Info(k.GetVirtuallyStakedAmount(ctx).String())
	logger.Info(cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32PrefixAccAddr, k.GetCurrentAddress(ctx)))
	logger.Info(fmt.Sprintf(strconv.FormatUint(k.GetAccountState(ctx, k.GetCurrentAddress(ctx)).GetSequence(), 10)))
	k.ProcessAllMintingStoreValue(ctx)
	k.ProcessProposals(ctx)
	k.ProcessAllTxAndDetails(ctx)
	k.ProcessAllUndelegateSuccess(ctx)
	k.ProcessAllSignature(ctx)
	k.ProcessAllSlashingEvents(ctx)
}
