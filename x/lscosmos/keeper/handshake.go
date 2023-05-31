package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// OnChanOpenInit implements the IBCModule interface
func (k Keeper) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	hostAccounts := k.GetHostAccounts(ctx)

	// Require portID is the portID module is bound to
	if portID != hostAccounts.DelegatorAccountPortID() &&
		portID != hostAccounts.RewardsAccountPortID() {
		return "", errorsmod.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected either of %s or %s",
			portID, hostAccounts.DelegatorAccountPortID(), hostAccounts.RewardsAccountPortID())
	}
	var versionData icatypes.Metadata
	if err := icatypes.ModuleCdc.UnmarshalJSON([]byte(version), &versionData); err != nil {
		return "", err
	}
	if versionData.Version != icatypes.Version {
		return "", errorsmod.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", versionData.Version, icatypes.Version)
	}

	return version, nil
}

// OnChanOpenTry implements the IBCModule interface
func (k Keeper) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	// Controller Auth Module does not do OnChanOpenTry
	return "", nil
}

// OnChanOpenAck implements the IBCModule interface
func (k Keeper) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	hostAccounts := k.GetHostAccounts(ctx)
	if portID != hostAccounts.DelegatorAccountPortID() &&
		portID != hostAccounts.RewardsAccountPortID() {
		return errorsmod.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected either of %s or %s",
			portID, hostAccounts.DelegatorAccountPortID(), hostAccounts.RewardsAccountPortID())
	}

	var counterpartyVersionData icatypes.Metadata
	if err := icatypes.ModuleCdc.UnmarshalJSON([]byte(counterpartyVersion), &counterpartyVersionData); err != nil {
		return err
	}

	if counterpartyVersionData.Version != icatypes.Version {
		return errorsmod.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: %s, expected %s", counterpartyVersion, icatypes.Version)
	}
	//TODO more checks, capability, channelID??

	hostChainParams := k.GetHostChainParams(ctx)

	if !k.GetModuleState(ctx) {
		if portID == hostAccounts.DelegatorAccountPortID() {
			delegationAddress, delegationAddrfound := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, hostChainParams.ConnectionID, hostAccounts.DelegatorAccountPortID())
			if delegationAddrfound {
				if err := k.SetHostChainDelegationAddress(ctx, delegationAddress); err != nil {
					return err
				}
				if err := k.icaControllerKeeper.RegisterInterchainAccount(ctx, hostChainParams.ConnectionID, hostAccounts.RewardsAccountOwnerID, ""); err != nil {
					return errorsmod.Wrap(err, "Could not register ica reward Address")
				}

			}
		}
		if portID == hostAccounts.RewardsAccountPortID() {
			rewardAddress, rewardAddrFound := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, hostChainParams.ConnectionID, hostAccounts.RewardsAccountPortID())
			delegationAddress := k.GetDelegationState(ctx).HostChainDelegationAddress
			if rewardAddrFound {
				_ = k.SetHostChainRewardAddressIfEmpty(ctx, types.NewHostChainRewardAddress(rewardAddress))
				setWithdrawAddrMsg := &distributiontypes.MsgSetWithdrawAddress{
					DelegatorAddress: delegationAddress,
					WithdrawAddress:  rewardAddress,
				}
				err := k.GenerateAndExecuteICATx(ctx, hostChainParams.ConnectionID, hostAccounts.DelegatorAccountOwnerID, []proto.Message{setWithdrawAddrMsg})
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	k.Logger(ctx).Info(fmt.Sprintf("Recreating ICA channel with channelID: %s, portID: %s, counterpartyID: %s", channelID, portID, counterpartyChannelID))

	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (k Keeper) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (k Keeper) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for channels
	return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (k Keeper) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface
func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// Controller Auth Module does not do OnRecvPacket
	return nil
}

// OnAcknowledgementPacket implements the IBCModule interface
func (k Keeper) OnAcknowledgementPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {

	// TODO add checks for ports, channels
	hostChainParams := k.GetHostChainParams(ctx)

	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet acknowledgement: %v", err)
	}

	var icaPacket icatypes.InterchainAccountPacketData
	if err := icatypes.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), &icaPacket); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		k.Logger(ctx).Info(fmt.Sprintln("ICA tx ack failed with ack:", ack.String()))
		err := k.resetToPreICATx(ctx, icaPacket)
		if err != nil {
			return err
		}
	case *channeltypes.Acknowledgement_Result:
		// this line is used by starport scaffolding # oracle/packet/module/ack
		err := k.handleSuccessfulAck(ctx, ack, icaPacket, hostChainParams)
		if err != nil {
			return err
		}
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
	//Return nil here

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAck, ack.String()),
		),
	)

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttribute(types.AttributeKeyAckSuccess, fmt.Sprintln(ack.Success())),
			),
		)
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttribute(types.AttributeKeyAckError, resp.Error),
			),
		)
	}

	return nil
}

// OnTimeoutPacket implements the IBCModule interface
func (k Keeper) OnTimeoutPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// this line is used by starport scaffolding # oracle/packet/module/ack

	var icaPacket icatypes.InterchainAccountPacketData
	if err := icatypes.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), &icaPacket); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	err := k.resetToPreICATx(ctx, icaPacket)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTimeout,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)
	return nil
}

// handleSuccessfulAck handles successful acknowledgements.
func (k Keeper) handleSuccessfulAck(ctx sdk.Context, ack channeltypes.Acknowledgement, icaPacket icatypes.InterchainAccountPacketData, hostChainParams types.HostChainParams) error {
	txMsgData := &sdk.TxMsgData{}
	if err := k.cdc.Unmarshal(ack.GetResult(), txMsgData); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	msgs, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot Deserialise icapacket data: %v", err)
	}

	// Dispatch packet
	msgsCount := 0
	expectedMsgType := sdk.MsgTypeURL(msgs[0])
	for i, msg := range msgs {
		var data []byte
		if len(txMsgData.Data) == 0 {
			data = txMsgData.GetMsgResponses()[i].Value
		} else {
			data = txMsgData.Data[i].Data
		}
		response, err := k.handleAckMsgData(ctx, data, msg, hostChainParams)
		if err != nil {
			return err
		}
		k.Logger(ctx).Info("message response in ICS-27 packet response", "response", response)
		if expectedMsgType == sdk.MsgTypeURL(msgs[i]) {
			msgsCount++
		}

		// assert all msgs are of same type.
		if len(msgs) == msgsCount {
			switch expectedMsgType {
			case sdk.MsgTypeURL(&distributiontypes.MsgWithdrawDelegatorReward{}):
				rewardAddr := k.GetHostChainRewardAddress(ctx)

				_, rewardAccAddr, err := bech32.DecodeAndConvert(rewardAddr.Address)
				if err != nil {
					return err
				}
				balanceQuery := banktypes.CreatePrefixedAccountStoreKey(rewardAccAddr, []byte(hostChainParams.BaseDenom))

				// total rewards balance withdrawn
				k.icqKeeper.MakeRequest(
					ctx,
					hostChainParams.ConnectionID,
					hostChainParams.ChainID,
					types.BankStoreQuery,
					balanceQuery,
					sdk.NewInt(int64(-1)),
					types.ModuleName,
					RewardsAccountBalance,
					0,
				)
			case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
				previousEpochNumber := types.PreviousUnbondingEpoch(k.epochKeeper.GetEpochInfo(ctx, types.UndelegationEpochIdentifier).CurrentEpoch)
				//May be also match amount with previous epoch incase host chain is down for multiple entire epoch duration. (or add epochnumber in memo ~ not clean, or store (sequenceNumber,epoch of the ica txn) )
				previousEpochUnbondings := k.GetUnbondingEpochCValue(ctx, previousEpochNumber)
				err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.UndelegationModuleAccount, types.ModuleName, sdk.NewCoins(previousEpochUnbondings.STKBurn))
				if err != nil {
					return err
				}
				err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(previousEpochUnbondings.STKBurn))
				if err != nil {
					return err
				}

				//update completionTime
				var msgUndelegateResponse stakingtypes.MsgUndelegateResponse
				if err := k.cdc.Unmarshal(data, &msgUndelegateResponse); err != nil {
					return err
				}
				k.UpdateCompletionTimeForUndelegationEpoch(ctx, previousEpochNumber, msgUndelegateResponse.CompletionTime.Add(types.UndelegationCompletionTimeBuffer))
			default:

			}
		}
	}
	if len(msgs) != msgsCount {
		k.SetModuleState(ctx, false) //Disable module, we assert single type of msg throughout the tx.
		k.Logger(ctx).Error(fmt.Sprintf("%s module has been disabled due to different msg types in a ica txn", types.ModuleName))
		return nil
	}

	return nil
}

// handleAckMsgData handles successful response.
func (k Keeper) handleAckMsgData(ctx sdk.Context, data []byte, msg sdk.Msg, hostChainParams types.HostChainParams) (string, error) {
	switch sdk.MsgTypeURL(msg) {
	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		parsedMsg, ok := msg.(*stakingtypes.MsgDelegate)
		if !ok {
			return "", errorsmod.Wrapf(sdkerrors.ErrInvalidType, "unable to unmarshal msg of type %s", sdk.MsgTypeURL(msg))
		}
		var msgResponse stakingtypes.MsgDelegateResponse
		if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
			return "", errorsmod.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal delegate response message: %s", err.Error())
		}
		// Add delegation state
		k.AddHostAccountDelegation(ctx, types.NewHostAccountDelegation(parsedMsg.ValidatorAddress, parsedMsg.Amount))
		k.RemoveICADelegateFromTransientStore(ctx, parsedMsg.Amount)

		return msgResponse.String(), nil

	case sdk.MsgTypeURL(&distributiontypes.MsgSetWithdrawAddress{}):
		var msgResponse distributiontypes.MsgSetWithdrawAddressResponse
		if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
			return "", errorsmod.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal set withdraw address response message: %s", err.Error())
		}
		k.SetModuleState(ctx, true)
		return msgResponse.String(), nil
	case sdk.MsgTypeURL(&distributiontypes.MsgWithdrawDelegatorReward{}):
		var msgResponse distributiontypes.MsgWithdrawDelegatorRewardResponse
		if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
			return "", errorsmod.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal withdraw delegator reward response message: %s", err.Error())
		}
		return msgResponse.String(), nil
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		parsedMsg, ok := msg.(*banktypes.MsgSend)
		if !ok {
			return "", errorsmod.Wrapf(sdkerrors.ErrInvalidType, "unable to unmarshal msg of type %s", sdk.MsgTypeURL(msg))
		}
		var msgResponse banktypes.MsgSendResponse
		if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
			return "", errorsmod.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}
		//is from rewardaddr to delegationaddr?
		rewardAddress := k.GetHostChainRewardAddress(ctx)
		delegationState := k.GetDelegationState(ctx)
		if rewardAddress.Address == parsedMsg.FromAddress && delegationState.HostChainDelegationAddress == parsedMsg.ToAddress {
			amountOfBaseDenom := parsedMsg.Amount.AmountOf(hostChainParams.BaseDenom)
			if amountOfBaseDenom.GT(sdk.ZeroInt()) {
				cValue := k.GetCValue(ctx)

				k.AddBalanceToDelegationState(ctx, sdk.NewCoin(hostChainParams.BaseDenom, amountOfBaseDenom))

				//Mint autocompounding fee, use old cValue as we mint tokens at previous cValue.
				pstakeFeeAmount := hostChainParams.PstakeParams.PstakeRestakeFee.MulInt(amountOfBaseDenom)
				protocolFee, _ := k.ConvertTokenToStk(ctx, sdk.NewDecCoinFromDec(hostChainParams.BaseDenom, pstakeFeeAmount), cValue)

				err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(protocolFee))
				if err != nil {
					return "", types.ErrMintFailed
				}

				//Send protocol fee to protocol pool
				err = k.SendProtocolFee(ctx, sdk.NewCoins(protocolFee), types.ModuleName, hostChainParams.PstakeParams.PstakeFeeAddress)
				if err != nil {
					return "", types.ErrFailedDeposit
				}
			}
		}
		return msgResponse.String(), nil
	case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
		parsedMsg, ok := msg.(*stakingtypes.MsgUndelegate)
		if !ok {
			return "", errorsmod.Wrapf(sdkerrors.ErrInvalidType, "unable to unmarshal msg of type %s", sdk.MsgTypeURL(msg))
		}
		var msgResponse stakingtypes.MsgUndelegateResponse
		if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
			return "", errorsmod.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal undelegate response message: %s", err.Error())
		}
		k.Logger(ctx).Info(fmt.Sprintf("Started unbonding for val: %s, amount: %s", parsedMsg.ValidatorAddress, parsedMsg.Amount))
		//burn stkatom (DONE OUTSIDE THE LOOP), remove from delegations, add unbonding entry completion time
		err := k.SubtractHostAccountDelegation(ctx, types.NewHostAccountDelegation(parsedMsg.ValidatorAddress, parsedMsg.Amount))
		if err != nil {
			return "", err
		}

		return msgResponse.String(), nil
	case sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}):
		var msgResponse ibctransfertypes.MsgTransferResponse
		if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
			return "", errorsmod.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}
		k.Logger(ctx).Info(fmt.Sprintf("Initiated IBC transfer from %s to %s with msg: %s", hostChainParams.ChainID, ctx.ChainID(), msg))
		// handle rest in ibc hooks.
		return msgResponse.String(), nil

	default:
		return "", nil
	}
}

// resetToPreICATx is called when ICA execution fails
func (k Keeper) resetToPreICATx(ctx sdk.Context, icaPacket icatypes.InterchainAccountPacketData) error {
	hostChainParams := k.GetHostChainParams(ctx)

	msgs, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot Deserialise icapacket data: %v", err)
	}
	// Dispatch packet
	msgsCount := 0
	expectedMsgType := sdk.MsgTypeURL(msgs[0])
	for _, msg := range msgs {
		err := k.handleResetMsgs(ctx, msg, hostChainParams)
		if err != nil {
			return err
		}
		if expectedMsgType == sdk.MsgTypeURL(msg) {
			msgsCount++
		}
		// assert all msgs are of same type.
		if len(msgs) == msgsCount && expectedMsgType == sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}) {
			previousEpochNumber := types.PreviousUnbondingEpoch(k.epochKeeper.GetEpochInfo(ctx, types.UndelegationEpochIdentifier).CurrentEpoch)
			err := k.RemoveHostAccountUndelegation(ctx, previousEpochNumber)
			if err != nil {
				return err
			}
			k.FailUnbondingEpochCValue(ctx, previousEpochNumber, sdk.NewCoin(hostChainParams.MintDenom, sdk.ZeroInt()))
			k.Logger(ctx).Info(fmt.Sprintf("Failed unbonding msgs: %s, for undelegationEpoch: %v", msgs, previousEpochNumber))
		}

		k.Logger(ctx).Info("ICA msg timed out, ", "msg", msg)
	}
	if msgsCount != len(msgs) {
		k.SetModuleState(ctx, false) //Disable module, we assert single type of msg throughout the tx.
		k.Logger(ctx).Error(fmt.Sprintf("%s module has been disabled due to different msg types in a ica txn", types.ModuleName))
		return nil
	}
	return nil
}

// handleResetMsgs is a helper function for handling reset messages in resetToPreICATx
func (k Keeper) handleResetMsgs(ctx sdk.Context, msg sdk.Msg, _ types.HostChainParams) error {
	switch sdk.MsgTypeURL(msg) {
	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		parsedMsg, ok := msg.(*stakingtypes.MsgDelegate)
		if !ok {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "unable to unmarshal msg of type %s", sdk.MsgTypeURL(msg))
		}
		// Add to host-balance, because delegate txn timed out.
		k.AddBalanceToDelegationState(ctx, parsedMsg.Amount)
		k.RemoveICADelegateFromTransientStore(ctx, parsedMsg.Amount)
		return nil
	case sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}):
		parsedMsg, ok := msg.(*ibctransfertypes.MsgTransfer)
		if !ok {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "unable to unmarshal msg of type %s", sdk.MsgTypeURL(msg))
		}
		removedTransientUndelegationTransfer, err := k.RemoveUndelegationTransferFromTransientStore(ctx, parsedMsg.Token)
		if err != nil {
			ctx.Logger().Error("Failed to do ICA + IBC transfer from host chain to controller chain", "Err: ", err)
		}
		k.AddHostAccountUndelegation(ctx, types.HostAccountUndelegation{
			EpochNumber:             removedTransientUndelegationTransfer.EpochNumber,
			TotalUndelegationAmount: parsedMsg.Token,
			CompletionTime:          ctx.BlockTime(),
			UndelegationEntries:     nil,
		})

		return nil
	default:
		return nil
	}
}
