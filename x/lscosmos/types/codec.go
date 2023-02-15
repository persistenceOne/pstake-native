package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers the necessary x/lscosmos interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MinDepositAndFeeChangeProposal{}, "cosmos/MinDepositAndFeeChangeProposal", nil)
	cdc.RegisterConcrete(&PstakeFeeAddressChangeProposal{}, "cosmos/PstakeFeeAddressChangeProposal", nil)
	cdc.RegisterConcrete(&AllowListedValidatorSetChangeProposal{}, "cosmos/AllowListedValidatorSetChangeProposal", nil)
	cdc.RegisterConcrete(&MsgLiquidStake{}, "cosmos/MsgLiquidStake", nil)
	cdc.RegisterConcrete(&MsgLiquidUnstake{}, "cosmos/MsgLiquidUnstake", nil)
	cdc.RegisterConcrete(&MsgRedeem{}, "cosmos/MsgRedeem", nil)
	cdc.RegisterConcrete(&MsgClaim{}, "cosmos/MsgClaim", nil)
	cdc.RegisterConcrete(&MsgRecreateICA{}, "cosmos/MsgRecreateICA", nil)
	cdc.RegisterConcrete(&MsgJumpStart{}, "cosmos/MsgJumpStart", nil)
	cdc.RegisterConcrete(&MsgChangeModuleState{}, "cosmos/MsgChangeModuleState", nil)
}

// RegisterInterfaces registers the x/lscosmos interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgLiquidStake{},
		&MsgLiquidUnstake{},
		&MsgRedeem{},
		&MsgClaim{},
		&MsgRecreateICA{},
		&MsgJumpStart{},
		&MsgChangeModuleState{},
	) // add the structs that implements sdk.Msg interface

	registry.RegisterImplementations((*govtypes.Content)(nil),
		// add the stucts that implements govTypes.Content interface
		&MinDepositAndFeeChangeProposal{},
		&PstakeFeeAddressChangeProposal{},
		&AllowListedValidatorSetChangeProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/lscosmos module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/lscosmos and
	// defined at the application level.
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	AminoCdc  = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	Amino.Seal()
}
