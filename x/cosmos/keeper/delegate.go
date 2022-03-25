package keeper

import (
	"fmt"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// Generate an event for delegating on cosmos chain once tokens are minted on native side
func (k Keeper) generateDelegateOutgoingEvent(ctx sdk.Context, keyAndValue cosmosTypes.KeyAndValueForMinting) {
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

	//fetches validator set for delegation on cosmos chain
	validatorSet := k.GetParams(ctx).ValidatorSetCosmosChain

	//create messages for delegation on cosmos chain
	var delegateMsgs []*codecTypes.Any
	for _, validator := range validatorSet {
		amount := sdk.NewCoin(cosmosTypes.StakeDenom, keyAndValue.Value.Amount.AmountOf(cosmosTypes.MintDenom).ToDec().Mul(validator.Weight).TruncateInt())
		msg := stakingTypes.MsgDelegate{
			DelegatorAddress: "cosmos15vm0p2x990762txvsrpr26ya54p5qlz9xqlw5z",
			ValidatorAddress: validator.Address,
			Amount:           amount,
		}
		msgAny, err := codecTypes.NewAnyWithValue(&msg)
		if err != nil {
			panic(err)
		}
		delegateMsgs = append(delegateMsgs, msgAny)
	}

	tx := cosmosTypes.CosmosTx{
		Tx: sdkTx.Tx{
			Body: &sdkTx.TxBody{
				Messages:      delegateMsgs,
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

	// set acknowledgment flag true for future reference (not any yet)
	k.setAcknowledgmentFlagTrue(ctx, keyAndValue.Key)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			cosmosTypes.EventTypeOutgoing,
			sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(nextID)),
		),
	)
	//Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
	k.setNewTxnInOutgoingPool(ctx, nextID, tx)
}
