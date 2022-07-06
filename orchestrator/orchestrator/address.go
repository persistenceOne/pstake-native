package orchestrator

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	stdlog "log"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdkcryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

func SetSDKConfigPrefix(prefix string) {
	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(prefix, prefix+sdkTypes.PrefixPublic)
	configuration.SetBech32PrefixForValidator(prefix, prefix+sdkTypes.PrefixValidator+sdkTypes.PrefixOperator)
	configuration.SetBech32PrefixForConsensusNode(prefix+sdkTypes.PrefixValidator+sdkTypes.PrefixConsensus, prefix+sdkTypes.PrefixValidator+sdkTypes.PrefixConsensus+sdkTypes.PrefixPublic)

}

func SignNativeTx(seed string, native *NativeChain, clientCtx client.Context, msg sdkTypes.Msg) ([]byte, error) {
	// Build the factory CLI
	// Create a new TxBuilder.

	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(400000)

	privKey, _ := GetPivKeyAddress(native.AccountPrefix, native.CoinType, seed)

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}
	SetSDKConfigPrefix(native.AccountPrefix)
	ac, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, msg.GetSigners()[0])
	if err != nil {
		return nil, err
	}
	stdlog.Println(ac, seq, err)

	sig := signing.SignatureV2{PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  clientCtx.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: seq,
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	signerData := xauthsigning.SignerData{
		ChainID:       native.ChainID,
		AccountNumber: ac,
		Sequence:      seq,
	}
	sigv2, err := tx.SignWithPrivKey(
		clientCtx.TxConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, privKey, clientCtx.TxConfig, seq)
	if err != nil {
		return nil, err
	}

	err = txBuilder.SetSignatures(sigv2)
	if err != nil {
		return nil, err
	}

	stdlog.Println(txBuilder.GetTx(), "Signed Tx")
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil

}

func GetSDKPivKeyAndAddress(Seed string) (sdkcryptotypes.PrivKey, sdkTypes.AccAddress) {

	privKey := secp256k1.GenPrivKeyFromSecret([]byte(Seed))

	pubkey := privKey.PubKey()

	address, err := sdkTypes.AccAddressFromHex(pubkey.Address().String())
	stdlog.Println(address.String())
	if err != nil {
		panic(err)
	}
	return privKey, address
}

func GetPivKeyAddress(prefix string, cointype uint32, mnemonic string) (sdkcryptotypes.PrivKey, string) {

	kb, err := keyring.New("pstake", keyring.BackendMemory, "", nil)

	keyringAlgos, _ := kb.SupportedAlgorithms()

	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)

	hdPath := hd.CreateHDPath(cointype, 0, 0)

	derivedPriv, err := algo.Derive()(mnemonic, "", hdPath.String())

	privKey := algo.Generate()(derivedPriv)

	//addrString, err := sdkTypes.Bech32ifyAddressBytes(prefix, privKey.PubKey().Address())
	addrString, err := Bech32ifyAddressBytes(prefix, sdkTypes.AccAddress(privKey.PubKey().Address()))
	if err != nil {
		panic(err)
	}
	return privKey, addrString

}

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

func GetSignBytesForCosmos(seed string, chain *CosmosChain, clientCtx client.Context, OutgoingTx sdkTx.Tx, signerAddress string) ([]byte, error) {
	privkey, _ := GetPivKeyAddress(chain.AccountPrefix, chain.CoinType, seed)

	SetSDKConfigPrefix(chain.AccountPrefix)
	signerAddr, err := AccAddressFromBech32(signerAddress, chain.AccountPrefix)
	if err != nil {
		return nil, err
	}

	ac, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, signerAddr)

	if err != nil {
		return nil, err
	}

	nativeProtoCodec := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

	err = OutgoingTx.UnpackInterfaces(nativeProtoCodec)
	if err != nil {
		panic(err)
	}

	txBuilder := authTx.WrapTx(&OutgoingTx)

	SignBytes, err := clientCtx.TxConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		xauthsigning.SignerData{
			ChainID:       chain.ChainID,
			AccountNumber: ac,
			Sequence:      seq,
		}, txBuilder.GetTx())

	if err != nil {
		return nil, err
	}

	signature, err := privkey.Sign(SignBytes)
	if err != nil {
		return nil, err
	}

	if !privkey.PubKey().VerifySignature(SignBytes, signature) {
		panic(fmt.Sprintf("signature verification failed err: %v", err))
	}

	return signature, nil

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
