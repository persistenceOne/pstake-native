package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/liquidstakeibc interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgRegisterHostChain{}, "pstake/MsgRegisterHostChain")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateHostChain{}, "pstake/MsgUpdateHostChain")
	legacy.RegisterAminoMsg(cdc, &MsgLiquidStake{}, "pstake/MsgLiquidStake")
	legacy.RegisterAminoMsg(cdc, &MsgLiquidStakeLSM{}, "pstake/MsgLiquidStakeLSM")
	legacy.RegisterAminoMsg(cdc, &MsgLiquidUnstake{}, "pstake/MsgLiquidUnstake")
	legacy.RegisterAminoMsg(cdc, &MsgRedeem{}, "pstake/MsgRedeem")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "pstake/MsgUpdateParams")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterHostChain{},
		&MsgUpdateHostChain{},
		&MsgLiquidStake{},
		&MsgLiquidStakeLSM{},
		&MsgLiquidUnstake{},
		&MsgRedeem{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
