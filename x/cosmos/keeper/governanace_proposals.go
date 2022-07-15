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

/*
HandleChangeMultisigProposal handles the multisig change proposal
Multisig change proposal is used to change the multsig account on cosmos side with the given oracle addresses in
the proposal and does the following actions :
1. Aggregate all the keys and form a multisig address
2. Handle the transaction queue (by making grant, feegrant and revoke transactions for changing authorizations)
*/
func HandleChangeMultisigProposal(ctx sdk.Context, k Keeper, p *cosmosTypes.ChangeMultisigProposal) error {
	oldAccountAddress := k.GetCurrentAddress(ctx)
	oldAccount := k.GetAccountState(ctx, oldAccountAddress)
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
	multisigAcc := k.GetAccountState(ctx, multisigAccAddress)
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
		k.SetAccountState(ctx, multisigAcc)
	}

	k.SetCurrentAddress(ctx, multisigAccAddress)

	k.handleTransactionQueue(ctx, oldAccount)
	return nil
}

/*
HandleEnableModuleProposal handles the proposal to enable module by setting the moduleEnable flag in params
Initially checks if module is already enabled. After verifying the account addresses supplied with the proposal a new
multisig address is formed and set. Post all these step the module is enabled by calling enableModule function
*/
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
	multisigAcc := k.GetAccountState(ctx, multisigAccAddress)

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
		Sequence:      p.SequenceNumber,
	}
	err = multisigAcc.SetPubKey(multisigPubkey)
	if err != nil {
		return err
	}
	k.SetAccountState(ctx, multisigAcc)

	// set new multisig address as the current address for transaction signing
	k.SetCurrentAddress(ctx, multisigAccAddress)

	custodialAccAddress, err := cosmosTypes.AccAddressFromBech32(p.CustodialAddress, cosmosTypes.Bech32PrefixAccAddr)
	if err != nil {
		return err
	}
	custodialAddressString, err := cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32PrefixAccAddr, custodialAccAddress)
	if err != nil {
		return err
	}

	// set custodial address, chainID and enable the module.
	k.setCustodialAddress(ctx, custodialAddressString)
	k.setCosmosChainID(ctx, p.ChainID)
	k.setBondDenom(ctx, p.Denom)
	k.enableModule(ctx)

	multisigAddress, err := cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32PrefixAccAddr, multisigAccAddress)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			cosmosTypes.EventModuleEnableProposal,
			sdk.NewAttribute(cosmosTypes.AttributeMultisigAddress, multisigAddress),
		),
	)

	return nil
}

/*
HandleChangeCosmosValidatorWeightsProposal handles the proposal for change of weights of cosmos side validator set
by verifying the supplied data and then setting it in DB
*/
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
	k.SetCosmosValidatorSet(ctx, p.WeightedAddresses)
	return nil
}

/*
HandleChangeOracleValidatorWeightsProposal handles the proposal for change of weights of oracle validator set
by verifying the supplied data and then setting it in DB
*/
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
	k.SetOracleValidatorSet(ctx, valAddresses, p.WeightedAddresses)
	return nil
}

// handleTransactionQueue helper function for handling transaction queue in HandleChangeMultisigProposal
func (k Keeper) handleTransactionQueue(ctx sdk.Context, oldAccount authTypes.AccountI) {
	//step 1 : move all pending and active transactions to an array
	list := k.GetAllFromTransactionQueue(ctx) //gets a map of all transactions

	//add grant and revoke transactions in an order to queue
	// for granting access to new multisig account and revoke from previous account

	// grant from old account
	grantTransactionID := k.AddGrantTransactions(ctx, oldAccount)

	// feegrant transaction from old account
	feegrantTransactionID := k.AddFeegrantTransaction(ctx, oldAccount)

	// revoke transaction from new account
	revokeTransactionID := k.AddRevokeTransactions(ctx, oldAccount)

	// set the above transactions in order in transaction queue
	k.SetNewInTransactionQueue(ctx, grantTransactionID)
	k.SetNewInTransactionQueue(ctx, feegrantTransactionID)
	k.SetNewInTransactionQueue(ctx, revokeTransactionID)

	// append all the remaining transactions to the queue with modified txIDs and remove old from the transaction queue
	k.shiftListOfTransactionsToNewIDs(ctx, list)
}
