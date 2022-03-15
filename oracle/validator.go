package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/persistenceOne/pStake-native/oracle/constants"
	"github.com/persistenceOne/pStake-native/oracle/utils"
)

func SignMintTx(encCfg params.EncodingConfig, clientCtx client.Context, msg sdk.Msg) ([]byte, error) {
	// Build the factory CLI
	// Create a new TxBuilder.

	txBuilder := encCfg.TxConfig.NewTxBuilder()

	privKey, _ := utils.GetSDKPivKeyAndAddress()
	accSeqs := []uint64{0}

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}

	sig := signing.SignatureV2{PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	signerData := xauthsigning.SignerData{
		ChainID: constants.NativeChainID,
	}
	sigv2, err := tx.SignWithPrivKey(
		encCfg.TxConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, privKey, encCfg.TxConfig, accSeqs[0])
	if err != nil {
		return nil, err
	}

	err = txBuilder.SetSignatures(sigv2)
	if err != nil {
		return nil, err
	}

	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
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
