package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		QueryDepositsCmd(),
		QueryUnbondingCmd(),
		QueryUserUnbondingCmd(),
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

// QueryDepositsCmd returns all user deposits.
func QueryDepositsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits",
		Short: "Query deposit records",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query all deposits: $ %s query liquidstakeibc deposits`,
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

			req := &types.QueryDepositsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Deposits(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// QueryUnbondingCmd returns an unbonding record.
func QueryUnbondingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding [epoch-number] [host-denom]",
		Short: "Query an unbonding record",
		Args:  cobra.ExactArgs(2),
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query an unbonding record: $ %s query liquidstakeibc unbonding [epoch-number] [host-denom]`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			epochNumber, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.Unbonding(
				context.Background(),
				&types.QueryUnbondingRequest{
					EpochNumber: epochNumber,
					HostDenom:   args[1],
				},
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

// QueryUserUnbondingCmd returns a user unbonding record.
func QueryUserUnbondingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-unbonding [delegator-address] [epoch-number] [host-denom]",
		Short: "Query a user unbonding record",
		Args:  cobra.ExactArgs(2),
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query a user unbonding record: $ %s query liquidstakeibc user-unbonding [delegator-address] [epoch-number] [host-denom]`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			epochNumber, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.UserUnbondings(
				context.Background(),
				&types.QueryUserUnbondingsRequest{
					Address:     args[0],
					EpochNumber: epochNumber,
					HostDenom:   args[2],
				},
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
