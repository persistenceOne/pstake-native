package keeper

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// TODO change
const HALT_HEIGHT = int64(12080000)
const CHAIN_ID = "core-1"
const PROTOCOL_ACC = "persistence12d7ett36q9vmtzztudt48f9rtyxlayflz5gun3" //TODO REPLACE

var FAILED_EPOCHS = []int64{164, 216, 228}
var FORCE_TO_FAIL_EPOCHS = []int64{216}
var DELETED_EPOCHS = []int64{232}

var UNBONDING_CREATION_HEIGHT_EPOCHS = map[int64]int64{232: 15826532, 236: 15881865}

func (k Keeper) Fork(ctx sdk.Context) error {
	hostChainParams := k.GetHostChainParams(ctx)

	// get the missing undelegation data for epoch 236 by querying cosmoshub chain via rest endpoints
	missingUndelegations := ParseHostAccountUnbondings(hostChainParams.MintDenom, hostChainParams.BaseDenom)

	// re-create deleted undelegations using queried data and mark them as not failed
	for _, epoch := range DELETED_EPOCHS {
		// mark the undelegation record as not failed
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, epoch)
		unbondingEpochCValue.IsFailed = false
		k.SetUnbondingEpochCValue(ctx, unbondingEpochCValue)

		// safe check to make sure there is no existing undelegation
		// this will never trigger, but let's be extra cautious here
		_, err := k.GetHostAccountUndelegationForEpoch(ctx, epoch)
		if err == nil {
			continue
		}

		// retrieve the missing undelegation from the queried data
		undelegation, ok := missingUndelegations[UNBONDING_CREATION_HEIGHT_EPOCHS[epoch]]
		if !ok {
			k.Logger(ctx).Error("Could not find missing undelegation from queried data for epoch %s", strconv.Itoa(int(epoch)))
		}

		// update the undelegation total amount to stkATOM as the amount returned from the query is in ATOM
		undelegation.TotalUndelegationAmount = unbondingEpochCValue.STKBurn
		k.AddHostAccountUndelegation(ctx, undelegation)
	}

	// packets for these epochs have already been relayed, which means they need to be manually be failed.
	for _, epoch := range FORCE_TO_FAIL_EPOCHS {
		err := k.RemoveHostAccountUndelegation(ctx, epoch)
		if err != nil {
			return err
		}

		k.FailUnbondingEpochCValue(ctx, epoch, sdk.NewCoin(hostChainParams.MintDenom, sdk.ZeroInt()))
	}

	// calculate the amount of stkATOM that will be needed for the claims
	stkAtomExpectedInUnbondingAcc := k.calculateTotalUndelegatingStkAtomAmount(ctx)

	// retrieve the balance of the undelegation module account
	currentStkAtomBalance := k.bankKeeper.GetBalance(
		ctx,
		authtypes.NewModuleAddress(types.UndelegationModuleAccount),
		hostChainParams.MintDenom,
	)

	// if the needed balance is greater than the balance in the account, it needs to be funded
	if stkAtomExpectedInUnbondingAcc.GT(currentStkAtomBalance.Amount) {
		// calculate the amount to be funded
		diff := stkAtomExpectedInUnbondingAcc.Sub(currentStkAtomBalance.Amount)

		// fund the undelegation account using the protocol account
		protocolAddr := sdk.MustAccAddressFromBech32(PROTOCOL_ACC)
		stkBalanceWithProtocol := k.bankKeeper.GetBalance(ctx, protocolAddr, hostChainParams.MintDenom)
		if stkBalanceWithProtocol.Amount.GT(diff) {
			err := k.bankKeeper.SendCoinsFromAccountToModule(
				ctx,
				protocolAddr,
				types.UndelegationModuleAccount,
				sdk.NewCoins(sdk.NewCoin(hostChainParams.MintDenom, diff)), // TODO: Would it be wise to add a small margin ?
			)
			if err != nil {
				return err
			}
		} else {
			panic(fmt.Sprintf("Protocol doesn't have enought stkATOM to fill deposits"))
		}

		for _, epoch := range DELETED_EPOCHS {
			unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, epoch)

			// find deficit of stkATOMs from the delegator unbonding entries
			totalUnstaked := k.getTotalDelegatorStkAmountUndelegatingForEpoch(ctx, epoch)
			if unbondingEpochCValue.STKBurn.Amount.GT(totalUnstaked) {
				deficitForThisEpoch := unbondingEpochCValue.STKBurn.Amount.Sub(totalUnstaked)

				// create a new record with the delegator address as the protocol account and difference as the amount.
				k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
					DelegatorAddress: PROTOCOL_ACC,
					EpochNumber:      epoch,
					Amount:           sdk.NewCoin(hostChainParams.MintDenom, deficitForThisEpoch),
				})

				diff = diff.Sub(deficitForThisEpoch)
			}
		}
	}

	return nil
}

func (k Keeper) getTotalDelegatorStkAmountUndelegatingForEpoch(ctx sdk.Context, epochNumber int64) sdk.Int {
	totalStkAmount := sdk.ZeroInt()
	for _, delegatorUnbondingEntry := range k.IterateAllDelegatorUnbondingEpochEntry(ctx) {
		if delegatorUnbondingEntry.EpochNumber == epochNumber {
			totalStkAmount = totalStkAmount.Add(delegatorUnbondingEntry.Amount.Amount)
		}
	}
	return totalStkAmount
}

func (k Keeper) calculateTotalUndelegatingStkAtomAmount(ctx sdk.Context) sdk.Int {
	// get the current unbonding epoch
	currUnbondingEpoch := types.CurrentUnbondingEpoch(
		k.epochKeeper.GetEpochInfo(ctx, types.UndelegationEpochIdentifier).CurrentEpoch,
	)

	allUnbondingEpochCValues := k.IterateAllUnbondingEpochCValues(ctx)

	totalStkAtomToBeInUnbondingAcc := sdk.ZeroInt()
	for _, delegatorUnbondingEntry := range k.IterateAllDelegatorUnbondingEpochEntry(ctx) {
		// the current epoch has no unbondingEpochCvalue, so just add their stkATOMs
		if delegatorUnbondingEntry.EpochNumber == currUnbondingEpoch {
			totalStkAtomToBeInUnbondingAcc = totalStkAtomToBeInUnbondingAcc.Add(delegatorUnbondingEntry.Amount.Amount)
			continue
		}

		// retrieve the unbonding c value for the delegator unbonding entry
		unbondingEpochCValue, found := findUnbondingCValue(allUnbondingEpochCValues, delegatorUnbondingEntry.EpochNumber)
		if !found {
			// should never occur as per the code logic, but if not found, just continue.
			k.Logger(ctx).Error(
				"Unbonding epoch not found",
				"epoch_number",
				delegatorUnbondingEntry.EpochNumber,
				"delegator_unbonding_entry",
				delegatorUnbondingEntry,
			)
			continue
		}

		// if the unbonding c value is matured, the claim will be in ATOM, so nothing to do
		if unbondingEpochCValue.IsMatured {
			continue
		}

		// get the host account unbonding for the epoch
		hostAccUnbonding, found := findHostAccUnbonding(
			k.GetDelegationState(ctx).HostAccountUndelegations,
			delegatorUnbondingEntry.EpochNumber,
		)

		// filter the unbonding requests that are in the maturing state but have no maturing time set
		// this should filter out [164, 216, 228] epochs
		if !unbondingEpochCValue.IsFailed && found && hostAccUnbonding.CompletionTime.Equal(time.Time{}) {
			totalStkAtomToBeInUnbondingAcc = totalStkAtomToBeInUnbondingAcc.Add(delegatorUnbondingEntry.Amount.Amount)
			continue
		}

		// if the unbonding has failed, the stkATOMs need to be claimed, so we can add the amount
		// NOTE: the re-creation of the deleted epochs needs to happen before this
		if unbondingEpochCValue.IsFailed {
			totalStkAtomToBeInUnbondingAcc = totalStkAtomToBeInUnbondingAcc.Add(delegatorUnbondingEntry.Amount.Amount)
			continue
		}
	}

	return totalStkAtomToBeInUnbondingAcc
}

func findUnbondingCValue(
	allUnbondingCValues []types.UnbondingEpochCValue,
	epoch int64,
) (types.UnbondingEpochCValue, bool) {
	for _, unbondingEpochCValue := range allUnbondingCValues {
		if epoch == unbondingEpochCValue.EpochNumber {
			return unbondingEpochCValue, true
		}
	}

	return types.UnbondingEpochCValue{}, false
}

func findHostAccUnbonding(
	allHostAccUndelegations []types.HostAccountUndelegation,
	epoch int64,
) (types.HostAccountUndelegation, bool) {
	for _, hostAccUndelgation := range allHostAccUndelegations {
		if epoch == hostAccUndelgation.EpochNumber {
			return hostAccUndelgation, true
		}
	}

	return types.HostAccountUndelegation{}, false
}

func (k Keeper) GetUnstakingEpochForPacket(ctx sdk.Context, packetMsgs []sdk.Msg) (int64, error) {
	if len(packetMsgs) == 0 {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrLogic, "packet should have more than 0 messages")
	}

	delegationState := k.GetDelegationState(ctx)
	unbondings := delegationState.HostAccountUndelegations
	unbondingEpochFound := make([]int64, len(packetMsgs))
PacketLoop:
	for i, packetMsg := range packetMsgs {
		undelegateMsg, ok := packetMsg.(*stakingtypes.MsgUndelegate)
		if !ok {
			return 0, sdkerrors.Wrapf(
				sdkerrors.ErrInvalidType,
				"unable to unmarshal msg of type %s, expected type MsgUndelegate",
				sdk.MsgTypeURL(packetMsg),
			)
		}
		for _, unbonding := range unbondings {
			if unbonding.CompletionTime.Equal(time.Time{}) {
				for _, ubd := range unbonding.UndelegationEntries {
					if undelegateMsg.ValidatorAddress == ubd.ValidatorAddress &&
						undelegateMsg.Amount.IsEqual(ubd.Amount) {
						unbondingEpochFound[i] = unbonding.EpochNumber
						continue PacketLoop
					}
				}
			}
		}
	}
	allElementsSame := func(a []int64) bool {
		for i := 1; i < len(a); i++ {
			if a[i] != a[0] {
				return false
			}
		}
		return true
	}
	if allElementsSame(unbondingEpochFound) {
		return unbondingEpochFound[0], nil
	}
	return 0, sdkerrors.ErrNotFound.Wrapf("Unstaking epoch not found for packet msgs %v", packetMsgs)
}

func (k Keeper) ClaimFailed(ctx sdk.Context, unbondingEntry types.DelegatorUnbondingEpochEntry) error {
	delegatorAddress := sdk.MustAccAddressFromBech32(unbondingEntry.DelegatorAddress)
	unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, unbondingEntry.EpochNumber)

	if unbondingEpochCValue.IsFailed {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.UndelegationModuleAccount,
			delegatorAddress,
			sdk.NewCoins(unbondingEntry.Amount),
		)
		if err != nil {
			return err
		}

		// remove entry from unbonding epoch entry
		k.RemoveDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEntry.EpochNumber)
	}

	return nil
}

// to fetch the data from the cosmoshub chain:
// curl -X GET -H "Content-Type: application/json" -H "x-cosmos-block-height: 15884400" "https://rest.cosmos.audit.one/cosmos/staking/v1beta1/delegators/cosmos13t4996czrgft9gw43epuwauccrldu5whx6uprjdmvsmuf7ylg8yqcxgzk3/unbonding_delegations" -o undelegations.json
// creationheight 15826532 => epoch 232, 15881865 => 236

type UnbondingEntry struct {
	CreationHeight          string    `json:"creation_height"`
	CompletionTime          time.Time `json:"completion_time"`
	InitialBalance          string    `json:"initial_balance"`
	Balance                 string    `json:"balance"`
	UnbondingId             string    `json:"unbonding_id"`
	UnbondingOnHoldRefCount string    `json:"unbonding_on_hold_ref_count"`
}

type UnbondingResponse struct {
	DelegatorAddress string           `json:"delegator_address"`
	ValidatorAddress string           `json:"validator_address"`
	Entries          []UnbondingEntry `json:"entries"`
}

type Unbondings struct {
	UnbondingResponses []UnbondingResponse `json:"unbonding_responses"`
}

func ParseHostAccountUnbondings(mintDenom string, baseDenom string) map[int64]types.HostAccountUndelegation {
	// create a map to quickly access each undelegation epoch entry and initialise it
	hostAccountUndelegationsMap := make(map[int64]types.HostAccountUndelegation)

	// 232
	hostAccountUndelegationsMap[15826532] = types.HostAccountUndelegation{
		EpochNumber:             232,
		TotalUndelegationAmount: sdk.NewCoin(mintDenom, sdk.ZeroInt()),
		CompletionTime:          time.Time{},
		UndelegationEntries:     make([]types.UndelegationEntry, 0),
	}

	// 236
	hostAccountUndelegationsMap[15881865] = types.HostAccountUndelegation{
		EpochNumber:             236,
		TotalUndelegationAmount: sdk.NewCoin(mintDenom, sdk.ZeroInt()),
		CompletionTime:          time.Time{},
		UndelegationEntries:     make([]types.UndelegationEntry, 0),
	}
	// read the file contents and unmarshal them
	contents, err := os.ReadFile("undelegations.json")
	if err != nil {
		panic(err)
	}

	var unbondings Unbondings
	if err := json.Unmarshal(contents, &unbondings); err != nil {
		panic(err)
	}

	for _, unbonding := range unbondings.UnbondingResponses {
		for _, unbondingEntry := range unbonding.Entries {
			// parse the epoch number
			ubdCreationHeight, err := strconv.ParseInt(unbondingEntry.CreationHeight, 10, 64)
			if err != nil {
				panic(err)
			}

			// get the undelegation object from the map
			hostAccountUndelegation := hostAccountUndelegationsMap[ubdCreationHeight]

			// parse the balance
			balance, ok := sdk.NewIntFromString(unbondingEntry.Balance)
			if !ok {
				panic(err)
			}

			// append the undelegation entry
			hostAccountUndelegation.UndelegationEntries = append(
				hostAccountUndelegation.UndelegationEntries,
				types.UndelegationEntry{
					ValidatorAddress: unbonding.ValidatorAddress,
					Amount:           sdk.NewCoin(baseDenom, balance),
				},
			)

			// if time has not been set for that undelegation entry, set it
			if hostAccountUndelegation.CompletionTime.Equal(time.Time{}) {
				hostAccountUndelegation.CompletionTime = unbondingEntry.CompletionTime
			}

			// save the object back to the map
			hostAccountUndelegationsMap[ubdCreationHeight] = hostAccountUndelegation
		}
	}

	return hostAccountUndelegationsMap
}
