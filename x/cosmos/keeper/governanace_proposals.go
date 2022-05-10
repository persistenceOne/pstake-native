package keeper

import (
	"bytes"
	multisig2 "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"sort"
)

func HandleChangeMultisigProposal(ctx sdk.Context, k Keeper, p *cosmosTypes.ChangeMultisigProposal) error {
	oldAccountAddress := k.getCurrentAddress(ctx)
	oldAccount := k.getAccountState(ctx, oldAccountAddress)
	_, ok := oldAccount.GetPubKey().(*multisig2.LegacyAminoPubKey)
	if !ok {
		return cosmosTypes.ErrInvalidMultisigPubkey
	}

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
			return cosmosTypes.ErrOrchAddressNotFound
		}

		// checks for singular orchestrator for a validator.
		if _, ok := valAddrMap[valAddr.String()]; ok {
			return cosmosTypes.ErrMoreMultisigAccountsBelongToOneValidator
		}
		valAddrMap[valAddr.String()] = orcastratorAddress

		account := k.authKeeper.GetAccount(ctx, orchestratorAccAddress)
		if account == nil {
			return cosmosTypes.ErrOrchAddressNotFound
		}
		if _, ok := account.GetPubKey().(multisig.PubKey); ok {
			return cosmosTypes.ErrOrcastratorPubkeyIsMultisig
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
		cosmosAddr, err := cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32Prefix, multisigAccAddress)
		if err != nil {
			return err
		}
		if cosmosAddr == "" {
			return cosmosTypes.ErrInvalidCustodialAddress
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

	k.handleTransactionQueue(ctx, oldAccount)
	return nil
}

func (k Keeper) handleTransactionQueue(ctx sdk.Context, oldAccount authTypes.AccountI) {
	var list []TransactionQueue
	//step 1 : move all pending and active transactions to an array
	list = k.getAllFromTransactionQueue(ctx) //gets a map of all transactions

	//step 2 : add grant and revoke transactions in an order to queue

	// grant from old account
	grantTransactionID := k.addGrantTransactions(ctx, oldAccount)
	// feegrant transaction from old account
	feegrantTransactionID := k.addFeegrantTransaction(ctx, oldAccount)
	// revoke transaction from new account
	revokeTransactionID := k.addRevokeTransactions(ctx, oldAccount)
	k.setNewInTransactionQueue(ctx, grantTransactionID)
	k.setNewInTransactionQueue(ctx, feegrantTransactionID)
	k.setNewInTransactionQueue(ctx, revokeTransactionID)

	//step 3 : append all the remaining transactions to the queue with modified txIDs and remove old from the transaction queue
	k.shiftListOfTransactionsToNewIDs(ctx, list)
	//step 4 : do transaction processing as normal transactions
}
