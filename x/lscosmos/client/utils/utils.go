package utils

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

type RegisterCosmosChainProposalReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	// TODO update
	ModuleEnabled         bool                        `json:"module_enabled" yaml:"module_enabled"`
	ConnectionID          string                      `json:"connectionID" yaml:"connectionID"`
	TransferChannel       string                      `json:"transfer_channel" yaml:"transfer_channel"`
	TransferPort          string                      `json:"transfer_port" yaml:"transfer_port"`
	BaseDenom             string                      `json:"base_denom" yaml:"base_denom"`
	MintDenom             string                      `json:"mint_denom" yaml:"mint_denom"`
	MinDeposit            string                      `json:"min_deposit" yaml:"min_deposit"`
	AllowListedValidators types.AllowListedValidators `json:"allow_listed_validators" yaml:"allow_listed_validators"`
	PstakeDepositFee      string                      `json:"pstake_deposit_fee" yaml:"pstake_deposit_fee"`
	PstakeRestakeFee      string                      `json:"pstake_restake_fee" yaml:"pstake_restake_fee"`
	PstakeUnstakeFee      string                      `json:"pstake_unstake_fee" yaml:"pstake_unstake_fee"`
	Proposer              sdk.AccAddress              `json:"proposer" yaml:"proposer"`
	Deposit               sdk.Coins                   `json:"deposit" yaml:"deposit"`
}

type RegisterCosmosChainProposalJSON struct {
	Title                 string                      `json:"title" yaml:"title"`
	Description           string                      `json:"description" yaml:"description"`
	ModuleEnabled         bool                        `json:"module_enabled" yaml:"module_enabled"`
	ConnectionID          string                      `json:"connectionID" yaml:"connectionID"`
	TransferChannel       string                      `json:"transfer_channel" yaml:"transfer_channel"`
	TransferPort          string                      `json:"transfer_port" yaml:"transfer_port"`
	BaseDenom             string                      `json:"base_denom" yaml:"base_denom"`
	MintDenom             string                      `json:"mint_denom" yaml:"mint_denom"`
	MinDeposit            string                      `json:"min_deposit" yaml:"min_deposit"`
	AllowListedValidators types.AllowListedValidators `json:"allow_listed_validators" yaml:"allow_listed_validators"`
	PstakeDepositFee      string                      `json:"pstake_deposit_fee" yaml:"pstake_deposit_fee"`
	PstakeRestakeFee      string                      `json:"pstake_restake_fee" yaml:"pstake_restake_fee"`
	PstakeUnstakeFee      string                      `json:"pstake_unstake_fee" yaml:"pstake_unstake_fee"`
	Deposit               string                      `json:"deposit" yaml:"deposit"`
}

func NewRegisterChainJSON(title, description string, moduleEnabled bool, connectionID, transferChannel, transferPort,
	baseDenom, mintDenom, minDeposit string, allowListedValidators types.AllowListedValidators, pstakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee, deposit string) RegisterCosmosChainProposalJSON {
	return RegisterCosmosChainProposalJSON{
		Title:                 title,
		Description:           description,
		ModuleEnabled:         moduleEnabled,
		ConnectionID:          connectionID,
		TransferChannel:       transferChannel,
		TransferPort:          transferPort,
		BaseDenom:             baseDenom,
		MintDenom:             mintDenom,
		MinDeposit:            minDeposit,
		AllowListedValidators: allowListedValidators,
		PstakeDepositFee:      pstakeDepositFee,
		PstakeRestakeFee:      pstakeRestakeFee,
		PstakeUnstakeFee:      pstakeUnstakeFee,
		Deposit:               deposit,
	}
}

// ParseRegisterCosmosChainProposalJSON reads and parses a RegisterCosmosChainProposalJSON from
// file.
func ParseRegisterCosmosChainProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (RegisterCosmosChainProposalJSON, error) {
	proposal := RegisterCosmosChainProposalJSON{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
