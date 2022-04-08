package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	//cdc.RegisterConcrete(&MsgSend{}, "cosmos-sdk/MsgSend", nil)
	cdc.RegisterConcrete(&MsgSetOrchestrator{}, "cosmos/MsgSetOrchestrator", nil)
	cdc.RegisterConcrete(&MsgWithdrawStkAsset{}, "cosmos/MsgWithdrawStkAsset", nil)
	cdc.RegisterConcrete(&MsgMintTokensForAccount{}, "cosmos/MsgMintTokensForAccount", nil)
	cdc.RegisterConcrete(&MsgMakeProposal{}, "cosmos/MsgMakeProposal", nil)
	cdc.RegisterConcrete(&MsgVote{}, "cosmos/MsgVote", nil)
	cdc.RegisterConcrete(&MsgVoteWeighted{}, "cosmos/MsgVoteWeighted", nil)
	cdc.RegisterConcrete(&MsgSignedTx{}, "cosmos/MsgSignedTx", nil)
	cdc.RegisterConcrete(&MsgTxStatus{}, "cosmos/MsgTxStatus", nil)
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
