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
	proposal(ctx, k)
}

func minting(ctx sdk.Context, k Keeper) {
	var list []types.KeyAndValueForMinting
	listNew, err := k.GetAllMintAddressAndAmount(ctx, list)
	if err != nil {
		panic("error in fetching address and amount list")
	}
	listWithRatio := k.FetchFromMintPoolTx(ctx, listNew)

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

		//TODO Delete txn once Acknowledgment is received if the amount is delegated successfully
	}
}

func proposal(ctx sdk.Context, k Keeper) {
	list := k.GetAllKeyAndValueForProposal(ctx)
	fmt.Println("------------", list, "------------")

	for _, element := range list {
		if element.Value.Ratio > 0.66 && !element.Value.ProposalPosted {
			err := k.CreateProposal(ctx, element)
			fmt.Println("Created Proposal")
			if err != nil {
				panic("Error in generating proposal" + err.Error())
			}
		}
	}
}
