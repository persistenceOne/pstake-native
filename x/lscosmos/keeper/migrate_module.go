package keeper

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

func (k Keeper) Migrate(ctx sdk.Context) error {
	hcparams := k.GetHostChainParams(ctx)
	cValue := k.GetCValue(ctx)
	delegationState := k.GetDelegationState(ctx)
	hostChainRewardAddress := k.GetHostChainRewardAddress(ctx)
	hostAccounts := k.GetHostAccounts(ctx)
	allowlistedVals := k.GetAllowListedValidators(ctx)
	// port all stores
	//	HostChainParamsKey

	pending, err := k.CheckPendingICATxs(ctx)
	if pending {
		return types.ErrModuleMigrationFailed.Wrapf("There are pending ica txs/ channels closed, err: %s", err)
	}

	consensusState, err := k.liquidStakeIBCKeeper.GetLatestConsensusState(ctx, hcparams.ConnectionID)
	if err != nil {
		k.Logger(ctx).Error("could not retrieve client state", "host_chain", hcparams.ChainID)
		return err
	}

	// set validators
	var validators []*liquidstakeibctypes.Validator
	for _, delval := range delegationState.HostAccountDelegations {
		allowlistedVal := types.AllowListedValidator{
			ValidatorAddress: delval.ValidatorAddress,
			TargetWeight:     sdk.ZeroDec(),
		}
		for _, av := range allowlistedVals.AllowListedValidators {
			if delval.ValidatorAddress == av.ValidatorAddress {
				allowlistedVal = av
				continue
			}
		}
		validators = append(validators, &liquidstakeibctypes.Validator{
			OperatorAddress: delval.ValidatorAddress,
			Status:          stakingtypes.BondStatusBonded,
			Weight:          allowlistedVal.TargetWeight,
			DelegatedAmount: delval.Amount.Amount,
			TotalAmount:     sdk.OneInt(),
			DelegatorShares: sdk.OneDec(),
			UnbondingEpoch:  0,
		})
	}
	newhc := &liquidstakeibctypes.HostChain{
		ChainId:      hcparams.ChainID,
		ConnectionId: hcparams.ConnectionID,
		ChannelId:    hcparams.TransferChannel,
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
		Validators: validators,
	}

	// save the host chain
	k.liquidStakeIBCKeeper.SetHostChain(ctx, newhc)

	// for updating valset in liquidstakeibc.
	for _, v := range newhc.Validators {
		err = k.liquidStakeIBCKeeper.QueryHostChainValidator(ctx, newhc, v.OperatorAddress)
		if err != nil {
			return errorsmod.Wrapf(
				liquidstakeibctypes.ErrFailedICQRequest,
				"error submitting validators icq: %s",
				err.Error(),
			)
		}
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
		k.Logger(ctx).Info("ibcTransientStore ICADelegate Amount IsPositive ")
		msgSendToRewards := &banktypes.MsgSend{
			FromAddress: newhc.DelegationAccount.Address,
			ToAddress:   newhc.RewardsAccount.Address,
			Amount:      sdk.NewCoins(ibcTransientStore.ICADelegate),
		}
		_, err := k.liquidStakeIBCKeeper.GenerateAndExecuteICATx(
			ctx,
			newhc.ConnectionId,
			newhc.DelegationAccount.Owner,
			[]proto.Message{msgSendToRewards},
		)
		if err != nil {
			return err
		}
	}
	if len(ibcTransientStore.UndelegatonCompleteIBCTransfer) != 0 {
		return types.ErrModuleMigrationFailed.Wrapf("no transient ica ibc transfer balance expected, clear ibc packets and try again")
	}

	//	UnbondingEpochCValueKey : -> Migrates to liquidstakeibctypes UnbondingKey
	// DO reject/ fail currentEpoch unstaking requests, and previous that might have failed.
	err = k.FailNotStartedUnbondings(ctx)
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
		claimableCoin, _ := sdk.NewDecCoinFromDec(newhc.HostDenom, claimableAmount).TruncateDecimal()

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

	// set fee address and params to liquidstakeibc
	k.liquidStakeIBCKeeper.SetParams(ctx, liquidstakeibctypes.Params{
		AdminAddress:     hcparams.PstakeParams.PstakeFeeAddress,
		FeeAddress:       hcparams.PstakeParams.PstakeFeeAddress,
		UpperCValueLimit: sdk.MustNewDecFromStr(liquidstakeibctypes.DefaultUpperCValueLimit),
		LowerCValueLimit: sdk.MustNewDecFromStr(liquidstakeibctypes.DefaultLowerCValueLimit),
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

// FailNotStartedUnbondings fails the current unbondings and the previous for which no packet was received.
func (k Keeper) FailNotStartedUnbondings(ctx sdk.Context) error {
	delegationState := k.GetDelegationState(ctx)
	for _, ubd := range delegationState.HostAccountUndelegations {
		if ubd.CompletionTime.Equal(time.Time{}) {
			err := k.FailUnbonding(ctx, ubd.EpochNumber)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) FailUnbonding(ctx sdk.Context, epochnumber int64) error {
	hostAccountUndelegationForEpoch, err := k.GetHostAccountUndelegationForEpoch(ctx, epochnumber)
	if err != nil {
		return err
	}
	err = k.RemoveHostAccountUndelegation(ctx, epochnumber)
	if err != nil {
		return err
	}
	k.FailUnbondingEpochCValue(ctx, epochnumber, hostAccountUndelegationForEpoch.TotalUndelegationAmount)
	return nil
}
