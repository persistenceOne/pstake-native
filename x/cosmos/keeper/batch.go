package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
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

func (k Keeper) setNewTxnInOutgoingPool(ctx sdk.Context, txID uint64, tx cosmosTypes.CosmosTx) {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(cosmosTypes.OutgoingTXPoolKey))
	key := cosmosTypes.UInt64Bytes(txID)
	bz, err := k.cdc.Marshal(&tx)
	if err != nil {
		panic(err)
	}
	outgoingStore.Set(key, bz)
}

func (k Keeper) updateStatusOnceProcessed(cts sdk.Context, txID uint64, tx cosmosTypes.CosmosTx) {
	//TODO
}

//gets txn details by ID
func (k Keeper) getTxnFromOutgoingPoolByID(ctx sdk.Context, txID uint64) (cosmosTypes.QueryOutgoingTxByIDResponse, error) {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(cosmosTypes.OutgoingTXPoolKey))
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
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(cosmosTypes.OutgoingTXPoolKey))
	key := cosmosTypes.UInt64Bytes(txID)
	outgoingStore.Delete(key)
}

func (k Keeper) setOutgoingTxnSignaturesAndEmitEvent(ctx sdk.Context, tx cosmosTypes.CosmosTx, txID uint64) error {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(cosmosTypes.OutgoingTXPoolKey))
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

//Sets txBytes once received from Orchestrator after signing.
func (k Keeper) setTxDetailsSignedByOrchestrator(ctx sdk.Context, txID uint64, txHash string, tx sdkTx.Tx) error {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(cosmosTypes.OutgoingTXPoolKey))
	key := cosmosTypes.UInt64Bytes(txID)
	var cosmosTx cosmosTypes.CosmosTx
	if outgoingStore.Has(key) {
		err := k.cdc.Unmarshal(outgoingStore.Get(key), &cosmosTx)
		if err != nil {
			return err
		}

		cosmosTx.TxHash = txHash
		cosmosTx.Tx = tx

		bz, err := cosmosTx.Marshal()
		if err != nil {
			return err
		}

		outgoingStore.Set(key, bz)
	}
	return nil
}

func (k Keeper) getAllTxInOutgoingPool(ctx sdk.Context) (details []txIDAndDetailsInOutgoingPool, err error) {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(cosmosTypes.OutgoingTXPoolKey))
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
/*
TODO : Add key and value structure
*/

type TxHashAndDetails struct {
	TxHash  string
	Details cosmosTypes.TxHashValue
}

// Set details corresponding to a particular txHash and update details if already present
func (k Keeper) setTxHashAndDetails(ctx sdk.Context, orchAddress sdk.AccAddress, txID uint64, txHash string, status string, accountNumber uint64, sequenceNumber uint64) {
	txHashAndTxIDStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.HashAndIDStore)
	key := []byte(txHash)
	if txHashAndTxIDStore.Has(key) {
		var txHashValue cosmosTypes.TxHashValue
		err := k.cdc.Unmarshal(txHashAndTxIDStore.Get(key), &txHashValue)
		if err != nil {
			panic("error in unmarshalling txHashValue")
		}
		if !txHashValue.Find(orchAddress.String()) {
			txHashValue.OrchestratorAddresses = append(txHashValue.OrchestratorAddresses, orchAddress.String())
			txHashValue.Status = append(txHashValue.Status, status)
			txHashValue.Counter++
			txHashValue.Ratio = float32(txHashValue.Counter) / float32(k.getTotalValidatorOrchestratorCount(ctx))
			bz, err := k.cdc.Marshal(&txHashValue)
			if err != nil {
				panic("error in marshaling txHashValue")
			}
			txHashAndTxIDStore.Set(key, bz)
		}
	} else {
		ratio := float32(1) / float32(k.getTotalValidatorOrchestratorCount(ctx))
		newTxHashValue := cosmosTypes.NewTxHashValue(txID, orchAddress, ratio, status, ctx.BlockHeight(), ctx.BlockHeight()+cosmosTypes.StorageWindow, accountNumber, sequenceNumber)
		bz, err := k.cdc.Marshal(&newTxHashValue)
		if err != nil {
			panic("error in marshaling txHashValue")
		}
		txHashAndTxIDStore.Set(key, bz)
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

//______________________________________________________________________________________________

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
		if value.Active == true {
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
func (k Keeper) getAllFromTransactionQueue(ctx sdk.Context) (txIDAndStatusMap map[uint64]cosmosTypes.OutgoingQueueValue) {
	transactionQueueStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyTransactionQueue)

	//iterate through all the transactions present in queue and add to map
	iterator := transactionQueueStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := cosmosTypes.UInt64FromBytes(iterator.Key())

		var value cosmosTypes.OutgoingQueueValue
		k.cdc.MustUnmarshal(iterator.Value(), &value)

		txIDAndStatusMap[key] = value
	}
	return txIDAndStatusMap
}

// Emits event for transaction to be picked up by oracles to be signed
func (k Keeper) emitEventForActiveTransaction(ctx sdk.Context, txID uint64) {
	//TODO : increment counter

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			cosmosTypes.EventTypeOutgoing,
			sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	)
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
		cosmosTxDetails.Tx.AuthInfo.Fee.GasLimit = cosmosTxDetails.Tx.AuthInfo.Fee.GasLimit * 2
	}

	//set it back again in outgoing txn
	k.setNewTxnInOutgoingPool(ctx, txID, cosmosTxDetails)

	//remove txHash and mapping
	k.removeTxHashAndDetails(ctx, txHash)
}

func (k Keeper) ProcessAllTxAndDetails(ctx sdk.Context) error {
	// fetch active transaction in the queue
	txID := k.getActiveFromTransactionQueue(ctx)

	//if txID returned is 0, then emit a new transaction
	if txID == 0 {
		nextID := k.getNextFromTransactionQueue(ctx)
		k.emitEventForActiveTransaction(ctx, nextID)
		return nil
	}

	// get all txHash and details aggregated
	txDetails, err := k.getAllTxHashAndDetails(ctx)
	if err != nil {
		return err
	}

	for _, tx := range txDetails {
		// avoid processing inactive transaction
		if tx.Details.TxID != txID {
			continue
		}

		// find majority status string : as spam might be possible
		majorityStatus := FindMajority(tx.Details.Status)

		// get tx from outgoing pool
		cosmosTx, err := k.getTxnFromOutgoingPoolByID(ctx, tx.Details.TxID)
		if err != nil {
			return err
		}

		custodialAddress, err := cosmosTypes.AccAddressFromBech32(k.GetParams(ctx).CustodialAddress, cosmosTypes.Bech32Prefix)
		if err != nil {
			return err
		}
		//TODO : remove bug
		multisigAccount := k.authKeeper.GetAccount(ctx, custodialAddress)
		if multisigAccount == nil {
			return cosmosTypes.ErrMultiSigAddressNotFound
		}

		txHashValue, err := k.getTxHashAndDetails(ctx, tx.TxHash)
		if err != nil {
			return err
		}

		// process tx if majority status is present
		if tx.Details.Ratio < cosmosTypes.MinimumRatioForMajority {
			return nil
		}

		if majorityStatus == "gas failure" {
			// retry txn with given failure
			k.retryTransactionWithFailure(ctx, cosmosTx, tx.Details.TxID, tx.TxHash, majorityStatus)
			k.emitEventForActiveTransaction(ctx, txID)
		} else if majorityStatus == "success" {
			// process txn success and perform success actions
			msgs := cosmosTx.CosmosTxDetails.Tx.GetMsgs()
			for _, msg := range msgs {
				execMsgs := msg.(*authz.MsgExec).Msgs
				for _, im := range execMsgs {
					//Only first element is checked as event transactions will always be grouped as one type of message
					switch im.GetCachedValue().(type) {
					case *stakingTypes.MsgDelegate:
						//TODO : update C value
						err = k.processStakingSuccessTxns(ctx, tx.Details.TxID)
						cosmosTx.CosmosTxDetails.Status = "success"
						k.updateStatusOnceProcessed(ctx, tx.Details.TxID, cosmosTx.CosmosTxDetails)
					case *stakingTypes.MsgUndelegate:
						//TODO : update C value
						err = k.setEpochAndValidatorDetailsForAllUndelegations(ctx, tx.Details.TxID)
						cosmosTx.CosmosTxDetails.Status = "success"
						k.updateStatusOnceProcessed(ctx, tx.Details.TxID, cosmosTx.CosmosTxDetails)
						if err != nil {
							return err
						}
					case *types.MsgSend:
						// TODO : update C value
					}
					break
				}
				k.removeFromTransactionQueue(ctx, txID)
				nextID := k.getNextFromTransactionQueue(ctx)
				k.emitEventForActiveTransaction(ctx, nextID)
			}
		} else if majorityStatus == "sequence mismatch" {
			// retry txn with the given failure
			k.retryTransactionWithFailure(ctx, cosmosTx, txID, tx.TxHash, majorityStatus)
			k.emitEventForActiveTransaction(ctx, txID)
		}

		// set sequence number in any case of status, so it stays up to date
		err = multisigAccount.SetSequence(txHashValue.SequenceNumber)
		if err != nil {
			return err
		}

		//set account number in any case of status, so it stays up to date
		err = multisigAccount.SetAccountNumber(txHashValue.AccountNumber)
		if err != nil {
			return err
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
	return nil
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
