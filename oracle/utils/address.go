package utils

import (
	"fmt"
	"github.com/persistenceOne/pstake-native/oracle/oracle"

	//"github.com/cosmos/cosmos-sdk/crypto/hd"
	//"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdkcryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

func SignCosmosTx(seed string, chain *oracle.CosmosChain, clientCtx client.Context, msg sdk.Msg) ([]byte, error) {
	// Build the factory CLI
	// Create a new TxBuilder.

	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(400000)

	privKey, _ := GetSDKPivKeyAndAddress(seed)
	accSeqs := []uint64{0}

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}

	sig := signing.SignatureV2{PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  clientCtx.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	ac, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, msg.GetSigners()[0])
	fmt.Println(ac, seq, err)
	signerData := xauthsigning.SignerData{
		ChainID:       chain.ChainID,
		AccountNumber: ac,
		Sequence:      seq,
	}
	sigv2, err := tx.SignWithPrivKey(
		clientCtx.TxConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, privKey, clientCtx.TxConfig, accSeqs[0])
	if err != nil {
		return nil, err
	}

	err = txBuilder.SetSignatures(sigv2)
	if err != nil {
		return nil, err
	}

	fmt.Println(txBuilder.GetTx(), "Signed Tx")
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}

func SignNativeTx(seed string, native *oracle.NativeChain, clientCtx client.Context, msg sdk.Msg) ([]byte, error) {
	// Build the factory CLI
	// Create a new TxBuilder.

	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(400000)

	privKey, _ := GetSDKPivKeyAndAddress(seed)
	accSeqs := []uint64{0}

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}

	sig := signing.SignatureV2{PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  clientCtx.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	ac, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, msg.GetSigners()[0])
	fmt.Println(ac, seq, err)
	signerData := xauthsigning.SignerData{
		ChainID:       native.ChainID,
		AccountNumber: ac,
		Sequence:      seq,
	}
	sigv2, err := tx.SignWithPrivKey(
		clientCtx.TxConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, privKey, clientCtx.TxConfig, accSeqs[0])
	if err != nil {
		return nil, err
	}

	err = txBuilder.SetSignatures(sigv2)
	if err != nil {
		return nil, err
	}

	fmt.Println(txBuilder.GetTx(), "Signed Tx")
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil

}

func GetSDKPivKeyAndAddress(Seed string) (sdkcryptotypes.PrivKey, sdk.AccAddress) {

	privKey := secp256k1.GenPrivKeyFromSecret([]byte(Seed))

	pubkey := privKey.PubKey()

	address, err := sdk.AccAddressFromHex(pubkey.Address().String())
	fmt.Println(address.String())
	if err != nil {
		panic(err)
	}
	return privKey, address
}
