package keeper

import (
	"time"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// adds grant transactions to the outgoing pool and returns txID for reference
func (k Keeper) addGrantTransactions(ctx sdk.Context, oldAccount authTypes.AccountI) uint64 {
	// generate ID for Grant Transaction
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

	custodialAddress, err := cosmosTypes.AccAddressFromBech32(k.GetParams(ctx).CustodialAddress, cosmosTypes.Bech32Prefix)
	if err != nil {
		panic(err)
	}

	var accessMsgsAny []*codecTypes.Any

	//create all generic authorizations
	grantGenericAuthorization := authz.NewGenericAuthorization("/cosmos.authz.v1beta1.MsgGrant")
	revokeAuthGenericAuthorization := authz.NewGenericAuthorization("/cosmos.authz.v1beta1.MsgRevoke")
	delegateGenericAuthorization := authz.NewGenericAuthorization("/cosmos.staking.v1beta1.MsgDelegate")
	undelegateGenericAuthorization := authz.NewGenericAuthorization("/cosmos.staking.v1beta1.MsgUndelegate")
	feegrantGenericAuthorization := authz.NewGenericAuthorization("/cosmos.feegrant.v1beta1.BasicAllowance")
	sendGenericAuthorization := authz.NewGenericAuthorization("/cosmos.bank.v1beta1.MsgSend")
	voteGenereicAuthorization := authz.NewGenericAuthorization("/cosmos.gov.v1beta1.MsgVoteWeighted")

	// generate a grant msg of type any for granting "grant" to new account
	grantGrantMsgAny := k.generateGrantMsgAny(ctx, custodialAddress, grantGenericAuthorization)

	// generate a grant msg of type any for granting "revoke" to new account
	revokeGrantMsgAny := k.generateGrantMsgAny(ctx, custodialAddress, revokeAuthGenericAuthorization)

	// generate a grant msg of type any for granting "revoke" to new account
	delegateGrantMsgAny := k.generateGrantMsgAny(ctx, custodialAddress, delegateGenericAuthorization)

	// generate a grant msg of type any for granting "undelegate" to new account
	undelegateGrantMsgAny := k.generateGrantMsgAny(ctx, custodialAddress, undelegateGenericAuthorization)

	// generate a grant msg of type any for granting "feegrant" to new account
	feegrantGrantMsgAny := k.generateGrantMsgAny(ctx, custodialAddress, feegrantGenericAuthorization)

	// generate a grant msg of type any for granting "send" to new account
	sendGrantMsgAny := k.generateGrantMsgAny(ctx, custodialAddress, sendGenericAuthorization)

	// generate a grant msg of type any for granting "send" to new account
	voteGrantMsgAny := k.generateGrantMsgAny(ctx, custodialAddress, voteGenereicAuthorization)

	// append all messages to the Any type array
	accessMsgsAny = append(
		accessMsgsAny,
		grantGrantMsgAny,
		revokeGrantMsgAny,
		delegateGrantMsgAny,
		undelegateGrantMsgAny,
		feegrantGrantMsgAny,
		sendGrantMsgAny,
		voteGrantMsgAny,
	)

	// add the above grant messages to the execMsg and generate execMsg
	execMsg := authz.MsgExec{
		Grantee: oldAccount.GetAddress().String(),
		Msgs:    accessMsgsAny,
	}

	// convert execMsg to Any type to embed in transaction
	execMsgAny, err := codecTypes.NewAnyWithValue(&execMsg)
	if err != nil {
		panic(err)
	}

	tx := cosmosTypes.CosmosTx{
		Tx: sdkTx.Tx{
			Body: &sdkTx.TxBody{
				Messages:      []*codecTypes.Any{execMsgAny},
				Memo:          "",
				TimeoutHeight: 0,
			},
			AuthInfo: &sdkTx.AuthInfo{
				SignerInfos: nil,
				Fee: &sdkTx.Fee{
					Amount:   nil,
					GasLimit: 400000,
					Payer:    "",
				},
			},
			Signatures: nil,
		},
		EventEmitted:      false,
		Status:            "",
		TxHash:            "",
		ActiveBlockHeight: ctx.BlockHeight() + cosmosTypes.StorageWindow,
		SignerAddress:     oldAccount.GetAddress().String(),
	}

	// sets transaction in outgoing pool with the given tx ID
	k.SetNewTxnInOutgoingPool(ctx, nextID, tx)

	return nextID
}

// generates grant messages with given authorization and then convert to type Any
func (k Keeper) generateGrantMsgAny(ctx sdk.Context, custodialAddress sdk.AccAddress, authorization authz.Authorization) *codecTypes.Any {
	// generate a grant msg of type any for granting given authorization to new account
	grantMsgAny, err := authz.NewMsgGrant(
		custodialAddress,
		k.GetAccountState(ctx, k.GetCurrentAddress(ctx)).GetAddress(), // todo : fix acc address issue as it converts to persistence address
		authorization,
		time.Unix(0, 0),
	)
	if err != nil {
		panic(err)
	}

	// generate any type for grant message
	any, err := codecTypes.NewAnyWithValue(grantMsgAny)
	if err != nil {
		panic(err)
	}

	return any
}

// adds feegrant transaction to the outgoing pool and returns txID for reference
func (k Keeper) addFeegrantTransaction(ctx sdk.Context, oldAccount authTypes.AccountI) uint64 {
	// generate ID for Feegrant Transaction
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

	var feegrantMsgsAny []*codecTypes.Any

	// generate a BasicAllowance message with no limit on spend
	basic := feegrant.BasicAllowance{
		SpendLimit: nil,
	}

	grant := &basic

	custodialAddress, err := cosmosTypes.AccAddressFromBech32(k.GetParams(ctx).CustodialAddress, cosmosTypes.Bech32Prefix)
	if err != nil {
		panic(err)
	}

	// generate a MsgGrantAllowance with the given grant and allowance
	feegrantMsg, err := feegrant.NewMsgGrantAllowance(
		grant,
		custodialAddress,
		k.GetAccountState(ctx, k.GetCurrentAddress(ctx)).GetAddress(), // todo : fix acc address issue as it converts to persistence address
	)
	if err != nil {
		panic(err)
	}

	// generate Any type value from the above generated MsgGrantAllowance to set in MsgExec transaction
	feegrantMsgAny, err := codecTypes.NewAnyWithValue(feegrantMsg)
	if err != nil {
		panic(err)
	}

	feegrantMsgsAny = append(feegrantMsgsAny, feegrantMsgAny)

	// generate a MsgExec to be sent out in transaction
	execMsg := authz.MsgExec{
		Grantee: oldAccount.GetAddress().String(),
		Msgs:    feegrantMsgsAny,
	}

	// convert execMsg to Any type to set in outgoing transaction
	execMsgAny, err := codecTypes.NewAnyWithValue(&execMsg)
	if err != nil {
		panic(err)
	}

	tx := cosmosTypes.CosmosTx{
		Tx: sdkTx.Tx{
			Body: &sdkTx.TxBody{
				Messages:      []*codecTypes.Any{execMsgAny},
				Memo:          "",
				TimeoutHeight: 0,
			},
			AuthInfo: &sdkTx.AuthInfo{
				SignerInfos: nil,
				Fee: &sdkTx.Fee{
					Amount:   nil,
					GasLimit: 400000,
					Payer:    "",
				},
			},
			Signatures: nil,
		},
		EventEmitted:      false,
		Status:            "",
		TxHash:            "",
		ActiveBlockHeight: ctx.BlockHeight() + cosmosTypes.StorageWindow,
		SignerAddress:     oldAccount.GetAddress().String(),
	}

	// set new transaction in outgoing pool with the given tx ID
	k.SetNewTxnInOutgoingPool(ctx, nextID, tx)

	return nextID
}

// todo check logic for revoke as oldAccount is not involved
// adds revoke transaction to outgoing pool and returns txID for reference
func (k Keeper) addRevokeTransactions(ctx sdk.Context, _ authTypes.AccountI) uint64 {
	// generate ID for Revoke Transaction
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

	custodialAddress, err := cosmosTypes.AccAddressFromBech32(k.GetParams(ctx).CustodialAddress, cosmosTypes.Bech32Prefix)
	if err != nil {
		panic(err)
	}

	var revokeMsgsAny []*codecTypes.Any

	// declare all revoke message authTypes
	grantAuthorized := "/cosmos.authz.v1beta1.MsgGrant"
	revokeAuthorized := "/cosmos.authz.v1beta1.MsgRevoke"
	delegateAuthorized := "/cosmos.staking.v1beta1.MsgDelegate"
	undelegateAuthorized := "/cosmos.staking.v1beta1.MsgUndelegate"
	feegrantAuthorized := "/cosmos.feegrant.v1beta1.BasicAllowance"
	sendAuthorized := "/cosmos.bank.v1beta1.MsgSend"
	voteAuthorized := "/cosmos.gov.v1beta1.MsgVoteWeighted"

	// generate all revoke messages
	revokeGrantAuthorizationAny := k.generateRevokeMsgAny(ctx, custodialAddress, grantAuthorized)
	revokeRevokeAuthorizationAny := k.generateRevokeMsgAny(ctx, custodialAddress, revokeAuthorized)
	revokeDelegateAuthorizationAny := k.generateRevokeMsgAny(ctx, custodialAddress, delegateAuthorized)
	revokeUndelegateAuthorizationAny := k.generateRevokeMsgAny(ctx, custodialAddress, undelegateAuthorized)
	revokeFeegrantAuthorizationAny := k.generateRevokeMsgAny(ctx, custodialAddress, feegrantAuthorized)
	revokeSendAuthorizationAny := k.generateRevokeMsgAny(ctx, custodialAddress, sendAuthorized)
	revokeVoteAuthorizationAny := k.generateRevokeMsgAny(ctx, custodialAddress, voteAuthorized)

	revokeMsgsAny = append(
		revokeMsgsAny,
		revokeGrantAuthorizationAny,
		revokeRevokeAuthorizationAny,
		revokeDelegateAuthorizationAny,
		revokeUndelegateAuthorizationAny,
		revokeFeegrantAuthorizationAny,
		revokeSendAuthorizationAny,
		revokeVoteAuthorizationAny,
	)

	cosmosAddrr, err := cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32PrefixAccAddr, k.GetCurrentAddress(ctx))
	if err != nil {
		panic(err)
	}
	execMsg := authz.MsgExec{
		Grantee: cosmosAddrr,
		Msgs:    revokeMsgsAny,
	}

	// convert msgExec to type Any to be set in outgoing transaction
	execMsgAny, err := codecTypes.NewAnyWithValue(&execMsg)
	if err != nil {
		panic(err)
	}

	tx := cosmosTypes.CosmosTx{
		Tx: sdkTx.Tx{
			Body: &sdkTx.TxBody{
				Messages:      []*codecTypes.Any{execMsgAny},
				Memo:          "",
				TimeoutHeight: 0,
			},
			AuthInfo: &sdkTx.AuthInfo{
				SignerInfos: nil,
				Fee: &sdkTx.Fee{
					Amount:   nil,
					GasLimit: 400000,
					Payer:    "",
				},
			},
			Signatures: nil,
		},
		EventEmitted:      false,
		Status:            "",
		TxHash:            "",
		ActiveBlockHeight: ctx.BlockHeight() + cosmosTypes.StorageWindow,
		SignerAddress:     cosmosAddrr,
	}

	// set new transaction in outgoing pool with the given tx ID
	k.SetNewTxnInOutgoingPool(ctx, nextID, tx)

	return nextID
}

// generates revoke messages with given msgAuthorized and then convert to type Any
func (k Keeper) generateRevokeMsgAny(ctx sdk.Context, custodialAddress sdk.AccAddress, msgAuthorized string) *codecTypes.Any {
	// generate revoke message with given msgAuthorized
	revokeMsg := authz.NewMsgRevoke(
		k.GetAccountState(ctx, k.GetCurrentAddress(ctx)).GetAddress(),
		custodialAddress,
		msgAuthorized,
	)

	// generate Any type value for revoke message
	revokeMsgAny, err := codecTypes.NewAnyWithValue(&revokeMsg)
	if err != nil {
		panic(err)
	}

	return revokeMsgAny
}

// helper function to be used in case of shifting the existing list of transactions before proposal to an ID after
// the authorization change transaction
func (k Keeper) shiftListOfTransactionsToNewIDs(ctx sdk.Context, transactionQueue []TransactionQueue) {
	for _, tq := range transactionQueue {
		// generate ID for regenerating Transaction
		nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

		// fetch tx details from db for the given txID
		txDetails, err := k.GetTxnFromOutgoingPoolByID(ctx, tq.txID)
		if err != nil {
			panic(err)
		}

		cosmosAddrr, err := cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32PrefixAccAddr, k.GetCurrentAddress(ctx))
		if err != nil {
			panic(err)
		}
		// reset few details of the outgoing transaction for accepting new signatures
		txDetails.CosmosTxDetails.Tx.AuthInfo.SignerInfos = nil
		txDetails.CosmosTxDetails.Tx.Signatures = nil
		txDetails.CosmosTxDetails.Status = ""
		txDetails.CosmosTxDetails.TxHash = ""
		txDetails.CosmosTxDetails.EventEmitted = false
		txDetails.CosmosTxDetails.ActiveBlockHeight = ctx.BlockHeight() + cosmosTypes.StorageWindow
		txDetails.CosmosTxDetails.SignerAddress = cosmosAddrr

		// set this transaction in outgoing pool with new ID
		k.SetNewTxnInOutgoingPool(ctx, nextID, txDetails.CosmosTxDetails)
		k.setNewInTransactionQueue(ctx, nextID)

		// remove old transaction from outgoing pool
		k.removeTxnDetailsByID(ctx, tq.txID)
		k.removeFromTransactionQueue(ctx, tq.txID)
	}
}
