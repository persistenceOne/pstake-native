package types

import (
	"errors"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx"
	signing2 "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
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

	return bz, nil
}

// ValAddressFromBech32 creates a ValAddress from a Bech32 string.
func ValAddressFromBech32(address, prefix string) (valAddr sdkTypes.ValAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdkTypes.ValAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := sdkTypes.GetFromBech32(address, prefix)
	if err != nil {
		return nil, err
	}

	err = sdkTypes.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// Bech32ifyAddressBytes returns a bech32 representation of address bytes.
// Returns an empty sting if the byte slice is 0-length. Returns an error if the bech32 conversion
// fails or the prefix is empty.
func Bech32ifyAddressBytes(prefix string, address sdkTypes.AccAddress) (string, error) {
	if address.Empty() {
		return "", nil
	}
	if len(address.Bytes()) == 0 {
		return "", nil
	}
	if len(prefix) == 0 {
		return "", errors.New("prefix cannot be empty")
	}
	return bech32.ConvertAndEncode(prefix, address.Bytes())
}

// Bech32ifyValAddressBytes returns a bech32 representation of valAddress bytes.
// Returns an empty sting if the byte slice is 0-length. Returns an error if the bech32 conversion
// fails or the prefix is empty.
func Bech32ifyValAddressBytes(prefix string, address sdkTypes.ValAddress) (string, error) {
	if address.Empty() {
		return "", nil
	}
	if len(address.Bytes()) == 0 {
		return "", nil
	}
	if len(prefix) == 0 {
		return "", errors.New("prefix cannot be empty")
	}
	return bech32.ConvertAndEncode(prefix, address.Bytes())
}
