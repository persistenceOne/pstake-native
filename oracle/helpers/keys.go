package helpers

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/relayer/relayer"
	relayer2 "github.com/persistenceOne/pStake-native/oracle/oracle"
)

// KeyOutput contains mnemonic and address of key
type KeyOutput struct {
	Mnemonic string `json:"mnemonic" yaml:"mnemonic"`
	Address  string `json:"address" yaml:"address"`
}

// KeyAddOrRestore is a helper function for add key and restores key when mnemonic is passed
func KeyAddOrRestore(chain relayer2.CosmosChain, keyName string, coinType uint32, mnemonic ...string) (KeyOutput, error) {
	var mnemonicStr string
	var err error

	if len(mnemonic) > 0 {
		mnemonicStr = mnemonic[0]
	} else {
		mnemonicStr, err = relayer.CreateMnemonic()
		if err != nil {
			return KeyOutput{}, err
		}
	}

	info, err := chain.KeyBase.NewAccount(keyName, mnemonicStr, "", hd.CreateHDPath(coinType, 0, 0).String(), hd.Secp256k1)
	if err != nil {
		return KeyOutput{}, err
	}

	done := chain.UseSDKContext()
	ko := KeyOutput{Mnemonic: mnemonicStr, Address: info.GetAddress().String()}
	done()

	return ko, nil
}

func KeyAddOrRestoreNative(chain relayer2.NativeChain, keyName string, coinType uint32, mnemonic ...string) (KeyOutput, error) {
	var mnemonicStr string
	var err error

	if len(mnemonic) > 0 {
		mnemonicStr = mnemonic[0]
	} else {
		mnemonicStr, err = relayer.CreateMnemonic()
		if err != nil {
			return KeyOutput{}, err
		}
	}

	info, err := chain.KeyBase.NewAccount(keyName, mnemonicStr, "", hd.CreateHDPath(coinType, 0, 0).String(), hd.Secp256k1)
	if err != nil {
		return KeyOutput{}, err
	}

	done := chain.UseSDKContext()
	ko := KeyOutput{Mnemonic: mnemonicStr, Address: info.GetAddress().String()}
	done()

	return ko, nil
}
