package cli

import (
	"fmt"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/persistenceOne/pstake-native/x/cosmos/client/utils"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/spf13/cobra"
)

// Proposal flags
const (
	flagStatus = "status"
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
		CmdSendNewProposal(),
		NewCmdVote(),
		NewCmdWeightedVote(),
		NewCmdTxStatusCmd(),
		NewWithdrawCmd(),
		NewRewardsClaimedCmd(),
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

			coins, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			chainID := args[3]

			txHash := args[4]

			blockHeight, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				return err
			}

			msg := cosmosTypes.NewMsgMintTokensForAccount(toAddr, orchAddress, coins, txHash, chainID, blockHeight)

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

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSendNewProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-proposal [title] [description] [orchestrator-address] [proposal-id] [chain-id] [block-height]",
		Short: "Allows orchestrator to send any proposal created on cosmos chain.",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title := args[0]

			description := args[1]

			orchAddress, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			proposalID, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			chainID := args[4]

			blockHeight, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				return err
			}

			votingStartTime := time.Now()
			votingEndTime := votingStartTime.Add(time.Minute * 2)

			msg := cosmosTypes.NewMsgMakeProposal(title, description, orchAddress, chainID, blockHeight, proposalID, votingStartTime, votingEndTime)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewCmdVote implements creating a new vote command.
func NewCmdVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [option]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for an active proposal, options: yes/no/no_with_veto/abstain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a vote for an active proposal. You can
find the proposal-id by running "%s query gov proposals".

Example:
$ %s tx gov vote 1 yes --from mykey
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// Get voting address
			from := clientCtx.GetFromAddress()

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// Find out which vote option user chose
			byteVoteOption, err := cosmosTypes.VoteOptionFromString(utils.NormalizeVoteOption(args[1]))
			if err != nil {
				return err
			}

			// Build vote message and run basic validation
			msg := cosmosTypes.NewMsgVote(from, proposalID, byteVoteOption)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewCmdWeightedVote implements creating a new weighted vote command.
func NewCmdWeightedVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "weighted-vote [proposal-id] [weighted-options]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for an active proposal, options: yes/no/no_with_veto/abstain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a vote for an active proposal. You can
find the proposal-id by running "%s query gov proposals".

Example:
$ %s tx gov weighted-vote 1 yes=0.6,no=0.3,abstain=0.05,no_with_veto=0.05 --from mykey
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Get voter address
			from := clientCtx.GetFromAddress()

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// Figure out which vote options user chose
			options, err := cosmosTypes.WeightedVoteOptionsFromString(utils.NormalizeWeightedVoteOptions(args[1]))
			if err != nil {
				return err
			}

			// Build vote message and run basic validation
			msg := cosmosTypes.NewMsgVoteWeighted(from, proposalID, options)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCmdTxStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx-status [orchestrator-address] [tx-hash] [status] [account-number] [sequence-number] [balance]",
		Args:  cobra.ExactArgs(6),
		Short: "Send status for transaction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit status for transaction relayed to cosmos chain.
Only "success" or "failure" accepted as status.`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			orchAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			txHash := args[1]

			status := args[2]

			accountNumber, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			sequenceNumber, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			balance, err := sdk.ParseCoinsNormalized(args[5])
			if err != nil {
				return err
			}

			// todo parse validator details in json file
			msg := cosmosTypes.NewMsgTxStatus(orchAddress, status, txHash, accountNumber, sequenceNumber, balance, []cosmosTypes.ValidatorDetails{})

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewWithdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [from-address] [to-address] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Withdraw transaction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit destination address on cosmos chain for uatom withdrawal`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			toAddress, err := cosmosTypes.AccAddressFromBech32(args[1], cosmosTypes.Bech32Prefix)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := cosmosTypes.NewMsgWithdrawStkAsset(fromAddress, toAddress, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewRewardsClaimedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards-claimed [orchestrator-address] [amount_claimed] [chainID] [block-height]",
		Args:  cobra.ExactArgs(4),
		Short: "Rewards claimed transaction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit amount claimed on other chain to be re staked`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			orchAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			chainID := args[2]

			blockHeight, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			msg := cosmosTypes.NewMsgRewardsClaimedOnCosmosChain(orchAddr, amount, chainID, blockHeight)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewEnableModuleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable-module [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a module enable proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a module-enable proposal along with an initial deposit.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := utils.ParseEnableModuleProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := cosmosTypes.NewEnableModuleProposal(proposal.Title, proposal.Description, proposal.Threshold, proposal.AccountNumber)

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

func NewChangeMultisigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "change-multisig [proposal-file] ",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a multisig change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a change-multisig proposal along with an initial deposit.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := utils.ParseChangeMultisigProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := cosmosTypes.NewChangeMultisigProposal(proposal.Title, proposal.Description, proposal.Threshold, proposal.OrchestratorAddresses, proposal.AccountNumber)

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

func NewChangeCosmosValidatorWeightsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "change-cosmos-validator-weights [proposal-file] ",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a cosmos validator weights change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a cosmos validator weights proposal along with an initial deposit.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("entering cosmos validator set proposal")
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := utils.ParseChangeCosmosValidatorWeightsProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			var weightedAddresses []cosmosTypes.WeightedAddressCosmos

			for _, weightedAddress := range proposal.WeightedAddresses {
				weight, err := sdk.NewDecFromStr(weightedAddress.Weight)
				if err != nil {
					return err
				}
				weightedAddresses = append(
					weightedAddresses,
					cosmosTypes.WeightedAddressCosmos{
						Address: weightedAddress.ValAddress,
						Weight:  weight,
					})
			}

			from := clientCtx.GetFromAddress()
			content := cosmosTypes.NewChangeCosmosValidatorWeightsProposal(proposal.Title, proposal.Description, weightedAddresses)

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

func NewChangeOracleValidatorWeightsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "change-oracle-validator-weights [proposal-file] ",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a oracle validator weights change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a oracle validator weights proposal along with an initial deposit.`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := utils.ParseChangeOracleValidatorWeightsProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			var weightedAddresses []cosmosTypes.WeightedAddress

			for _, weightedAddress := range proposal.WeightedAddresses {
				weight, err := sdk.NewDecFromStr(weightedAddress.Weight)
				if err != nil {
					return err
				}
				weightedAddresses = append(
					weightedAddresses,
					cosmosTypes.WeightedAddress{
						Address: weightedAddress.ValAddress,
						Weight:  weight,
					})
			}

			from := clientCtx.GetFromAddress()
			content := cosmosTypes.NewChangeOracleValidatorWeightsProposal(proposal.Title, proposal.Description, weightedAddresses)

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
