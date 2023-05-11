package cli

// DONTCOVER
// client is excluded from test coverage in MVP version

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/persistenceOne/pstake-native/v2/x/lspersistence/types"
)

// GetTxCmd returns a root CLI command handler for all x/lspersistence transaction commands.
func GetTxCmd() *cobra.Command {
	liquidstakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Liquid-staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	liquidstakingTxCmd.AddCommand(
		NewLiquidStakeCmd(),
		NewLiquidUnstakeCmd(),
		NewUpdateParamsCmd(),
	)

	return liquidstakingTxCmd
}

// NewLiquidStakeCmd implements the liquid stake coin command handler.
func NewLiquidStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-stake [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Liquid-stake coin",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Liquid-stake coin. 
			
Example:
$ %s tx %s liquid-stake 1000stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liquidStaker := clientCtx.GetFromAddress()

			stakingCoin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidStake(liquidStaker, stakingCoin)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewLiquidUnstakeCmd implements the liquid unstake coin command handler.
func NewLiquidUnstakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-unstake [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Liquid-unstake coin",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Liquid-unstake coin. 
			
Example:
$ %s tx %s liquid-unstake 500stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liquidStaker := clientCtx.GetFromAddress()

			unstakingCoin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidUnstake(liquidStaker, unstakingCoin)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUpdateParamsCmd implements the liquid unstake coin command handler.
func NewUpdateParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [params.json]",
		Args:  cobra.ExactArgs(1),
		Short: "Update-params for lspersistence",
		Long: strings.TrimSpace(
			fmt.Sprintf(`update-params param-file. 
			
Example:
$ %s tx %s update-params ~/params.json --from mykey

Example params.json 
{
  "liquid_bond_denom": "stk/uxprt",
  "whitelisted_validators": [
    {
      "validator_address": "persistencevaloper1hcqg5wj9t42zawqkqucs7la85ffyv08lmnhye9",
      "target_weight": "10"
    }
  ],
  "stake_fee_rate": "0.000000000000000000",
  "unstake_fee_rate": "0.000000000000000000",
  "redemption_fee_rate": "0.025000000000000000",
  "restake_fee_rate": "0.050000000000000000",
  "min_liquid_staking_amount": "10000",
  "admin_address": "persistence1kk3vjcjsvy3kd6389lavdkt5f2h5k3d2ry296c",
  "fee_address": "persistence1kk3vjcjsvy3kd6389lavdkt5f2h5k3d2ry296c"
}

Example for msg for submit-proposal v0.46 onwards
{
  "@type": "/pstake.lspersistence.v1beta1.MsgUpdateParams",
  "authority": "persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu",
  "params": {
    "liquid_bond_denom": "stk/uxprt",
    "whitelisted_validators": [
      {
        "validator_address": "persistencevaloper1hcqg5wj9t42zawqkqucs7la85ffyv08lmnhye9",
        "target_weight": "10"
      }
    ],
    "stake_fee_rate": "0.000000000000000000",
    "unstake_fee_rate": "0.000000000000000000",
    "redemption_fee_rate": "0.025000000000000000",
    "restake_fee_rate": "0.050000000000000000",
    "min_liquid_staking_amount": "10000",
    "admin_address": "persistence1kk3vjcjsvy3kd6389lavdkt5f2h5k3d2ry296c",
    "fee_address": "persistence1kk3vjcjsvy3kd6389lavdkt5f2h5k3d2ry296c"
  }
}
`,
				version.AppName, types.ModuleName,
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
