package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/app/params"
	"github.com/persistenceOne/pstake-native/v2/x/lspersistence/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/lspersistence/types"
)

// Simulation operation weights constants.
//
//nolint:gosec
const (
	OpWeightSimulateAddWhitelistValidatorsProposal    = "op_weight_add_whitelist_validators_proposal"
	OpWeightSimulateUpdateWhitelistValidatorsProposal = "op_weight_update_whitelist_validators_proposal"
	OpWeightSimulateDeleteWhitelistValidatorsProposal = "op_weight_delete_whitelist_validators_proposal"
	OpWeightCompleteRedelegationUnbonding             = "op_weight_complete_redelegation_unbonding"
	OpWeightTallyWithLiquidStaking                    = "op_weight_tally_with_liquid_staking"
	MaxWhitelistValidators                            = 10
)

// ProposalContents defines the module weighted proposals' contents for mocking param changes, other actions with keeper
func ProposalContents(ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper, k keeper.Keeper) []simtypes.WeightedProposalContent { //nolint:staticcheck
	return []simtypes.WeightedProposalContent{ //nolint:staticcheck
		simulation.NewWeightedProposalContent(
			OpWeightSimulateAddWhitelistValidatorsProposal,
			params.DefaultWeightAddWhitelistValidatorsProposal,
			SimulateAddWhitelistValidatorsProposal(sk, k),
		),
		simulation.NewWeightedProposalContent(
			OpWeightSimulateUpdateWhitelistValidatorsProposal,
			params.DefaultWeightUpdateWhitelistValidatorsProposal,
			SimulateUpdateWhitelistValidatorsProposal(sk, k),
		),
		simulation.NewWeightedProposalContent(
			OpWeightSimulateDeleteWhitelistValidatorsProposal,
			params.DefaultWeightDeleteWhitelistValidatorsProposal,
			SimulateDeleteWhitelistValidatorsProposal(sk, k),
		),
		simulation.NewWeightedProposalContent(
			OpWeightCompleteRedelegationUnbonding,
			params.DefaultWeightCompleteRedelegationUnbonding,
			SimulateCompleteRedelegationUnbonding(sk),
		),
	}
}

// SimulateAddWhitelistValidatorsProposal generates random add whitelisted validator param change proposal content.
func SimulateAddWhitelistValidatorsProposal(sk types.StakingKeeper, k keeper.Keeper) simtypes.ContentSimulatorFn { //nolint:staticcheck
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content { //nolint:staticcheck
		params := k.GetParams(ctx)

		vals := sk.GetBondedValidatorsByPower(ctx)

		wm := params.WhitelistedValsMap()
		for i := 0; i < len(vals) && len(params.WhitelistedValidators) < MaxWhitelistValidators; i++ {
			val, _ := keeper.RandomValidator(r, sk, ctx)
			if _, ok := wm[val.OperatorAddress]; !ok {
				params.WhitelistedValidators = append(params.WhitelistedValidators,
					types.WhitelistedValidator{
						ValidatorAddress: val.OperatorAddress,
						TargetWeight:     genTargetWeight(r),
					})
				// manually set params for simulation
				k.SetParams(ctx, params)
				break
			}
		}
		return nil
	}
}

// SimulateUpdateWhitelistValidatorsProposal generates random update whitelisted validator param change proposal content.
func SimulateUpdateWhitelistValidatorsProposal(sk types.StakingKeeper, k keeper.Keeper) simtypes.ContentSimulatorFn { //nolint:staticcheck
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content { //nolint:staticcheck
		params := k.GetParams(ctx)

		targetVal, found := keeper.RandomActiveLiquidValidator(r, ctx, k, sk)
		if found {
			for i := range params.WhitelistedValidators {
				if params.WhitelistedValidators[i].ValidatorAddress == targetVal.OperatorAddress {
					params.WhitelistedValidators[i].TargetWeight = genTargetWeight(r)
					// manually set params for simulation
					k.SetParams(ctx, params)
					break
				}
			}
		}
		return nil
	}
}

// SimulateDeleteWhitelistValidatorsProposal generates random delete whitelisted validator param change proposal content.
func SimulateDeleteWhitelistValidatorsProposal(sk types.StakingKeeper, k keeper.Keeper) simtypes.ContentSimulatorFn { //nolint:staticcheck
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content { //nolint:staticcheck
		params := k.GetParams(ctx)

		targetVal, found := keeper.RandomActiveLiquidValidator(r, ctx, k, sk)
		if found {
			remove := func(slice []types.WhitelistedValidator, s int) []types.WhitelistedValidator {
				return append(slice[:s], slice[s+1:]...)
			}

			for i := range params.WhitelistedValidators {
				if params.WhitelistedValidators[i].ValidatorAddress == targetVal.OperatorAddress {
					params.WhitelistedValidators[i].TargetWeight = genTargetWeight(r)
					params.WhitelistedValidators = remove(params.WhitelistedValidators, i)
					k.SetParams(ctx, params)
					break
				}
			}
		}
		return nil
	}
}

// SimulateCompleteRedelegationUnbonding mocking complete redelegations, unbondings by BlockValidatorUpdates of staking keeper.
func SimulateCompleteRedelegationUnbonding(sk types.StakingKeeper) simtypes.ContentSimulatorFn { //nolint:staticcheck
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content { //nolint:staticcheck
		reds := sk.GetAllRedelegations(ctx, types.LiquidStakingProxyAcc, nil, nil)
		ubds := sk.GetAllUnbondingDelegations(ctx, types.LiquidStakingProxyAcc)
		if len(reds) != 0 || len(ubds) != 0 {
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 100).WithBlockTime(ctx.BlockTime().Add(stakingtypes.DefaultUnbondingTime))
			sk.BlockValidatorUpdates(ctx)
		}
		return nil
	}
}
