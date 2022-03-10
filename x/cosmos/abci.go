package cosmos

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	// TODO implement EndBlocker
	// TODO : Add MintTokens on Majority
	minting(ctx, k)
}

func minting(ctx sdk.Context, k Keeper) {
	var list []types.KeyAndValueForMinting
	listNew, err := k.GetAllMintAddressAndAmount(ctx, list)
	if err != nil {
		panic("error in fetching address and amount list")
	}
	listWithRatio := k.FetchFromMintPoolTx(ctx, listNew)
	fmt.Println("-----------", listWithRatio, "-----------")

	for _, addressToMintTokens := range listWithRatio {
		if addressToMintTokens.Ratio > types.MinimumRatioForMajority && !addressToMintTokens.Value.Minted {
			err = k.MintTokensOnMajority(ctx, addressToMintTokens.Key, addressToMintTokens.Value)
			if err != nil {
				panic("unable to mint tokens")
			}
		}

		if addressToMintTokens.Value.NativeBlockHeight+types.StorageWindow < ctx.BlockHeight() {
			k.DeleteMintedAddressAndAmountKeys(ctx, addressToMintTokens.Key)
			destinationAddress, err := sdk.AccAddressFromBech32(addressToMintTokens.Value.DestinationAddress)
			if err != nil {
				panic("error in converting address to AccAddress")
			}
			k.DeleteFromMintPoolTx(ctx, destinationAddress, addressToMintTokens.Value.Amount, addressToMintTokens.Key.TxHash)
		}

		//TODO Delete txn once Acknowledgment is received that the amount is delegated successfully
	}
}
