/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for the cosmos module.
func GetQueryCmd() *cobra.Command {
	cosmosQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the cosmos module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cosmosQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetOutgoingTxByID(),
	)

	return cosmosQueryCmd
}

// GetCmdQueryParams implements a command to return the current cosmos parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current cosmos parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryParamsRequest{}
			res, err := queryClient.QueryParams(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetOutgoingTxByID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txByID [tx-id]",
		Short: "Query the outgoing tx by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			txID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			txByID := &types.QueryOutgoingTxByIDRequest{
				TxID: txID,
			}

			res, err := queryClient.QueryTxByID(context.Background(), txByID)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.CosmosTxDetails)
		},
	}
	return cmd
}
