package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
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
) error {

	// Require portID is the portID module is bound to
	if portID != types.DelegationAccountPortID &&
		portID != types.RewardAccountPortID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected either of %s or %s",
			portID, types.DelegationAccountPortID, types.RewardAccountPortID)
	}
	var versionData icatypes.Metadata
	if err := icatypes.ModuleCdc.UnmarshalJSON([]byte(version), &versionData); err != nil {
		return err
	}
	if versionData.Version != icatypes.Version {
		return sdkerrors.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", versionData.Version, icatypes.Version)
	}

	// Claim channel capability passed back by IBC module
	if err := k.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return err
	}

	return nil
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
	if portID != types.DelegationAccountPortID &&
		portID != types.RewardAccountPortID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected either of %s or %s",
			portID, types.DelegationAccountPortID, types.RewardAccountPortID)
	}

	var counterpartyVersionData icatypes.Metadata
	if err := icatypes.ModuleCdc.UnmarshalJSON([]byte(counterpartyVersion), &counterpartyVersionData); err != nil {
		return err
	}

	if counterpartyVersionData.Version != icatypes.Version {
		return sdkerrors.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: %s, expected %s", counterpartyVersion, types.Version)
	}
	//TODO more checks, capability, channelID??

	hostChainParams := k.GetHostChainParams(ctx)

	if portID == types.DelegationAccountPortID {
		delegationAddress, delegationAddrfound := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, hostChainParams.ConnectionID, types.DelegationAccountPortID)
		if delegationAddrfound {
			if err := k.SetHostChainDelegationAddress(ctx, delegationAddress); err != nil {
				return err
			}
			if err := k.icaControllerKeeper.RegisterInterchainAccount(ctx, hostChainParams.ConnectionID, types.RewardModuleAccount); err != nil {
				return sdkerrors.Wrap(err, "Could not register ica reward Address")
			}

		}
	}
	if portID == types.RewardAccountPortID {
		rewardAddress, rewardAddrFound := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, hostChainParams.ConnectionID, types.RewardAccountPortID)
		delegationAddress := k.GetDelegationState(ctx).HostChainDelegationAddress
		if rewardAddrFound {
			_ = k.SetHostChainRewardAddressIfEmpty(ctx, types.NewHostChainRewardAddress(rewardAddress))
			setWithdrawAddrMsg := &distributiontypes.MsgSetWithdrawAddress{
				DelegatorAddress: delegationAddress,
				WithdrawAddress:  rewardAddress,
			}
			err := k.GenerateAndExecuteICATx(ctx, hostChainParams.ConnectionID, types.DelegationAccountPortID, []sdk.Msg{setWithdrawAddrMsg})
			if err != nil {
				return err
			}
		}
	}
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
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
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
	_, ok := k.lscosmosScopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(modulePacket.GetSourcePort(), modulePacket.GetSourceChannel()))
	if !ok {
		return sdkerrors.Wrapf(capabilitytypes.ErrCapabilityNotOwned, "capability not found for port: %s channel: %s in module: %s", modulePacket.GetSourcePort(), modulePacket.GetSourceChannel(), types.ModuleName)
	}

	// TODO add checks for capabilities, ports, channels
	hostChainParams := k.GetHostChainParams(ctx)

	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet acknowledgement: %v", err)
	}
	if !ack.Success() {
		return sdkerrors.Wrapf(channeltypes.ErrInvalidAcknowledgement, "acknowledgement failed")
	}
	// this line is used by starport scaffolding # oracle/packet/module/ack
	txMsgData := &sdk.TxMsgData{}
	if err := k.cdc.Unmarshal(ack.GetResult(), txMsgData); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	icaPacket := &icatypes.InterchainAccountPacketData{}
	if err := icatypes.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), icaPacket); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}
	msgs, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot Deserialise icapacket data: %v", err)
	}
	var eventType string

	// Dispatch packet
	switch len(txMsgData.Data) {
	case 0:
		// TODO: handle for sdk 0.46.x
		return nil
	default:
		for i, msgData := range txMsgData.Data {
			response, err := k.handleAckMsgData(ctx, msgData, msgs[i], hostChainParams)
			if err != nil {
				return err
			}
			k.Logger(ctx).Info("message response in ICS-27 packet response", "response", response)

			// if the packet has withdrawrewards msgs
			if i == 0 && msgData.MsgType == sdk.MsgTypeURL(&distributiontypes.MsgWithdrawDelegatorReward{}) {
				rewardAddr := k.GetHostChainRewardAddress(ctx)

				balanceQuery := banktypes.QueryBalanceRequest{Address: rewardAddr.Address, Denom: hostChainParams.BaseDenom}
				bz, err := k.cdc.Marshal(&balanceQuery)
				if err != nil {
					return err
				}

				// total rewards balance withdrawn
				k.icqKeeper.MakeRequest(
					ctx,
					hostChainParams.ConnectionID,
					hostChainParams.ChainID,
					"cosmos.bank.v1beta1.Query/Balance",
					bz,
					sdk.NewInt(int64(-1)),
					types.ModuleName,
					RewardsAccountBalance,
					0,
				)
			}
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAck, fmt.Sprintf("%v", ack)),
		),
	)

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				eventType,
				sdk.NewAttribute(types.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				eventType,
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
	_, ok := k.lscosmosScopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(modulePacket.GetSourcePort(), modulePacket.GetSourceChannel()))
	if !ok {
		return sdkerrors.Wrapf(capabilitytypes.ErrCapabilityNotOwned, "capability not found for port: %s channel: %s in module: %s", modulePacket.GetSourcePort(), modulePacket.GetSourceChannel(), types.ModuleName)
	}

	// TODO add checks for capabilities, ports, channels
	hostChainParams := k.GetHostChainParams(ctx)

	icaPacket := &icatypes.InterchainAccountPacketData{}
	if err := icatypes.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), icaPacket); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}
	msgs, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot Deserialise icapacket data: %v", err)
	}
	// Dispatch packet
	switch len(icaPacket.Data) {
	case 0:
		// TODO: handle for sdk 0.46.x
		return nil
	default:
		for _, msg := range msgs {
			response, err := k.handleTimeoutMsgData(ctx, msg, hostChainParams)
			if err != nil {
				return err
			}

			k.Logger(ctx).Info("message response in ICS-27 packet response", "response", response)
		}
	}

	return nil
}

func (k Keeper) handleAckMsgData(ctx sdk.Context, msgData *sdk.MsgData, msg sdk.Msg, hostChainParams types.HostChainParams) (string, error) {
	switch msgData.MsgType {
	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		parsedMsg, ok := msg.(*stakingtypes.MsgDelegate)
		if !ok {
			return "", sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "unable to unmarshal msg of type %s", msgData.MsgType)
		}
		var msgResponse stakingtypes.MsgDelegateResponse
		if err := k.cdc.Unmarshal(msgData.Data, &msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}
		// Add delegation state
		k.AddHostAccountDelegation(ctx, types.NewHostAccountDelegation(parsedMsg.ValidatorAddress, parsedMsg.Amount))

		return msgResponse.String(), nil

	case sdk.MsgTypeURL(&distributiontypes.MsgSetWithdrawAddress{}):
		var msgResponse distributiontypes.MsgSetWithdrawAddressResponse
		if err := k.cdc.Unmarshal(msgData.Data, &msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}
		k.SetModuleState(ctx, true)
		return msgResponse.String(), nil
	case sdk.MsgTypeURL(&distributiontypes.MsgWithdrawDelegatorReward{}):
		var msgResponse distributiontypes.MsgWithdrawDelegatorRewardResponse
		if err := k.cdc.Unmarshal(msgData.Data, &msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}
		return msgResponse.String(), nil
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		parsedMsg, ok := msg.(*banktypes.MsgSend)
		if !ok {
			return "", sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "unable to unmarshal msg of type %s", msgData.MsgType)
		}
		var msgResponse banktypes.MsgSendResponse
		if err := k.cdc.Unmarshal(msgData.Data, &msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}
		//is from rewardaddr to delegationaddr?
		rewardAddress := k.GetHostChainRewardAddress(ctx)
		delegationState := k.GetDelegationState(ctx)
		if rewardAddress.Address == parsedMsg.FromAddress && delegationState.HostChainDelegationAddress == parsedMsg.ToAddress {
			amountOfBaseDenom := parsedMsg.Amount.AmountOf(hostChainParams.BaseDenom)
			if amountOfBaseDenom.GT(sdk.ZeroInt()) {
				k.AddBalanceToDelegationState(ctx, sdk.NewCoin(hostChainParams.BaseDenom, amountOfBaseDenom))

				//Mint autocompounding fee
				pstakeFeeAmount := hostChainParams.PstakeRestakeFee.MulInt(amountOfBaseDenom).Mul(k.GetCValue(ctx))
				protocolFee, _ := sdk.NewDecCoinFromDec(hostChainParams.MintDenom, pstakeFeeAmount).TruncateDecimal()

				err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(protocolFee))
				if err != nil {
					return "", types.ErrMintFailed
				}

				//Send protocol fee to protocol pool
				err = k.SendProtocolFee(ctx, sdk.NewCoins(protocolFee), types.ModuleName, hostChainParams.PstakeFeeAddress)
				if err != nil {
					return "", types.ErrFailedDeposit
				}
			}
		}
		return msgResponse.String(), nil
	default:
		return "", nil
	}
}

func (k Keeper) handleTimeoutMsgData(ctx sdk.Context, msg sdk.Msg, _ types.HostChainParams) (string, error) {
	switch sdk.MsgTypeURL(msg) {
	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		parsedMsg, ok := msg.(*stakingtypes.MsgDelegate)
		if !ok {
			return "", sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "unable to unmarshal msg of type %s", sdk.MsgTypeURL(msg))
		}
		// Add to host-balance, because delegate txn timed out.
		k.AddBalanceToDelegationState(ctx, parsedMsg.Amount)
		return msg.String(), nil
	default:
		return "", nil
	}
}
