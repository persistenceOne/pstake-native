package orchestrator

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/persistenceOne/pstake-native/app"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

func (n *NativeChain) MakeEncodingConfig() params.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := n.NewProtoCodec(interfaceRegistry, n.AccountPrefix)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}

	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	app.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	app.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

func (c *CosmosChain) MakeEncodingConfig() params.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := c.NewProtoCodec(interfaceRegistry, c.AccountPrefix)
	txCfg := tx.NewTxConfig(marshaler, []signingtypes.SignMode{signingtypes.SignMode_SIGN_MODE_DIRECT, signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON})

	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}

	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	simapp.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	simapp.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

type ProtoCodec struct {
	interfaceRegistry types.InterfaceRegistry
	useContext        func() func()
}

var _ codec.Codec = &ProtoCodec{}
var _ codec.ProtoCodecMarshaler = &ProtoCodec{}

func (c *CosmosChain) NewProtoCodec(interfaceRegistry types.InterfaceRegistry, accountPrefix string) *ProtoCodec {
	return &ProtoCodec{interfaceRegistry: interfaceRegistry, useContext: c.UseSDKContext}
}

func (n *NativeChain) NewProtoCodec(interfaceRegistry types.InterfaceRegistry, accountPrefix string) *ProtoCodec {
	return &ProtoCodec{interfaceRegistry: interfaceRegistry, useContext: n.UseSDKContext}
}

func (pc *ProtoCodec) Marshal(o codec.ProtoMarshaler) ([]byte, error) {
	defer pc.useContext()()
	return o.Marshal()
}

func (pc *ProtoCodec) MustMarshal(o codec.ProtoMarshaler) []byte {
	bz, err := pc.Marshal(o)
	if err != nil {
		panic(err)
	}

	return bz
}

func (pc *ProtoCodec) MarshalLengthPrefixed(o codec.ProtoMarshaler) ([]byte, error) {
	defer pc.useContext()()
	bz, err := pc.Marshal(o)
	if err != nil {
		return nil, err
	}

	var sizeBuf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(sizeBuf[:], uint64(o.Size()))
	return append(sizeBuf[:n], bz...), nil
}

func (pc *ProtoCodec) MustMarshalLengthPrefixed(o codec.ProtoMarshaler) []byte {
	bz, err := pc.MarshalLengthPrefixed(o)
	if err != nil {
		panic(err)
	}

	return bz
}

func (pc *ProtoCodec) Unmarshal(bz []byte, ptr codec.ProtoMarshaler) error {
	defer pc.useContext()()
	err := ptr.Unmarshal(bz)
	if err != nil {
		return err
	}
	err = types.UnpackInterfaces(ptr, pc)
	if err != nil {
		return err
	}
	return nil
}

func (pc *ProtoCodec) MustUnmarshal(bz []byte, ptr codec.ProtoMarshaler) {
	if err := pc.Unmarshal(bz, ptr); err != nil {
		panic(err)
	}
}

func (pc *ProtoCodec) UnmarshalLengthPrefixed(bz []byte, ptr codec.ProtoMarshaler) error {
	defer pc.useContext()()
	size, n := binary.Uvarint(bz)
	if n < 0 {
		return fmt.Errorf("invalid number of bytes read from length-prefixed encoding: %d", n)
	}

	if size > uint64(len(bz)-n) {
		return fmt.Errorf("not enough bytes to read; want: %v, got: %v", size, len(bz)-n)
	} else if size < uint64(len(bz)-n) {
		return fmt.Errorf("too many bytes to read; want: %v, got: %v", size, len(bz)-n)
	}

	bz = bz[n:]
	return pc.Unmarshal(bz, ptr)
}

func (pc *ProtoCodec) MustUnmarshalLengthPrefixed(bz []byte, ptr codec.ProtoMarshaler) {
	if err := pc.UnmarshalLengthPrefixed(bz, ptr); err != nil {
		panic(err)
	}
}

func (pc *ProtoCodec) MarshalJSON(o proto.Message) ([]byte, error) {
	done := pc.useContext()
	m, ok := o.(codec.ProtoMarshaler)
	if !ok {
		return nil, fmt.Errorf("cannot protobuf JSON encode unsupported type: %T", o)
	}
	bz, err := codec.ProtoMarshalJSON(m, pc.interfaceRegistry)
	if err != nil {
		return []byte{}, err
	}
	done()
	return bz, nil
}

func (pc *ProtoCodec) MustMarshalJSON(o proto.Message) []byte {
	bz, err := pc.MarshalJSON(o)
	if err != nil {
		panic(err)
	}

	return bz
}

func (pc *ProtoCodec) UnmarshalJSON(bz []byte, ptr proto.Message) error {
	defer pc.useContext()()
	m, ok := ptr.(codec.ProtoMarshaler)
	if !ok {
		return fmt.Errorf("cannot protobuf JSON decode unsupported type: %T", ptr)
	}

	err := jsonpb.Unmarshal(strings.NewReader(string(bz)), m)
	if err != nil {
		return err
	}

	return types.UnpackInterfaces(ptr, pc)
}

func (pc *ProtoCodec) MustUnmarshalJSON(bz []byte, ptr proto.Message) {
	if err := pc.UnmarshalJSON(bz, ptr); err != nil {
		panic(err)
	}
}

func (pc *ProtoCodec) MarshalInterface(i proto.Message) ([]byte, error) {
	defer pc.useContext()()
	if err := assertNotNil(i); err != nil {
		return nil, err
	}
	any, err := types.NewAnyWithValue(i)
	if err != nil {
		return nil, err
	}

	return pc.Marshal(any)
}

func (pc *ProtoCodec) UnmarshalInterface(bz []byte, ptr interface{}) error {
	any := &types.Any{}
	err := pc.Unmarshal(bz, any)
	if err != nil {
		return err
	}

	return pc.UnpackAny(any, ptr)
}

func (pc *ProtoCodec) MarshalInterfaceJSON(x proto.Message) ([]byte, error) {
	defer pc.useContext()()
	any, err := types.NewAnyWithValue(x)
	if err != nil {
		return nil, err
	}
	return pc.MarshalJSON(any)
}

func (pc *ProtoCodec) UnmarshalInterfaceJSON(bz []byte, iface interface{}) error {
	any := &types.Any{}
	err := pc.UnmarshalJSON(bz, any)
	if err != nil {
		return err
	}
	return pc.UnpackAny(any, iface)
}

func (pc *ProtoCodec) UnpackAny(any *types.Any, iface interface{}) error {
	defer pc.useContext()()
	return pc.interfaceRegistry.UnpackAny(any, iface)
}

func (pc *ProtoCodec) InterfaceRegistry() types.InterfaceRegistry {
	return pc.interfaceRegistry
}

func assertNotNil(i interface{}) error {
	if i == nil {
		return errors.New("can't marshal <nil> value")
	}
	return nil
}
