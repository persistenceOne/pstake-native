package utils

import (
	"fmt"
	//"github.com/cosmos/cosmos-sdk/crypto/hd"
	//"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdkcryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pStake-native/oracle/constants"
)

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

func SignMintTx(clientCtx client.Context, msg sdk.Msg) ([]byte, error) {
	// Build the factory CLI
	// Create a new TxBuilder.

	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(400000)

	privKey, _ := GetSDKPivKeyAndAddress()
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
		ChainID:       constants.NativeChainID,
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

	fmt.Println(txBuilder.GetTx(), "HELLO WORLD")
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	//txJSONBytes, err := encCfg.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	//if err != nil {
	//	return "",err
	//}
	//txJSON := string(txJSONBytes)

	return txBytes, nil
}

func GetSDKPivKeyAndAddress() (sdkcryptotypes.PrivKey, sdk.AccAddress) {
	//kb, err := keyring.New("orcTest", "test", "./", nil)
	//if err != nil {
	//	panic(err)
	//}
	//path := hd.CreateHDPath(118, 0, 0).String()
	//info, err := kb.NewAccount("orcTest", constants.Seed, "", path, hd.Secp256k1)

	privKey := secp256k1.GenPrivKeyFromSecret([]byte(constants.Seed))

	pubkey := privKey.PubKey()

	address, err := sdk.AccAddressFromHex(pubkey.Address().String())
	fmt.Println(address.String())
	//becaddress, err := bech32.ConvertAndEncode("cosmos", pubkey.Bytes())
	//fmt.Println(becaddress)
	if err != nil {
		panic(err)
	}
	return privKey, address
}

func TssSignMintTx(clientCtx client.Context, msg sdk.Msg, signature []byte) ([]byte, error) {
	// Build the factory CLI
	// Create a new TxBuilder.

	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(400000)

	privKey, _ := GetSDKPivKeyAndAddress()
	accSeqs := []uint64{0}

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}

	sig := signing.SignatureV2{PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  clientCtx.TxConfig.SignModeHandler().DefaultMode(),
			Signature: signature,
		},
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	ac, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, msg.GetSigners()[0])
	fmt.Println(ac, seq, err)
	signerData := xauthsigning.SignerData{
		ChainID:       constants.NativeChainID,
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

	fmt.Println(txBuilder.GetTx(), "HELLO WORLD")
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	//txJSONBytes, err := encCfg.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	//if err != nil {
	//	return "",err
	//}
	//txJSON := string(txJSONBytes)

	return txBytes, nil
}