package keeper

import (
	"context"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibcporttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
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

// LiquidStake defines a method for liquid staking tokens
func (m msgServer) LiquidStake(goCtx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// check if module is inactive or active
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	//GetParams
	hostChainParams := m.GetHostChainParams(ctx)

	//check for minimum deposit amount
	if msg.Amount.Amount.LT(hostChainParams.MinDeposit) {
		return nil, errorsmod.Wrapf(
			types.ErrMinDeposit, "expected amount more than %s, got %s", hostChainParams.MinDeposit, msg.Amount.Amount,
		)
	}

	expectedIBCPrefix := ibctransfertypes.GetDenomPrefix(hostChainParams.TransferPort, hostChainParams.TransferChannel)

	denomTraceStr, err := m.ibcTransferKeeper.DenomPathFromHash(ctx, msg.Amount.Denom)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "got error : %s", err)
	}
	denomTrace := ibctransfertypes.ParseDenomTrace(denomTraceStr)

	// Check if ibc path matches allowlisted path.
	if expectedIBCPrefix != denomTrace.GetPrefix() {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidDenomPath, "expected %s, got %s", expectedIBCPrefix, denomTrace.GetPrefix(),
		)
	}
	//Check if base denom is valid (uatom) , this can be programmed further to accommodate for liquid staked vouchers.
	if denomTrace.BaseDenom != hostChainParams.BaseDenom {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidDenom, "expected %s, got %s", hostChainParams.BaseDenom, denomTrace.BaseDenom,
		)
	}

	// get the delegator address from bech32 string
	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "got error : %s", err)
	}

	// amount of stk tokens to be minted. We calculate this before depositing any amount so as to not affect minting c-value.
	// We do not care about residue here because it won't be minted and bank.TotalSupply invariant should not be affected
	cValue := m.GetCValue(ctx)
	mintToken, _ := m.ConvertTokenToStk(ctx, sdktypes.NewDecCoinFromCoin(msg.Amount), cValue)

	//send the deposit to the deposit-module account
	depositAmount := sdktypes.NewCoins(msg.Amount)
	err = m.SendTokensToDepositModule(ctx, depositAmount, delegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrFailedDeposit, "failed to deposit tokens to module account %s, got error : %s", types.DepositModuleAccount, err,
		)
	}

	//Mint staked representative tokens in lscosmos module account
	err = m.bankKeeper.MintCoins(ctx, types.ModuleName, sdktypes.NewCoins(mintToken))
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrMintFailed, "failed to mint coins in module %s, got error %s", types.ModuleName, err,
		)
	}

	//Calculate protocol fee
	protocolFee := hostChainParams.PstakeParams.PstakeDepositFee
	protocolFeeAmount := protocolFee.MulInt(mintToken.Amount)
	// We do not care about residue, as to not break Total calculation invariant.
	protocolCoin, _ := sdktypes.NewDecCoinFromDec(hostChainParams.MintDenom, protocolFeeAmount).TruncateDecimal()

	//Send (mintedTokens - protocolTokens) to delegator address
	err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegatorAddress,
		sdktypes.NewCoins(mintToken.Sub(protocolCoin)))
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrMintFailed, "failed to send coins from module %s to account %s, got error : %s",
			types.ModuleName, delegatorAddress.String(), err,
		)
	}

	//Send protocol fee to protocol pool
	if protocolCoin.IsPositive() {
		err = m.SendProtocolFee(ctx, sdktypes.NewCoins(protocolCoin), types.ModuleName, hostChainParams.PstakeParams.PstakeFeeAddress)
		if err != nil {
			return nil, errorsmod.Wrapf(
				types.ErrFailedDeposit, "failed to send protocol fee to pstake fee address %s, got error : %s",
				hostChainParams.PstakeParams.PstakeFeeAddress, err,
			)
		}
	}
	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeLiquidStake,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmount, mintToken.String()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, mintToken.Sub(protocolCoin).String()),
			sdktypes.NewAttribute(types.AttributePstakeDepositFee, protocolCoin.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)
	return &types.MsgLiquidStakeResponse{}, nil
}

// LiquidUnstake defines a method for unstaking the liquid staked tokens
func (m msgServer) LiquidUnstake(goCtx context.Context, msg *types.MsgLiquidUnstake) (*types.MsgLiquidUnstakeResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// check if module is inactive or active
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	hostChainParams := m.GetHostChainParams(ctx)

	if msg.Amount.Denom != hostChainParams.MintDenom {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "Expected %s, got %s", hostChainParams.MintDenom, msg.Amount.Denom)
	}

	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	// take deposit into module acc
	err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, delegatorAddress, types.UndelegationModuleAccount, sdktypes.NewCoins(msg.Amount))
	if err != nil {
		return nil, err
	}
	// take pstake fees
	unstakeCoin := msg.Amount
	pstakeFeeAmt := hostChainParams.PstakeParams.PstakeUnstakeFee.MulInt(msg.Amount.Amount).TruncateInt()
	pstakeFee := sdktypes.NewCoin(msg.Amount.Denom, pstakeFeeAmt)
	if pstakeFeeAmt.IsPositive() {
		err = m.SendProtocolFee(ctx, sdktypes.NewCoins(pstakeFee), types.UndelegationModuleAccount, hostChainParams.PstakeParams.PstakeFeeAddress)
		if err != nil {
			return nil, err
		}
		unstakeCoin = msg.Amount.Sub(pstakeFee)
	}

	// Add entry to unbonding db
	epoch := m.epochKeeper.GetEpochInfo(ctx, types.UndelegationEpochIdentifier)
	unbondingEpochNumber := types.CurrentUnbondingEpoch(epoch.CurrentEpoch)
	m.AddDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEpochNumber, unstakeCoin)
	m.AddTotalUndelegationForEpoch(ctx, unbondingEpochNumber, unstakeCoin)

	// check is there are delegations worth the amount to be undelegated.
	// there are chances where the delegation epoch is not yet done so stkAtom are more than delegated amount
	// in this case users should just redeem tokens. (as tokens should be present as part of deposit tokens)
	delegationState := m.GetDelegationState(ctx)
	undelegations, err := m.GetHostAccountUndelegationForEpoch(ctx, unbondingEpochNumber)
	if err != nil {
		return nil, err
	}
	totalDelegations := delegationState.TotalDelegations(hostChainParams.BaseDenom)
	baseDenomUndelegations, _ := m.ConvertStkToToken(ctx, sdktypes.NewDecCoinFromCoin(undelegations.TotalUndelegationAmount), m.GetCValue(ctx))
	if totalDelegations.IsLT(sdktypes.NewCoin(hostChainParams.BaseDenom, baseDenomUndelegations.Amount)) {
		return nil, errorsmod.Wrapf(types.ErrHostChainDelegationsLTUndelegations, "Delegated amount: %s is less than total undelegations for the epoch: %s", totalDelegations, undelegations.TotalUndelegationAmount)
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeLiquidUnstake,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, msg.GetDelegatorAddress()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, msg.Amount.String()),
			sdktypes.NewAttribute(types.AttributePstakeUnstakeFee, pstakeFee.String()),
			sdktypes.NewAttribute(types.AttributeUnstakeAmount, unstakeCoin.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.GetDelegatorAddress()),
		)},
	)
	return &types.MsgLiquidUnstakeResponse{}, nil
}

// Redeem defines a method for redeeming liquid staked tokens instantly
func (m msgServer) Redeem(goCtx context.Context, msg *types.MsgRedeem) (*types.MsgRedeemResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// check if module is inactive or active
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	// take redeem address from msg address string
	redeemAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "got error : %s", err)
	}

	// get the ibc denom and host chain params
	ibcDenom := m.GetIBCDenom(ctx)
	hostChainParams := m.GetHostChainParams(ctx)

	// check msg amount denom
	if msg.Amount.Denom != hostChainParams.MintDenom {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "expected %s, got %s", hostChainParams.BaseDenom, msg.Amount.Denom)
	}

	// We do not care about residue, as to not break Total calculation invariant.
	// protocolCoin is the redemption fee
	protocolCoin, _ := sdktypes.NewDecCoinFromDec(
		hostChainParams.MintDenom,
		hostChainParams.PstakeParams.PstakeRedemptionFee.MulInt(msg.Amount.Amount),
	).TruncateDecimal()

	// send redeem tokens to module account from redeem account
	err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, redeemAddress, types.ModuleName, sdktypes.NewCoins(msg.Amount))
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrMintFailed, "failed to send coins from account %s to module %s, got error : %s",
			redeemAddress.String(), types.ModuleName, err,
		)
	}

	// send protocol fee to protocol pool
	if protocolCoin.IsPositive() {
		err = m.SendProtocolFee(ctx, sdktypes.NewCoins(protocolCoin), types.ModuleName, hostChainParams.PstakeParams.PstakeFeeAddress)
		if err != nil {
			return nil, errorsmod.Wrapf(
				types.ErrFailedDeposit, "failed to send protocol fee to pstake fee address %s, got error : %s",
				hostChainParams.PstakeParams.PstakeFeeAddress, err,
			)
		}
	}
	// convert redeem amount to ibc/allow-listed-denom amount (sub protocolCoin) based on the current c-value
	redeemStk := msg.Amount.Sub(protocolCoin)
	redeemToken, _ := m.ConvertStkToToken(ctx, sdktypes.NewDecCoinFromCoin(redeemStk), m.GetCValue(ctx))

	// get all deposit account balances
	allDepositBalances := m.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(types.DepositModuleAccount))
	delegationBalance := sdktypes.NewCoin(ibcDenom, allDepositBalances.AmountOf(ibcDenom))

	// check deposit account has sufficient funds
	if redeemToken.IsGTE(delegationBalance) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "expected tokens under %s, got %s for redeem", delegationBalance.String(), redeemToken.String())
	}

	// send the ibc/Denom token from module to the account
	err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.DepositModuleAccount, redeemAddress, sdktypes.NewCoins(redeemToken))
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrMintFailed, "failed to send coins from module %s to account %s, got error : %s",
			types.DepositModuleAccount, redeemAddress.String(), err,
		)
	}

	// burn the redeemStk token
	err = m.bankKeeper.BurnCoins(ctx, types.ModuleName, sdktypes.NewCoins(redeemStk))
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrBurnFailed, "failed to burn coins from module %s, got error %s", types.ModuleName, err,
		)
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeRedeem,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, redeemAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, redeemToken.String()),
			sdktypes.NewAttribute(types.AttributePstakeRedeemFee, protocolCoin.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)
	return &types.MsgRedeemResponse{}, nil
}

// Claim defines a method for claiming unstaked mature tokens or failed unbondings
func (m msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// check if module is inactive or active
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	// get AccAddress from bech32 string
	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	// get all the entries corresponding to the delegator address
	delegatorUnbondingEntries := m.IterateDelegatorUnbondingEpochEntry(ctx, delegatorAddress)

	// loop through all the epoch and send tokens if an entry has matured.
	for _, unbondingEntry := range delegatorUnbondingEntries {
		unbondingEpochCValue := m.GetUnbondingEpochCValue(ctx, unbondingEntry.EpochNumber)
		if unbondingEpochCValue.IsMatured {
			// get c value from the UnbondingEpochCValue struct
			// calculate claimable amount from un inverse c value
			claimableAmount := sdktypes.NewDecFromInt(unbondingEntry.Amount.Amount).Quo(unbondingEpochCValue.GetUnbondingEpochCValue())

			// calculate claimable coin and community coin to be sent to delegator account and community pool respectively
			claimableCoin, _ := sdktypes.NewDecCoinFromDec(m.GetIBCDenom(ctx), claimableAmount).TruncateDecimal()

			// send coin to delegator address from undelegation module account
			err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.UndelegationModuleAccount, delegatorAddress, sdktypes.NewCoins(claimableCoin))
			if err != nil {
				return nil, err
			}

			ctx.EventManager().EmitEvents(sdktypes.Events{
				sdktypes.NewEvent(
					types.EventTypeClaim,
					sdktypes.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
					sdktypes.NewAttribute(types.AttributeAmount, unbondingEntry.Amount.String()),
					sdktypes.NewAttribute(types.AttributeClaimedAmount, claimableAmount.String()),
				)},
			)

			// remove entry from unbonding epoch entry
			m.RemoveDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEntry.EpochNumber)
		}
		if unbondingEpochCValue.IsFailed {
			err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.UndelegationModuleAccount, delegatorAddress, sdktypes.NewCoins(unbondingEntry.Amount))
			if err != nil {
				return nil, err
			}

			// remove entry from unbonding epoch entry
			m.RemoveDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEntry.EpochNumber)
		}
	}

	// emit event
	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)

	return &types.MsgClaimResponse{}, nil
}

// JumpStart defines a method for jump-starting the module through fee address account.
func (m msgServer) JumpStart(goCtx context.Context, msg *types.MsgJumpStart) (*types.MsgJumpStartResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// check pstake fee address == from addr
	hostChainParams := m.GetHostChainParams(ctx)
	if msg.PstakeAddress != hostChainParams.PstakeParams.PstakeFeeAddress {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("msg.pstakeAddress should be equal to msg.PstakeParams.PstakeFeeAddress, got %s expected %s", msg.PstakeAddress, hostChainParams.PstakeParams.PstakeFeeAddress))
	}

	// check module disabled
	if m.GetModuleState(ctx) {
		return nil, types.ErrModuleAlreadyEnabled
	}
	// ensure if params were set before never allow the denoms to be reset/changed, so tvu can always be checked
	if !hostChainParams.IsEmpty() && hostChainParams.BaseDenom != msg.BaseDenom {
		return nil, types.ErrInvalidDenom
	}
	if types.ConvertBaseDenomToMintDenom(msg.BaseDenom) != msg.MintDenom {
		return nil, types.ErrInvalidMintDenom
	}
	// check if there is any mints.
	stkSupply := m.bankKeeper.GetSupply(ctx, msg.MintDenom)
	if stkSupply.Amount.GT(sdktypes.ZeroInt()) {
		return nil, errorsmod.Wrap(types.ErrModuleAlreadyEnabled, "Module cannot be reset once via admin once it has positive delegations.")
	}
	// reset db, no need to release capability
	m.SetDelegationState(ctx, types.DelegationState{})
	m.SetHostChainRewardAddress(ctx, types.HostChainRewardAddress{})
	if err := msg.HostAccounts.Validate(); err != nil {
		return nil, err
	}
	m.SetHostAccounts(ctx, msg.HostAccounts)

	// check fees limits
	if err := msg.PstakeParams.Validate(); err != nil {
		return nil, err
	}

	if msg.MinDeposit.LTE(sdktypes.ZeroInt()) {
		return nil, errorsmod.Wrap(types.ErrInvalidDeposit, "MinDeposit should be GT 0")
	}
	// do proposal things
	if msg.TransferPort != ibctransfertypes.PortID {
		return nil, errorsmod.Wrap(ibcporttypes.ErrInvalidPort, "Only acceptable TransferPort is \"transfer\"")
	}

	// checks for valid and active channel
	channel, found := m.channelKeeper.GetChannel(ctx, msg.TransferPort, msg.TransferChannel)
	if !found {
		return nil, errorsmod.Wrap(ibcchanneltypes.ErrChannelNotFound, fmt.Sprintf("channel for ibc transfer: %s not found", msg.TransferChannel))
	}
	if channel.State != ibcchanneltypes.OPEN {
		return nil, errorsmod.Wrapf(
			ibcchanneltypes.ErrInvalidChannelState,
			"channel state is not OPEN (got %s)", channel.State.String(),
		)
	}
	// TODO Understand capabilities and see if it has to be/ should be claimed in lsscopedkeeper. If it even matters.
	_, err := m.lscosmosScopedKeeper.NewCapability(ctx, host.ChannelCapabilityPath(msg.TransferPort, msg.TransferChannel))
	if err != nil {
		ctx.Logger().Info(fmt.Sprintf("err: %s, Capability already exists", err.Error()))
	}

	hostAccounts := m.GetHostAccounts(ctx)
	// This checks for channel being active
	err = m.icaControllerKeeper.RegisterInterchainAccount(ctx, msg.ConnectionID, hostAccounts.DelegatorAccountOwnerID, "")
	if err != nil {
		return nil, errorsmod.Wrap(err, "Could not register ica delegation Address")
	}

	newHostChainParams := types.NewHostChainParams(msg.ChainID, msg.ConnectionID, msg.TransferChannel,
		msg.TransferPort, msg.BaseDenom, msg.MintDenom, msg.PstakeParams.PstakeFeeAddress,
		msg.MinDeposit, msg.PstakeParams.PstakeDepositFee, msg.PstakeParams.PstakeRestakeFee,
		msg.PstakeParams.PstakeUnstakeFee, msg.PstakeParams.PstakeRedemptionFee)

	m.SetHostChainParams(ctx, newHostChainParams)
	m.SetAllowListedValidators(ctx, msg.AllowListedValidators)

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeJumpStart,
			sdktypes.NewAttribute(types.AttributePstakeAddress, msg.PstakeAddress),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.PstakeAddress),
		)},
	)

	return &types.MsgJumpStartResponse{}, nil
}

// RecreateICA defines a method for recreating closed ica channels
func (m msgServer) RecreateICA(goCtx context.Context, msg *types.MsgRecreateICA) (*types.MsgRecreateICAResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// check if module is inactive or active
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	hostAccounts := m.Keeper.GetHostAccounts(ctx)
	hostChainParams := m.Keeper.GetHostChainParams(ctx)

	msgAttributes := []sdktypes.Attribute{sdktypes.NewAttribute(types.AttributeFromAddress, msg.FromAddress)}

	_, ok := m.icaControllerKeeper.GetOpenActiveChannel(ctx, hostChainParams.ConnectionID, hostAccounts.DelegatorAccountPortID())
	if !ok {
		err := m.icaControllerKeeper.RegisterInterchainAccount(ctx, hostChainParams.ConnectionID, hostAccounts.DelegatorAccountOwnerID, "")
		if err != nil {
			return nil, errorsmod.Wrap(err, "Could not register ica delegation Address")
		}
		msgAttributes = append(msgAttributes, sdktypes.NewAttribute(types.AttributeRecreateDelegationICA, hostAccounts.DelegatorAccountPortID()))
	}
	_, ok = m.icaControllerKeeper.GetOpenActiveChannel(ctx, hostChainParams.ConnectionID, hostAccounts.RewardsAccountPortID())
	if !ok {
		err := m.icaControllerKeeper.RegisterInterchainAccount(ctx, hostChainParams.ConnectionID, hostAccounts.RewardsAccountOwnerID, "")
		if err != nil {
			return nil, errorsmod.Wrap(err, "Could not register ica reward Address")
		}
		msgAttributes = append(msgAttributes, sdktypes.NewAttribute(types.AttributeRecreateRewardsICA, hostAccounts.RewardsAccountPortID()))
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeRecreateICA,
			msgAttributes...,
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.FromAddress),
		)},
	)

	return &types.MsgRecreateICAResponse{}, nil

}

// ChangeModuleState defines an admin method for disabling or re-enabling module state
func (m msgServer) ChangeModuleState(goCtx context.Context, msg *types.MsgChangeModuleState) (*types.MsgChangeModuleStateResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	hostChainParams := m.Keeper.GetHostChainParams(ctx)
	if hostChainParams.IsEmpty() {
		return nil, types.ErrModuleNotInitialised
	}
	if hostChainParams.PstakeParams.PstakeFeeAddress != msg.PstakeAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("Only admin address is allowed to call this method, current admin address: %s", hostChainParams.PstakeParams.PstakeFeeAddress))
	}
	moduleState := m.Keeper.GetModuleState(ctx)
	if moduleState == msg.ModuleState {
		return nil, sdkerrors.Wrap(types.ErrModuleNotInitialised, fmt.Sprintf("currentState: %v", moduleState))
	}
	m.Keeper.SetModuleState(ctx, msg.ModuleState)

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeChangeModuleState,
			sdktypes.NewAttribute(types.AttributeChangedModuleState, strconv.FormatBool(msg.ModuleState)),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.PstakeAddress),
		)},
	)
	return &types.MsgChangeModuleStateResponse{}, nil

}

// ReportSlashing defines an admin method for reporting slashing on a validator
func (m msgServer) ReportSlashing(goCtx context.Context, msg *types.MsgReportSlashing) (*types.MsgReportSlashingResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// check if module is inactive or active
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	hostChainParams := m.Keeper.GetHostChainParams(ctx)
	if hostChainParams.IsEmpty() {
		return nil, types.ErrModuleNotInitialised
	}
	if hostChainParams.PstakeParams.PstakeFeeAddress != msg.PstakeAddress {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "Only admin address is allowed to call this method, current admin address: %s", hostChainParams.PstakeParams.PstakeFeeAddress)
	}

	delegationState := m.Keeper.GetDelegationState(ctx)
	_, hostAccountDelegations := m.Keeper.GetAllValidatorsState(ctx, hostChainParams.BaseDenom)
	exists := false
	for _, val := range hostAccountDelegations {
		if val.ValidatorAddress == msg.ValidatorAddress && val.Amount.IsPositive() {
			exists = true
			break
		}
	}
	if !exists {
		return nil, errorsmod.Wrapf(types.ErrNoHostChainDelegations, "No delegation found for validator: %s", msg.ValidatorAddress)
	}

	pending, err := m.Keeper.CheckPendingICATxs(ctx)
	if pending {
		return nil, err
	}

	delegationRequest := stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delegationState.HostChainDelegationAddress,
		ValidatorAddr: msg.ValidatorAddress,
	}
	bz, err := m.cdc.Marshal(&delegationRequest)
	if err != nil {
		return nil, errorsmod.Wrap(err, "Failed to Marshal delegationRequest")
	}
	m.Keeper.icqKeeper.MakeRequest(ctx, hostChainParams.ConnectionID, hostChainParams.ChainID, "cosmos.staking.v1beta1.Query/Delegation",
		bz, sdktypes.NewInt(int64(-1)), types.ModuleName, Delegation, 0)

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeReportSlashing,
			sdktypes.NewAttribute(types.AttributeValidatorAddress, msg.ValidatorAddress),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.PstakeAddress),
		)},
	)
	return &types.MsgReportSlashingResponse{}, nil
}
