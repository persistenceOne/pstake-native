package app

import (
	"github.com/cosmos/cosmos-sdk/std"

	"github.com/persistenceOne/pstake-native/v5/app/params"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v5/x/liquidstakeibc/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v5/x/lscosmos/types"
	ratesynctypes "github.com/persistenceOne/pstake-native/v5/x/ratesync/types"
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
