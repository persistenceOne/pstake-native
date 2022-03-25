package oracle

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	tssSign "github.com/binance-chain/tss-lib/ecdsa/signing"
	"github.com/binance-chain/tss-lib/test"
	"github.com/binance-chain/tss-lib/tss"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	tx2 "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/btcd/btcec"
	"runtime"
	"sync/atomic"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sign "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto"
	"math/big"
	"os"
	"testing"
	"time"
)

const (
	Bech32MainPrefix = "persistence"

	CoinType           = 750
	FullFundraiserPath = "44'/750'/0'/0/0"

	Bech32PrefixAccAddr  = Bech32MainPrefix
	Bech32PrefixAccPub   = Bech32MainPrefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

// hashToInt converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
// This is borrowed from crypto/ecdsa.
func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	fmt.Println(orderBits, orderBytes, "orderbytes")
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

// Serialize signature to R || S.
// R, S are padded to 32 bytes respectively.
func serializeSig(R, S []byte) []byte {
	rBytes := R
	sBytes := S
	sigBytes := make([]byte, 64)
	// 0 pad the byte arrays from the left if they aren't big enough.
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}

func TestTss(t *testing.T) {
	testParticipants := test.TestParticipants

	keys, signPIDs, err := keygen.LoadKeygenTestFixturesRandomSet(11, testParticipants)

	pkX, pkY := keys[0].ECDSAPub.X(), keys[0].ECDSAPub.Y()
	pubKey := ecdsa.PublicKey{
		Curve: tss.EC(),
		X:     pkX,
		Y:     pkY,
	}
	pubkeyObject := (*btcec.PublicKey)(&pubKey)
	pk := pubkeyObject.SerializeCompressed()
	publicKey := &secp256k1.PubKey{Key: pk}

	//privKey := secp256k1.GenPrivKeyFromSecret([]byte("Hello"))
	//publicKey := privKey.PubKey()

	configuration := sdk.GetConfig()
	configuration.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	configuration.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	configuration.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	configuration.SetCoinType(CoinType)
	fromAddr := sdk.AccAddress(publicKey.Address())
	addrStr := fromAddr.String()
	fmt.Println("address!!!!", addrStr)
	msgSend := &banktypes.MsgSend{
		FromAddress: addrStr,
		ToAddress:   addrStr,
		Amount:      sdk.NewCoins(sdk.NewCoin("uxprt", sdk.NewInt(10))),
	}
	encodingConfig := simapp.MakeTestEncodingConfig()
	client, err := newRPCClient("https://rpc.testnet.persistence.one:443", 10*time.Second)
	if err != nil {
		panic(err)
	}
	var ctx = cosmosClient.Context{
		From: addrStr,
	}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(simapp.DefaultNodeHome).
		WithViper("").
		WithNodeURI("https://rpc.testnet.persistence.one:443").
		WithClient(client).WithFromAddress(fromAddr).WithChainID("test-core-1")
	cfg := simapp.MakeTestEncodingConfig()
	accnum, seq, err := ctx.AccountRetriever.GetAccountNumberSequence(ctx, fromAddr)

	txf := tx2.Factory{}.
		WithChainID("test-core-1").
		WithTxConfig(cfg.TxConfig).
		WithSignMode(cfg.TxConfig.SignModeHandler().DefaultMode()).
		WithGas(200000).WithAccountNumber(accnum).WithSequence(seq)

	if txf.SimulateAndExecute() {
		if ctx.Offline {
			panic(err)
		}

		_, adjusted, err := tx2.CalculateGas(ctx, txf, msgSend)
		if err != nil {
			panic(err)
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", tx2.GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	tx, err := tx2.BuildUnsignedTx(txf, sdk.Msg(msgSend))
	if err != nil {
		panic(err)
	}

	signMode := txf.SignMode()
	if signMode == sign.SignMode_SIGN_MODE_UNSPECIFIED {
		signMode = ctx.TxConfig.SignModeHandler().DefaultMode()
	}

	signerData := signing.SignerData{
		ChainID:       txf.ChainID(),
		AccountNumber: txf.AccountNumber(),
		Sequence:      txf.Sequence(),
	}
	fmt.Println(txf.AccountNumber(), txf.Sequence(), txf.ChainID())
	sigData := sign.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := sign.SignatureV2{
		PubKey:   publicKey,
		Data:     &sigData,
		Sequence: txf.Sequence(),
	}
	if err := tx.SetSignatures(sig); err != nil {
		panic(err)
	}
	bz, err := ctx.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, tx.GetTx())

	//bz, err := ctx.TxConfig.TxJSONEncoder()(tx.GetTx())
	fmt.Println(new(big.Int).SetBytes(crypto.Sha256(bz)), "bytes to int!!!!")
	if err != nil {
		panic(err)
	}

	//sigR, _ := privKey.Sign(bz)
	p2pCtx := tss.NewPeerContext(signPIDs)
	parties := make([]*tssSign.LocalParty, 0, len(signPIDs))
	updater := test.SharedPartyUpdater

	errCh := make(chan *tss.Error, len(signPIDs))
	outCh := make(chan tss.Message, len(signPIDs))
	endCh := make(chan common.SignatureData, len(signPIDs))
	var b common.SignatureData
	for i := 0; i < len(signPIDs); i++ {
		params := tss.NewParameters(tss.S256(), p2pCtx, signPIDs[i], len(signPIDs), 10)
		fmt.Println("params!!!", params.PartyCount())
		// TODO to figure new(big.Int).SetBytes(bz)
		bigInt := hashToInt(crypto.Sha256(bz), pubKey.Curve)
		fmt.Println(bigInt, new(big.Int).SetBytes(crypto.Sha256(bz)))
		P := tssSign.NewLocalParty(bigInt, params, keys[i], outCh, endCh).(*tssSign.LocalParty)
		parties = append(parties, P)
		go func(P *tssSign.LocalParty) {
			if err := P.Start(); err != nil {
				errCh <- err
			}
		}(P)
	}
	fmt.Println("signparty length!!!", len(parties))
	var ended int32
signing:
	for {
		fmt.Printf("ACTIVE GOROUTINES: %d\n", runtime.NumGoroutine())
		select {
		case err := <-errCh:
			common.Logger.Errorf("Error: %s", err)
			assert.FailNow(t, err.Error())
			break signing

		case msg := <-outCh:
			dest := msg.GetTo()
			if dest == nil {
				for _, P := range parties {
					if P.PartyID().Index == msg.GetFrom().Index {
						continue
					}
					go updater(P, msg, errCh)
				}
			} else {
				if dest[0].Index == msg.GetFrom().Index {
					t.Fatalf("party %d tried to send a message to itself (%d)", dest[0].Index, msg.GetFrom().Index)
				}
				go updater(parties[dest[0].Index], msg, errCh)
			}

		case b = <-endCh:
			atomic.AddInt32(&ended, 1)
			if atomic.LoadInt32(&ended) == int32(len(signPIDs)) {
				t.Logf("Done. Received signature data from %d participants", ended)
				break signing
			}
		}
	}

	sig = sign.SignatureV2{
		PubKey: publicKey,
		Data: &sign.SingleSignatureData{
			SignMode:  ctx.TxConfig.SignModeHandler().DefaultMode(),
			Signature: serializeSig(b.R, b.S),
		},
	}
	// Construct the SignatureV2 struct
	//sigData = sign.SingleSignatureData{
	//	SignMode:  signMode,
	//	Signature: sigR,
	//}
	//sig = sign.SignatureV2{
	//	PubKey:   publicKey,
	//	Data:     &sigData,
	//	Sequence: txf.Sequence(),
	//}
	fmt.Println("signature!!", sig)
	err = tx.SetSignatures(sig)
	txBytes, _ := ctx.TxConfig.TxEncoder()(tx.GetTx())

	// make signature sdk compliant

	// broadcast
	r, err := ctx.BroadcastTx(txBytes)
	fmt.Println(r, err)

}
