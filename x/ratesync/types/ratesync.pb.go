// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pstake/ratesync/v1beta1/ratesync.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	types "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type InstantiationState int32

const (
	// Not Initiated
	InstantiationState_INSTANTIATION_NOT_INITIATED InstantiationState = 0
	// Initiated
	InstantiationState_INSTANTIATION_INITIATED InstantiationState = 1
	// we should have an address
	InstantiationState_INSTANTIATION_COMPLETED InstantiationState = 2
)

var InstantiationState_name = map[int32]string{
	0: "INSTANTIATION_NOT_INITIATED",
	1: "INSTANTIATION_INITIATED",
	2: "INSTANTIATION_COMPLETED",
}

var InstantiationState_value = map[string]int32{
	"INSTANTIATION_NOT_INITIATED": 0,
	"INSTANTIATION_INITIATED":     1,
	"INSTANTIATION_COMPLETED":     2,
}

func (x InstantiationState) String() string {
	return proto.EnumName(InstantiationState_name, int32(x))
}

func (InstantiationState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_429540018f2469ab, []int{0}
}

type FeatureType int32

const (
	FeatureType_LIQUID_STAKE_IBC FeatureType = 0
	FeatureType_LIQUID_STAKE     FeatureType = 1
)

var FeatureType_name = map[int32]string{
	0: "LIQUID_STAKE_IBC",
	1: "LIQUID_STAKE",
}

var FeatureType_value = map[string]int32{
	"LIQUID_STAKE_IBC": 0,
	"LIQUID_STAKE":     1,
}

func (x FeatureType) String() string {
	return proto.EnumName(FeatureType_name, int32(x))
}

func (FeatureType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_429540018f2469ab, []int{1}
}

// HostChain defines the ratesync module's HostChain state.
type HostChain struct {
	// unique id
	ID                uint64           `protobuf:"varint,1,opt,name=i_d,json=iD,proto3" json:"i_d,omitempty"`
	ChainID           string           `protobuf:"bytes,2,opt,name=chain_i_d,json=chainID,proto3" json:"chain_i_d,omitempty"`
	ConnectionID      string           `protobuf:"bytes,3,opt,name=connection_i_d,json=connectionID,proto3" json:"connection_i_d,omitempty"`
	ICAAccount        types.ICAAccount `protobuf:"bytes,4,opt,name=i_c_a_account,json=iCAAccount,proto3" json:"i_c_a_account"`
	Features          Feature          `protobuf:"bytes,5,opt,name=features,proto3" json:"features"`
	TransferChannelID string           `protobuf:"bytes,6,opt,name=transfer_channel_i_d,json=transferChannelID,proto3" json:"transfer_channel_i_d,omitempty"`
	TransferPortID    string           `protobuf:"bytes,7,opt,name=transfer_port_i_d,json=transferPortID,proto3" json:"transfer_port_i_d,omitempty"`
}

func (m *HostChain) Reset()         { *m = HostChain{} }
func (m *HostChain) String() string { return proto.CompactTextString(m) }
func (*HostChain) ProtoMessage()    {}
func (*HostChain) Descriptor() ([]byte, []int) {
	return fileDescriptor_429540018f2469ab, []int{0}
}
func (m *HostChain) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HostChain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HostChain.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HostChain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HostChain.Merge(m, src)
}
func (m *HostChain) XXX_Size() int {
	return m.Size()
}
func (m *HostChain) XXX_DiscardUnknown() {
	xxx_messageInfo_HostChain.DiscardUnknown(m)
}

var xxx_messageInfo_HostChain proto.InternalMessageInfo

func (m *HostChain) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *HostChain) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

func (m *HostChain) GetConnectionID() string {
	if m != nil {
		return m.ConnectionID
	}
	return ""
}

func (m *HostChain) GetICAAccount() types.ICAAccount {
	if m != nil {
		return m.ICAAccount
	}
	return types.ICAAccount{}
}

func (m *HostChain) GetFeatures() Feature {
	if m != nil {
		return m.Features
	}
	return Feature{}
}

func (m *HostChain) GetTransferChannelID() string {
	if m != nil {
		return m.TransferChannelID
	}
	return ""
}

func (m *HostChain) GetTransferPortID() string {
	if m != nil {
		return m.TransferPortID
	}
	return ""
}

type Feature struct {
	// triggers on hooks
	LiquidStakeIBC LiquidStake `protobuf:"bytes,1,opt,name=liquid_stake_i_b_c,json=liquidStakeIBC,proto3" json:"liquid_stake_i_b_c"`
	// triggers on hour epoch
	LiquidStake LiquidStake `protobuf:"bytes,2,opt,name=liquid_stake,json=liquidStake,proto3" json:"liquid_stake"`
}

func (m *Feature) Reset()         { *m = Feature{} }
func (m *Feature) String() string { return proto.CompactTextString(m) }
func (*Feature) ProtoMessage()    {}
func (*Feature) Descriptor() ([]byte, []int) {
	return fileDescriptor_429540018f2469ab, []int{1}
}
func (m *Feature) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Feature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Feature.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Feature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Feature.Merge(m, src)
}
func (m *Feature) XXX_Size() int {
	return m.Size()
}
func (m *Feature) XXX_DiscardUnknown() {
	xxx_messageInfo_Feature.DiscardUnknown(m)
}

var xxx_messageInfo_Feature proto.InternalMessageInfo

func (m *Feature) GetLiquidStakeIBC() LiquidStake {
	if m != nil {
		return m.LiquidStakeIBC
	}
	return LiquidStake{}
}

func (m *Feature) GetLiquidStake() LiquidStake {
	if m != nil {
		return m.LiquidStake
	}
	return LiquidStake{}
}

type LiquidStake struct {
	FeatureType FeatureType `protobuf:"varint,1,opt,name=feature_type,json=featureType,proto3,enum=pstake.ratesync.v1beta1.FeatureType" json:"feature_type,omitempty"`
	// needs to be uploaded before hand
	CodeID uint64 `protobuf:"varint,2,opt,name=code_i_d,json=codeID,proto3" json:"code_i_d,omitempty"`
	// state of instantiation, do not support gov based instantiation. (need ICA
	// to be at least admin)
	Instantiation InstantiationState `protobuf:"varint,3,opt,name=instantiation,proto3,enum=pstake.ratesync.v1beta1.InstantiationState" json:"instantiation,omitempty"`
	// address of instantiated contract.
	ContractAddress string `protobuf:"bytes,4,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
	// allow * as default for all denoms in case of lsibc, or default bond denom
	// in case of ls.
	Denoms  []string `protobuf:"bytes,5,rep,name=denoms,proto3" json:"denoms,omitempty"`
	Enabled bool     `protobuf:"varint,6,opt,name=enabled,proto3" json:"enabled,omitempty"`
}

func (m *LiquidStake) Reset()         { *m = LiquidStake{} }
func (m *LiquidStake) String() string { return proto.CompactTextString(m) }
func (*LiquidStake) ProtoMessage()    {}
func (*LiquidStake) Descriptor() ([]byte, []int) {
	return fileDescriptor_429540018f2469ab, []int{2}
}
func (m *LiquidStake) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LiquidStake) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LiquidStake.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LiquidStake) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LiquidStake.Merge(m, src)
}
func (m *LiquidStake) XXX_Size() int {
	return m.Size()
}
func (m *LiquidStake) XXX_DiscardUnknown() {
	xxx_messageInfo_LiquidStake.DiscardUnknown(m)
}

var xxx_messageInfo_LiquidStake proto.InternalMessageInfo

func (m *LiquidStake) GetFeatureType() FeatureType {
	if m != nil {
		return m.FeatureType
	}
	return FeatureType_LIQUID_STAKE_IBC
}

func (m *LiquidStake) GetCodeID() uint64 {
	if m != nil {
		return m.CodeID
	}
	return 0
}

func (m *LiquidStake) GetInstantiation() InstantiationState {
	if m != nil {
		return m.Instantiation
	}
	return InstantiationState_INSTANTIATION_NOT_INITIATED
}

func (m *LiquidStake) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

func (m *LiquidStake) GetDenoms() []string {
	if m != nil {
		return m.Denoms
	}
	return nil
}

func (m *LiquidStake) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

// aim to keep this smaller than 256 MaxCharLen in ICA memo.
type ICAMemo struct {
	FeatureType FeatureType `protobuf:"varint,1,opt,name=feature_type,json=featureType,proto3,enum=pstake.ratesync.v1beta1.FeatureType" json:"feature_type,omitempty"`
	HostChainID uint64      `protobuf:"varint,2,opt,name=host_chain_i_d,json=hostChainID,proto3" json:"host_chain_i_d,omitempty"`
}

func (m *ICAMemo) Reset()         { *m = ICAMemo{} }
func (m *ICAMemo) String() string { return proto.CompactTextString(m) }
func (*ICAMemo) ProtoMessage()    {}
func (*ICAMemo) Descriptor() ([]byte, []int) {
	return fileDescriptor_429540018f2469ab, []int{3}
}
func (m *ICAMemo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ICAMemo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ICAMemo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ICAMemo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ICAMemo.Merge(m, src)
}
func (m *ICAMemo) XXX_Size() int {
	return m.Size()
}
func (m *ICAMemo) XXX_DiscardUnknown() {
	xxx_messageInfo_ICAMemo.DiscardUnknown(m)
}

var xxx_messageInfo_ICAMemo proto.InternalMessageInfo

func (m *ICAMemo) GetFeatureType() FeatureType {
	if m != nil {
		return m.FeatureType
	}
	return FeatureType_LIQUID_STAKE_IBC
}

func (m *ICAMemo) GetHostChainID() uint64 {
	if m != nil {
		return m.HostChainID
	}
	return 0
}

func init() {
	proto.RegisterEnum("pstake.ratesync.v1beta1.InstantiationState", InstantiationState_name, InstantiationState_value)
	proto.RegisterEnum("pstake.ratesync.v1beta1.FeatureType", FeatureType_name, FeatureType_value)
	proto.RegisterType((*HostChain)(nil), "pstake.ratesync.v1beta1.HostChain")
	proto.RegisterType((*Feature)(nil), "pstake.ratesync.v1beta1.Feature")
	proto.RegisterType((*LiquidStake)(nil), "pstake.ratesync.v1beta1.LiquidStake")
	proto.RegisterType((*ICAMemo)(nil), "pstake.ratesync.v1beta1.ICAMemo")
}

func init() {
	proto.RegisterFile("pstake/ratesync/v1beta1/ratesync.proto", fileDescriptor_429540018f2469ab)
}

var fileDescriptor_429540018f2469ab = []byte{
	// 692 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x41, 0x6f, 0x12, 0x41,
	0x14, 0x66, 0x29, 0x42, 0x19, 0x28, 0xa5, 0x93, 0xc6, 0x92, 0x36, 0xa1, 0x04, 0x1b, 0x43, 0x6b,
	0x64, 0x53, 0x8c, 0x27, 0x4f, 0xb0, 0x5b, 0xeb, 0xc4, 0x16, 0xda, 0x85, 0xc6, 0x44, 0x0f, 0x93,
	0x61, 0x98, 0x96, 0x89, 0x30, 0x83, 0xbb, 0x43, 0xb5, 0xff, 0xc2, 0x9f, 0xe0, 0xd1, 0xa3, 0x67,
	0x7f, 0x41, 0x8f, 0x3d, 0x7a, 0x30, 0xc6, 0xb4, 0x07, 0xff, 0x86, 0xd9, 0xd9, 0x05, 0x16, 0x9b,
	0x7a, 0x30, 0x5e, 0xc8, 0xbe, 0xef, 0xfb, 0xde, 0x37, 0xef, 0x3d, 0xde, 0x0c, 0x78, 0x38, 0xf2,
	0x14, 0x79, 0xcb, 0x4c, 0x97, 0x28, 0xe6, 0x5d, 0x08, 0x6a, 0x9e, 0xef, 0x76, 0x99, 0x22, 0xbb,
	0x53, 0xa0, 0x3a, 0x72, 0xa5, 0x92, 0x70, 0x2d, 0xd0, 0x55, 0xa7, 0x70, 0xa8, 0x5b, 0x5f, 0x3d,
	0x93, 0x67, 0x52, 0x6b, 0x4c, 0xff, 0x2b, 0x90, 0xaf, 0xd7, 0x42, 0xdb, 0x01, 0x7f, 0x37, 0xe6,
	0x3d, 0xfd, 0xcd, 0xbb, 0x33, 0xf3, 0x79, 0x38, 0xcc, 0x59, 0x21, 0x43, 0x2e, 0xa4, 0xa9, 0x7f,
	0x03, 0xa8, 0xfc, 0x3d, 0x0e, 0xd2, 0x2f, 0xa4, 0xa7, 0xac, 0x3e, 0xe1, 0x02, 0x2e, 0x83, 0x05,
	0x8e, 0x7b, 0x05, 0xa3, 0x64, 0x54, 0x12, 0x4e, 0x9c, 0xdb, 0x70, 0x1d, 0xa4, 0xa9, 0xcf, 0x60,
	0x1f, 0x8e, 0x97, 0x8c, 0x4a, 0xda, 0x49, 0x69, 0x00, 0xd9, 0x70, 0x0b, 0xe4, 0xa8, 0x14, 0x82,
	0x51, 0xc5, 0x65, 0x20, 0x58, 0xd0, 0x82, 0xec, 0x0c, 0x45, 0x36, 0x7c, 0x05, 0x96, 0x38, 0xa6,
	0x98, 0x60, 0x42, 0xa9, 0x1c, 0x0b, 0x55, 0x48, 0x94, 0x8c, 0x4a, 0xa6, 0xb6, 0x5d, 0x0d, 0xdb,
	0xfd, 0xa3, 0xd0, 0xb0, 0xfe, 0x2a, 0xb2, 0xea, 0xf5, 0x20, 0xa1, 0x91, 0xbe, 0xfc, 0xb1, 0x19,
	0xfb, 0xfc, 0xeb, 0xcb, 0x8e, 0xe1, 0x00, 0x3e, 0x85, 0xe1, 0x3e, 0x58, 0x3c, 0x65, 0x44, 0x8d,
	0x5d, 0xe6, 0x15, 0xee, 0x69, 0xcf, 0x52, 0xf5, 0x8e, 0x11, 0x56, 0x9f, 0x07, 0xc2, 0xa8, 0xd5,
	0x34, 0x19, 0x9a, 0x60, 0x55, 0xb9, 0x44, 0x78, 0xa7, 0xcc, 0xc5, 0xb4, 0x4f, 0x84, 0x60, 0x03,
	0xdd, 0x4d, 0x52, 0x77, 0xb3, 0x32, 0xe1, 0xac, 0x80, 0x42, 0x36, 0xdc, 0x06, 0x53, 0x10, 0x8f,
	0xa4, 0xab, 0xb4, 0x3a, 0xa5, 0xd5, 0xb9, 0x09, 0x71, 0x24, 0x5d, 0x85, 0xec, 0xf2, 0x57, 0x03,
	0xa4, 0xc2, 0xc3, 0xe1, 0x1b, 0x00, 0x83, 0x66, 0xb1, 0xae, 0x12, 0x73, 0xdc, 0xc5, 0x54, 0xcf,
	0x3a, 0x53, 0xdb, 0xba, 0xb3, 0xf4, 0x03, 0x9d, 0xd2, 0xf6, 0xc9, 0x68, 0xf9, 0xb9, 0xc1, 0x0c,
	0x47, 0x0d, 0x0b, 0x3a, 0x20, 0x1b, 0x35, 0xd7, 0xff, 0xd5, 0x3f, 0xd8, 0x66, 0x22, 0xb6, 0xe5,
	0x4f, 0x71, 0x90, 0x89, 0xe8, 0xe0, 0x3e, 0xc8, 0x86, 0x43, 0xc3, 0xea, 0x62, 0xc4, 0x74, 0xe9,
	0xb9, 0xbf, 0x9c, 0x11, 0x36, 0xde, 0xb9, 0x18, 0x31, 0x27, 0x73, 0x3a, 0x0b, 0x60, 0x01, 0x2c,
	0x52, 0xd9, 0x63, 0xd3, 0xa5, 0x4a, 0x38, 0x49, 0x3f, 0x46, 0x36, 0x3c, 0x06, 0x4b, 0x5c, 0x78,
	0x8a, 0x08, 0xc5, 0x89, 0xbf, 0x40, 0x7a, 0xa5, 0x72, 0xb5, 0x47, 0x77, 0x9e, 0x81, 0xa2, 0xea,
	0xb6, 0x22, 0x8a, 0x39, 0xf3, 0x0e, 0x70, 0x1b, 0xe4, 0xa9, 0x14, 0xca, 0x25, 0x54, 0x61, 0xd2,
	0xeb, 0xb9, 0xcc, 0xf3, 0xf4, 0x0e, 0xa6, 0x9d, 0xe5, 0x09, 0x5e, 0x0f, 0x60, 0x78, 0x1f, 0x24,
	0x7b, 0x4c, 0xc8, 0xa1, 0xbf, 0x50, 0x0b, 0x95, 0xb4, 0x13, 0x46, 0xb0, 0x00, 0x52, 0x4c, 0x90,
	0xee, 0x80, 0x05, 0x4b, 0xb1, 0xe8, 0x4c, 0xc2, 0xf2, 0x7b, 0x90, 0x42, 0x56, 0xfd, 0x90, 0x0d,
	0xe5, 0xff, 0x9b, 0xce, 0x03, 0x90, 0xeb, 0x4b, 0x4f, 0xe1, 0xf9, 0x8b, 0x97, 0x70, 0x32, 0xfd,
	0xc9, 0x3d, 0x45, 0xf6, 0x8e, 0x04, 0xf0, 0x76, 0xeb, 0x70, 0x13, 0x6c, 0xa0, 0x66, 0xbb, 0x53,
	0x6f, 0x76, 0x50, 0xbd, 0x83, 0x5a, 0x4d, 0xdc, 0x6c, 0x75, 0x30, 0x6a, 0x22, 0x3f, 0xdc, 0xb3,
	0xf3, 0x31, 0xb8, 0x01, 0xd6, 0xe6, 0x05, 0x33, 0xd2, 0xb8, 0x4d, 0x5a, 0xad, 0xc3, 0xa3, 0x83,
	0x3d, 0x9f, 0x8c, 0xef, 0x3c, 0x05, 0x99, 0x48, 0xc5, 0x70, 0x15, 0xe4, 0x0f, 0xd0, 0xf1, 0x09,
	0xb2, 0x71, 0xbb, 0x53, 0x7f, 0xb9, 0x87, 0x51, 0xc3, 0xca, 0xc7, 0x60, 0x1e, 0x64, 0xa3, 0x68,
	0xde, 0x68, 0x9c, 0x5c, 0x5e, 0x17, 0x8d, 0xab, 0xeb, 0xa2, 0xf1, 0xf3, 0xba, 0x68, 0x7c, 0xbc,
	0x29, 0xc6, 0xae, 0x6e, 0x8a, 0xb1, 0x6f, 0x37, 0xc5, 0xd8, 0xeb, 0x67, 0x67, 0x5c, 0xf5, 0xc7,
	0xdd, 0x2a, 0x95, 0x43, 0x73, 0xc4, 0x5c, 0x8f, 0x7b, 0x8a, 0x09, 0xca, 0x5a, 0x82, 0x99, 0xc1,
	0xc8, 0x1e, 0x0b, 0xa2, 0xf8, 0x39, 0x33, 0xcf, 0x6b, 0xe6, 0x87, 0xd9, 0xeb, 0xe9, 0xcf, 0xd6,
	0xeb, 0x26, 0xf5, 0xeb, 0xf5, 0xe4, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x05, 0x21, 0xe8, 0x15,
	0x5d, 0x05, 0x00, 0x00,
}

func (m *HostChain) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HostChain) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HostChain) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TransferPortID) > 0 {
		i -= len(m.TransferPortID)
		copy(dAtA[i:], m.TransferPortID)
		i = encodeVarintRatesync(dAtA, i, uint64(len(m.TransferPortID)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.TransferChannelID) > 0 {
		i -= len(m.TransferChannelID)
		copy(dAtA[i:], m.TransferChannelID)
		i = encodeVarintRatesync(dAtA, i, uint64(len(m.TransferChannelID)))
		i--
		dAtA[i] = 0x32
	}
	{
		size, err := m.Features.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintRatesync(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size, err := m.ICAAccount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintRatesync(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.ConnectionID) > 0 {
		i -= len(m.ConnectionID)
		copy(dAtA[i:], m.ConnectionID)
		i = encodeVarintRatesync(dAtA, i, uint64(len(m.ConnectionID)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ChainID) > 0 {
		i -= len(m.ChainID)
		copy(dAtA[i:], m.ChainID)
		i = encodeVarintRatesync(dAtA, i, uint64(len(m.ChainID)))
		i--
		dAtA[i] = 0x12
	}
	if m.ID != 0 {
		i = encodeVarintRatesync(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Feature) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Feature) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Feature) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.LiquidStake.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintRatesync(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.LiquidStakeIBC.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintRatesync(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *LiquidStake) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LiquidStake) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LiquidStake) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Enabled {
		i--
		if m.Enabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if len(m.Denoms) > 0 {
		for iNdEx := len(m.Denoms) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Denoms[iNdEx])
			copy(dAtA[i:], m.Denoms[iNdEx])
			i = encodeVarintRatesync(dAtA, i, uint64(len(m.Denoms[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.ContractAddress) > 0 {
		i -= len(m.ContractAddress)
		copy(dAtA[i:], m.ContractAddress)
		i = encodeVarintRatesync(dAtA, i, uint64(len(m.ContractAddress)))
		i--
		dAtA[i] = 0x22
	}
	if m.Instantiation != 0 {
		i = encodeVarintRatesync(dAtA, i, uint64(m.Instantiation))
		i--
		dAtA[i] = 0x18
	}
	if m.CodeID != 0 {
		i = encodeVarintRatesync(dAtA, i, uint64(m.CodeID))
		i--
		dAtA[i] = 0x10
	}
	if m.FeatureType != 0 {
		i = encodeVarintRatesync(dAtA, i, uint64(m.FeatureType))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ICAMemo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ICAMemo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ICAMemo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.HostChainID != 0 {
		i = encodeVarintRatesync(dAtA, i, uint64(m.HostChainID))
		i--
		dAtA[i] = 0x10
	}
	if m.FeatureType != 0 {
		i = encodeVarintRatesync(dAtA, i, uint64(m.FeatureType))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintRatesync(dAtA []byte, offset int, v uint64) int {
	offset -= sovRatesync(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *HostChain) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovRatesync(uint64(m.ID))
	}
	l = len(m.ChainID)
	if l > 0 {
		n += 1 + l + sovRatesync(uint64(l))
	}
	l = len(m.ConnectionID)
	if l > 0 {
		n += 1 + l + sovRatesync(uint64(l))
	}
	l = m.ICAAccount.Size()
	n += 1 + l + sovRatesync(uint64(l))
	l = m.Features.Size()
	n += 1 + l + sovRatesync(uint64(l))
	l = len(m.TransferChannelID)
	if l > 0 {
		n += 1 + l + sovRatesync(uint64(l))
	}
	l = len(m.TransferPortID)
	if l > 0 {
		n += 1 + l + sovRatesync(uint64(l))
	}
	return n
}

func (m *Feature) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.LiquidStakeIBC.Size()
	n += 1 + l + sovRatesync(uint64(l))
	l = m.LiquidStake.Size()
	n += 1 + l + sovRatesync(uint64(l))
	return n
}

func (m *LiquidStake) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.FeatureType != 0 {
		n += 1 + sovRatesync(uint64(m.FeatureType))
	}
	if m.CodeID != 0 {
		n += 1 + sovRatesync(uint64(m.CodeID))
	}
	if m.Instantiation != 0 {
		n += 1 + sovRatesync(uint64(m.Instantiation))
	}
	l = len(m.ContractAddress)
	if l > 0 {
		n += 1 + l + sovRatesync(uint64(l))
	}
	if len(m.Denoms) > 0 {
		for _, s := range m.Denoms {
			l = len(s)
			n += 1 + l + sovRatesync(uint64(l))
		}
	}
	if m.Enabled {
		n += 2
	}
	return n
}

func (m *ICAMemo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.FeatureType != 0 {
		n += 1 + sovRatesync(uint64(m.FeatureType))
	}
	if m.HostChainID != 0 {
		n += 1 + sovRatesync(uint64(m.HostChainID))
	}
	return n
}

func sovRatesync(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRatesync(x uint64) (n int) {
	return sovRatesync(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *HostChain) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRatesync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HostChain: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HostChain: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnectionID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ICAAccount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ICAAccount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Features", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Features.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferChannelID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TransferChannelID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferPortID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TransferPortID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRatesync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRatesync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Feature) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRatesync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Feature: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Feature: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidStakeIBC", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidStakeIBC.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidStake", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidStake.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRatesync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRatesync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *LiquidStake) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRatesync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LiquidStake: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LiquidStake: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeatureType", wireType)
			}
			m.FeatureType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeatureType |= FeatureType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CodeID", wireType)
			}
			m.CodeID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CodeID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Instantiation", wireType)
			}
			m.Instantiation = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Instantiation |= InstantiationState(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denoms", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRatesync
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRatesync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denoms = append(m.Denoms, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Enabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Enabled = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipRatesync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRatesync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ICAMemo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRatesync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ICAMemo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ICAMemo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeatureType", wireType)
			}
			m.FeatureType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeatureType |= FeatureType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HostChainID", wireType)
			}
			m.HostChainID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HostChainID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRatesync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRatesync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRatesync(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRatesync
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRatesync
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthRatesync
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRatesync
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRatesync
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRatesync        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRatesync          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRatesync = fmt.Errorf("proto: unexpected end of group")
)
