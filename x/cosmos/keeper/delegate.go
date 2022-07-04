package keeper

import (
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// Generate an event for delegating on cosmos chain once staking epoch is called
func (k Keeper) generateDelegateOutgoingEvent(ctx sdk.Context, validatorSet []ValAddressAmount) error {
	params := k.GetParams(ctx)

	//create chunks for delegation on cosmos chain
	chunkSlice := ChunkDelegateAndUndelegateSlice(validatorSet, params.ChunkSize)

	for _, chunk := range chunkSlice {
		nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

		var delegateMsgsAny []*codecTypes.Any
		for _, element := range chunk {
			msg := stakingTypes.MsgDelegate{
				DelegatorAddress: params.CustodialAddress,
				ValidatorAddress: element.Validator,
				Amount:           element.Amount,
			}
			anyMsg, err := codecTypes.NewAnyWithValue(&msg)
			if err != nil {
				return err
			}
			delegateMsgsAny = append(delegateMsgsAny, anyMsg)
		}

		cosmosAddrr, err := cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32PrefixAccAddr, k.GetCurrentAddress(ctx))
		if err != nil {
			return err
		}
		execMsg := authz.MsgExec{
			Grantee: cosmosAddrr,
			Msgs:    delegateMsgsAny,
		}

		execMsgAny, err := codecTypes.NewAnyWithValue(&execMsg)
		if err != nil {
			return err
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

		// set acknowledgment flag true for future reference (not any yet)

		// Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
		k.SetNewTxnInOutgoingPool(ctx, nextID, tx)

		// sets the transaction in transaction queue with the given tx ID
		k.setNewInTransactionQueue(ctx, nextID)
	}

	return nil
}
