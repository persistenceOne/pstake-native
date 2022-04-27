package keeper

import (
	"bytes"
	multisig2 "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"sort"
)

func HandleChangeMultisigProposal(ctx sdk.Context, k Keeper, p *types.ChangeMultisigProposal) error {
	//oldAccountAddress := k.getCurrentAddress(ctx)
	//oldAccount := k.getAccountState(ctx, oldAccountAddress)
	//oldPubKey, ok := oldAccount.GetPubKey().(*multisig2.LegacyAminoPubKey)
	//if !ok{return types.ErrInvalidMultisigPubkey}

	var multisigPubkeys []cryptotypes.PubKey

	// can remove this validation when we allow to have multiple keys with one validator
	// do not iterate over this, will cause non determinism.
	valAddrMap := make(map[string]string)
	for _, orcastratorAddress := range p.OrcastratorAddresses {
		//validate is orchestrator is actually correct
		orchestratorAccAddress, err := sdk.AccAddressFromBech32(orcastratorAddress)
		if err != nil {
			return err
		}
		_, valAddr, found, err := k.getAllValidartorOrchestratorMappingAndFindIfExist(ctx, orchestratorAccAddress)
		if err != nil {
			return err
		}
		if !found {
			return types.ErrOrchAddressNotFound
		}

		// checks for singular orchestrator for a validator.
		if _, ok := valAddrMap[valAddr.String()]; ok {
			return types.ErrMoreMultisigAccountsBelongToOneValidator
		}
		valAddrMap[valAddr.String()] = orcastratorAddress

		account := k.authKeeper.GetAccount(ctx, orchestratorAccAddress)
		if account == nil {
			return types.ErrOrchAddressNotFound
		}
		if _, ok := account.GetPubKey().(multisig.PubKey); ok {
			return types.ErrOrcastratorPubkeyIsMultisig
		}
		multisigPubkeys = append(multisigPubkeys, account.GetPubKey())
	}

	// sorts pubkey so that unique key is formed with same pubkeys
	sort.Slice(multisigPubkeys, func(i, j int) bool {
		return bytes.Compare(multisigPubkeys[i].Address(), multisigPubkeys[j].Address()) < 0
	})
	multisigPubkey := multisig2.NewLegacyAminoPubKey(int(p.Threshold), multisigPubkeys)
	multisigAccAddress := sdk.AccAddress(multisigPubkey.Address().Bytes())
	multisigAcc := k.getAccountState(ctx, multisigAccAddress)
	if multisigAcc == nil {
		//TODO add caching for this address string.
		cosmosAddr, err := types.Bech32ifyAddressBytes(types.Bech32Prefix, multisigAccAddress)
		if err != nil {
			return err
		}
		if cosmosAddr == "" {
			return types.ErrInvalidCustodialAddress
		}
		multisigAcc := &authTypes.BaseAccount{
			Address:       cosmosAddr,
			PubKey:        nil,
			AccountNumber: p.AccountNumber,
			Sequence:      0,
		}
		err = multisigAcc.SetPubKey(multisigPubkey)
		if err != nil {
			return err
		}
		k.setAccountState(ctx, multisigAcc)
	}

	k.setCurrentAddress(ctx, multisigAccAddress)

	//TODO add txns to outgoingTx, authzGrant(multipletxns), authzRevoke with old address.
	return nil
}
