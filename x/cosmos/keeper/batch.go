package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type TxHashAndDetails struct {
	TxHash  string
	Details cosmosTypes.TxHashValue
}

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
	bz, err := tx.Marshal()
	if err != nil {
		panic(err)
	}
	outgoingStore.Set(key, bz)
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
	err := cosmosTx.Unmarshal(bz)
	if err != nil {
		return cosmosTypes.QueryOutgoingTxByIDResponse{}, err
	}
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

//Sets txBytes once received from Orchestrator after signing.
func (k Keeper) setTxDetailsSignedByOrchestrator(ctx sdk.Context, txID uint64, txHash string, tx sdkTx.Tx) error {
	outgoingStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(cosmosTypes.OutgoingTXPoolKey))
	key := cosmosTypes.UInt64Bytes(txID)
	var cosmosTx cosmosTypes.CosmosTx
	if outgoingStore.Has(key) {
		err := cosmosTx.Unmarshal(outgoingStore.Get(key))
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
		err := tx.Unmarshal(iterator.Value())
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
// Set details corresponding to a particular txHash and update details if already present
func (k Keeper) setTxHashAndDetails(ctx sdk.Context, orchAddress sdk.AccAddress, txID uint64, txHash string, status string) {
	txHashAndTxIDStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.HashAndIDStore)
	key := []byte(txHash)
	if txHashAndTxIDStore.Has(key) {
		var txHashValue cosmosTypes.TxHashValue
		err := txHashValue.Unmarshal(txHashAndTxIDStore.Get(key))
		if err != nil {
			panic("error in unmarshalling txHashValue")
		}
		if !txHashValue.Find(orchAddress.String()) {
			txHashValue.OrchestratorAddresses = append(txHashValue.OrchestratorAddresses, orchAddress.String())
			txHashValue.Status = append(txHashValue.Status, status)
			txHashValue.Counter++
			txHashValue.Ratio = float32(txHashValue.Counter) / float32(k.getTotalValidatorOrchestratorCount(ctx))
			bz, err := txHashValue.Marshal()
			if err != nil {
				panic("error in marshaling txHashValue")
			}
			txHashAndTxIDStore.Set(key, bz)
		}
	} else {
		ratio := float32(1) / float32(k.getTotalValidatorOrchestratorCount(ctx))
		newTxHashValue := cosmosTypes.NewTxHashValue(txID, orchAddress, ratio, status, ctx.BlockHeight(), ctx.BlockHeight()+cosmosTypes.StorageWindow)
		bz, err := newTxHashValue.Marshal()
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
		err := txHashAndValue.Unmarshal(hashAndIDStore.Get(key))
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
		returnErr = value.Unmarshal(iterator.Value())
		if returnErr != nil {
			return nil, returnErr
		}
		list = append(list, TxHashAndDetails{string(iterator.Key()), value})
	}
	return list, nil
}

// RetryTransactionWithDoubleGas : retry txn with double gas
func (k Keeper) retryTransactionWithDoubleGas(ctx sdk.Context, txDetails cosmosTypes.QueryOutgoingTxByIDResponse, txID uint64, txHash string) {
	// doubles gas fees and emit a new event
	cosmosTxDetails := txDetails.CosmosTxDetails
	cosmosTxDetails.Tx.AuthInfo.Fee.GasLimit = cosmosTxDetails.Tx.AuthInfo.Fee.GasLimit * 2
	cosmosTxDetails.Tx.AuthInfo.SignerInfos = nil
	cosmosTxDetails.Tx.Signatures = nil
	cosmosTxDetails.TxHash = ""

	//create new ID for next txn and set it in kv store
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

	//emit new event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			cosmosTypes.EventTypeOutgoing,
			sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(nextID)),
		),
	)

	//set new outgoing txn
	k.setNewTxnInOutgoingPool(ctx, nextID, cosmosTxDetails)

	//remove old txn from db
	k.removeTxnDetailsByID(ctx, txID)

	//remove txHash and mapping
	k.removeTxHashAndDetails(ctx, txHash)
}

// ProcessAllTxAndDetails Process all the transaction details that are pending and retry if failed by less gas or delete them once they are past the avtive block limit
func (k Keeper) ProcessAllTxAndDetails(ctx sdk.Context) error {
	list, err := k.getAllTxHashAndDetails(ctx)
	if err != nil {
		panic(err)
	}
	for _, element := range list {
		majorityStatus := FindMajority(element.Details.Status)
		txDetails, err := k.getTxnFromOutgoingPoolByID(ctx, element.Details.TxID)
		//if err != nil {
		//k.removeTxHashAndDetails(ctx, element.TxHash)
		//TODO : Check if Signed tx is sent later than Status of the same.
		//return err
		//}
		//err check to see if details have been found or not
		if err == nil {
			if element.Details.Ratio > 0.80 {
				if majorityStatus == "failure" {
					k.retryTransactionWithDoubleGas(ctx, txDetails, element.Details.TxID, element.TxHash)
				}
				if majorityStatus == "success" {
					msgs := txDetails.CosmosTxDetails.Tx.GetMsgs()
					//Only first element is checked as event transactions will always be grouped as one type of message
					switch msgs[0].(type) {
					//TODO : Add cases for rewards claim, unbonding
					case *stakingTypes.MsgDelegate:
						k.updateCosmosValidatorStakingParams(ctx, msgs)
					case *types.MsgWithdrawDelegatorReward:
						k.emitStakingTxnForClaimedRewards(ctx, msgs)
					case *stakingTypes.MsgUndelegate:
						// TODO
					}
				}
			}
			if txDetails.CosmosTxDetails.ActiveBlockHeight >= ctx.BlockHeight() {
				//TODO : check for rewards claim txns
				k.removeTxnDetailsByID(ctx, element.Details.TxID)
				k.removeTxHashAndDetails(ctx, element.TxHash)
			}
		}

	}
	txDetailsList, err := k.getAllTxInOutgoingPool(ctx) //TODO Implement Rest
	if err != nil {
		panic(err)
	}
	for _, element := range txDetailsList {
		if element.txDetails.ActiveBlockHeight <= ctx.BlockHeight() {
			k.removeTxnDetailsByID(ctx, element.txID)
			element.txDetails.NativeBlockHeight = ctx.BlockHeight()
			element.txDetails.ActiveBlockHeight = ctx.BlockHeight() + cosmosTypes.StorageWindow
			nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))
			k.setNewTxnInOutgoingPool(ctx, nextID, element.txDetails)
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
