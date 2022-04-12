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
func (k Keeper) generateDelegateOutgoingEvent(ctx sdk.Context, keyAndValue cosmosTypes.KeyAndValueForMinting) error {
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

	params := k.GetParams(ctx)
	//fetches validator set for delegation on cosmos chain
	uatomDenom, err := params.GetBondDenomOf("uatom")
	if err != nil {
		return err
	}
	amount := sdk.NewCoin(uatomDenom, keyAndValue.Value.Amount.AmountOf(uatomDenom))
	validatorSet := k.fetchValidatorsToDelegate(ctx, amount)

	//create messages for delegation on cosmos chain
	var delegateMsgs []*codecTypes.Any
	for _, validator := range validatorSet {
		//amount := sdk.NewCoin(cosmosTypes.StakeDenom, keyAndValue.Value.Amount.AmountOf(cosmosTypes.MintDenom).ToDec().Mul(validator.Weight).TruncateInt())
		msg := stakingTypes.MsgDelegate{
			DelegatorAddress: "cosmos15vm0p2x990762txvsrpr26ya54p5qlz9xqlw5z",
			ValidatorAddress: validator.validator.String(),
			Amount:           validator.amount,
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
	return nil
}

func (k Keeper) setTotalDelegatedAmountTillDate(ctx sdk.Context, addToTotal sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz, err := addToTotal.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(cosmosTypes.KeyTotalDelegationTillDate), bz)
}

func (k Keeper) getTotalDelegatedAmountTillDate(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(cosmosTypes.KeyTotalDelegationTillDate))
	var amount sdk.Coin
	err := amount.Unmarshal(bz)
	if err != nil {
		panic(err)
	}
	return amount
}

func (k Keeper) emitStakingTxnForClaimedRewards(ctx sdk.Context, msgs []sdk.Msg) {
	//totalAmountInClaimMsgs := sdk.NewInt64Coin(k.GetParams(ctx).BondDenom, 0)
	//TODO : Ask which impl to go forwards with txn response for claimRewards and minting rewards for devs and validators
}
