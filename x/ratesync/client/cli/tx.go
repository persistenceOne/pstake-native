package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

var DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdMsgUpdateParams())
	cmd.AddCommand(CmdCreateChain())
	cmd.AddCommand(CmdUpdateChain())
	cmd.AddCommand(CmdDeleteChain())
	// this line is used by starport scaffolding # 1

	return cmd
}

func CmdMsgUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [params-file]",
		Short: "Broadcast message MsgUpdateParams",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var params types.Params

			paramsInFile, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			err = json.Unmarshal(paramsInFile, &params)
			if err != nil {
				return err
			}
			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgUpdateParams(authority.String(), params)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreateChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chain [path_to_file]",
		Short: "Create a new chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var hostChain types.HostChain

			hostChainInFile, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			err = json.Unmarshal(hostChainInFile, &hostChain)
			if err != nil {
				return fmt.Errorf("err unmarshalling json err: %v, should be of type %v", err, hostChain)
			}

			msg := types.NewMsgCreateHostChain(
				clientCtx.GetFromAddress().String(),
				hostChain,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-chain [index]",
		Short: "Update a chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var hostChain types.HostChain

			hostChainInFile, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			err = json.Unmarshal(hostChainInFile, &hostChain)
			if err != nil {
				return fmt.Errorf("err unmarshalling json err: %v, should be of type %v", err, hostChain)
			}

			msg := types.NewMsgUpdateHostChain(
				clientCtx.GetFromAddress().String(),
				hostChain,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-chain [id]",
		Short: "Delete a chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			id := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			idInt, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeleteHostChain(
				clientCtx.GetFromAddress().String(),
				idInt,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
