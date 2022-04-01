/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package cli

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	cosmosUtils "github.com/persistenceOne/pstake-native/x/cosmos/client/utils"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for the cosmos module.
func GetQueryCmd() *cobra.Command {
	cosmosQueryCmd := &cobra.Command{
		Use:                        cosmosTypes.ModuleName,
		Short:                      "Querying commands for the cosmos module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cosmosQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetOutgoingTxByID(),
		GetCmdQueryProposal(),
		GetCmdQueryProposals(),
		GetCmdQueryVote(),
		GetCmdQueryVotes(),
	)

	return cosmosQueryCmd
}

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
			queryClient := cosmosTypes.NewQueryClient(clientCtx)

			params := &cosmosTypes.QueryParamsRequest{}
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
			queryClient := cosmosTypes.NewQueryClient(clientCtx)

			txID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			txByID := &cosmosTypes.QueryOutgoingTxByIDRequest{
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

func GetCmdQueryProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a proposal. You can find the
proposal-id by running "%s query gov proposals".

Example:
$ %s query gov proposal 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := cosmosTypes.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			// Query the proposal
			res, err := queryClient.Proposal(
				cmd.Context(),
				&cosmosTypes.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Proposal)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposals with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated proposals that match optional filters:

Example:
$ %s query gov proposals --status (DepositPeriod|VotingPeriod|Passed|Rejected)
$ %s query gov proposals --page=2 --limit=100
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			strProposalStatus, _ := cmd.Flags().GetString(flagStatus)

			var proposalStatus cosmosTypes.ProposalStatus

			if len(strProposalStatus) != 0 {
				proposalStatus1, err := cosmosTypes.ProposalStatusFromString(cosmosUtils.NormalizeProposalStatus(strProposalStatus))
				proposalStatus = proposalStatus1
				if err != nil {
					return err
				}
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := cosmosTypes.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Proposals(
				cmd.Context(),
				&cosmosTypes.QueryProposalsRequest{
					ProposalStatus: proposalStatus,
					Pagination:     pageReq,
				},
			)
			if err != nil {
				return err
			}

			if len(res.GetProposals()) == 0 {
				return fmt.Errorf("no proposals found")
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status, status: deposit_period/voting_period/passed/rejected")
	flags.AddPaginationFlagsToCmd(cmd, "proposals")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [voter-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a single vote on a proposal given its identifier.

Example:
$ %s query gov vote 1 cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := cosmosTypes.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			ctx := cmd.Context()
			_, err = queryClient.Proposal(
				ctx,
				&cosmosTypes.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			voterAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.Vote(
				ctx,
				&cosmosTypes.QueryVoteRequest{ProposalId: proposalID, Voter: voterAddr.String()},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Vote)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query vote details for a single proposal by its identifier.

Example:
$ %[1]s query gov votes 1
$ %[1]s query gov votes 1 --page=2 --limit=100
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := cosmosTypes.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			ctx := cmd.Context()
			_, err = queryClient.Proposal(
				ctx,
				&cosmosTypes.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Votes(
				ctx,
				&cosmosTypes.QueryVotesRequest{ProposalId: proposalID, Pagination: pageReq},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "votes")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
