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
		Aliases:                    []string{"liquidstake", "lsibc"},
		Short:                      "Querying commands for the pstake liquid staking ibc module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		QueryParamsCmd(),
		QueryHostChainsCmd(),
		QueryDepositsCmd(),
		QueryUnbondingsCmd(),
		QueryUserUnbondingsCmd(),
		QueryValidatorUnbondingsCmd(),
		QueryDepositAccountBalanceCmd(),
		QueryExchangeRateCmd(),
		QueryUnbondingCmd(),
	)

	return cmd
}

// QueryParamsCmd returns the current module params.
func QueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current liquidstakeibc parameters",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query the current liquidstakeibc parameters: $ %s query liquidstakeibc params`,
				version.AppName,
			),
		),
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

			res, err := queryClient.HostChains(cmd.Context(), &types.QueryHostChainsRequest{})
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
		Use:   "deposits [chain-id]",
		Short: "Query deposit records for a host chain",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query all deposits: $ %s query liquidstakeibc deposits [chain-id]`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Deposits(cmd.Context(), &types.QueryDepositsRequest{ChainId: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// QueryUnbondingsCmd returns all unbonding records for a host chain.
func QueryUnbondingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbondings [chain-id]",
		Short: "Query all unbonding records for a host chain",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query an unbonding record: $ %s query liquidstakeibc unbondings [chain-id]`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Unbondings(context.Background(), &types.QueryUnbondingsRequest{ChainId: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// QueryUnbondingCmd returns an unbonding record for a host chain and an epoch.
func QueryUnbondingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding [chain-id] [epoch]",
		Short: "Query an unbonding record for a host chain and an epoch",
		Args:  cobra.ExactArgs(2),
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query an unbonding record: $ %s query liquidstakeibc unbonding [chain-id] [epoch]`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			epoch, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.Unbonding(
				context.Background(),
				&types.QueryUnbondingRequest{ChainId: args[0], Epoch: epoch})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// QueryUserUnbondingsCmd returns all user unbondings.
func QueryUserUnbondingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-unbondings [delegator-address]",
		Short: "Query all user unbonding records",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query a user unbonding record: $ %s query liquidstakeibc user-unbondings [delegator-address]`,
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

			res, err := queryClient.UserUnbondings(
				context.Background(),
				&types.QueryUserUnbondingsRequest{
					Address: args[0],
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

// QueryValidatorUnbondingsCmd returns all validator unbondings for a host chain.
func QueryValidatorUnbondingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-unbondings [chain-id]",
		Short: "Query all validator unbonding records for a host chain",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query all validator unbondings for a host chain: $ %s query liquidstakeibc validator-unbondings [chain-id]`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ValidatorUnbondings(
				context.Background(),
				&types.QueryValidatorUnbondingRequest{
					ChainId: args[0],
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

// QueryDepositAccountBalanceCmd returns the host chain deposit account balance.
func QueryDepositAccountBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-account-balance [chain-id]",
		Short: "Query deposit records for a host chain",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query a host chain deposit account balance: $ %s query liquidstakeibc deposit-account-balance [chain-id]`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DepositAccountBalance(
				cmd.Context(),
				&types.QueryDepositAccountBalanceRequest{ChainId: args[0]},
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

// QueryExchangeRateCmd returns the host chain exchange rate between the host token and the stk token.
func QueryExchangeRateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exchange-rate [chain-id]",
		Short: "Query the exchange rate of a host chain",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Query the exchange rate of a host chain: $ %s query liquidstakeibc exchange-rate [chain-id]`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ExchangeRate(cmd.Context(), &types.QueryExchangeRateRequest{ChainId: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
