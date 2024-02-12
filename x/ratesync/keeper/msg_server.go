package keeper

import (
	"context"
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateHostChain(goCtx context.Context, msg *types.MsgCreateHostChain) (*types.MsgCreateHostChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	// Checks if the msg creator is the same as the current owner
	if msg.Authority != k.authority && msg.Authority != params.Admin {
		return nil, errorsmod.Wrapf(sdkerrors.ErrorInvalidSigner, "tx signer is not a module authority")
	}

	// get the host chain id
	chainID, err := k.GetChainID(ctx, msg.HostChain.ConnectionID)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "chain id not found for connection \"%s\": \"%s\"", msg.HostChain.ConnectionID, err)
	}
	if chainID != msg.HostChain.ChainID {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidChainID, "chain id does not match connection-chainID input \"%s\": found\"%s\"", msg.HostChain.ChainID, chainID)
	}

	id := k.IncrementHostChainID(ctx)
	msg.HostChain.ID = id

	if msg.HostChain.ICAAccount.Owner == "" {
		msg.HostChain.ICAAccount.Owner = types.DefaultPortOwner(id)
	} // else handled in msg.ValidateBasic()
	// register ratesyn ICA
	if msg.HostChain.ICAAccount.ChannelState == liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATING {
		err = k.icaControllerKeeper.RegisterInterchainAccount(ctx, msg.HostChain.ConnectionID, msg.HostChain.ICAAccount.Owner, "")
		if err != nil {
			return nil, errorsmod.Wrapf(
				types.ErrRegisterFailed,
				"error registering %s ratesync ica with owner: %s, err:%s",
				chainID, msg.HostChain.ICAAccount.Owner,
				err.Error(),
			)
		}
	} // else handled in validate basic (not allowed to create new host chain with previous ICA as portID is default and suffixed by ID

	channel, found := k.ibcKeeper.ChannelKeeper.GetChannel(ctx, msg.HostChain.TransferPortID, msg.HostChain.TransferChannelID)
	if !found || channel.State != channeltypes.OPEN {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrNotFound,
			"error creating %s ratesync with channel: %s, port: %s",
			chainID, msg.HostChain.TransferChannelID, msg.HostChain.TransferPortID,
		)
	}

	k.SetHostChain(
		ctx,
		msg.HostChain,
	)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateHostChain,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
			sdk.NewAttribute(types.AttributeChainID, msg.HostChain.ChainID),
			sdk.NewAttribute(types.AttributeConnectionID, msg.HostChain.ConnectionID),
			sdk.NewAttribute(types.AttributeID, fmt.Sprintf("%v", id)),
		),
	})
	return &types.MsgCreateHostChainResponse{ID: id}, nil
}

func (k msgServer) UpdateHostChain(goCtx context.Context, msg *types.MsgUpdateHostChain) (*types.MsgUpdateHostChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	// Checks if the msg creator is the same as the current owner
	if msg.Authority != k.authority && msg.Authority != params.Admin {
		return nil, errorsmod.Wrapf(sdkerrors.ErrorInvalidSigner, "tx signer is not a module authority")
	}

	// Check if the value exists
	oldHC, isFound := k.GetHostChain(
		ctx,
		msg.HostChain.ID,
	)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "id not set, hostchain does not exist")
	}

	// only allow enable disable feature && instantiate.
	// to change chain-id etc, add delete and create new hostchain with same details
	if msg.HostChain.ConnectionID != oldHC.ConnectionID {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "invalid connectionID, connectionID cannot be updated, "+
			"connectionID mismatch got %s, found %s", msg.HostChain.ConnectionID, oldHC.ConnectionID)
	}

	if msg.HostChain.TransferChannelID != oldHC.TransferChannelID {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "invalid channelID, channelID cannot be updated, "+
			"channelID mismatch got %s, found %s", msg.HostChain.TransferChannelID, oldHC.TransferChannelID)
	}

	if msg.HostChain.TransferPortID != oldHC.TransferPortID {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "invalid portID, portID cannot be updated, "+
			"portID mismatch got %s, found %s", msg.HostChain.TransferPortID, oldHC.TransferPortID)
	}
	channel, found := k.ibcKeeper.ChannelKeeper.GetChannel(ctx, oldHC.TransferPortID, oldHC.TransferChannelID)
	if !found || channel.State != channeltypes.OPEN {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrNotFound,
			"error creating %s ratesync with channel: %s, port: %s",
			oldHC.ChainID, msg.HostChain.TransferChannelID, msg.HostChain.TransferPortID,
		)
	}
	if oldHC.ICAAccount.ChannelState != liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATED {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "invalid ICAAccount state, should already be active")
	}
	if msg.HostChain.ICAAccount.ChannelState != oldHC.ICAAccount.ChannelState ||
		msg.HostChain.ICAAccount.Address != oldHC.ICAAccount.Address ||
		msg.HostChain.ICAAccount.Owner != oldHC.ICAAccount.Owner ||
		!msg.HostChain.ICAAccount.Balance.IsEqual(oldHC.ICAAccount.Balance) {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "invalid ICAAccount, ICA account cannot be updated, "+
			"ICAAccount mismatch got %s, found %s", msg.HostChain.ICAAccount, oldHC.ICAAccount)
	}

	updateStr := ""
	isOneUpdated := false
	saveUpdate := func(updates string) (bool, string) {
		return true, updates
	}

	chainID, err := k.GetChainID(ctx, msg.HostChain.ConnectionID)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "chain id not found for connection \"%s\": \"%s\"", msg.HostChain.ConnectionID, err)
	}
	if chainID != msg.HostChain.ChainID {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidChainID, "chain id does not match connection-chainID input \"%s\": found\"%s\"", msg.HostChain.ChainID, chainID)
	}
	if msg.HostChain.ChainID != oldHC.ChainID {
		oldHC.ChainID = msg.HostChain.ChainID
		isOneUpdated, updateStr = saveUpdate(fmt.Sprintf("updates host chain chainID %v to %v \n", oldHC.ChainID, msg.HostChain.ChainID))
	}

	// allow only one feature update per tx.
	if !isOneUpdated && !msg.HostChain.Features.LiquidStakeIBC.Equals(oldHC.Features.LiquidStakeIBC) {
		if oldHC.Features.LiquidStakeIBC.Instantiation == types.InstantiationState_INSTANTIATION_NOT_INITIATED {
			// allow to add details and instantiate or just save if trying to recover.
			switch msg.HostChain.Features.LiquidStakeIBC.Instantiation {
			case types.InstantiationState_INSTANTIATION_NOT_INITIATED:
				// just update oldhc, validate basic will take care of mismatch states.
				oldHC.Features.LiquidStakeIBC = msg.HostChain.Features.LiquidStakeIBC
			case types.InstantiationState_INSTANTIATION_INITIATED:
				// update oldhc, generate and execute wasm instantiate
				oldHC.Features.LiquidStakeIBC = msg.HostChain.Features.LiquidStakeIBC

				err := k.InstantiateLiquidStakeContract(ctx, oldHC.ICAAccount, oldHC.Features.LiquidStakeIBC, oldHC.ID, oldHC.ConnectionID, channel.Counterparty.ChannelId, channel.Counterparty.PortId)
				if err != nil {
					return nil, err
				}
			case types.InstantiationState_INSTANTIATION_COMPLETED:
				// just update oldhc, validate basic will take care of mismatch states.
				oldHC.Features.LiquidStakeIBC = msg.HostChain.Features.LiquidStakeIBC
			}
		}
		if !slices.Equal(oldHC.Features.LiquidStakeIBC.Denoms, msg.HostChain.Features.LiquidStakeIBC.Denoms) {
			oldHC.Features.LiquidStakeIBC.Denoms = msg.HostChain.Features.LiquidStakeIBC.Denoms
		}
		isOneUpdated, updateStr = saveUpdate(fmt.Sprintf("updates LiquidStakeIBC feature from %v to %v \n", oldHC.Features.LiquidStakeIBC, msg.HostChain.Features.LiquidStakeIBC))
	}
	if !isOneUpdated && !msg.HostChain.Features.LiquidStake.Equals(oldHC.Features.LiquidStake) {
		if oldHC.Features.LiquidStake.Instantiation == types.InstantiationState_INSTANTIATION_NOT_INITIATED {
			// allow to add details and instantiate or just save if trying to recover.
			switch msg.HostChain.Features.LiquidStake.Instantiation {
			case types.InstantiationState_INSTANTIATION_NOT_INITIATED:
				// just update oldhc, validate basic will take care of mismatch states.
				oldHC.Features.LiquidStake = msg.HostChain.Features.LiquidStake
			case types.InstantiationState_INSTANTIATION_INITIATED:
				// update oldhc, generate and execute wasm instantiate
				oldHC.Features.LiquidStake = msg.HostChain.Features.LiquidStake

				err := k.InstantiateLiquidStakeContract(ctx, oldHC.ICAAccount, oldHC.Features.LiquidStake, oldHC.ID, oldHC.ConnectionID, channel.Counterparty.ChannelId, channel.Counterparty.PortId)
				if err != nil {
					return nil, err
				}
			case types.InstantiationState_INSTANTIATION_COMPLETED:
				// just update oldhc, validate basic will take care of mismatch states.
				oldHC.Features.LiquidStake = msg.HostChain.Features.LiquidStake
			}
		}
		if !slices.Equal(oldHC.Features.LiquidStake.Denoms, msg.HostChain.Features.LiquidStake.Denoms) {
			oldHC.Features.LiquidStake.Denoms = msg.HostChain.Features.LiquidStake.Denoms
		}
		//nolint: ineffassign,staticcheck // it will be required if more features are added.
		isOneUpdated, updateStr = saveUpdate(fmt.Sprintf("updates LiquidStake feature from %v to %v", oldHC.Features.LiquidStake, msg.HostChain.Features.LiquidStake))
	}
	err = oldHC.Features.ValdidateBasic()
	if err != nil {
		return nil, err
	}

	k.SetHostChain(ctx, oldHC)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateHostChain,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
			sdk.NewAttribute(types.AttributeChainID, oldHC.ChainID),
			sdk.NewAttribute(types.AttributeConnectionID, oldHC.ConnectionID),
			sdk.NewAttribute(types.AttributeID, fmt.Sprintf("%v", oldHC.ID)),
			sdk.NewAttribute(types.AttributeUpdates, updateStr),
		),
	})
	return &types.MsgUpdateHostChainResponse{}, nil
}

func (k msgServer) DeleteHostChain(goCtx context.Context, msg *types.MsgDeleteHostChain) (*types.MsgDeleteHostChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	// Checks if the msg creator is the same as the current owner
	if msg.Authority != k.authority && msg.Authority != params.Admin {
		return nil, errorsmod.Wrapf(sdkerrors.ErrorInvalidSigner, "tx signer is not a module authority")
	}

	// Check if the value exists
	hc, isFound := k.GetHostChain(
		ctx,
		msg.ID,
	)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "id not set")
	}

	// check pending packets, do not allow to delete if packets are pending.
	portID := types.MustICAPortIDFromOwner(hc.ICAAccount.Owner)
	channelID, ok := k.icaControllerKeeper.GetOpenActiveChannel(ctx, hc.ConnectionID, portID)
	if !ok {
		return nil, errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "PortID: %s, connectionID: %s", portID, hc.ConnectionID)
	}
	nextSendSeq, ok := k.ibcKeeper.ChannelKeeper.GetNextSequenceSend(ctx, portID, channelID)
	if !ok {
		return nil, errorsmod.Wrapf(channeltypes.ErrSequenceSendNotFound, "PortID: %s, channelID: %s", portID, channelID)
	}
	nextAckSeq, ok := k.ibcKeeper.ChannelKeeper.GetNextSequenceAck(ctx, portID, channelID)
	if !ok {
		return nil, errorsmod.Wrapf(channeltypes.ErrSequenceAckNotFound, "PortID: %s, channelID: %s", portID, channelID)
	}
	if nextSendSeq != nextAckSeq {
		return nil, errorsmod.Wrapf(channeltypes.ErrPacketSequenceOutOfOrder, "PortID: %s, channelID: %s, NextSendSequence: %v, NextAckSequence: %v", portID, channelID, nextSendSeq, nextAckSeq)
	}

	k.RemoveHostChain(
		ctx,
		msg.ID,
	)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeleteHostChain,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
			sdk.NewAttribute(types.AttributeChainID, hc.ChainID),
			sdk.NewAttribute(types.AttributeConnectionID, hc.ConnectionID),
			sdk.NewAttribute(types.AttributeID, fmt.Sprintf("%v", hc.ID)),
		),
	})

	return &types.MsgDeleteHostChainResponse{}, nil
}

// UpdateParams defines a method for updating the module params
func (k msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)

	// authority needs to be either the gov module account (for proposals)
	// or the module admin account (for normal txs)
	if msg.Authority != k.authority && msg.Authority != params.Admin {
		return nil, errorsmod.Wrapf(sdkerrors.ErrorInvalidSigner, "tx signer is not a module authority")
	}

	k.SetParams(ctx, msg.Params)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateParams,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
			sdk.NewAttribute(types.AttributeKeyUpdatedParams, msg.Params.String()),
		),
	})

	return &types.MsgUpdateParamsResponse{}, nil
}
