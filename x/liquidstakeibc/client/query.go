package client

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// NewQueryCmd returns the parent command for all x/liquidstakeibc CLi query commands.
func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the pstake liquid staking ibc module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		QueryParamsCmd(),
		QueryHostChainsCmd(),
	)

	return cmd
}

// QueryParamsCmd returns the current module params.
func QueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current liquidstakeibc parameters",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query the current liquidstakeibc parameters:

$ <appd> query liquidstakeibc params
`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// QueryHostChainsCmd returns the registered host chains.
func QueryHostChainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host-chains",
		Short: "Query registered host chains",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query the current registered host chains: $ %s query liquidstakeibc host-chains`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryHostChainsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.HostChains(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
