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
	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sign "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/tendermint/tendermint/crypto"
	"runtime"
	"sync/atomic"

	"math/big"
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

func getParticipantPartyIDs() tss.UnSortedPartyIDs  {
	return nil
}


func SignTss(clientCtx client.Context, bz []byte, keys []keygen.LocalPartySaveData, signPIDs tss.SortedPartyIDs, errCh chan *tss.Error, outCh chan tss.Message, endCh chan common.SignatureData) (sign.SignatureV2){
	pkX, pkY := keys[0].ECDSAPub.X(), keys[0].ECDSAPub.Y()
	pubKey := ecdsa.PublicKey{
		Curve: tss.EC(),
		X:     pkX,
		Y:     pkY,
	}

	pubkeyObject := (*btcec.PublicKey)(&pubKey)
	pk := pubkeyObject.SerializeCompressed()
	publicKey := &secp256k1.PubKey{Key: pk}
	p2pCtx := tss.NewPeerContext(signPIDs)
	parties := make([]*tssSign.LocalParty, 0, len(signPIDs))
	var b common.SignatureData
	updater := test.SharedPartyUpdater

	for i := 0; i < len(signPIDs); i++ {
		params := tss.NewParameters(tss.S256(), p2pCtx, signPIDs[i], len(signPIDs), 10)
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

	var ended int32

signing:
	for {
		fmt.Printf("ACTIVE GOROUTINES: %d\n", runtime.NumGoroutine())
		select {
		case err := <-errCh:
			common.Logger.Errorf("Error: %s", err)
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
				}
				go updater(parties[dest[0].Index], msg, errCh)
			}

		case b = <-endCh:
			atomic.AddInt32(&ended, 1)
			if atomic.LoadInt32(&ended) == int32(len(signPIDs)) {
				break signing
			}
		}
	}

	sig := sign.SignatureV2{
		PubKey: publicKey,
		Data: &sign.SingleSignatureData{
			SignMode:  clientCtx.TxConfig.SignModeHandler().DefaultMode(),
			Signature: serializeSig(b.R, b.S),
		},
	}

	return sig
}
