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

	execMsg := authz.MsgExec{
		Grantee: oldAccount.GetAddress().String(),
		Msgs:    accessMsgsAny,
	}

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
					GasLimit: 200000,
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

	k.setNewTxnInOutgoingPool(ctx, nextID, tx)

	return nextID
}

func (k Keeper) generateGrantMsgAny(ctx sdk.Context, custodialAddress sdk.AccAddress, authorization *authz.GenericAuthorization) *codecTypes.Any {
	// generate a grant msg of type any for granting given authorization to new account
	grantMsgAny, err := authz.NewMsgGrant(
		custodialAddress,
		k.getAccountState(ctx, k.getCurrentAddress(ctx)).GetAddress(),
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

func (k Keeper) addFeegrantTransaction(ctx sdk.Context, oldAccount authTypes.AccountI) uint64 {
	// generate ID for Feegrant Transaction
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

	var feegrantMsgsAny []*codecTypes.Any

	basic := feegrant.BasicAllowance{
		SpendLimit: nil,
	}

	var grant feegrant.FeeAllowanceI
	grant = &basic

	custodialAddress, err := cosmosTypes.AccAddressFromBech32(k.GetParams(ctx).CustodialAddress, cosmosTypes.Bech32Prefix)
	if err != nil {
		panic(err)
	}

	feegrantMsg, err := feegrant.NewMsgGrantAllowance(
		grant,
		custodialAddress,
		k.getAccountState(ctx, k.getCurrentAddress(ctx)).GetAddress(),
	)
	if err != nil {
		panic(err)
	}

	feegrantMsgAny, err := codecTypes.NewAnyWithValue(feegrantMsg)
	if err != nil {
		panic(err)
	}

	feegrantMsgsAny = append(feegrantMsgsAny, feegrantMsgAny)

	execMsg := authz.MsgExec{
		Grantee: oldAccount.GetAddress().String(),
		Msgs:    feegrantMsgsAny,
	}

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
					GasLimit: 200000,
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

	k.setNewTxnInOutgoingPool(ctx, nextID, tx)

	return nextID
}

func (k Keeper) addRevokeTransactions(ctx sdk.Context, oldAccount authTypes.AccountI) uint64 {
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

	execMsg := authz.MsgExec{
		Grantee: k.getCurrentAddress(ctx).String(),
		Msgs:    revokeMsgsAny,
	}

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
					GasLimit: 200000,
					Payer:    "",
				},
			},
			Signatures: nil,
		},
		EventEmitted:      false,
		Status:            "",
		TxHash:            "",
		ActiveBlockHeight: ctx.BlockHeight() + cosmosTypes.StorageWindow,
		SignerAddress:     k.getCurrentAddress(ctx).String(),
	}

	k.setNewTxnInOutgoingPool(ctx, nextID, tx)

	return nextID
}

func (k Keeper) generateRevokeMsgAny(ctx sdk.Context, custodialAddress sdk.AccAddress, msgAuthorized string) *codecTypes.Any {
	// generate revoke message with given msgAuthorized
	revokeMsg := authz.NewMsgRevoke(
		k.getAccountState(ctx, k.getCurrentAddress(ctx)).GetAddress(),
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

func (k Keeper) shiftListOfTransactionsToNewIDs(ctx sdk.Context, transactionQueue []TransactionQueue) {
	for _, tq := range transactionQueue {
		// generate ID for regenerating Transaction
		nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

		// fetch tx details from db for the given txID
		txDetails, err := k.getTxnFromOutgoingPoolByID(ctx, tq.txID)
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
		txDetails.CosmosTxDetails.SignerAddress = k.getCurrentAddress(ctx).String()

		// set this transaction in outgoing pool with new ID
		k.setNewTxnInOutgoingPool(ctx, nextID, txDetails.CosmosTxDetails)
		k.setNewInTransactionQueue(ctx, nextID)

		// remove old transaction from outgoing pool
		k.removeTxnDetailsByID(ctx, tq.txID)
		k.removeFromTransactionQueue(ctx, tq.txID)
	}
}
