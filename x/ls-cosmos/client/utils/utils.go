package utils

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type RegisterCosmosChainProposalReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	// todo : check for ibc connection type
	IBCConnection        string         `json:"ibc_connection" yaml:"ibc_connection"`
	TokenTransferChannel string         `json:"token_transfer_channel" yaml:"token_transfer_channel"`
	TokenTransferPort    string         `json:"token_transfer_port" yaml:"token_transfer_port"`
	BaseDenom            string         `json:"base_denom" yaml:"base_denom"`
	MintDenom            string         `json:"mint_denom" yaml:"mint_denom"`
	Proposer             sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Deposit              sdk.Coins      `json:"deposit" yaml:"deposit"`
}

type RegisterCosmosChainProposalJSON struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	// todo : check for ibc connection type
	IBCConnection        string `json:"ibc_connection" yaml:"ibc_connection"`
	TokenTransferChannel string `json:"token_transfer_channel" yaml:"token_transfer_channel"`
	TokenTransferPort    string `json:"token_transfer_port" yaml:"token_transfer_port"`
	BaseDenom            string `json:"base_denom" yaml:"base_denom"`
	MintDenom            string `json:"mint_denom" yaml:"mint_denom"`
	Deposit              string `json:"deposit" yaml:"deposit"`
}

func NewRegisterChainJSON(Title, Description, IBCConnection, TokenTransferChannel, TokenTransferPort, BaseDenom, MintDenom, Deposit string) RegisterCosmosChainProposalJSON {
	return RegisterCosmosChainProposalJSON{
		Title:                Title,
		Description:          Description,
		IBCConnection:        IBCConnection,
		TokenTransferChannel: TokenTransferChannel,
		TokenTransferPort:    TokenTransferPort,
		BaseDenom:            BaseDenom,
		MintDenom:            MintDenom,
		Deposit:              Deposit,
	}
}

// ParseRegisterCosmosChainProposalJSON reads and parses a RegisterCosmosChainProposalJSON from
// file.
func ParseRegisterCosmosChainProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (RegisterCosmosChainProposalJSON, error) {
	proposal := RegisterCosmosChainProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
