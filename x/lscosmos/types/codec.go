package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&RegisterHostChainProposal{}, "cosmos/RegisterHostChainProposal", nil)
	cdc.RegisterConcrete(&MinDepositAndFeeChangeProposal{}, "cosmos/MinDepositAndFeeChangeProposal", nil)
	cdc.RegisterConcrete(&PstakeFeeAddressChangeProposal{}, "cosmos/PstakeFeeAddressChangeProposal", nil)
	cdc.RegisterConcrete(&AllowListedValidatorSetChangeProposal{}, "cosmos/AllowListedValidatorSetChangeProposal", nil)
	cdc.RegisterConcrete(&MsgLiquidStake{}, "cosmos/MsgLiquidStake", nil)
	cdc.RegisterConcrete(&MsgJuice{}, "cosmos/MsgJuice", nil)
	cdc.RegisterConcrete(&MsgLiquidUnstake{}, "cosmos/MsgLiquidUnstake", nil)
	cdc.RegisterConcrete(&MsgRedeem{}, "cosmos/MsgRedeem", nil)
	cdc.RegisterConcrete(&MsgClaim{}, "cosmos/MsgClaim", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgLiquidStake{},
		&MsgJuice{},
		&MsgLiquidUnstake{},
		&MsgRedeem{},
		&MsgClaim{},
	) // add the structs that implements sdk.Msg interface

	registry.RegisterImplementations((*govtypes.Content)(nil),
		// add the stucts that implements govTypes.Content interface
		&RegisterHostChainProposal{},
		&MinDepositAndFeeChangeProposal{},
		&PstakeFeeAddressChangeProposal{},
		&AllowListedValidatorSetChangeProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
