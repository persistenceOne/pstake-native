package oracle

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"github.com/binance-chain/tss-lib/common"
	_"github.com/binance-chain/tss-lib/crypto"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/binance-chain/tss-lib/crypto/vss"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	tssSign "github.com/binance-chain/tss-lib/ecdsa/signing"
	"github.com/binance-chain/tss-lib/test"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sign "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/persistenceOne/pStake-native/oracle/constants"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"

	"math/big"
)

const (
	testFixtureDirFormat  = "%s/../../test/_ecdsa_fixtures"
	testFixtureFileFormat = "keygen_data_%d.json"
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

func KeyGenTss(id, moniker string, uniqueKey  *big.Int, pIDs tss.SortedPartyIDs, outCh chan tss.Message, errCh chan *tss.Error, endCh chan keygen.LocalPartySaveData){
	p2pCtx := tss.NewPeerContext(pIDs)
	thisParty := tss.NewPartyID(id, moniker, uniqueKey)

	parties := make([]*keygen.LocalParty, 0, len(pIDs))
	params := tss.NewParameters(tss.S256(), p2pCtx, thisParty, len(parties), constants.Threshold)
	updater := test.SharedPartyUpdater
	party := keygen.NewLocalParty(params, outCh, endCh)
	go func() {
		err := party.Start()
		if err != nil {
			errCh <- err
		}
	}()
	var ended int32
keygen:
	for {
		select {
		case err := <-errCh:
			common.Logger.Errorf("Error: %s", err)
			break keygen

		case msg := <-outCh:
			dest := msg.GetTo()
			if dest == nil { // broadcast!
				for _, P := range parties {
					if P.PartyID().Index == msg.GetFrom().Index {
						continue
					}
					go updater(P, msg, errCh)
				}
			} else { // point-to-point!
				if dest[0].Index == msg.GetFrom().Index {
					return
				}
				go updater(parties[dest[0].Index], msg, errCh)
			}

		case save := <-endCh:
			index, _ := save.OriginalIndex()
			tryWriteTestFixtureFile(index, save)

			atomic.AddInt32(&ended, 1)
			if atomic.LoadInt32(&ended) == int32(len(pIDs)) {
				// combine shares for each Pj to get u
				u := new(big.Int)
				for range parties {
					pShares := make(vss.Shares, 0)
					for _, P := range parties {
						share := keygen.KGRound2Message1.GetShare()
						shareStruct := &vss.Share{
							Threshold: constants.Threshold,
							ID:        P.PartyID().KeyInt(),
							Share:     new(big.Int).SetBytes(share),
						}
						pShares = append(pShares, shareStruct)
					}
					uj, _ := pShares[:constants.Threshold+1].ReConstruct(tss.S256())

					u = new(big.Int).Add(u, uj)
				}
				break keygen
			}
		}
	}

}

func SignTss(clientCtx client.Context, bz []byte, keys []keygen.LocalPartySaveData, signPIDs tss.SortedPartyIDs, errCh chan *tss.Error, outCh chan tss.Message, endCh chan common.SignatureData) sign.SignatureV2 {
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
		bigInt := hashToInt(tcrypto.Sha256(bz), pubKey.Curve)
		fmt.Println(bigInt, new(big.Int).SetBytes(tcrypto.Sha256(bz)))
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


func tryWriteTestFixtureFile(index int, data keygen.LocalPartySaveData) {
	fixtureFileName := makeTestFixtureFilePath(index)

	fi, err := os.Stat(fixtureFileName)
	if !(err == nil && fi != nil && !fi.IsDir()) {
		fd, _ := os.OpenFile(fixtureFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		bz, _ := json.Marshal(&data)
		_, err = fd.Write(bz)
	}
}

func makeTestFixtureFilePath(partyIndex int) string {
	_, callerFileName, _, _ := runtime.Caller(0)
	srcDirName := filepath.Dir(callerFileName)
	fixtureDirName := fmt.Sprintf(testFixtureDirFormat, srcDirName)
	return fmt.Sprintf("%s/"+testFixtureFileFormat, fixtureDirName, partyIndex)
}