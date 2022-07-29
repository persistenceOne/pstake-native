package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/pstake-native/x/ls-cosmos/client/utils"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
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
		NewRegisterCosmosChainCmd(),
		NewLiquidStakeCmd(),
	)

	return cmd
}

func NewRegisterCosmosChainCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "register-cosmos-chain [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a register cosmos chain proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a register cosmos chain proposal along with an initial deposit
The proposal details must be supplied via a JSON file. For values that contains objects,
only non-empty fields will be updated.

IMPORTANT : The values for the fields in this proposal are not validated, so it is very
important that any value change is valid.

Example:
$ %s tx gov submit-proposal register-cosmos-chain <path/to/proposal.json> --from <key_or_address>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := utils.ParseRegisterCosmosChainProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewRegisterCosmosChainProposal(
				proposal.Title,
				proposal.Description,
				proposal.IBCConnection,
				proposal.TokenTransferChannel,
				proposal.TokenTransferPort,
				proposal.BaseDenom,
				proposal.MintDenom,
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
		Use:   "liquid-stake [amount(whitelisted-ibcDenom coin)] [mint-address] [deposit-address] ",
		Short: `Liquid Stake ibc/Atom to stkAtom`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			mintAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			depositAddress, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidStake(amount, mintAddress, depositAddress)

			return tx.GenerateOrBroadcastTxCLI(clientctx, cmd.Flags(), msg)

		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
