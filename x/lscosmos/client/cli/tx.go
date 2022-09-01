package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/pstake-native/x/lscosmos/client/utils"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		NewRegisterHostChainCmd(),
		NewLiquidStakeCmd(),
	)

	return cmd
}

func NewRegisterHostChainCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "register-host-chain [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a register host chain proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a register host chain proposal along with an initial deposit
The proposal details must be supplied via a JSON file. For values that contains objects,
only non-empty fields will be updated.

IMPORTANT : The values for the fields in this proposal are not validated, so it is very
important that any value change is valid.

Example Proposal :
{
	"title": "register host chain proposal",
	"description": "this proposal register host chain params in the chain",
	"connection_i_d": "test connection",
	"transfer_channel": "test-channel-1",
	"transfer_port": "test-transfer",
	"base_denom": "uatom",
	"mint_denom": "ustkatom",
	"min_deposit": "5",
	"pstake_deposit_fee": "0.1",
	"pstake_restake_fee": "0.1",
	"pstake_unstake_fee": "0.1",
	"deposit": "100stake"
}

Example:
$ %s tx gov submit-proposal register-host-chain <path/to/proposal.json> --from <key_or_address> --fees <1000stake> --gas <200000>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := utils.ParseRegisterHostChainProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			minDeposit, ok := sdk.NewIntFromString(proposal.MinDeposit)
			if !ok {
				return types.ErrInvalidIntParse
			}
			depositFee, err := sdk.NewDecFromStr(proposal.PstakeDepositFee)
			if err != nil {
				return err
			}

			restakeFee, err := sdk.NewDecFromStr(proposal.PstakeRestakeFee)
			if err != nil {
				return err
			}
			unstakeFee, err := sdk.NewDecFromStr(proposal.PstakeUnstakeFee)
			if err != nil {
				return err
			}

			content := types.NewRegisterHostChainProposal(
				proposal.Title,
				proposal.Description,
				proposal.ModuleEnabled,
				proposal.ChainID,
				proposal.ConnectionID,
				proposal.TransferChannel,
				proposal.TransferPort,
				proposal.BaseDenom,
				proposal.MintDenom,
				minDeposit,
				proposal.AllowListedValidators,
				depositFee,
				restakeFee,
				unstakeFee,
			)

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

func NewLiquidStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-stake [amount(whitelisted-ibcDenom coin)]",
		Short: `Liquid Stake ibc/Atom to stkAtom`,
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
