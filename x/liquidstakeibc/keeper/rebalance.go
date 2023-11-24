package keeper

import (
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

type delegation struct {
	validator        string
	ideal            math.Int
	delegation       math.Int
	diff             math.Int
	validatorDetails types.Validator
}

// Rebalance tries to make redelegate transactions to host-chain to balance the delegations as per the weights.
func (k Keeper) Rebalance(ctx sdk.Context, epoch int64) []proto.Message {

	hcs := k.GetAllHostChains(ctx)
	for _, hc := range hcs {
		// skip unbonding epoch, as we do not want to redelegate tokens that might be going through unbond txn in same epoch.
		// nothing bad will happen even if we do as long as unbonding txns are triggered before redelegations.
		if !types.IsUnbondingEpoch(hc.UnbondingFactor, epoch) {
			k.Logger(ctx).Info("redelegation epoch co-incides with unbonding epoch, skipping it")
			continue
		}
		msgs := k.GenerateRedelegateMsgs(ctx, *hc)
		// send one msg per ica
		for _, msg := range msgs {
			ibcSeq, err := k.GenerateAndExecuteICATx(ctx, hc.ConnectionId, hc.DelegationAccount.Owner, []proto.Message{msg})
			if err != nil {
				k.Logger(ctx).Error("Failed to submit ica redelegate txns with", "err:", err)
			}
			k.SetRedelegationTx(ctx, &types.RedelegateTx{
				ChainId:       hc.ChainId,
				IbcSequenceId: ibcSeq,
				State:         types.RedelegateTx_REDELEGATE_SENT,
			})
		}
	}
	return nil
}

func (k Keeper) GenerateRedelegateMsgs(ctx sdk.Context, hc types.HostChain) []proto.Message {
	var AcceptableDelta = hc.Params.RedelegationAcceptableDelta
	var MaxRedelegationEntries = hc.Params.MaxEntries
	sum := math.ZeroInt()
	for _, validator := range hc.Validators {
		sum = sum.Add(validator.DelegatedAmount)
	}

	idealDelegationList := make([]delegation, len(hc.Validators))
	sum2 := math.ZeroInt()
	for i, validator := range hc.Validators {
		idealAmt := validator.Weight.MulInt(sum).TruncateInt()
		// last element
		if i == len(hc.Validators)-1 {
			idealAmt = sum.Sub(sum2)
		}
		sum2 = sum2.Add(idealAmt)
		idealDelegationList = append(idealDelegationList,
			delegation{
				validator:        validator.OperatorAddress,
				ideal:            idealAmt,
				delegation:       validator.DelegatedAmount,
				diff:             validator.DelegatedAmount.Sub(idealAmt),
				validatorDetails: *validator,
			})
	}
	// negative diffs first, so ascending
	idealDelegationList = sortDelegationListAsc(idealDelegationList)
	revIdealList := make([]delegation, len(idealDelegationList))
	copy(revIdealList, idealDelegationList)
	// positive diffs first (descending)
	Reverse(revIdealList)
	redelegations, ok := k.GetRedelegations(ctx, hc.ChainId)
	if !ok {
		redelegations = &types.Redelegations{
			ChainID:       hc.ChainId,
			Redelegations: []*stakingtypes.Redelegation{},
		}
	}

	var msgs []proto.Message
L1:
	for i := range revIdealList {
		if revIdealList[i].diff.LT(AcceptableDelta) {
			break L1
		}
		// RedelegationExistsToValidator: This is not updated inside the loop (with newer msgs), so some ICA redelegate txns might fail, and it is ok.
		if !k.RedelegationExistsToValidator(redelegations.Redelegations, revIdealList[i].validator) {
			//re-sort idealDelegationAsc
			idealDelegationList = sortDelegationListAsc(idealDelegationList)
		L2:
			for j := range idealDelegationList {
				if revIdealList[i].validator == idealDelegationList[j].validator {
					break L1
				}
				if revIdealList[i].diff.LT(AcceptableDelta) || idealDelegationList[j].diff.IsPositive() {
					break L2
				}
				if !idealDelegationList[j].validatorDetails.Delegable || idealDelegationList[j].validatorDetails.Status != stakingtypes.Bonded.String() {
					continue L2
				}

				// RedelegationFromAToB: This is not updated inside the loop (with newer msgs), so some ICA redelegate txns might fail, and it is ok.
				_, numEntries := k.RedelegationFromAToB(redelegations.Redelegations, revIdealList[i].validator, idealDelegationList[j].validator)
				if numEntries < MaxRedelegationEntries {
					redelegationAmt := math.MinInt(revIdealList[i].diff.Abs(), idealDelegationList[j].diff.Abs())
					redelegateMsg := &stakingtypes.MsgBeginRedelegate{
						DelegatorAddress:    hc.DelegationAccount.Address,
						ValidatorSrcAddress: revIdealList[i].validator,
						ValidatorDstAddress: idealDelegationList[j].validator,
						Amount:              sdk.NewCoin(hc.HostDenom, redelegationAmt),
					}
					msgs = append(msgs, redelegateMsg)
					revIdealList[i].diff = revIdealList[i].diff.Sub(redelegationAmt)
					idealDelegationList[j].diff = idealDelegationList[j].diff.Add(redelegationAmt)
				}
			}
		}
	}
	return msgs
}

func (k Keeper) RedelegationExistsToValidator(redelegations []*stakingtypes.Redelegation, toValoper string) bool {
	for _, redelegation := range redelegations {
		if redelegation.ValidatorDstAddress == toValoper && len(redelegation.Entries) > 0 {
			return true
		}
	}
	return false
}

func (k Keeper) RedelegationFromAToB(redelegations []*stakingtypes.Redelegation, fromValoper, toValoper string) (bool, uint32) {
	for _, redelegation := range redelegations {
		if redelegation.ValidatorDstAddress == toValoper && redelegation.ValidatorSrcAddress == fromValoper {
			return true, uint32(len(redelegation.Entries))
		}
	}
	return false, 0
}

func sortDelegationListAsc(idealDelegationList []delegation) []delegation {
	sort.SliceStable(idealDelegationList, func(i, j int) bool {
		switch {
		case idealDelegationList[i].diff.LT(idealDelegationList[j].diff):
			return true
		case idealDelegationList[i].diff.GT(idealDelegationList[j].diff):
			return false
		default:
			return idealDelegationList[i].validator < idealDelegationList[j].validator
		}
	})
	return idealDelegationList
}

// remove when go updates to 1.21, and use slices package.
// Reverse reverses the elements of the slice in place.
func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
