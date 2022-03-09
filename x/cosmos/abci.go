package cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/keeper"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	// TODO implement EndBlocker
	// TODO : Add MintTokens on Majority
	minting(ctx, k)
}

func minting(ctx sdk.Context, k Keeper) {
	var list []keeper.KeyAndValueForMinting
	listNew, err := k.GetAllMintAddressAndAmount(ctx, list)
	if err != nil {
		panic("error in fetching address and amount list")
	}
	listWithRatio := k.FetchFromMintPoolTx(ctx, listNew)

	for _, addressToMintTokens := range listWithRatio {
		err = k.MintTokensOnMajority(ctx, addressToMintTokens.Key, addressToMintTokens.Value)
		if err != nil {
			panic("unable to mint tokens")
		}
	}

}
