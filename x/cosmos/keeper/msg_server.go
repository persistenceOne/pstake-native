package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) cosmosTypes.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) SetOrchestrator(c context.Context, msg *cosmosTypes.MsgSetOrchestrator) (*cosmosTypes.MsgSetOrchestratorResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	ctx := sdkTypes.UnwrapSDKContext(c)
	validator, e1 := sdkTypes.ValAddressFromBech32(msg.Validator)
	orchestrator, e2 := sdkTypes.AccAddressFromBech32(msg.Orchestrator)
	if e1 != nil || e2 != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	_, foundExistingOrchestratorKey := k.GetOrchestratorValidator(ctx, orchestrator)

	if k.Keeper.stakingKeeper.Validator(ctx, validator) == nil {
		return nil, sdkErrors.Wrap(stakingTypes.ErrNoValidatorFound, validator.String())
	} else if foundExistingOrchestratorKey {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrResetDelegateKeys, validator.String())
	}

	delegateKeys := k.GetDelegateKeys(ctx)
	for i := range delegateKeys {
		if delegateKeys[i].Orchestrator == orchestrator.String() {
			return nil, sdkErrors.Wrap(err, "Duplicate Orchestrator Key")
		}
	}
	// set the orchestrator address
	k.SetOrchestratorValidator(ctx, validator, orchestrator)

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, msg.Type()),
			sdkTypes.NewAttribute(cosmosTypes.AttributeKeySetOperatorAddr, orchestrator.String()),
		),
	)

	return &cosmosTypes.MsgSetOrchestratorResponse{}, nil
}

func (k msgServer) Send(c context.Context, msg *cosmosTypes.MsgSendWithFees) (*cosmosTypes.MsgSendWithFeesResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	from, err := sdkTypes.AccAddressFromBech32(msg.MessageSend.FromAddress)
	if err != nil {
		return nil, err
	}
	to, err := sdkTypes.AccAddressFromBech32(msg.MessageSend.ToAddress)
	if err != nil {
		return nil, err
	}

	msgAny, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	//TODO denom check
	if ctx.IsZero() || sdkTypes.VerifyAddressFormat(from) != nil || sdkTypes.VerifyAddressFormat(to) != nil ||
		!msg.MessageSend.Amount.IsValid() || !msg.Fees.IsValid() {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrInvalid, "arguments")
	}
	//TODO what to do with amount till txn is confirmed? : sample below
	//totalAmount := msg.MessageSend.Amount.Add(msg.Fees)
	//if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, cosmosTypes.ModuleName, totalAmount); err != nil {
	//	return nil, err
	//}

	txID, err := k.AddToOutgoingPool(ctx, from, msgAny)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, msg.Type()),
			sdkTypes.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	)
	return &cosmosTypes.MsgSendWithFeesResponse{}, nil
}

func (k msgServer) Vote(c context.Context, msg *cosmosTypes.MsgVoteWithFees) (*cosmosTypes.MsgVoteWithFeesResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	voter, err := sdkTypes.AccAddressFromBech32(msg.MessageVote.Voter)
	if err != nil {
		return nil, err
	}
	msgAny, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	//TODO checks
	if ctx.IsZero() || sdkTypes.VerifyAddressFormat(voter) != nil || !msg.Fees.IsValid() {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrInvalid, "arguments")
	}
	//TODO what to do with amount till txn is confirmed?

	txID, err := k.AddToOutgoingPool(ctx, voter, msgAny)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, msg.Type()),
			sdkTypes.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	)
	return &cosmosTypes.MsgVoteWithFeesResponse{}, nil
}

func (k msgServer) Delegate(c context.Context, msg *cosmosTypes.MsgDelegateWithFees) (*cosmosTypes.MsgDelegateWithFeesResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	// TODO
	delegator_address, err := sdkTypes.AccAddressFromBech32(msg.MessageDelegate.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	validator_address, err := sdkTypes.AccAddressFromBech32(msg.MessageDelegate.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	msgAny, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	//TODO checks
	if ctx.IsZero() || sdkTypes.VerifyAddressFormat(delegator_address) != nil || sdkTypes.VerifyAddressFormat(validator_address) != nil ||
		!msg.MessageDelegate.Amount.IsValid() || !msg.Fees.IsValid() {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrInvalid, "arguments")
	}
	//TODO what to do with amount till txn is confirmed?

	txID, err := k.AddToOutgoingPool(ctx, delegator_address, msgAny)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, msg.Type()),
			sdkTypes.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	)
	return &cosmosTypes.MsgDelegateWithFeesResponse{}, nil
}

func (k msgServer) Undelegate(c context.Context, msg *cosmosTypes.MsgUndelegateWithFees) (*cosmosTypes.MsgUndelegateWithFeesResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	//TODO
	delegator_address, err := sdkTypes.AccAddressFromBech32(msg.MessageUndelegate.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	validator_address, err := sdkTypes.AccAddressFromBech32(msg.MessageUndelegate.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	msgAny, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	//TODO checks
	if ctx.IsZero() || sdkTypes.VerifyAddressFormat(delegator_address) != nil || sdkTypes.VerifyAddressFormat(validator_address) != nil ||
		!msg.MessageUndelegate.Amount.IsValid() || !msg.Fees.IsValid() {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrInvalid, "arguments")
	}
	//TODO what to do with amount till txn is confirmed?

	txID, err := k.AddToOutgoingPool(ctx, delegator_address, msgAny)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, msg.Type()),
			sdkTypes.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	)
	return &cosmosTypes.MsgUndelegateWithFeesResponse{}, nil
}

func (k msgServer) MintTokensForAccount(c context.Context, msg *cosmosTypes.MsgMintTokensForAccount) (*cosmosTypes.MsgMintTokensForAccountResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	destinationAddress, err := sdkTypes.AccAddressFromBech32(msg.AddressFromMemo)
	if err != nil {
		return nil, err
	}

	orchestratorAddress, err := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, err
	}

	if ctx.IsZero() || sdkTypes.VerifyAddressFormat(destinationAddress) != nil || sdkTypes.VerifyAddressFormat(orchestratorAddress) != nil ||
		!msg.Amount.IsValid() {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrInvalid, "arguments")
	}

	coinsAmount := msg.Amount.AmountOf("uatom")
	coinString := coinsAmount.String() + cosmosTypes.MintDenom
	newAmount, err := sdkTypes.ParseCoinsNormalized(coinString)
	if err != nil {
		return nil, err
	}

	k.setMintAddressAndAmount(ctx, msg.ChainID, msg.BlockHeight, msg.TxHash, destinationAddress, newAmount)

	_, found := k.GetOrchestratorValidator(ctx, orchestratorAddress)
	if found {
		err = k.addToMintingPoolTx(ctx, destinationAddress, orchestratorAddress, msg.Amount)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, msg.Type()),
			sdkTypes.NewAttribute("set_minting_txn", orchestratorAddress.String()),
		),
	)
	return &cosmosTypes.MsgMintTokensForAccountResponse{}, nil
}
