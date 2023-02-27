package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/client/utils"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
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
		NewLiquidStakeCmd(),
		NewLiquidUnstakeCmd(),
		NewRedeemCmd(),
		NewClaimCmd(),
		NewJumpStartCmd(),
		NewRecreateICACmd(),
		NewChangeModuleStateCmd(),
	)

	return cmd
}

func NewLiquidStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-stake [amount(allow-listed-ibcDenom coin)]",
		Short: `Liquid Stake ibc/Atom to stk/Atom`,
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

func NewMinDepositAndFeeChangeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pstake-lscosmos-min-deposit-and-fee-change",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a minimum deposit and fee change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a min-deposit and fee change proposal along with an initial deposit
The proposal details must be supplied via a JSON file. For values that contains objects,
only non-empty fields will be updated.

IMPORTANT : The values for the fields in this proposal are not validated, so it is very
important that any value change is valid.

Example Proposal :
{
	"title": "min-deposit and fee change proposal",
	"description": "this proposal changes min-deposit and protocol fee on chain",
	"min_deposit": "5",
	"pstake_deposit_fee": "0.1",
	"pstake_restake_fee": "0.1",
	"pstake_unstake_fee": "0.1",
	"pstake_redemption_fee": "0.1",
	"deposit": "100stake"
}

Example:
$ %s tx gov submit-proposal pstake-lscosmos-min-deposit-and-fee-change  <path/to/proposal.json> --from <key_or_address> --fees <1000stake> --gas <200000>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := utils.ParseMinDepositAndFeeChangeProposalJSON(clientCtx.LegacyAmino, args[0])
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
			redemptionFee, err := sdk.NewDecFromStr(proposal.PstakeRedemptionFee)
			if err != nil {
				return err
			}

			content := types.NewMinDepositAndFeeChangeProposal(
				proposal.Title,
				proposal.Description,
				minDeposit,
				depositFee,
				restakeFee,
				unstakeFee,
				redemptionFee,
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

func NewPstakeFeeAddressChangeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pstake-lscosmos-change-pstake-fee-address [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a pstake fee address change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a pstake fee address change proposal along with an initial deposit
The proposal details must be supplied via a JSON file. For values that contains objects,
only non-empty fields will be updated.

IMPORTANT : The values for the fields in this proposal are not validated, so it is very
important that any value change is valid.

Example Proposal :
{
	"title": "change pstake fee address",
	"description": "this proposal changes pstake fee address in the chain",
	"pstake_fee_address" : "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9"
	"deposit": "100stake"
}

Example:
$ %s tx gov submit-proposal pstake-lscosmos-change-pstake-fee-address <path/to/proposal.json> --from <key_or_address> --fees <1000stake> --gas <200000>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := utils.ParsePstakeFeeAddressChangeProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			pstakeFeeAddress, err := sdk.AccAddressFromBech32(proposal.PstakeFeeAddress)
			if err != nil {
				return err
			}

			content := types.NewPstakeFeeAddressChangeProposal(
				proposal.Title,
				proposal.Description,
				pstakeFeeAddress.String(),
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

func NewAllowListedValidatorSetChangeProposalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pstake-lscosmos-change-allow-listed-validator-set [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a AllowListed Validator set change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a AllowListed Validator set change proposal along with an initial deposit
The proposal details must be supplied via a JSON file. For values that contains objects,
only non-empty fields will be updated.

IMPORTANT : The values for the fields in this proposal are not validated, so it is very
important that any value change is valid.

Example Proposal :
{
	"title": "change pstake fee address",
	"description": "this proposal changes pstake fee address in the chain",
	"allow_listed_validators": {
   		 "allow_listed_validators": [
      {
        "validator_address": "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
        "target_weight": "1"
      }
    ]
  },
"deposit": "100stake"
}

Example:
$ %s tx gov submit-proposal pstake-lscosmos-change-allow-listed-validator-set <path/to/proposal.json> --from <key_or_address> --fees <1000stake> --gas <200000>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := utils.ParseAllowListedValidatorSetChangeProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewAllowListedValidatorSetChangeProposal(
				proposal.Title,
				proposal.Description,
				proposal.AllowListedValidators,
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

func NewLiquidUnstakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-unstake [amount(stk/Atom)]",
		Short: `Liquid Unstake stkAtom to ibc/Atom`,
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
			msg := types.NewMsgLiquidUnstake(delegatorAddress, amount)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewRedeemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem [amount(stkDenom)]",
		Short: `Instantly redeem skt tokens`,
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

			redeemAddress := clientctx.GetFromAddress()
			msg := types.NewMsgRedeem(redeemAddress, amount)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim",
		Short: `Claim matured tokens`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delegatorAddress := clientctx.GetFromAddress()
			msg := types.NewMsgClaim(delegatorAddress)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewRecreateICACmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recreate-ica",
		Short: "Recreate ICA accounts",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromAddress := clientctx.GetFromAddress()
			msg := types.NewMsgRecreateICA(fromAddress)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewJumpStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jump-start",
		Short: "jump start the module using allowlisted address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pstakeAddress := clientctx.GetFromAddress()

			msgDetails, err := utils.ParseJumpstartTxnJSON(clientctx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			minDeposit, ok := sdk.NewIntFromString(msgDetails.MinDeposit)
			if !ok {
				return types.ErrInvalidIntParse
			}
			depositFee, err := sdk.NewDecFromStr(msgDetails.PstakeParams.PstakeDepositFee)
			if err != nil {
				return err
			}
			restakeFee, err := sdk.NewDecFromStr(msgDetails.PstakeParams.PstakeRestakeFee)
			if err != nil {
				return err
			}
			unstakeFee, err := sdk.NewDecFromStr(msgDetails.PstakeParams.PstakeUnstakeFee)
			if err != nil {
				return err
			}
			redemptionFee, err := sdk.NewDecFromStr(msgDetails.PstakeParams.PstakeRedemptionFee)
			if err != nil {
				return err
			}
			pstakeParams := types.PstakeParams{
				PstakeDepositFee:    depositFee,
				PstakeRestakeFee:    restakeFee,
				PstakeUnstakeFee:    unstakeFee,
				PstakeRedemptionFee: redemptionFee,
				PstakeFeeAddress:    pstakeAddress.String(),
			}

			if types.ConvertBaseDenomToMintDenom(msgDetails.BaseDenom) != msgDetails.MintDenom {
				return types.ErrInvalidMintDenom
			}

			msg := types.NewMsgJumpStart(pstakeAddress, msgDetails.ChainID, msgDetails.ConnectionID, msgDetails.TransferChannel,
				msgDetails.TransferPort, msgDetails.BaseDenom, msgDetails.MintDenom, minDeposit, msgDetails.AllowListedValidators,
				pstakeParams, msgDetails.HostAccounts)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewChangeModuleStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-module-state",
		Short: "Admin functionality to disable/enable the functionality incase of failures",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			moduleState, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			fromAddress := clientctx.GetFromAddress()
			msg := types.NewMsgChangeModuleState(fromAddress, moduleState)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
