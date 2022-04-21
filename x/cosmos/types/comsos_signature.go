package types

import (
	"errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	signing2 "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"strings"
)

// VerifySignature Multisig only supports Amino Signing, hence the code will only check for amino signing
func VerifySignature(pubkey cryptotypes.PubKey, signerData signing.SignerData, sigData signing2.SingleSignatureData, transaction tx.Tx) error {
	aminoSignModeHandler := legacytx.NewStdTxSignModeHandler()

	return signing.VerifySignature(pubkey, signerData, &sigData, aminoSignModeHandler, &transaction)
}

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func AccAddressFromBech32(address, prefix string) (addr sdkTypes.AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdkTypes.AccAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := sdkTypes.GetFromBech32(address, prefix)
	if err != nil {
		return nil, err
	}

	err = sdkTypes.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdkTypes.AccAddress(bz), nil
}
