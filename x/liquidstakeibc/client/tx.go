package client

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// NewTxCmd returns a root CLI command handler for all liquid staking transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Aliases:                    []string{"liquidstake", "lsibc"},
		Short:                      "Pstake liquid staking ibc transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewRegisterHostChainCmd(),
		NewUpdateHostChainCmd(),
		NewLiquidStakeCmd(),
		NewLiquidStakeCmdLSM(),
		NewLiquidUnstakeCmd(),
		NewRedeemCmd(),
		NewUpdateParamsCmd(),
	)

	return txCmd
}

// NewRegisterHostChainCmd implements the command to register a host chain.
func NewRegisterHostChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-host-chain [connection-id] [channel-id] [port-id] [deposit-fee] [restake-fee] [unstake-fee] [redemption-fee] [host-denom] [minimum-deposit] [unbonding-factor] [autocompound-factor]",
		Args:  cobra.ExactArgs(11),
		Short: "Register a host chain",
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Submit a register host chain transaction: $ %s tx liquidstakeibc register-host-chain connection-0 channel-0 transfer 0.00 0.05 0.00 0.005 uatom 1 4 20`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			minimumDeposit, ok := sdk.NewIntFromString(args[8])
			if !ok {
				return fmt.Errorf("unable to parse string to sdk.Int")
			}

			unbondingFactor, err := strconv.ParseInt(args[9], 10, 64)
			if err != nil {
				return fmt.Errorf("unable to parse string to int64")
			}

			autocompoundFactor, err := strconv.ParseInt(args[10], 10, 64)
			if err != nil {
				return fmt.Errorf("unable to parse string to int64")
			}

			msg := types.NewMsgRegisterHostChain(
				args[0],
				args[1],
				args[2],
				args[3],
				args[4],
				args[5],
				args[6],
				args[7],
				minimumDeposit,
				unbondingFactor,
				clientCtx.FromAddress.String(),
				autocompoundFactor,
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

// NewUpdateHostChainCmd implements the command to update a host chain.
func NewUpdateHostChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-host-chain [chain-id] [updates]",
		Args:  cobra.ExactArgs(2),
		Short: "Update a host chain",
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Submit an update host chain transaction: 
$ %s tx liquidstakeibc update-host-chain gaia-1 '[
    {
        "key": "active",
        "value": "true"
    },
    {
        "key": "set_withdraw_address",
        "value": ""
    },
    {
        "key": "flags",
        "value": "{\"lsm\": true}"
    },
    {
        "key": "add_validator",
        "value": "{\"operator_address\": \"cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt\", \"status\": \"BOND_STATUS_BONDED\", \"weight\": \"1\", \"delegated_amount\": \"0\", \"exchange_rate\": \"0\", \"unbonding_epoch\": 0}"
    }
]'`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			updates := make([]*types.KVUpdate, 0)
			if err = json.Unmarshal([]byte(args[1]), &updates); err != nil {
				return err
			}

			msg := types.NewMsgUpdateHostChain(
				args[0],
				clientCtx.FromAddress.String(),
				updates,
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

func NewLiquidStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-stake [amount]",
		Short: `Liquid Stake tokens from a registered host chain into stk tokens`,
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Submit a liquid stake transaction: $ %s tx liquidstakeibc liquid-stake 100000000ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
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

func NewLiquidStakeCmdLSM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-stake-lsm [delegations]",
		Short: `Liquid Stake an existing delegation from a registered host chain into stk tokens`,
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Submit a liquid stake LSM transaction: 
$ %s tx liquidstakeibc liquid-stake-lsm '[
    {
        "amount": "100000000",
        "denom": "ibc/7976C604E31F2C1062F1BF20175FC14E08FC855C4BECDBA4E1274646914FCB7C"
    }
]'`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			delegatorAddress := clientctx.GetFromAddress()
			msg := types.NewMsgLiquidStakeLSM(coins, delegatorAddress)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewLiquidUnstakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-unstake [amount]",
		Short: `Unstake stk tokens from a registered host chain`,
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Submit a liquid unstake transaction: $ %s tx liquidstakeibc liquid-unstake 100000000stk/uatom`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
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
			msg := types.NewMsgLiquidUnstake(amount, delegatorAddress)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewRedeemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem [amount]",
		Short: `Instantly redeem stk tokens from a registered host chain`,
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Submit a redeem transaction: $ %s tx liquidstakeibc redeem 50000000stk/uatom`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
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
			msg := types.NewMsgRedeem(amount, delegatorAddress)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUpdateParamsCmd implements the command to update the module params.
func NewUpdateParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [params-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Update the module params",
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`Submit an update params transaction: $ %s tx liquidstakeibc update-params /params-file.json

Params file contents:

{
  "messages": [{
    "@type": "/pstake.liquidstakeibc.v1beta1.MsgUpdateParams",
    "authority": "persistence10d07y265gmmuvt4z0w9aw880jnsr700j5w4kch",
    "params": {
      "admin_address": "persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr",
      "fee_address": "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld"
    }
  }],
  "deposit": "10000000uxprt",
  "proposer": "persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu",
  "title": "Update module addresses",
  "summary": "Updates both the admin and the fee address of the module",
  "metadata": ""
}`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			msg := types.NewMsgUpdateParams(authority, params)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
