package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers all the necessary types and interfaces for the
// cosmos module.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSetOrchestrator{}, "cosmos/MsgSetOrchestrator", nil)
	cdc.RegisterConcrete(&MsgWithdrawStkAsset{}, "cosmos/MsgWithdrawStkAsset", nil)
	cdc.RegisterConcrete(&MsgMintTokensForAccount{}, "cosmos/MsgMintTokensForAccount", nil)
	cdc.RegisterConcrete(&MsgMakeProposal{}, "cosmos/MsgMakeProposal", nil)
	cdc.RegisterConcrete(&MsgVote{}, "cosmos/MsgVote", nil)
	cdc.RegisterConcrete(&MsgVoteWeighted{}, "cosmos/MsgVoteWeighted", nil)
	cdc.RegisterConcrete(&MsgSignedTx{}, "cosmos/MsgSignedTx", nil)
	cdc.RegisterConcrete(&MsgTxStatus{}, "cosmos/MsgTxStatus", nil)
	cdc.RegisterConcrete(&MsgUndelegateSuccess{}, "cosmos/MsgUndelegateSuccess", nil)
	cdc.RegisterConcrete(&MsgSetSignature{}, "cosmos/MsgSetSignature", nil)
	cdc.RegisterConcrete(&EnableModuleProposal{}, "cosmos/EnableModuleProposal", nil)
	cdc.RegisterConcrete(&ChangeMultisigProposal{}, "cosmos/ChangeMultisigProposal", nil)
	cdc.RegisterConcrete(&ChangeCosmosValidatorWeightsProposal{}, "cosmos/ChangeCosmosValidatorWeightsProposal", nil)
	cdc.RegisterConcrete(&ChangeOracleValidatorWeightsProposal{}, "cosmos/ChangeOracleValidatorWeightsProposal", nil)
}

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetOrchestrator{},
		&MsgWithdrawStkAsset{},
		&MsgMintTokensForAccount{},
		&MsgMakeProposal{},
		&MsgVote{},
		&MsgVoteWeighted{},
		&MsgSignedTx{},
		&MsgTxStatus{},
		&MsgSetSignature{},
		&MsgUndelegateSuccess{},
	)

	registry.RegisterImplementations((*govtypes.Content)(nil),
		&EnableModuleProposal{},
		&ChangeMultisigProposal{},
		&ChangeCosmosValidatorWeightsProposal{},
		&ChangeOracleValidatorWeightsProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptoCodec.RegisterCrypto(amino)
	amino.Seal()
}
