package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"slices"
	"sort"
)

type delegation struct {
	validator        string
	ideal            math.Int
	delegation       math.Int
	diff             math.Int
	validatorDetails types.Validator
}

func (k Keeper) Rebalance(ctx sdk.Context, epoch int64) []proto.Message {

	hcs := k.GetAllHostChains(ctx)
	for _, hc := range hcs {
		// skip unbonding epoch, as we do not want to redelegate tokens that might be going through unbond txn in same epoch.
		// nothing bad will happen even if we do as long as unbonding txns are triggered before redelegations.
		if !types.IsUnbondingEpoch(hc.UnbondingFactor, epoch) {
			k.Logger(ctx).Info("redelegation epoch co-incides with unbonding epoch, skipping it")
			continue
		}
		var AcceptableDelta = hc.Params.RedelegationAcceptableDelta
		var MaxRedelegationEntries = hc.Params.MaxEntries
		sum := math.ZeroInt()
		for _, validator := range hc.Validators {
			sum = sum.Add(validator.DelegatedAmount)
		}

		var idealDelegationList []delegation
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
		// negative diffs first, so descending
		idealDelegationList = k.SortDelegationListAsc(idealDelegationList)
		revIdealList := idealDelegationList
		// positive diffs first (descending)
		slices.Reverse(revIdealList)
		redelegations, _ := k.GetRedelegations(ctx, hc.ChainId)

		var msgs []proto.Message
	L1:
		for i, _ := range revIdealList {
			if revIdealList[i].diff.LT(AcceptableDelta) {
				break L1
			}
			if !k.RedelegationExistsToValidator(redelegations.Redelegations, revIdealList[i].validator) {
				//re-sort idealDelegationAsc
				idealDelegationList = k.SortDelegationListAsc(idealDelegationList)
			L2:
				for j, _ := range idealDelegationList {
					if revIdealList[i].validator == idealDelegationList[j].validator {
						break L1
					}
					if !idealDelegationList[j].validatorDetails.Delegable {
						continue L2
					}
					if idealDelegationList[j].diff.IsPositive() {
						break L2
					}
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
						if revIdealList[i].diff.LT(AcceptableDelta) {
							break L2
						}
					}
				}
			}
		}
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

func (k Keeper) SortDelegationListAsc(idealDelegationList []delegation) []delegation {
	sort.SliceStable(idealDelegationList, func(i, j int) bool {
		if idealDelegationList[i].diff.LT(idealDelegationList[j].diff) {
			return true
		} else if idealDelegationList[i].diff.GT(idealDelegationList[j].diff) {
			return false
		} else {
			if idealDelegationList[i].validator < idealDelegationList[j].validator {
				return true
			} else {
				return false
			}
		}
	})
	return idealDelegationList
}
