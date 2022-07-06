package keeper_test

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	multisig2 "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/cosmos/keeper"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func GetSDKPivKeyAndAddressR(prefix string, cointype uint32, mnemonic string) (cryptotypes.PrivKey, string) {

	kb, err := keyring.New("pstake", keyring.BackendMemory, "", nil)

	keyringAlgos, _ := kb.SupportedAlgorithms()

	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)

	hdPath := hd.CreateHDPath(cointype, 0, 0)

	derivedPriv, err := algo.Derive()(mnemonic, "", hdPath.String())

	privKey := algo.Generate()(derivedPriv)

	//addrString, err := sdk.Bech32ifyAddressBytes(prefix, privKey.PubKey().Address())
	addrString, err := sdk.Bech32ifyAddressBytes(prefix, privKey.PubKey().Address())
	if err != nil {
		panic(err)
	}
	fmt.Println(addrString)
	return privKey, addrString

}

func multisig(OrcastratorAddresses []string, k keeper.Keeper, t *testing.T, ctx sdk.Context, accountNumber uint64, threshold int64) *authTypes.BaseAccount {
	var multisigPubkeys []cryptotypes.PubKey

	// can remove this validation when we allow to have multiple keys with one validator
	// do not iterate over this, will cause non determinism.
	for _, orcastratorAddress := range OrcastratorAddresses {
		//validate is orchestrator is actually correct
		orchestratorAccAddress, err := sdk.AccAddressFromBech32(orcastratorAddress)
		require.NoError(t, nil, err)

		account := k.authKeeper.GetAccount(ctx, orchestratorAccAddress)
		multisigPubkeys = append(multisigPubkeys, account.GetPubKey())
	}

	// sorts pubkey so that unique key is formed with same pubkeys
	sort.Slice(multisigPubkeys, func(i, j int) bool {
		return bytes.Compare(multisigPubkeys[i].Address(), multisigPubkeys[j].Address()) < 0
	})
	multisigPubkey := multisig2.NewLegacyAminoPubKey(int(threshold), multisigPubkeys)
	multisigAccAddress := sdk.AccAddress(multisigPubkey.Address().Bytes())
	multisigAcc := k.GetAccountState(ctx, multisigAccAddress)
	if multisigAcc == nil {
		//TODO add caching for this address string.
		cosmosAddr, err := cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32Prefix, multisigAccAddress)
		require.NoError(t, nil, err)
		multisigAcc := &authTypes.BaseAccount{
			Address:       cosmosAddr,
			PubKey:        nil,
			AccountNumber: accountNumber,
			Sequence:      0,
		}
		err = multisigAcc.SetPubKey(multisigPubkey)
		k.SetAccountState(ctx, multisigAcc)
		return multisigAcc
	}
	return nil
}

func TestKeeper_SetAccountState(t *testing.T) {
	_, app, ctx := helpers.CreateTestApp()
	cosmosKeeper := app.CosmosKeeper

	orcastratorAddress := "persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu"
	prvKey, err := GetSDKPivKeyAndAddressR("persistence", 118, "together chief must vocal account off apart dinosaur move canvas spring whisper improve cruise idea earn reflect flash goat illegal mistake blood earn ridge")
	require.NoError(t, nil, err)
	acc := &authTypes.BaseAccount{
		Address:       orcastratorAddress,
		PubKey:        nil,
		AccountNumber: 1,
		Sequence:      0,
	}
	acc.SetPubKey(prvKey.PubKey())
	cosmosKeeper.authKeeper.SetAccount(ctx, acc)
	baseAccount := multisig([]string{orcastratorAddress}, cosmosKeeper, t, ctx, 0, 1)

	address, _ := cosmosTypes.AccAddressFromBech32(baseAccount.Address, "cosmos")
	cosmosKeeper.GetAccountState(ctx, address)
}
