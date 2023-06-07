package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

func (k Keeper) Migrate(ctx sdk.Context) error {
	hcparams := k.GetHostChainParams(ctx)
	cValue := k.GetCValue(ctx)
	delegationState := k.GetDelegationState(ctx)
	hostChainRewardAddress := k.GetHostChainRewardAddress(ctx)
	hostAccounts := k.GetHostAccounts(ctx)

	// port all stores
	//	HostChainParamsKey

	consensusState, err := k.liquidStakeIBCKeeper.GetLatestConsensusState(ctx, hcparams.ConnectionID)
	if err != nil {
		k.Logger(ctx).Error("could not retrieve client state", "host_chain", hcparams.ChainID)
		return err
	}

	newhc := &liquidstakeibctypes.HostChain{
		ChainId:      hcparams.ChainID,
		ConnectionId: hcparams.ConnectionID,
		ChannelId:    hcparams.ChainID,
		PortId:       hcparams.TransferPort,
		Params: &liquidstakeibctypes.HostChainLSParams{
			DepositFee:    hcparams.PstakeParams.PstakeDepositFee,
			RestakeFee:    hcparams.PstakeParams.PstakeRestakeFee,
			UnstakeFee:    hcparams.PstakeParams.PstakeUnstakeFee,
			RedemptionFee: hcparams.PstakeParams.PstakeRedemptionFee,
		},
		HostDenom:       hcparams.BaseDenom,
		MinimumDeposit:  hcparams.MinDeposit,
		CValue:          cValue,
		NextValsetHash:  consensusState.NextValidatorsHash,
		UnbondingFactor: types.UndelegationEpochNumberFactor,
		Active:          false, // <- disable the module and update it with MsgUpdateHostChain
		DelegationAccount: &liquidstakeibctypes.ICAAccount{
			Address:      delegationState.HostChainDelegationAddress,
			Balance:      sdk.NewCoin(hcparams.BaseDenom, delegationState.HostDelegationAccountBalance.AmountOf(hcparams.BaseDenom)),
			Owner:        hostAccounts.DelegatorAccountOwnerID,
			ChannelState: liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATED, //add in created state, if closed, the begin blocker RecreateICA will check it out
		},
		RewardsAccount: &liquidstakeibctypes.ICAAccount{
			Address:      hostChainRewardAddress.Address,
			Balance:      sdk.NewInt64Coin(hcparams.BaseDenom, 0),
			Owner:        hostAccounts.RewardsAccountOwnerID,
			ChannelState: liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATED,
		},
	}

	// save the host chain
	k.liquidStakeIBCKeeper.SetHostChain(ctx, newhc)

	// for updating valset.
	err = k.liquidStakeIBCKeeper.QueryHostChainValidators(ctx, newhc, stakingtypes.QueryValidatorsRequest{})
	if err != nil {
		return errorsmod.Wrapf(
			liquidstakeibctypes.ErrFailedICQRequest,
			"error submitting validators icq: %s",
			err.Error(),
		)
	}

	//	AllowListedValidatorsKey <- will be set at MsgUpdateHostChain

	//	DelegationStateKey <- will be set by valSetQuery callback.

	//	HostChainRewardAddressKey <- already sent in newhc

	//	IBCTransientStoreKey
	ibcTransientStore := k.GetIBCTransientStore(ctx) //TODO need a better way.
	if len(ibcTransientStore.IBCTransfer) != 0 {
		return types.ErrModuleMigrationFailed.Wrapf("no transient ibc transfer balance expected, clear ibc packets and try again")
	}
	if ibcTransientStore.ICADelegate.Amount.IsPositive() {
		return types.ErrModuleMigrationFailed.Wrapf("no transient ICADelegate balance expected, clear ibc packets and try again")
	}
	if len(ibcTransientStore.UndelegatonCompleteIBCTransfer) != 0 {
		return types.ErrModuleMigrationFailed.Wrapf("no transient ica ibc transfer balance expected, clear ibc packets and try again")
	}

	//	UnbondingEpochCValueKey : -> Migrates to liquidstakeibctypes UnbondingKey
	// DO reject/ fail currentEpoch unstaking requests.
	err = k.FailCurrentUnbonding(ctx)
	if err != nil {
		return err
	}

	// DO auto claim of previous and delete the unbonding epoch.
	err = k.ClaimAll(ctx)
	if err != nil {
		return err
	}
	// Add pending epochValues to store. this should leave us with only pending to mature
	delegationState = k.GetDelegationState(ctx) //refresh delegation state after deletion of curr epoch entry.
	for _, validatorUnbondings := range delegationState.HostAccountUndelegations {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, validatorUnbondings.EpochNumber)
		if !unbondingEpochCValue.IsMatured && !unbondingEpochCValue.IsFailed {
			k.liquidStakeIBCKeeper.SetUnbonding(ctx, &liquidstakeibctypes.Unbonding{
				ChainId:       newhc.ChainId,
				EpochNumber:   validatorUnbondings.EpochNumber,
				MatureTime:    validatorUnbondings.CompletionTime,
				BurnAmount:    unbondingEpochCValue.STKBurn,
				UnbondAmount:  unbondingEpochCValue.AmountUnbonded,
				IbcSequenceId: "",
				State:         liquidstakeibctypes.Unbonding_UNBONDING_MATURING, //UNBONDING_MATURING
			})

		} else {
			return types.ErrModuleMigrationFailed.Wrapf("failed to migrate, epochNumber %v", validatorUnbondings.EpochNumber)
		}
	}

	//	DelegatorUnbondingEpochEntryKey
	delegatorUnbondingEpochEntries := k.IterateAllDelegatorUnbondingEpochEntry(ctx) //refresh undelegation user entries after claimall
	for _, delegatorUnbondingEpochEntry := range delegatorUnbondingEpochEntries {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, delegatorUnbondingEpochEntry.EpochNumber)
		claimableAmount := sdk.NewDecFromInt(delegatorUnbondingEpochEntry.Amount.Amount).Quo(unbondingEpochCValue.GetUnbondingEpochCValue())
		claimableCoin, _ := sdk.NewDecCoinFromDec(k.GetIBCDenom(ctx), claimableAmount).TruncateDecimal()

		k.liquidStakeIBCKeeper.SetUserUnbonding(ctx, &liquidstakeibctypes.UserUnbonding{
			ChainId:      newhc.ChainId,
			EpochNumber:  delegatorUnbondingEpochEntry.EpochNumber,
			Address:      delegatorUnbondingEpochEntry.DelegatorAddress,
			StkAmount:    delegatorUnbondingEpochEntry.Amount,
			UnbondAmount: claimableCoin,
		})
	}

	//	HostAccountsKey <- done in newhc
	//  ModuleEnableKey, disable module temporarily
	k.SetModuleState(ctx, false)

	// Migrate accounts
	// DepositModuleAccount
	depositAccBalance := k.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(types.DepositModuleAccount))
	//DELEGATION add to liquidstakeibc delegate
	currEpoch := k.epochKeeper.GetEpochInfo(ctx, types.DelegationEpochIdentifier)
	k.liquidStakeIBCKeeper.SetDeposit(ctx, &liquidstakeibctypes.Deposit{
		ChainId:       newhc.ChainId,
		Amount:        sdk.NewCoin(newhc.IBCDenom(), depositAccBalance.AmountOf(newhc.IBCDenom())),
		Epoch:         sdk.NewInt(currEpoch.CurrentEpoch),
		State:         liquidstakeibctypes.Deposit_DEPOSIT_PENDING,
		IbcSequenceId: "",
	})
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.DepositModuleAccount, liquidstakeibctypes.DepositModuleAccount, depositAccBalance)
	if err != nil {
		return err
	}
	// DelegationModuleAccount
	delegationAccBalance := k.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(types.DelegationModuleAccount))
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.DelegationModuleAccount, liquidstakeibctypes.DepositModuleAccount, delegationAccBalance)
	if err != nil {
		return err
	}
	// RewardModuleAccount
	rewardAccBalance := k.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(types.RewardModuleAccount))
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.RewardModuleAccount, liquidstakeibctypes.DepositModuleAccount, rewardAccBalance)
	if err != nil {
		return err
	}
	// UndelegationModuleAccount
	undelegationAccBalance := k.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(types.UndelegationModuleAccount))
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.UndelegationModuleAccount, liquidstakeibctypes.UndelegationModuleAccount, undelegationAccBalance)
	if err != nil {
		return err
	}
	// RewardBoosterModuleAccount
	rewardBoosterAccBalance := k.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(types.RewardBoosterModuleAccount))
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.RewardBoosterModuleAccount, liquidstakeibctypes.DepositModuleAccount, rewardBoosterAccBalance)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ClaimAll(ctx sdk.Context) error {
	// get all the entries corresponding to the delegator address
	delegatorUnbondingEntries := k.IterateAllDelegatorUnbondingEpochEntry(ctx)

	// loop through all the epoch and send tokens if an entry has matured.
	for _, unbondingEntry := range delegatorUnbondingEntries {
		delegatorAddress := sdk.MustAccAddressFromBech32(unbondingEntry.DelegatorAddress)
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, unbondingEntry.EpochNumber)
		if unbondingEpochCValue.IsMatured {
			// get c value from the UnbondingEpochCValue struct
			// calculate claimable amount from un inverse c value
			claimableAmount := sdk.NewDecFromInt(unbondingEntry.Amount.Amount).Quo(unbondingEpochCValue.GetUnbondingEpochCValue())

			// calculate claimable coin and community coin to be sent to delegator account and community pool respectively
			claimableCoin, _ := sdk.NewDecCoinFromDec(k.GetIBCDenom(ctx), claimableAmount).TruncateDecimal()

			// send coin to delegator address from undelegation module account
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.UndelegationModuleAccount, delegatorAddress, sdk.NewCoins(claimableCoin))
			if err != nil {
				return err
			}

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeClaim,
					sdk.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
					sdk.NewAttribute(types.AttributeAmount, unbondingEntry.Amount.String()),
					sdk.NewAttribute(types.AttributeClaimedAmount, claimableAmount.String()),
				)},
			)

			// remove entry from unbonding epoch entry
			k.RemoveDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEntry.EpochNumber)
		}
		if unbondingEpochCValue.IsFailed {
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.UndelegationModuleAccount, delegatorAddress, sdk.NewCoins(unbondingEntry.Amount))
			if err != nil {
				return err
			}

			// remove entry from unbonding epoch entry
			k.RemoveDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEntry.EpochNumber)
		}
	}
	return nil
}

func (k Keeper) FailCurrentUnbonding(ctx sdk.Context) error {
	currEpochInfo := k.epochKeeper.GetEpochInfo(ctx, types.UndelegationEpochIdentifier)
	currentUnbondingEpochNumber := types.CurrentUnbondingEpoch(currEpochInfo.CurrentEpoch)
	hostAccountUndelegationForEpoch, err := k.GetHostAccountUndelegationForEpoch(ctx, currentUnbondingEpochNumber)
	if err != nil {
		return err
	}
	err = k.RemoveHostAccountUndelegation(ctx, currentUnbondingEpochNumber)
	if err != nil {
		return err
	}
	k.FailUnbondingEpochCValue(ctx, currentUnbondingEpochNumber, hostAccountUndelegationForEpoch.TotalUndelegationAmount)
	return nil
}
