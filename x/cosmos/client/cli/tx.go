package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/spf13/cobra"
	"strconv"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        cosmosTypes.ModuleName,
		Short:                      "Cosmos transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewIncomingTxnCmd(),
		CmdSetOrchestratorAddress(),
	)

	return txCmd
}

func NewIncomingTxnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "incoming [destination_address] [orchestrator_address] [amount] [chain_id(cosmos)] [txHash(cosmos)] [block_height(cosmos)]",
		Short: `Send txn from cosmos side to native side`,
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			//cmd.Flags().Set(flags.FlagFrom, args[2])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			toAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			orchAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			chainId := args[3]

			txHash := args[4]

			blockHeight, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				return err
			}

			msg := cosmosTypes.NewMsgMintTokensForAccount(toAddr, orchAddress, coins, chainId, txHash, blockHeight)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSetOrchestratorAddress() *cobra.Command {
	//nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:   "set-orchestrator-address [validator-address] [orchestrator-address]",
		Short: "Allows validators to delegate their voting responsibilities to a given key.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := cosmosTypes.MsgSetOrchestrator{
				Validator:    args[0],
				Orchestrator: args[1],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// Send it
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
