package keeper

import (
	"bytes"
	"fmt"
	"sort"

	multisig2 "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
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
		valAddr, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orchestratorAccAddress)
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

func HandleEnableModuleProposal(ctx sdk.Context, k Keeper, p *cosmosTypes.EnableModuleProposal) error {
	// check if module already enabled
	if k.GetParams(ctx).ModuleEnabled {
		return cosmosTypes.ErrModuleAlreadyEnabled
	}

	// check if all the validators have orchestrator address set or not
	_, err := k.checkAllValidatorsHaveOrchestrators(ctx)
	if err != nil {
		return err
	}

	// make multisig from the orchestrators keys
	var multisigPubkeys []cryptotypes.PubKey

	// can remove this validation when we allow to have multiple keys with one validator
	// do not iterate over this, will cause non determinism.
	valAddrMap := make(map[string]string)
	for _, orcastratorAddress := range p.OrchestratorAddresses {
		//validate is orchestrator is actually correct
		orchestratorAccAddress, err := sdk.AccAddressFromBech32(orcastratorAddress)
		if err != nil {
			return err
		}
		valAddr, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orchestratorAccAddress)
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
		multisigAcc = &authTypes.BaseAccount{
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

	// set new multisig address as the current address for transaction signing
	k.setCurrentAddress(ctx, multisigAccAddress)

	k.enableModule(ctx)

	return nil
}

func HandleChangeCosmosValidatorWeightsProposal(ctx sdk.Context, k Keeper, p *cosmosTypes.ChangeCosmosValidatorWeightsProposal) error {
	// step 1 : check total sum of weights is 1
	err := cosmosTypes.ValidateValidatorSetCosmosChain(p.WeightedAddresses)
	if err != nil {
		return err
	}
	// step 2 : check if all the validator addresses are correct
	for _, va := range p.WeightedAddresses {
		_, err := cosmosTypes.ValAddressFromBech32(va.Address, cosmosTypes.Bech32PrefixValAddr)
		if err != nil {
			return err
		}
	}
	// step 3 : update the weights and other values
	k.setCosmosValidatorSet(ctx, p.WeightedAddresses)
	return nil
}

func HandleChangeOracleValidatorWeightsProposal(ctx sdk.Context, k Keeper, p *cosmosTypes.ChangeOracleValidatorWeightsProposal) error {
	// step 1 : check total sum of weights is 1
	err := cosmosTypes.ValidateValidatorSetNativeChain(p.WeightedAddresses)
	if err != nil {
		return err
	}
	// step 2 : check if all the validator are correct if already present in kv store
	valAddresses := []sdk.ValAddress{}
	for _, va := range p.WeightedAddresses {
		valAddress, err := sdk.ValAddressFromBech32(va.Address)
		if err != nil {
			return err
		}
		valAddresses = append(valAddresses, valAddress)
	}
	// step 3 : update the weights and other values
	if len(valAddresses) != len(p.WeightedAddresses) {
		return fmt.Errorf("validator addresses and weight are not equally mapped")
	}
	k.setOracleValidatorSet(ctx, valAddresses, p.WeightedAddresses)
	return nil
}

func (k Keeper) handleTransactionQueue(ctx sdk.Context, oldAccount authTypes.AccountI) {
	//step 1 : move all pending and active transactions to an array
	list := k.getAllFromTransactionQueue(ctx) //gets a map of all transactions

	//add grant and revoke transactions in an order to queue
	// for granting access to new multisig account and revoke from previous account

	// grant from old account
	grantTransactionID := k.addGrantTransactions(ctx, oldAccount)

	// feegrant transaction from old account
	feegrantTransactionID := k.addFeegrantTransaction(ctx, oldAccount)

	// revoke transaction from new account
	revokeTransactionID := k.addRevokeTransactions(ctx, oldAccount)

	// set the above transactions in order in transaction queue
	k.setNewInTransactionQueue(ctx, grantTransactionID)
	k.setNewInTransactionQueue(ctx, feegrantTransactionID)
	k.setNewInTransactionQueue(ctx, revokeTransactionID)

	// append all the remaining transactions to the queue with modified txIDs and remove old from the transaction queue
	k.shiftListOfTransactionsToNewIDs(ctx, list)
}
