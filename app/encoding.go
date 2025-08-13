package app

import (
	"github.com/cosmos/cosmos-sdk/std"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v4/x/liquidstakeibc/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v4/x/lscosmos/types"
	ratesynctypes "github.com/persistenceOne/pstake-native/v4/x/ratesync/types"

	"github.com/persistenceOne/pstake-native/v4/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	ratesynctypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ratesynctypes.RegisterCodec(encodingConfig.Amino)

	liquidstakeibctypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	liquidstakeibctypes.RegisterLegacyAminoCodec(encodingConfig.Amino)

	lscosmostypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	lscosmostypes.RegisterLegacyAminoCodec(encodingConfig.Amino)

	return encodingConfig
}
