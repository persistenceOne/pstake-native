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
var RESET_EPOCHS = []int64{232}

var UNBONDING_CREATION_HEIGHT_EPOCHS = map[int64]int64{232: 15826532, 236: 15881865}

func (k Keeper) Fork(ctx sdk.Context) error {
	hostChainParams := k.GetHostChainParams(ctx)

	// Readd deleted undelegations. // only for RESET_EPOCHS
	undelegationsMap := ParseHostAccountUnbondings(hostChainParams.MintDenom, hostChainParams.BaseDenom)

	// reveert failed unbondingsCValue for deleted undelegations.
	for _, epoch := range RESET_EPOCHS {
		// reveert failed unbondingsCValue
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, epoch)
		unbondingEpochCValue.IsFailed = false
		k.SetUnbondingEpochCValue(ctx, unbondingEpochCValue)

		// Readd deleted undelegations
		_, err := k.GetHostAccountUndelegationForEpoch(ctx, epoch)
		if err == nil {
			continue
		}
		undelegations, ok := undelegationsMap[UNBONDING_CREATION_HEIGHT_EPOCHS[epoch]]
		if !ok {
			k.Logger(ctx).Error("Undelegations not found to reset")
		}
		undelegations.TotalUndelegationAmount = unbondingEpochCValue.STKBurn
		k.AddHostAccountUndelegation(ctx, undelegations)
	}

	// Do claims for all matured undelegations users (claim atoms).
	//allDelegatorUnbondingEntries := k.IterateAllDelegatorUnbondingEpochEntry(ctx)
	//for _, delegatorUnbondingEntry := range allDelegatorUnbondingEntries {
	//	err := k.ClaimMatured(ctx, delegatorUnbondingEntry)
	//	if err != nil {
	//		return err
	//	}
	//}
	// CheckDelegatorUnbondings Calculate for claims of stkatom. failed unbondings
	stkAtomExpectedInUnbondingAcc := k.CheckDelegatorUnbondings(ctx)
	currentStkAtomBalance := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.UndelegationModuleAccount), hostChainParams.MintDenom)
	if stkAtomExpectedInUnbondingAcc.GT(currentStkAtomBalance.Amount) {
		diff := stkAtomExpectedInUnbondingAcc.Sub(currentStkAtomBalance.Amount)
		// move tokens from protocol account to -> module account
		protocolAddr := sdk.MustAccAddressFromBech32(PROTOCOL_ACC)
		stkbalanceWithProtocol := k.bankKeeper.GetBalance(ctx, protocolAddr, hostChainParams.MintDenom)
		if stkbalanceWithProtocol.Amount.GT(diff) {
			err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, protocolAddr, types.UndelegationModuleAccount, sdk.NewCoins(sdk.NewCoin(hostChainParams.MintDenom, diff)))
			if err != nil {
				return err
			}
		} else {
			panic(fmt.Sprintf("Protocol doesn't have enought stkatom to fill deposits"))
		}

		for _, epoch := range RESET_EPOCHS {
			// Find deficit of stkatoms from epoch of delegator unbonding entries. (232 and future ones)
			unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, epoch)
			_, totalUnstaked := k.GetAllEpochDelegatorsUnbondingsForCValue(ctx, epoch)
			if unbondingEpochCValue.STKBurn.Amount.GT(totalUnstaked) {
				deficitForThisEpoch := unbondingEpochCValue.STKBurn.Amount.Sub(totalUnstaked)

				// add atom claim to protocol account
				k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
					DelegatorAddress: PROTOCOL_ACC,
					EpochNumber:      epoch,
					Amount:           sdk.NewCoin(hostChainParams.MintDenom, deficitForThisEpoch),
				})
				diff = diff.Sub(deficitForThisEpoch)
			}
		}
	}

	// Claim failed unbondings
	//allDelegatorUnbondingEntries = k.IterateAllDelegatorUnbondingEpochEntry(ctx)
	//for _, delegatorUnbondingEntry := range allDelegatorUnbondingEntries {
	//	err := k.ClaimFailed(ctx, delegatorUnbondingEntry)
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

// Returns entreis and total
func (k Keeper) GetAllEpochDelegatorsUnbondingsForCValue(ctx sdk.Context, epochNumber int64) ([]types.DelegatorUnbondingEpochEntry, sdk.Int) {
	totalStkAmount := sdk.ZeroInt()
	allDelegatorUnbondingEntries := k.IterateAllDelegatorUnbondingEpochEntry(ctx)
	var epochUnbondings []types.DelegatorUnbondingEpochEntry
	for _, delegatorUnbondingEntry := range allDelegatorUnbondingEntries {
		if delegatorUnbondingEntry.EpochNumber == epochNumber {
			epochUnbondings = append(epochUnbondings, delegatorUnbondingEntry)
			totalStkAmount = totalStkAmount.Add(delegatorUnbondingEntry.Amount.Amount)
		}
	}
	return epochUnbondings, totalStkAmount
}

func (k Keeper) CheckDelegatorUnbondings(ctx sdk.Context) sdk.Int {
	//hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)
	allUnbondingEpochCValues := k.IterateAllUnbondingEpochCValues(ctx)
	allDelegatorUnbondingEntries := k.IterateAllDelegatorUnbondingEpochEntry(ctx)
	currEpoch := k.epochKeeper.GetEpochInfo(ctx, types.UndelegationEpochIdentifier)
	currUnbondingEpoch := types.CurrentUnbondingEpoch(currEpoch.CurrentEpoch)
	totalStkAtomToBeInUnbondingAcc := sdk.ZeroInt()
	for _, delegatorUnbondingEntry := range allDelegatorUnbondingEntries {
		//Current epoch will not have unbondingEpochCvalue, so just add their stkatoms
		if delegatorUnbondingEntry.EpochNumber == currUnbondingEpoch {
			totalStkAtomToBeInUnbondingAcc = totalStkAtomToBeInUnbondingAcc.Add(delegatorUnbondingEntry.Amount.Amount)
			continue
		}
		//
		unbondingEpochCValue, found := findUnbondingCValue(allUnbondingEpochCValues, delegatorUnbondingEntry.EpochNumber)
		if !found {
			// should never occur as per the code logic, but if not found, just continue.
			k.Logger(ctx).Error("Unbonding epoch not found for ", "epochnumber", delegatorUnbondingEntry.EpochNumber, "delegatorunbondingentry", delegatorUnbondingEntry)
			continue
		}
		// if is matured, their claims are atoms, not stkatoms.
		if unbondingEpochCValue.IsMatured {
			continue
		}
		hostAccUnbonding, found := findHostAccUnbonding(delegationState.HostAccountUndelegations, delegatorUnbondingEntry.EpochNumber)
		//filter the unbonding requests we might have stored that are not in the unstaking period (should filter 164,216,228 epochs)
		if !unbondingEpochCValue.IsFailed && found && hostAccUnbonding.CompletionTime.Equal(time.Time{}) {
			totalStkAtomToBeInUnbondingAcc = totalStkAtomToBeInUnbondingAcc.Add(delegatorUnbondingEntry.Amount.Amount)
			continue
		}
		// if unstaking is failed, their stkatoms should be available for claim,
		//NOTE the reset of Db has to occur before this function orelse 232 will come here
		if unbondingEpochCValue.IsFailed {
			totalStkAtomToBeInUnbondingAcc = totalStkAtomToBeInUnbondingAcc.Add(delegatorUnbondingEntry.Amount.Amount)
			continue
		}
	}
	return totalStkAtomToBeInUnbondingAcc
}

func findUnbondingCValue(allunbondingCValues []types.UnbondingEpochCValue, epoch int64) (types.UnbondingEpochCValue, bool) {
	for _, unbondingEpochCvalue := range allunbondingCValues {
		if epoch == unbondingEpochCvalue.EpochNumber {
			return unbondingEpochCvalue, true
		}
	}
	return types.UnbondingEpochCValue{}, false
}
func findHostAccUnbonding(allHostAccUndelegations []types.HostAccountUndelegation, epoch int64) (types.HostAccountUndelegation, bool) {
	for _, hostAccUndelgation := range allHostAccUndelegations {
		if epoch == hostAccUndelgation.EpochNumber {
			return hostAccUndelgation, true
		}
	}
	return types.HostAccountUndelegation{}, false
}

// func (k Keeper) ClaimMatured(ctx sdk.Context, unbondingEntry types.DelegatorUnbondingEpochEntry) error {
//
//		delegatorAddress := sdk.MustAccAddressFromBech32(unbondingEntry.DelegatorAddress)
//		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, unbondingEntry.EpochNumber)
//		if unbondingEpochCValue.IsMatured {
//			// get c value from the UnbondingEpochCValue struct
//			// calculate claimable amount from un inverse c value
//			claimableAmount := sdk.NewDecFromInt(unbondingEntry.Amount.Amount).Quo(unbondingEpochCValue.GetUnbondingEpochCValue())
//
//			// calculate claimable coin and community coin to be sent to delegator account and community pool respectively
//			claimableCoin, _ := sdk.NewDecCoinFromDec(k.GetIBCDenom(ctx), claimableAmount).TruncateDecimal()
//
//			// send coin to delegator address from undelegation module account
//			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.UndelegationModuleAccount, delegatorAddress, sdk.NewCoins(claimableCoin))
//			if err != nil {
//				return err
//			}
//
//			ctx.EventManager().EmitEvents(sdk.Events{
//				sdk.NewEvent(
//					types.EventTypeClaim,
//					sdk.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
//					sdk.NewAttribute(types.AttributeAmount, unbondingEntry.Amount.String()),
//					sdk.NewAttribute(types.AttributeClaimedAmount, claimableAmount.String()),
//				)},
//			)
//
//			// remove entry from unbonding epoch entry
//			k.RemoveDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEntry.EpochNumber)
//		}
//		return nil
//	}
func (k Keeper) ClaimFailed(ctx sdk.Context, unbondingEntry types.DelegatorUnbondingEpochEntry) error {
	delegatorAddress := sdk.MustAccAddressFromBech32(unbondingEntry.DelegatorAddress)
	unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, unbondingEntry.EpochNumber)
	if unbondingEpochCValue.IsFailed {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.UndelegationModuleAccount, delegatorAddress, sdk.NewCoins(unbondingEntry.Amount))
		if err != nil {
			return err
		}

		// remove entry from unbonding epoch entry
		k.RemoveDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEntry.EpochNumber)
	}
	return nil
}

func (k Keeper) CheckUnstakingEpochForPacket(ctx sdk.Context, packetMsgs []sdk.Msg) (int64, error) {
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
			return 0, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "unable to unmarshal msg of type %s, expected type MsgUndelegate", sdk.MsgTypeURL(packetMsg))
		}
		for _, unbonding := range unbondings {
			if unbonding.CompletionTime.Equal(time.Time{}) {
				for _, ubd := range unbonding.UndelegationEntries {
					if undelegateMsg.ValidatorAddress == ubd.ValidatorAddress && undelegateMsg.Amount.IsEqual(ubd.Amount) {
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

// ////////// fetch data from json:
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
