package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/gogoproto/proto"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibclocalhosttypes "github.com/cosmos/ibc-go/v7/modules/light-clients/09-localhost"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	epochsKeeper        types.EpochsKeeper
	icaControllerKeeper types.ICAControllerKeeper
	ibcKeeper           *ibckeeper.Keeper
	ibcTransferKeeper   types.IBCTransferKeeper
	icqKeeper           types.ICQKeeper

	paramSpace paramtypes.Subspace

	msgRouter *baseapp.MsgServiceRouter

	hooks types.LiquidStakeIBCHooks

	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	epochsKeeper types.EpochsKeeper,
	icaControllerKeeper types.ICAControllerKeeper,
	ibcKeeper *ibckeeper.Keeper,
	ibcTransferKeeper types.IBCTransferKeeper,
	icqKeeper types.ICQKeeper,

	paramSpace paramtypes.Subspace,

	msgRouter *baseapp.MsgServiceRouter,

	authority string,
) Keeper {
	return Keeper{
		cdc:                 cdc,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		epochsKeeper:        epochsKeeper,
		icaControllerKeeper: icaControllerKeeper,
		ibcKeeper:           ibcKeeper,
		ibcTransferKeeper:   ibcTransferKeeper,
		icqKeeper:           icqKeeper,
		storeKey:            storeKey,
		paramSpace:          paramSpace,
		msgRouter:           msgRouter,
		hooks:               nil,
		authority:           authority,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams gets the total set of liquidstakeibc parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the total set of liquidstakeibc parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bytes := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bytes)
}

// GetDepositModuleAccount returns deposit module account interface
func (k *Keeper) GetDepositModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.DepositModuleAccount)
}

// GetUndelegationModuleAccount returns undelegation module account interface
func (k *Keeper) GetUndelegationModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.UndelegationModuleAccount)
}

// SendProtocolFee to the community pool
func (k *Keeper) SendProtocolFee(ctx sdk.Context, protocolFee sdk.Coins, moduleAccount, feeAddress string) error {
	addr, err := sdk.AccAddressFromBech32(feeAddress)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAccount, addr, protocolFee)
	if err != nil {
		return err
	}
	return nil
}

// GetClientState retrieves the client state given a connection id
func (k *Keeper) GetClientState(ctx sdk.Context, connectionID string) (exported.ClientState, error) {
	conn, found := k.ibcKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		return nil, fmt.Errorf("invalid connection id, \"%s\" not found", connectionID)
	}

	clientState, found := k.ibcKeeper.ClientKeeper.GetClientState(ctx, conn.ClientId)
	if !found {
		return nil, fmt.Errorf("client id \"%s\" not found for connection \"%s\"", conn.ClientId, connectionID)
	}

	return clientState, nil
}

// GetChainID gets the id of the host chain given a connection id
func (k *Keeper) GetChainID(ctx sdk.Context, connectionID string) (string, error) {
	clientState, err := k.GetClientState(ctx, connectionID)
	if err != nil {
		return "", fmt.Errorf("client state not found for connection \"%s\": \"%s\"", connectionID, err.Error())
	}

	switch clientType := clientState.(type) {
	case *ibctmtypes.ClientState:
		return clientType.ChainId, nil
	case *ibclocalhosttypes.ClientState:
		return ctx.ChainID(), nil
	default:
		return "", fmt.Errorf("unexpected type of client, cannot determine chain-id: clientType: %s, connectionid: %s", clientState.ClientType(), connectionID)
	}
}

// GetPortID constructs a port id given the port owner
func (k *Keeper) GetPortID(owner string) string {
	return fmt.Sprintf("%s%s", icatypes.ControllerPortPrefix, owner)
}

// RegisterICAAccount registers an ICA
func (k *Keeper) RegisterICAAccount(ctx sdk.Context, connectionID, owner string) error {
	return k.icaControllerKeeper.RegisterInterchainAccount(
		ctx,
		connectionID,
		owner,
		"",
	)
}

// SetWithdrawAddress sends a MsgSetWithdrawAddress to set the withdrawal address to the rewards account
func (k *Keeper) SetWithdrawAddress(ctx sdk.Context, hc *types.HostChain) error {
	msgSetWithdrawAddress := &distributiontypes.MsgSetWithdrawAddress{
		DelegatorAddress: hc.DelegationAccount.Address,
		WithdrawAddress:  hc.RewardsAccount.Address,
	}

	_, err := k.GenerateAndExecuteICATx(
		ctx,
		hc.ConnectionId,
		hc.DelegationAccount.Owner,
		[]proto.Message{msgSetWithdrawAddress},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsICAChannelActive checks if an ICA channel is active
func (k *Keeper) IsICAChannelActive(ctx sdk.Context, hc *types.HostChain, portID string) bool {
	_, isActive := k.icaControllerKeeper.GetOpenActiveChannel(ctx, hc.ConnectionId, portID)
	return isActive
}

func (k *Keeper) GetEpochNumber(ctx sdk.Context, epoch string) int64 {
	return k.epochsKeeper.GetEpochInfo(ctx, epoch).CurrentEpoch
}

func (k *Keeper) SendICATransfer(
	ctx sdk.Context,
	hc *types.HostChain,
	amount sdk.Coin,
	sender string,
	receiver string,
	portOwner string,
) (string, error) {
	channel, found := k.ibcKeeper.ChannelKeeper.GetChannel(ctx, hc.PortId, hc.ChannelId)
	if !found {
		return "", fmt.Errorf(
			"could not retrieve channel for host chain %s while sending ICA transfer",
			hc.ChainId,
		)
	}

	timeoutTimestamp := uint64(ctx.BlockTime().UnixNano() + (types.IBCTimeoutTimestamp).Nanoseconds())

	// prepare the msg transfer to bring the undelegation back
	msgTransfer := ibctransfertypes.NewMsgTransfer(
		channel.Counterparty.PortId,
		channel.Counterparty.ChannelId,
		amount,
		sender,
		receiver,
		clienttypes.ZeroHeight(),
		timeoutTimestamp,
		"",
	)

	// execute the transfers
	sequenceID, err := k.GenerateAndExecuteICATx(
		ctx,
		hc.ConnectionId,
		portOwner,
		[]proto.Message{msgTransfer},
	)
	if err != nil {
		return "", fmt.Errorf(
			"could not send ICA transfer for host chain %s",
			hc.ChainId,
		)
	}

	return sequenceID, nil
}

func (k *Keeper) UpdateCValues(ctx sdk.Context) {
	hostChains := k.GetAllHostChains(ctx)

	for _, hc := range hostChains {

		// total stk tokens minted
		mintedAmount := k.bankKeeper.GetSupply(ctx, hc.MintDenom()).Amount

		// total tokenized staked amount
		tokenizedStakedAmount := k.GetLSMDepositAmountUntokenized(ctx, hc.ChainId)

		// amount staked by the module in any of the validators of the host chain
		stakedAmount := hc.GetHostChainTotalDelegations()

		// amount that is in the staking flow and hasn't left Persistence yet
		amountOnPersistence := k.GetDepositAmountOnPersistence(ctx, hc.ChainId)

		// amount that is in the staking flow and has arrived to the host chain, but hasn't been staked yet
		amountOnHostChain := k.GetDepositAmountOnHostChain(ctx, hc.ChainId)

		// amount unbonded from a validator that has been in the Unbonding state for more than 4 unbonding epochs
		totalUnbondingAmount := k.GetAllValidatorUnbondedAmount(ctx, hc)

		// total amount staked
		liquidStakedAmount := tokenizedStakedAmount.Add(stakedAmount).Add(amountOnPersistence).Add(amountOnHostChain).Add(totalUnbondingAmount)

		var cValue sdk.Dec
		if mintedAmount.IsZero() || liquidStakedAmount.IsZero() {
			cValue = sdk.OneDec()
		} else {
			cValue = sdk.NewDecFromInt(mintedAmount).Quo(sdk.NewDecFromInt(liquidStakedAmount))
		}

		k.Logger(ctx).Info(
			fmt.Sprintf(
				"Updated CValue for %s. Total minted amount: %v. Total liquid staked amount: %v. Composed of %v staked tokens, %v tokens on Persistence, %v tokens on the host chain, %v tokens from a validator total unbonding. New c_value: %v - Old c_value: %v",
				hc.ChainId,
				mintedAmount,
				liquidStakedAmount,
				stakedAmount,
				amountOnPersistence,
				amountOnHostChain,
				totalUnbondingAmount,
				cValue,
				hc.CValue,
			),
		)

		hc.LastCValue = hc.CValue
		hc.CValue = cValue
		k.SetHostChain(ctx, hc)

		if err := k.Hooks().PostCValueUpdate(ctx, hc.MintDenom(), hc.HostDenom, hc.CValue); err != nil {
			k.Logger(ctx).Error("PostCValueUpdate hook failed with ", "err:", err)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCValueUpdate,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeModuleMintedAmount, sdk.NewCoin(hc.MintDenom(), mintedAmount).String()),
				sdk.NewAttribute(types.AttributeModuleLSMTokenizedAmount, sdk.NewCoin(hc.HostDenom, tokenizedStakedAmount).String()),
				sdk.NewAttribute(types.AttributeModuleStakedAmount, sdk.NewCoin(hc.HostDenom, stakedAmount).String()),
				sdk.NewAttribute(types.AttributeModuleAmountOnPersistence, sdk.NewCoin(hc.HostDenom, amountOnPersistence).String()),
				sdk.NewAttribute(types.AttributeModuleAmountOnHostChain, sdk.NewCoin(hc.HostDenom, amountOnHostChain).String()),
				sdk.NewAttribute(types.AttributeModuleUnbondingAmount, sdk.NewCoin(hc.HostDenom, totalUnbondingAmount).String()),
				sdk.NewAttribute(types.AttributeOldCValue, hc.LastCValue.String()),
				sdk.NewAttribute(types.AttributeNewCValue, hc.CValue.String()),
			),
		)

		// if the c value is out of bounds, disable the chain
		if !k.CValueWithinLimits(hc) {
			hc.Active = false
			k.SetHostChain(ctx, hc)

			k.Logger(ctx).Error(fmt.Sprintf("C value out of limits !!! Disabling chain %s with c value %v.", hc.ChainId, hc.CValue))
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeChainDisabled,
					sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
					sdk.NewAttribute(types.AttributeOldCValue, hc.LastCValue.String()),
					sdk.NewAttribute(types.AttributeNewCValue, hc.CValue.String()),
				),
			)

			continue
		}

		k.RecalculateCValueLimits(ctx, hc, mintedAmount, liquidStakedAmount)
	}
}

func (k *Keeper) CValueWithinLimits(hc *types.HostChain) bool {
	return hc.CValue.LT(hc.Params.UpperCValueLimit) && hc.CValue.GT(hc.Params.LowerCValueLimit)
}

func (k *Keeper) RecalculateCValueLimits(ctx sdk.Context, hc *types.HostChain, mintedAmount, lsAmount math.Int) {
	// if there has been no activity on the chain yet, leave the limits as they are
	if mintedAmount.IsZero() || lsAmount.IsZero() {
		return
	}

	// calculate twice the amount of rewards which will be transferred next rewards epoch
	diff := sdk.NewDecFromInt(hc.GetHostChainTotalDelegations()).
		Mul(hc.AutoCompoundFactor).Mul(sdk.NewDec(types.CValueDynamicLowerDiff))

	// calculate the new lower and upper limit
	newLowerLimit := sdk.NewDecFromInt(mintedAmount).Quo(sdk.NewDecFromInt(lsAmount).Add(diff))
	newUpperLimit := hc.CValue.Add(hc.CValue.Sub(newLowerLimit).Mul(sdk.NewDec(types.CValueDynamicUpperDiff)))

	// update the limits on the host chain
	hc.Params.LowerCValueLimit = newLowerLimit
	hc.Params.UpperCValueLimit = newUpperLimit
	k.SetHostChain(ctx, hc)

	k.Logger(ctx).Info(
		fmt.Sprintf("Updated C Value limits for %s. Current C Value %s, new lower limit %s, new upper limit %s",
			hc.ChainId,
			hc.CValue,
			hc.Params.LowerCValueLimit,
			hc.Params.UpperCValueLimit,
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCValueLimitsUpdated,
			sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
			sdk.NewAttribute(types.AttributeLowerLimit, hc.Params.LowerCValueLimit.String()),
			sdk.NewAttribute(types.AttributeUpperLimit, hc.Params.UpperCValueLimit.String()),
		),
	)
}

func (k *Keeper) CalculateAutocompoundLimit(autocompoundFactor sdk.Dec) sdk.Dec {
	return autocompoundFactor.Quo(sdk.NewDec(types.Percentage)).Quo(sdk.NewDec(types.DaysInYear))
}

// Hooks gets the hooks for liquidstakeibc *Keeper {
func (k *Keeper) Hooks() types.LiquidStakeIBCHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiLiquidStakeIBCHooks{}
	}

	return k.hooks
}

// Set the liquidstakeibc hooks
func (k *Keeper) SetHooks(hooks types.LiquidStakeIBCHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set hooks twice")
	}

	k.hooks = hooks

	return k
}
