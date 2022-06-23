package oracle

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

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

func SignCosmosTx(seed string, chain *CosmosChain, clientCtx client.Context, msg sdk.Msg) ([]byte, error) {
	// Build the factory CLI
	// Create a new TxBuilder.

	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(400000)

	privKey, _ := GetSDKPivKeyAndAddressR(chain.AccountPrefix, chain.CoinType, seed)
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

func SignNativeTx(seed string, native *NativeChain, clientCtx client.Context, msg sdk.Msg) ([]byte, error) {
	// Build the factory CLI
	// Create a new TxBuilder.

	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(400000)

	privKey, _ := GetSDKPivKeyAndAddressR(native.AccountPrefix, native.CoinType, seed)
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

func GetSDKPivKeyAndAddressR(prefix string, cointype uint32, mnemonic string) (sdkcryptotypes.PrivKey, sdk.AccAddress) {

	kb, err := keyring.New("pstake", keyring.BackendMemory, "", nil)

	keyringAlgos, _ := kb.SupportedAlgorithms()

	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)

	hdPath := hd.CreateHDPath(cointype, 0, 0)

	derivedPriv, err := algo.Derive()(mnemonic, "", hdPath.String())

	privKey := algo.Generate()(derivedPriv)

	addrString, err := sdk.Bech32ifyAddressBytes(prefix, privKey.PubKey().Address())
	if err != nil {
		panic(err)
	}
	return privKey, sdk.AccAddress(addrString)

}

func GetSignature(seed string, chain *CosmosChain, clientCtx client.Context, msg sdk.Msg) ([]byte, error) {
	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(400000)

	privKey, _ := GetSDKPivKeyAndAddressR(chain.AccountPrefix, chain.CoinType, seed)
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

	signature := sigv2.Data.(*signing.SingleSignatureData).Signature

	return signature, nil
}
