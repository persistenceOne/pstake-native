package keeper

import (
	"fmt"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) generateWithdrawRewardsEvent(ctx sdk.Context, withdrawMsgs []*codecTypes.Any) {
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

	tx := cosmosTypes.CosmosTx{
		Tx: sdkTx.Tx{
			Body: &sdkTx.TxBody{
				Messages:      withdrawMsgs,
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
	//Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
	k.setNewTxnInOutgoingPool(ctx, nextID, tx)
}
