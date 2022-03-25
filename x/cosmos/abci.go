package cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	// TODO implement EndBlocker
	minting(ctx, k)
	proposal(ctx, k)
	checkTransactions(ctx, k)
}

func minting(ctx sdk.Context, k Keeper) {
	err := k.ProcessAllMintingTransactions(ctx)
	logger := k.Logger(ctx)
	if err != nil {
		logger.Info(err.Error())
	}
}

func proposal(ctx sdk.Context, k Keeper) {
	//TODO Add missing actions
	err := k.ProcessProposals(ctx)
	logger := k.Logger(ctx)
	if err != nil {
		logger.Info(err.Error())
	}
}

// For querying transactions (sent to cosmos side) status and once majority is reached then check if success or failure.
// If failure then the next steps regarding that
func checkTransactions(ctx sdk.Context, k Keeper) {
	err := k.ProcessAllTxAndDetails(ctx)
	logger := k.Logger(ctx)
	if err != nil {
		logger.Info(err.Error())
	}
}
