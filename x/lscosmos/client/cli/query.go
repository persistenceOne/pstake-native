package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group lscosmos queries under a subcommand
	cmd := &cobra.Command{
		Use:                        queryRoute,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryHostChainParams(),
		CmdQueryDelegationState(),
		CmdQueryAllowListedValidators(),
		CmdQueryCValue(),
		CmdQueryModuleState(),
		CmdQueryIBCTransientStore(),
		CmdQueryUnclaimed(),
		CmdQueryFailedUnbondings(),
		CmdQueryPendingUnbondings(),
		CmdQueryUnbondingEpoch(),
		CmdQueryHostAccountUndelegation(),
		CmdQueryDelegatorUnbodingEpochEntry(),
		CmdQueryRewardsBoosterAccount(),
		CmdQueryHostAccounts(),
		CmdQueryDepositModuleAccount(),
		CmdDelegatorUnbondingEpochEntries(),
	)

	return cmd
}

func CmdQueryHostChainParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host-chain-params",
		Short: "shows host chain parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.HostChainParams(context.Background(), &types.QueryHostChainParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryDelegationState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation-state",
		Short: "shows delegation state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DelegationState(context.Background(), &types.QueryDelegationStateRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAllowListedValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allow-listed-validators",
		Short: "shows allow listed validators",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AllowListedValidators(context.Background(), &types.QueryAllowListedValidatorsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryCValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "c-value",
		Short: "shows current c-value of the protocol",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CValue(context.Background(), &types.QueryCValueRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryModuleState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module-state",
		Short: "shows current module state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ModuleState(context.Background(), &types.QueryModuleStateRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryIBCTransientStore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc-transient-store",
		Short: "shows amount in ibc-transient-store",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.IBCTransientStore(context.Background(), &types.QueryIBCTransientStoreRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryUnclaimed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unclaimed [delegator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "shows unclaimed amounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Unclaimed(context.Background(), &types.QueryUnclaimedRequest{DelegatorAddress: delegatorAddress.String()})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryFailedUnbondings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "failed-unbondings [delegator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "shows failed unbondings request",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			_, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.FailedUnbondings(context.Background(), &types.QueryFailedUnbondingsRequest{DelegatorAddress: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPendingUnbondings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-unbondings [delegator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "shows pending unbondings",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			_, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.PendingUnbondings(context.Background(), &types.QueryPendingUnbondingsRequest{DelegatorAddress: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryUnbondingEpoch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-epoch [epoch-number]",
		Args:  cobra.ExactArgs(1),
		Short: "Shows unbonding epoch details",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			epochNumber, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.UnbondingEpochCValue(context.Background(), &types.QueryUnbondingEpochCValueRequest{EpochNumber: epochNumber})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryHostAccountUndelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host-account-undelegation [epoch-number]",
		Args:  cobra.ExactArgs(1),
		Short: "Shows host account undelegation for the given epoch number",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			epochNumber, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.HostAccountUndelegation(context.Background(), &types.QueryHostAccountUndelegationRequest{EpochNumber: epochNumber})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryDelegatorUnbodingEpochEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-unbonding-epoch-entry [delegator-address] [epoch-number]",
		Args:  cobra.ExactArgs(2),
		Short: "Shows host account undelegation for the given epoch number",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			epochNumber, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorUnbondingEpochEntry(context.Background(), &types.QueryDelegatorUnbondingEpochEntryRequest{DelegatorAddress: delegatorAddress.String(), EpochNumber: epochNumber})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryRewardsBoosterAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards-booster-account",
		Short: "shows the balance of rewards booster account",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.RewardsBoosterAccount(context.Background(), &types.QueryRewardBoosterAccountRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryHostAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host-accounts",
		Short: "shows the host accounts (delegation and reward) ica owner id",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.HostAccounts(context.Background(), &types.QueryHostAccountsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryDepositModuleAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-module-account",
		Short: "shows the balance of deposit module account",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DepositModuleAccount(context.Background(), &types.QueryDepositModuleAccountRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdDelegatorUnbondingEpochEntries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-unbonding-epoch-entries [delegator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "Shows host account undelegations for all the epoch number the entry has been made.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorUnbondingEpochEntries(context.Background(), &types.QueryAllDelegatorUnbondingEpochEntriesRequest{DelegatorAddress: delegatorAddress.String()})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
