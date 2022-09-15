package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
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
