package types

import sdk "github.com/cosmos/cosmos-sdk/types"

//TODO : Impl functions related to batch
func NewIncomingMintTx(orchAddress sdk.AccAddress, counter uint64) IncomingMintTx {
	return IncomingMintTx{
		OrchAddresses: []string{orchAddress.String()},
		Counter:       counter,
	}
}

func NewAddressAndAmount(destinationAddress sdk.AccAddress, amount sdk.Coins) AddressAndAmount {
	return AddressAndAmount{
		DestinationAddress: destinationAddress.String(),
		Amount:             amount,
	}
}

func NewChainIDHeightAndTxHash(chainID string, blockHeight int64, txHash string) ChainIDHeightAndTxHash {
	return ChainIDHeightAndTxHash{
		ChainID:     chainID,
		BlockHeight: blockHeight,
		TxHash:      txHash,
	}
}
