package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	KeyConnectionId   string = "connection_id"
	KeyProtocolDenom  string = "prot_denom"
	KeyBaseDenom      string = "base_denom"
	KeyMinimumDeposit string = "min_deposit"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the liquidstakeibc MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RegisterHostChain adds a new host chain to the protocol
func (k msgServer) RegisterHostChain(
	goCtx context.Context,
	msg *types.MsgRegisterHostChain,
) (*types.MsgRegisterHostChainResponse, error) {
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", k.authority, msg.Authority)
	}

	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// Get the host chain client state
	conn, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, msg.ConnectionId)
	if !found {
		return nil, fmt.Errorf("invalid connection id, \"%s\" not found", msg.ConnectionId)
	}
	clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, conn.ClientId)
	if !found {
		return nil, fmt.Errorf(
			"client id \"%s\" not found for connection \"%s\"",
			conn.ClientId,
			msg.ConnectionId,
		)
	}
	client, ok := clientState.(*ibctmtypes.ClientState)
	if !ok {
		return nil, fmt.Errorf(
			"invalid client state for client \"%s\" on connection \"%s\"",
			conn.ClientId,
			msg.ConnectionId,
		)
	}

	// Check if host chain is already registered
	_, found = k.GetHostChain(sdktypes.UnwrapSDKContext(ctx), client.ChainId)
	if found {
		return nil, fmt.Errorf("invalid chain id \"%s\", host chain already registered", client.ChainId)
	}

	hs := &types.HostChain{
		ChainId:        client.ChainId,
		ConnectionId:   msg.ConnectionId,
		LocalDenom:     msg.LocalDenom,
		HostDenom:      msg.HostDenom,
		MinimumDeposit: msg.MinimumDeposit,
	}

	k.SetHostChain(ctx, hs)

	return &types.MsgRegisterHostChainResponse{}, nil
}

// UpdateHostChain updates a registered host chain
func (k msgServer) UpdateHostChain(
	goCtx context.Context,
	msg *types.MsgUpdateHostChain,
) (*types.MsgUpdateHostChainResponse, error) {
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", k.authority, msg.Authority)
	}

	ctx := sdktypes.UnwrapSDKContext(goCtx)

	hs, found := k.GetHostChain(ctx, msg.ChainId)
	if !found {
		return nil, fmt.Errorf("invalid chain id \"%s\", host chain is not registered", msg.ChainId)
	}

	for _, update := range msg.Updates {
		switch update.Key {
		case KeyConnectionId:
			// TODO: Update connection + re-create ICA
		case KeyProtocolDenom:
			hs.HostDenom = update.Value
			if err := sdktypes.ValidateDenom(update.Value); err != nil {
				return nil, err
			}
		case KeyBaseDenom:
			hs.LocalDenom = update.Value
			if err := sdktypes.ValidateDenom(update.Value); err != nil {
				return nil, err
			}
		case KeyMinimumDeposit:
			minimumDeposit, ok := sdktypes.NewIntFromString(update.Value)
			if !ok {
				return nil, fmt.Errorf("unable to parse string to sdk.Int")
			}
			hs.MinimumDeposit = minimumDeposit
			if minimumDeposit.LT(sdktypes.NewInt(0)) {
				return nil, fmt.Errorf("invalid minimum deposit value less than zero")
			}
		}
	}

	k.SetHostChain(ctx, &hs)

	return &types.MsgUpdateHostChainResponse{}, nil
}

// LiquidStake defines a method for liquid staking tokens
func (k msgServer) LiquidStake(
	goCtx context.Context,
	msg *types.MsgLiquidStake,
) (*types.MsgLiquidStakeResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// TODO: Check if module is active
	//// check if module is inactive or active
	//if !k.GetModuleState(ctx) {
	//	return nil, types.ErrModuleDisabled
	//}

	// retrieve the host chain
	hostChain, found := k.GetHostChainFromLocalDenom(ctx, msg.Amount.Denom)
	if !found {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidHostChain,
			"host chain with ibc denom %s not registered",
			msg.Amount.Denom,
		)
	}

	// check for minimum deposit amount
	if msg.Amount.Amount.LT(hostChain.MinimumDeposit) {
		return nil, errorsmod.Wrapf(
			types.ErrMinDeposit,
			"expected amount more than %s, got %s",
			hostChain.MinimumDeposit,
			msg.Amount.Amount,
		)
	}

	// TODO: Get IBC prefix from the IBC connection and compare it with what has been provided.
	//expectedIBCPrefix := ibctransfertypes.GetDenomPrefix(hostChainParams.TransferPort, hostChainParams.TransferChannel)
	//
	//denomTraceStr, err := k.ibcTransferKeeper.DenomPathFromHash(ctx, msg.Amount.Denom)
	//if err != nil {
	//	return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "got error : %s", err)
	//}
	//denomTrace := ibctransfertypes.ParseDenomTrace(denomTraceStr)
	//
	//// Check if ibc path matches allowlisted path.
	//if expectedIBCPrefix != denomTrace.GetPrefix() {
	//	return nil, errorsmod.Wrapf(
	//		types.ErrInvalidDenomPath, "expected %s, got %s", expectedIBCPrefix, denomTrace.GetPrefix(),
	//	)
	//}
	////Check if base denom is valid (uatom) , this can be programmed further to accommodate for liquid staked vouchers.
	//if denomTrace.BaseDenom != hostChainParams.BaseDenom {
	//	return nil, errorsmod.Wrapf(
	//		types.ErrInvalidDenom, "expected %s, got %s", hostChainParams.BaseDenom, denomTrace.BaseDenom,
	//	)
	//}

	// get the delegator address from the bech32 string
	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "error parsing delegator address: %s", err)
	}

	// amount of stk tokens to be minted
	mintDenom := "stk" + hostChain.HostDenom
	mintAmount := sdktypes.NewDecCoinFromCoin(msg.Amount).Amount.Mul(hostChain.CValue) // TODO: CValue needs to be recalculated and saved on every LS/LU/Redeem/Slash
	mintToken, _ := sdktypes.NewDecCoinFromDec(mintDenom, mintAmount).TruncateDecimal()

	// send the deposit to the deposit-module account
	depositAmount := sdktypes.NewCoins(msg.Amount)
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, delegatorAddress, types.DepositModuleAccount, depositAmount)
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrFailedDeposit,
			"failed to deposit tokens to module account %s: %s",
			types.DepositModuleAccount,
			err,
		)
	}

	// mint stk tokens in the module account
	err = k.BankKeeper.MintCoins(ctx, types.ModuleName, sdktypes.NewCoins(mintToken))
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrMintFailed,
			"failed to mint coins in module %s: %s",
			types.ModuleName, err,
		)
	}

	// retrieve the module params
	params := k.GetParams(ctx)

	// calculate protocol fee
	protocolFeeAmount := params.DepositFee.MulInt(mintToken.Amount)
	protocolFee, _ := sdktypes.NewDecCoinFromDec(mintDenom, protocolFeeAmount).TruncateDecimal()

	// send stk tokens to the delegator address
	err = k.BankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		delegatorAddress,
		sdktypes.NewCoins(mintToken.Sub(protocolFee)),
	)
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrMintFailed,
			"failed to send coins from module %s to account %s: %s",
			types.ModuleName,
			delegatorAddress.String(),
			err,
		)
	}

	// send the protocol fee to the protocol pool
	if protocolFee.IsPositive() {
		err = k.SendProtocolFee(ctx, sdktypes.NewCoins(protocolFee), types.ModuleName, params.FeeAddress)
		if err != nil {
			return nil, errorsmod.Wrapf(
				types.ErrFailedDeposit,
				"failed to send protocol fee to pStake fee address %s: %s",
				params.FeeAddress,
				err,
			)
		}
	}
	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeLiquidStake,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmount, mintToken.String()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, mintToken.Sub(protocolFee).String()),
			sdktypes.NewAttribute(types.AttributePstakeDepositFee, protocolFee.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)
	return &types.MsgLiquidStakeResponse{}, nil
}
