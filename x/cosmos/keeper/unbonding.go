package keeper

import (
	"fmt"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) generateUnbondingOutgoingEvent(ctx sdk.Context, listOfValidatorsAndUnbondingAmount []ValAddressAndAmountForStakingAndUnstaking, epochNumber int64) {
	params := k.GetParams(ctx)

	chunkMsgs := ChunkUndelegationSlice(listOfValidatorsAndUnbondingAmount, params.ChunkSize)

	for _, chunk := range chunkMsgs {
		nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

		var undelegateMsgsAny []*codecTypes.Any
		var undelegategMsgs []stakingTypes.MsgUndelegate
		for _, element := range chunk {
			msg := stakingTypes.MsgUndelegate{
				DelegatorAddress: params.CustodialAddress,
				ValidatorAddress: element.validator.String(),
				Amount:           element.amount,
			}
			anyMsg, err := codecTypes.NewAnyWithValue(&msg)
			if err != nil {
				panic(err)
			}
			undelegateMsgsAny = append(undelegateMsgsAny, anyMsg)
			undelegategMsgs = append(undelegategMsgs, msg)
		}

		tx := cosmosTypes.CosmosTx{
			Tx: sdkTx.Tx{
				Body: &sdkTx.TxBody{
					Messages:      undelegateMsgsAny,
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
			EventEmitted:      true,
			Status:            "",
			TxHash:            "",
			NativeBlockHeight: ctx.BlockHeight(),
			ActiveBlockHeight: ctx.BlockHeight() + cosmosTypes.StorageWindow,
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				cosmosTypes.EventTypeOutgoing,
				sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(nextID)),
			),
		)

		err := k.setIDInEpochPoolForWithdrawals(ctx, nextID, undelegategMsgs, params.CustodialAddress, epochNumber)
		if err != nil {
			panic(err)
		}
		//Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
		k.setNewTxnInOutgoingPool(ctx, nextID, tx)
	}
}

func (k Keeper) setIDInEpochPoolForWithdrawals(ctx sdk.Context, txID uint64, undelegateMsgs []stakingTypes.MsgUndelegate, custodialAddress string, epochNumber int64) error {
	unbondingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingUnbondStore)
	key := cosmosTypes.UInt64Bytes(txID)
	value := cosmosTypes.NewValuOutgoingUnbondStore(undelegateMsgs, epochNumber)
	bz, err := value.Marshal()
	if err != nil {
		return err
	}
	unbondingEpochStore.Set(key, bz)
	return nil
}

func ChunkUndelegationSlice(slice []ValAddressAndAmountForStakingAndUnstaking, chunkSize int64) (chunks [][]ValAddressAndAmountForStakingAndUnstaking) {
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if int64(len(slice)) < chunkSize {
			chunkSize = int64(len(slice))
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
