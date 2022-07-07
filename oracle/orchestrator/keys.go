package orchestrator

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdkcryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createMemoryKeyFromMnemonic(mnemonic string) (sdkcryptotypes.PrivKey, sdk.AccAddress, error) {
	kb, err := keyring.New("pstake", keyring.BackendMemory, "", nil)
	if err != nil {
		return nil, nil, err
	}

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return nil, nil, err
	}

	derivedPriv, err := algo.Derive()(mnemonic, "", "m/44'/750'/0'/0/0")
	if err != nil {
		return nil, nil, err
	}

	privKey := algo.Generate()(derivedPriv)

	account, err := kb.NewAccount("oraclekey", mnemonic, "", "m/44'/750'/0'/0/0", algo)

	bytes, err := sdk.Bech32ifyAddressBytes("persistence", account.GetAddress())
	if err != nil {
		return nil, nil, err
	}

	return privKey, sdk.AccAddress(bytes), nil
}
