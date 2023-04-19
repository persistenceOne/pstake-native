package client

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// NewTxCmd returns a root CLI command handler for all liquid staking transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Pstake liquid staking ibc transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewRegisterHostChainCmd(),
		NewUpdateHostChainCmd(),
		NewLiquidStakeCmd(),
	)

	return txCmd
}

// NewRegisterHostChainCmd implements the command to register a host chain.
// TODO: Remove this when tagging version. Users should not be able to register chains.
func NewRegisterHostChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-host-chain [connection-id] [host-denom] [local-denom] [minimum-deposit]",
		Args:  cobra.ExactArgs(4),
		Short: "Register a host chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			minimumDeposit, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("unable to parse string to sdk.Int")
			}

			msg := types.NewMsgRegisterHostChain(
				args[0],
				args[1],
				args[2],
				minimumDeposit,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

// NewUpdateHostChainCmd implements the command to update a host chain.
// TODO: Remove this when tagging version. Users should not be able to update chains.
func NewUpdateHostChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-host-chain [chain-id] [chain-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Update a host chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			updates := make([]*types.KVUpdate, 0)
			if err = json.Unmarshal([]byte(args[1]), updates); err != nil {
				return err
			}

			msg := types.NewMsgUpdateHostChain(
				args[0],
				updates,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func NewLiquidStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-stake [amount]",
		Short: `Liquid Stake tokens from a registered host chain into stk tokens`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			delegatorAddress := clientctx.GetFromAddress()
			msg := types.NewMsgLiquidStake(amount, delegatorAddress)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
