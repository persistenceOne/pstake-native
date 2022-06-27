package keeper

import (
	"fmt"
	"reflect"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

//______________________________________________________________________________________________
/*
TODO : Add Key and value structure as comment
*/

type txIDAndDetailsInOutgoingPool struct {
	txID      uint64
	txDetails cosmosTypes.CosmosTx
}

// SetNewTxnInOutgoingPool sets new transaction in outgoing pool with the given transaction details
func (k Keeper) SetNewTxnInOutgoingPool(ctx sdk.Context, txID uint64, tx cosmosTypes.CosmosTx) {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.OutgoingTXPoolKey)
	key := cosmosTypes.UInt64Bytes(txID)
	bz, err := k.cdc.Marshal(&tx)
	if err != nil {
		panic(err)
	}
	outgoingStore.Set(key, bz)
}

// updateStatusOnceProcessed updates the status of the transactions once they are done with processing and successful
func (k Keeper) updateStatusOnceProcessed(ctx sdk.Context, txID uint64, status string) {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.OutgoingTXPoolKey)
	key := cosmosTypes.UInt64Bytes(txID)

	var cosmosTx cosmosTypes.CosmosTx
	k.cdc.MustUnmarshal(outgoingStore.Get(key), &cosmosTx)

	cosmosTx.Status = status
	outgoingStore.Set(key, k.cdc.MustMarshal(&cosmosTx))
}

// sets event emitted to true once even is emitted
func (k Keeper) setEventEmittedFlag(ctx sdk.Context, txID uint64, flag bool) {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.OutgoingTXPoolKey)
	key := cosmosTypes.UInt64Bytes(txID)

	var cosmosTx cosmosTypes.CosmosTx
	k.cdc.MustUnmarshal(outgoingStore.Get(key), &cosmosTx)

	cosmosTx.EventEmitted = flag
	outgoingStore.Set(key, k.cdc.MustMarshal(&cosmosTx))
}

// gets txn details corresponding to the given ID
func (k Keeper) getTxnFromOutgoingPoolByID(ctx sdk.Context, txID uint64) (cosmosTypes.QueryOutgoingTxByIDResponse, error) {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.OutgoingTXPoolKey)
	key := cosmosTypes.UInt64Bytes(txID)
	bz := outgoingStore.Get(key)
	if bz == nil {
		return cosmosTypes.QueryOutgoingTxByIDResponse{}, cosmosTypes.ErrTxnNotPresentInOutgoingPool
	}
	var cosmosTx cosmosTypes.CosmosTx
	k.cdc.MustUnmarshal(bz, &cosmosTx)
	return cosmosTypes.QueryOutgoingTxByIDResponse{
		CosmosTxDetails: cosmosTx,
	}, nil
}

// Deletes txn Details by ID
func (k Keeper) removeTxnDetailsByID(ctx sdk.Context, txID uint64) {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.OutgoingTXPoolKey)
	key := cosmosTypes.UInt64Bytes(txID)
	outgoingStore.Delete(key)
}

func (k Keeper) setOutgoingTxnSignaturesAndEmitEvent(ctx sdk.Context, tx cosmosTypes.CosmosTx, txID uint64) error {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.OutgoingTXPoolKey)
	key := cosmosTypes.UInt64Bytes(txID)

	//calculate and set txHash
	txBytes, err := k.cdc.Marshal(&tx.Tx)
	if err != nil {
		return err
	}
	txHash := cosmosTypes.BytesToHexUpper(txBytes)
	tx.TxHash = txHash

	bz, err := k.cdc.Marshal(&tx)
	if err != nil {
		return err
	}
	outgoingStore.Set(key, bz)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			cosmosTypes.EventTypeSignedOutgoing,
			sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	)

	k.removeFromOutgoingSignaturePool(ctx, txID)
	return nil
}

func (k Keeper) getAllTxInOutgoingPool(ctx sdk.Context) (details []txIDAndDetailsInOutgoingPool, err error) {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.OutgoingTXPoolKey)
	iterator := outgoingStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := cosmosTypes.UInt64FromBytes(iterator.Key())
		var tx cosmosTypes.CosmosTx
		err := k.cdc.Unmarshal(iterator.Value(), &tx)
		if err != nil {
			return nil, err
		}
		details = append(details, txIDAndDetailsInOutgoingPool{
			txID:      key,
			txDetails: tx,
		})
	}
	return details, nil
}

//______________________________________________________________________________________________

type TxHashAndDetails struct {
	TxHash  string
	Details cosmosTypes.TxHashValue
}

// Set details corresponding to a particular txHash and update details if already present
func (k Keeper) setTxHashAndDetails(ctx sdk.Context, msg cosmosTypes.MsgTxStatus, validatorAddress sdk.ValAddress) {
	txHashStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.HashAndIDStore)
	key := []byte(msg.TxHash)
	totalValidatorCount := k.GetTotalValidatorOrchestratorCount(ctx)

	if !txHashStore.Has(key) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newTxHashValue := cosmosTypes.NewTxHashValue(msg, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow, validatorAddress)
		txHashStore.Set(key, k.cdc.MustMarshal(&newTxHashValue))
		return
	}

	var txHashValue cosmosTypes.TxHashValue
	k.cdc.MustUnmarshal(txHashStore.Get(key), &txHashValue)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotTxStatus(txHashValue, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newTxHashValue := cosmosTypes.NewTxHashValue(msg, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow, validatorAddress)
		txHashStore.Set(key, k.cdc.MustMarshal(&newTxHashValue))
		return
	}

	if !txHashValue.Find(validatorAddress.String()) {
		txHashValue.UpdateValues(validatorAddress.String(), totalValidatorCount)
		txHashStore.Set(key, k.cdc.MustMarshal(&txHashValue))
		return
	}
}

//Fetch details mapped to particular hash
func (k Keeper) getTxHashAndDetails(ctx sdk.Context, txHash string) (cosmosTypes.TxHashValue, error) {
	hashAndIDStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.HashAndIDStore)
	key := []byte(txHash)
	if hashAndIDStore.Has(key) {
		var txHashAndValue cosmosTypes.TxHashValue
		err := k.cdc.Unmarshal(hashAndIDStore.Get(key), &txHashAndValue)
		if err != nil {
			return cosmosTypes.TxHashValue{}, err
		}
		return txHashAndValue, nil
	}
	return cosmosTypes.TxHashValue{}, nil
}

// Removes all the details mapped to txHash
func (k Keeper) removeTxHashAndDetails(ctx sdk.Context, txHash string) {
	hashAndIDStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.HashAndIDStore)
	key := []byte(txHash)
	hashAndIDStore.Delete(key)
}

// Fetches the list of all details mapped to txHash
func (k Keeper) getAllTxHashAndDetails(ctx sdk.Context) (list []TxHashAndDetails, returnErr error) {
	hashAndIDStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.HashAndIDStore)
	iterator := hashAndIDStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var value cosmosTypes.TxHashValue
		returnErr = k.cdc.Unmarshal(iterator.Value(), &value)
		if returnErr != nil {
			return nil, returnErr
		}
		list = append(list, TxHashAndDetails{string(iterator.Key()), value})
	}
	return list, nil
}

func StoreValueEqualOrNotTxStatus(storeValue cosmosTypes.TxHashValue, msgValue cosmosTypes.MsgTxStatus) bool {
	if storeValue.TxStatus.TxHash != msgValue.TxHash {
		return false
	}
	if storeValue.TxStatus.AccountNumber != msgValue.AccountNumber {
		return false
	}
	if storeValue.TxStatus.SequenceNumber != msgValue.SequenceNumber {
		return false
	}
	if !storeValue.TxStatus.Balance.IsEqual(msgValue.Balance) {
		return false
	}

	valueValidatorMap := make(map[string]cosmosTypes.ValidatorDetails)
	for _, vd := range storeValue.TxStatus.ValidatorDetails {
		valueValidatorMap[vd.ValidatorAddress] = vd
	}

	validatorDetailsMap := make(map[string]cosmosTypes.ValidatorDetails)
	for _, vd := range msgValue.ValidatorDetails {
		validatorDetailsMap[vd.ValidatorAddress] = vd
	}

	return reflect.DeepEqual(valueValidatorMap, validatorDetailsMap)
}

//______________________________________________________________________________________________

type TransactionQueue struct {
	txID   uint64
	status cosmosTypes.OutgoingQueueValue
}

// Set new transaction in transaction queue with value 0 (pending)
func (k Keeper) setNewInTransactionQueue(ctx sdk.Context, txID uint64) {
	transactionQueueStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyTransactionQueue)
	key := cosmosTypes.UInt64Bytes(txID)
	if transactionQueueStore.Has(key) {
		panic(fmt.Errorf("transaction present in queue"))
	}

	// true : active transaction, false : inactive transaction
	value := cosmosTypes.NewOutgoingQueueValue(false, 0)
	bz := k.cdc.MustMarshal(&value)
	transactionQueueStore.Set(key, bz)
}

// Get active transaction from the tx queue : returns 0 if no active transaction in queue
func (k Keeper) getActiveFromTransactionQueue(ctx sdk.Context) uint64 {
	transactionQueueStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyTransactionQueue)

	// returns the first transaction which is active : supposed to be the first transaction in the list
	iterator := transactionQueueStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var value cosmosTypes.OutgoingQueueValue
		k.cdc.MustUnmarshal(iterator.Value(), &value)
		if value.Active {
			return cosmosTypes.UInt64FromBytes(iterator.Key())
		}
	}

	// if txID returned is 0 : there is no active transaction
	return 0
}

func (k Keeper) incrementRetryCounterInTransactionQueue(ctx sdk.Context, txID uint64) {
	transactionQueueStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyTransactionQueue)
	key := cosmosTypes.UInt64Bytes(txID)

	if transactionQueueStore.Has(key) {
		var value cosmosTypes.OutgoingQueueValue
		k.cdc.MustUnmarshal(transactionQueueStore.Get(key), &value)

		// disable module if the retry counter has reached the max count
		if value.RetryCounter >= k.GetParams(ctx).RetryLimit {
			k.disableModule(ctx)
		}

		//increment retry counter
		value.RetryCounter++

		bz := k.cdc.MustMarshal(&value)
		transactionQueueStore.Set(key, bz)
	}
}

// Fetches the next transaction to be sent out and mark it active
// called after deleting the active transaction which has been successful
func (k Keeper) getNextFromTransactionQueue(ctx sdk.Context) uint64 {
	transactionQueueStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyTransactionQueue)

	//start iteration through the store and return the first key found in the store
	//as the keys stored are in ascending order
	iterator := transactionQueueStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := cosmosTypes.UInt64FromBytes(iterator.Key())
		value := cosmosTypes.NewOutgoingQueueValue(false, 0)
		bz := k.cdc.MustMarshal(&value)
		transactionQueueStore.Set(iterator.Key(), bz)
		return key
	}

	// if txID returned is zero : there are 0 pending transactions
	return 0
}

// Removes the transaction corresponding to the given txID
// called once the transaction is successful and all action required after its success are complete
func (k Keeper) removeFromTransactionQueue(ctx sdk.Context, txID uint64) {
	transactionQueueStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyTransactionQueue)
	transactionQueueStore.Delete(cosmosTypes.UInt64Bytes(txID))
}

// Gets the list of all transaction in the outgoing queue which are being sent out or yet to be sent out
func (k Keeper) getAllFromTransactionQueue(ctx sdk.Context) (txIDAndStatus []TransactionQueue) {
	transactionQueueStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyTransactionQueue)

	//iterate through all the transactions present in queue and add to map
	iterator := transactionQueueStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := cosmosTypes.UInt64FromBytes(iterator.Key())

		var value cosmosTypes.OutgoingQueueValue
		k.cdc.MustUnmarshal(iterator.Value(), &value)

		txIDAndStatus = append(txIDAndStatus, TransactionQueue{txID: key, status: value})
	}
	return txIDAndStatus
}

// Emits event for transaction to be picked up by oracles to be signed
func (k Keeper) emitEventForActiveTransaction(ctx sdk.Context, txID uint64) {
	k.incrementRetryCounterInTransactionQueue(ctx, txID)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			cosmosTypes.EventTypeOutgoing,
			sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	)
	k.setEventEmittedFlag(ctx, txID, true)
}

//______________________________________________________________________________________________

// RetryTransactionWithDoubleGas : retry txn with double gas
func (k Keeper) retryTransactionWithFailure(ctx sdk.Context, txDetails cosmosTypes.QueryOutgoingTxByIDResponse, txID uint64, txHash string, failure string) {

	// doubles gas fees and emit a new event
	cosmosTxDetails := txDetails.CosmosTxDetails

	cosmosTxDetails.Tx.AuthInfo.SignerInfos = nil
	cosmosTxDetails.Tx.Signatures = nil
	cosmosTxDetails.TxHash = ""

	// double gas in case of gas failure
	if failure == "gas failure" {
		//cosmosTxDetails.Tx.AuthInfo.Fee.GasLimit == cosmosTypes.GasLimit &&
		//2*cosmosTxDetails.Tx.AuthInfo.Fee.GasLimit < cosmosTypes.GasLimit // TODO
		// TODO : test case when transaction fails even after reaching max_gas limit
		cosmosTxDetails.Tx.AuthInfo.Fee.GasLimit *= 2
	}

	//set it back again in outgoing txn
	k.SetNewTxnInOutgoingPool(ctx, txID, cosmosTxDetails)

	//remove txHash and mapping
	k.removeTxHashAndDetails(ctx, txHash)
}

func (k Keeper) ProcessAllTxAndDetails(ctx sdk.Context) {
	// fetch active transaction in the queue
	txID := k.getActiveFromTransactionQueue(ctx)

	//if txID returned is 0, then emit a new transaction
	if txID == 0 {
		nextID := k.getNextFromTransactionQueue(ctx)
		if nextID == 0 {
			return
		}
		k.emitEventForActiveTransaction(ctx, nextID)
		return
	}

	// get all txHash and details aggregated
	txDetails, err := k.getAllTxHashAndDetails(ctx)
	if err != nil {
		panic(err)
	}

	queryResponse, err := k.getTxnFromOutgoingPoolByID(ctx, txID)
	if err != nil {
		panic(err)
	}

	for _, tx := range txDetails {

		// avoid processing inactive transaction
		if tx.TxHash != queryResponse.CosmosTxDetails.TxHash {
			continue
		}

		// find majority status string : as spam might be possible
		majorityStatus := FindMajority(tx.Details.Status)

		// get tx from outgoing pool
		cosmosTx, err := k.getTxnFromOutgoingPoolByID(ctx, txID)
		if err != nil {
			panic(err)
		}

		custodialAddress, err := cosmosTypes.AccAddressFromBech32(k.GetParams(ctx).CustodialAddress, cosmosTypes.Bech32Prefix)
		if err != nil {
			panic(err)
		}

		multisigAccount := k.getAccountState(ctx, custodialAddress)
		if multisigAccount == nil {
			panic(cosmosTypes.ErrMultiSigAddressNotFound)
		}

		txHashValue, err := k.getTxHashAndDetails(ctx, tx.TxHash)
		if err != nil {
			panic(err)
		}

		// process tx if majority status is present
		if tx.Details.Ratio.LT(cosmosTypes.MinimumRatioForMajority) {
			panic(err)
		}

		// TODO : deal with keeper failure and insufficient balance
		switch majorityStatus {
		case cosmosTypes.GasFailure:
			// retry txn with given failure
			k.retryTransactionWithFailure(ctx, cosmosTx, txID, tx.TxHash, majorityStatus)
			k.emitEventForActiveTransaction(ctx, txID)
		case cosmosTypes.Success:
			// process txn success and perform success actions
			msgs := cosmosTx.CosmosTxDetails.Tx.GetMsgs()
			for _, msg := range msgs {
				execMsgs := msg.(*authz.MsgExec).Msgs
				for _, im := range execMsgs {
					//Only first element is checked as event transactions will always be grouped as one type of message
					switch im.GetCachedValue().(type) {
					case *stakingTypes.MsgDelegate:
						k.updateStatusOnceProcessed(ctx, txID, "success")
						k.SubFromVirtuallyStaked(ctx, GetAmountFromMessage(execMsgs))
						k.AddToStaked(ctx, GetAmountFromMessage(execMsgs))
					case *stakingTypes.MsgUndelegate:
						k.setEpochAndValidatorDetailsForAllUndelegations(ctx, txID)
						k.updateStatusOnceProcessed(ctx, txID, "success")
						k.SubFromVirtuallyUnbonded(ctx, GetAmountFromMessage(execMsgs))
						k.SubFromStaked(ctx, GetAmountFromMessage(execMsgs))
					}
					break
				}
				k.removeFromTransactionQueue(ctx, txID)
			}
		case cosmosTypes.SequenceMismatch:
			// retry txn with the given failure
			k.retryTransactionWithFailure(ctx, cosmosTx, txID, tx.TxHash, majorityStatus)
			k.emitEventForActiveTransaction(ctx, txID)
		}

		// TODO : handle balance, bonded tokens and unbonding tokens value
		k.setCosmosBalance(ctx, tx.Details.TxStatus.Balance)
		for _, delegation := range tx.Details.TxStatus.ValidatorDetails {
			valAddress, err := cosmosTypes.ValAddressFromBech32(delegation.ValidatorAddress, cosmosTypes.Bech32PrefixValAddr)
			if err != nil {
				panic(err)
			}
			k.UpdateDelegationCosmosValidator(ctx, valAddress, delegation.BondedTokens, delegation.UnbondingTokens)
		}

		// set sequence number in any case of status, so it stays up to date
		err = multisigAccount.SetSequence(txHashValue.TxStatus.SequenceNumber)
		if err != nil {
			panic(err)
		}

		//set account number in any case of status, so it stays up to date
		err = multisigAccount.SetAccountNumber(txHashValue.TxStatus.AccountNumber)
		if err != nil {
			panic(err)
		}
	}

	txDetailsList, err := k.getAllTxInOutgoingPool(ctx)
	if err != nil {
		panic(err)
	}
	for _, tx := range txDetailsList {
		//remove transaction if active block limit is reached and status is set to success
		if tx.txDetails.ActiveBlockHeight <= ctx.BlockHeight() && tx.txDetails.Status == "success" {
			k.removeTxnDetailsByID(ctx, tx.txID)
			k.removeFromOutgoingSignaturePool(ctx, tx.txID)
			k.removeTxHashAndDetails(ctx, tx.txDetails.TxHash)
			k.removeFromTransactionQueue(ctx, tx.txID)
		}
	}
}

//______________________________________________________________________________________________

// FindMajority Find the majority element in any string slice
func FindMajority(inputArr []string) string {
	var m string //store majority element if exists
	i := 0       //counter
	for _, element := range inputArr {
		// If counter `i` becomes 0, set the current candidate
		// to `nums[j]` and reset the counter to 1
		if i == 0 {
			m = element
			i = 1
		} else {
			// If the counter is non-zero, increment or decrement it
			// according to whether `nums[j]` is a current candidate
			if m == element {
				i++
			} else {
				i--
			}
		}
	}
	return m //return majority element
}

func GetAmountFromMessage(execMsgs []*codecTypes.Any) sdk.Coin {
	tempAmnt := sdk.NewInt64Coin("uatom", 0)
	for _, m := range execMsgs {
		switch m.GetCachedValue().(type) {
		case *stakingTypes.MsgDelegate:
			tempAmnt = tempAmnt.Add(m.GetCachedValue().(*stakingTypes.MsgDelegate).Amount)
		case *stakingTypes.MsgUndelegate:
			tempAmnt = tempAmnt.Add(m.GetCachedValue().(*stakingTypes.MsgDelegate).Amount)
		case *bankTypes.MsgSend:
			tempAmnt = tempAmnt.Add(m.GetCachedValue().(*stakingTypes.MsgDelegate).Amount)
		}
	}
	return tempAmnt
}
