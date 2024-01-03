package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

// GetQueryCmd returns a root CLI command handler for all x/liquidstake query commands.
func GetQueryCmd() *cobra.Command {
	liquidValidatorQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Aliases:                    []string{"ls"},
		Short:                      "Querying commands for the liquidstake module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	liquidValidatorQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryLiquidValidators(),
		GetCmdQueryStates(),
	)

	return liquidValidatorQueryCmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the values set as liquidstake parameters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as liquidstake parameters.

Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(
				cmd.Context(),
				&types.QueryParamsRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryLiquidValidators implements the query liquidValidators command.
func GetCmdQueryLiquidValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-validators",
		Args:  cobra.NoArgs,
		Short: "Query all liquid validators",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Queries all liquid validators.

Example:
$ %s query %s liquid-validators
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			if err != nil {
				return err
			}

			res, err := queryClient.LiquidValidators(
				cmd.Context(),
				&types.QueryLiquidValidatorsRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryStates implements the query states command.
func GetCmdQueryStates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "states",
		Args:  cobra.NoArgs,
		Short: "Query states",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Queries states about net amount, mint rate.

Example:
$ %s query %s states
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.States(
				cmd.Context(),
				&types.QueryStatesRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
