// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmos/v1beta1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	tx "github.com/cosmos/cosmos-sdk/types/tx"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type StakeAmount struct {
	Denom  string                                   `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,3,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount" yaml:"amount"`
}

func (m *StakeAmount) Reset()         { *m = StakeAmount{} }
func (m *StakeAmount) String() string { return proto.CompactTextString(m) }
func (*StakeAmount) ProtoMessage()    {}
func (*StakeAmount) Descriptor() ([]byte, []int) {
	return fileDescriptor_9abdaadb2561c892, []int{0}
}
func (m *StakeAmount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StakeAmount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StakeAmount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StakeAmount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StakeAmount.Merge(m, src)
}
func (m *StakeAmount) XXX_Size() int {
	return m.Size()
}
func (m *StakeAmount) XXX_DiscardUnknown() {
	xxx_messageInfo_StakeAmount.DiscardUnknown(m)
}

var xxx_messageInfo_StakeAmount proto.InternalMessageInfo

func (m *StakeAmount) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *StakeAmount) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

type DelegationCosmos struct {
	Address string      `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Stake   StakeAmount `protobuf:"bytes,2,opt,name=stake,proto3" json:"stake" yaml:"stake"`
}

func (m *DelegationCosmos) Reset()         { *m = DelegationCosmos{} }
func (m *DelegationCosmos) String() string { return proto.CompactTextString(m) }
func (*DelegationCosmos) ProtoMessage()    {}
func (*DelegationCosmos) Descriptor() ([]byte, []int) {
	return fileDescriptor_9abdaadb2561c892, []int{1}
}
func (m *DelegationCosmos) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DelegationCosmos) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DelegationCosmos.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DelegationCosmos) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DelegationCosmos.Merge(m, src)
}
func (m *DelegationCosmos) XXX_Size() int {
	return m.Size()
}
func (m *DelegationCosmos) XXX_DiscardUnknown() {
	xxx_messageInfo_DelegationCosmos.DiscardUnknown(m)
}

var xxx_messageInfo_DelegationCosmos proto.InternalMessageInfo

func (m *DelegationCosmos) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *DelegationCosmos) GetStake() StakeAmount {
	if m != nil {
		return m.Stake
	}
	return StakeAmount{}
}

type IncomingTx struct {
	TxResponse types.TxResponse `protobuf:"bytes,1,opt,name=tx_response,json=txResponse,proto3" json:"tx_response" yaml:"tx_response"`
}

func (m *IncomingTx) Reset()         { *m = IncomingTx{} }
func (m *IncomingTx) String() string { return proto.CompactTextString(m) }
func (*IncomingTx) ProtoMessage()    {}
func (*IncomingTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_9abdaadb2561c892, []int{2}
}
func (m *IncomingTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IncomingTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IncomingTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IncomingTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IncomingTx.Merge(m, src)
}
func (m *IncomingTx) XXX_Size() int {
	return m.Size()
}
func (m *IncomingTx) XXX_DiscardUnknown() {
	xxx_messageInfo_IncomingTx.DiscardUnknown(m)
}

var xxx_messageInfo_IncomingTx proto.InternalMessageInfo

func (m *IncomingTx) GetTxResponse() types.TxResponse {
	if m != nil {
		return m.TxResponse
	}
	return types.TxResponse{}
}

type OutgoingTx struct {
	TxBody tx.TxBody `protobuf:"bytes,1,opt,name=tx_body,json=txBody,proto3" json:"tx_body" yaml:"tx_body"`
}

func (m *OutgoingTx) Reset()         { *m = OutgoingTx{} }
func (m *OutgoingTx) String() string { return proto.CompactTextString(m) }
func (*OutgoingTx) ProtoMessage()    {}
func (*OutgoingTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_9abdaadb2561c892, []int{3}
}
func (m *OutgoingTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OutgoingTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OutgoingTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OutgoingTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutgoingTx.Merge(m, src)
}
func (m *OutgoingTx) XXX_Size() int {
	return m.Size()
}
func (m *OutgoingTx) XXX_DiscardUnknown() {
	xxx_messageInfo_OutgoingTx.DiscardUnknown(m)
}

var xxx_messageInfo_OutgoingTx proto.InternalMessageInfo

func (m *OutgoingTx) GetTxBody() tx.TxBody {
	if m != nil {
		return m.TxBody
	}
	return tx.TxBody{}
}

type GenesisState struct {
	Params            Params             `protobuf:"bytes,1,opt,name=params,proto3" json:"params" yaml:"params"`
	CosmosDelegations []DelegationCosmos `protobuf:"bytes,2,rep,name=cosmos_delegations,json=cosmosDelegations,proto3" json:"cosmos_delegations" yaml:"cosmos_delegations"`
	IncomingTxn       []IncomingTx       `protobuf:"bytes,3,rep,name=incoming_txn,json=incomingTxn,proto3" json:"incoming_txn" yaml:"incoming_txn"`
	OutgoingTxn       OutgoingTx         `protobuf:"bytes,4,opt,name=outgoing_txn,json=outgoingTxn,proto3" json:"outgoing_txn" yaml:"outgoing_txn"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_9abdaadb2561c892, []int{4}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetCosmosDelegations() []DelegationCosmos {
	if m != nil {
		return m.CosmosDelegations
	}
	return nil
}

func (m *GenesisState) GetIncomingTxn() []IncomingTx {
	if m != nil {
		return m.IncomingTxn
	}
	return nil
}

func (m *GenesisState) GetOutgoingTxn() OutgoingTx {
	if m != nil {
		return m.OutgoingTxn
	}
	return OutgoingTx{}
}

func init() {
	proto.RegisterType((*StakeAmount)(nil), "comsos.v1beta1.StakeAmount")
	proto.RegisterType((*DelegationCosmos)(nil), "comsos.v1beta1.DelegationCosmos")
	proto.RegisterType((*IncomingTx)(nil), "comsos.v1beta1.IncomingTx")
	proto.RegisterType((*OutgoingTx)(nil), "comsos.v1beta1.OutgoingTx")
	proto.RegisterType((*GenesisState)(nil), "comsos.v1beta1.GenesisState")
}

func init() { proto.RegisterFile("cosmos/v1beta1/genesis.proto", fileDescriptor_9abdaadb2561c892) }

var fileDescriptor_9abdaadb2561c892 = []byte{
	// 607 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x53, 0xcf, 0x4e, 0xdb, 0x4c,
	0x10, 0x8f, 0xe1, 0x23, 0xe8, 0x5b, 0x53, 0xd4, 0xba, 0x14, 0x85, 0x80, 0x4c, 0xba, 0xed, 0x21,
	0x17, 0x6c, 0x41, 0xa5, 0x1e, 0x7a, 0xc3, 0xb4, 0x42, 0xed, 0x05, 0x64, 0x38, 0x54, 0x5c, 0xa2,
	0x75, 0xbc, 0x75, 0x57, 0xe0, 0x1d, 0x2b, 0xbb, 0x41, 0xce, 0xb1, 0x6f, 0xd0, 0x6b, 0x5f, 0xa1,
	0x4f, 0xc2, 0x91, 0x63, 0x4f, 0xb4, 0x82, 0x37, 0xe8, 0x13, 0x54, 0xde, 0x9d, 0x18, 0x63, 0x4e,
	0xc9, 0xf8, 0xf7, 0x6f, 0x76, 0x66, 0x97, 0x6c, 0x8d, 0x41, 0xe5, 0xa0, 0xc2, 0xcb, 0xdd, 0x84,
	0x6b, 0xb6, 0x1b, 0x66, 0x5c, 0x72, 0x25, 0x54, 0x50, 0x4c, 0x40, 0x83, 0xb7, 0x3a, 0x86, 0x5c,
	0x81, 0x0a, 0x10, 0xed, 0xaf, 0x65, 0x90, 0x81, 0x81, 0xc2, 0xea, 0x9f, 0x65, 0xf5, 0xb7, 0x32,
	0x80, 0xec, 0x82, 0x87, 0xac, 0x10, 0x21, 0x93, 0x12, 0x34, 0xd3, 0x02, 0x24, 0x7a, 0xf4, 0x37,
	0x10, 0x35, 0x55, 0x32, 0xfd, 0x12, 0x32, 0x39, 0x43, 0x68, 0xb3, 0x15, 0x6e, 0x4b, 0x04, 0x7d,
	0x04, 0x13, 0xa6, 0x78, 0x83, 0x21, 0x24, 0xe2, 0xaf, 0x9a, 0x38, 0x4b, 0xc6, 0xa2, 0x26, 0x55,
	0x05, 0x92, 0xfa, 0x48, 0xd2, 0x65, 0x8d, 0xea, 0xd2, 0x62, 0xf4, 0x87, 0x43, 0xdc, 0x13, 0xcd,
	0xce, 0xf9, 0x7e, 0x0e, 0x53, 0xa9, 0xbd, 0x35, 0xb2, 0x94, 0x72, 0x09, 0x79, 0xcf, 0x19, 0x38,
	0xc3, 0xff, 0x63, 0x5b, 0x78, 0x9a, 0x74, 0x99, 0xc1, 0x7b, 0x8b, 0x83, 0xc5, 0xa1, 0xbb, 0xb7,
	0x11, 0x60, 0x97, 0x55, 0xee, 0x7c, 0x30, 0xc1, 0x01, 0x08, 0x19, 0xed, 0x5f, 0xdd, 0x6c, 0x77,
	0xfe, 0xde, 0x6c, 0x3f, 0x99, 0xb1, 0xfc, 0xe2, 0x1d, 0xb5, 0x32, 0xfa, 0xf3, 0xf7, 0xf6, 0x30,
	0x13, 0xfa, 0xeb, 0x34, 0x09, 0xc6, 0x90, 0xe3, 0x19, 0xf1, 0x67, 0x47, 0xa5, 0xe7, 0xa1, 0x9e,
	0x15, 0x5c, 0x19, 0x07, 0x15, 0x63, 0x16, 0x9d, 0x92, 0xa7, 0xef, 0xf9, 0x05, 0xcf, 0xcc, 0x24,
	0x0f, 0x0c, 0xd7, 0xeb, 0x91, 0x65, 0x96, 0xa6, 0x13, 0xae, 0x14, 0x76, 0x38, 0x2f, 0xbd, 0x43,
	0xb2, 0xa4, 0xaa, 0x83, 0xf4, 0x16, 0x06, 0xce, 0xd0, 0xdd, 0xdb, 0x0c, 0x1e, 0xae, 0x2d, 0x68,
	0x9c, 0x32, 0x5a, 0xc3, 0x26, 0x57, 0x6c, 0x93, 0x46, 0x47, 0x63, 0xab, 0xa7, 0x40, 0xc8, 0x47,
	0x39, 0x86, 0x5c, 0xc8, 0xec, 0xb4, 0xf4, 0x18, 0x71, 0x75, 0x39, 0x9a, 0x70, 0x55, 0x80, 0x54,
	0xdc, 0x84, 0xba, 0x7b, 0xaf, 0x1f, 0x9c, 0xdf, 0x8c, 0x7a, 0x1e, 0x73, 0x5a, 0xc6, 0xc8, 0x8d,
	0xfa, 0x98, 0xe2, 0xd9, 0x94, 0x86, 0x0d, 0x8d, 0x89, 0xae, 0x79, 0xf4, 0x33, 0x21, 0x47, 0x53,
	0x9d, 0x81, 0x0d, 0xfc, 0x44, 0x96, 0x75, 0x39, 0x4a, 0x20, 0x9d, 0x61, 0x58, 0x3d, 0x6c, 0x5d,
	0x36, 0x52, 0x22, 0x48, 0x67, 0xd1, 0x3a, 0x26, 0xac, 0xd6, 0x09, 0x95, 0x8e, 0xc6, 0x5d, 0x6d,
	0x70, 0xfa, 0x6d, 0x91, 0xac, 0x1c, 0xda, 0xcb, 0x7c, 0xa2, 0x99, 0xe6, 0xde, 0x07, 0xd2, 0x2d,
	0xd8, 0x84, 0xe5, 0x0a, 0xbd, 0xd7, 0xdb, 0x53, 0x3a, 0x36, 0x68, 0xf4, 0xe2, 0xe1, 0x16, 0xad,
	0x86, 0xc6, 0x28, 0xf6, 0x26, 0xc4, 0xb3, 0x3d, 0x8d, 0xd2, 0x7a, 0x41, 0xaa, 0xb7, 0x60, 0xee,
	0xc6, 0xa0, 0x6d, 0xd9, 0xde, 0x61, 0xf4, 0x12, 0xcd, 0x37, 0xac, 0xf9, 0x63, 0x27, 0x1a, 0x3f,
	0xb3, 0x1f, 0xef, 0xa5, 0xca, 0x3b, 0x23, 0x2b, 0x02, 0xd7, 0x32, 0xd2, 0xa5, 0xc4, 0x9b, 0xd8,
	0x6f, 0xa7, 0xdd, 0xaf, 0x2e, 0xda, 0xc4, 0x9c, 0xe7, 0x36, 0xa7, 0xa9, 0xa6, 0xb1, 0x2b, 0x6a,
	0xa2, 0xac, 0xbc, 0x01, 0x37, 0x60, 0xbc, 0xff, 0x33, 0xc3, 0x79, 0xe4, 0x7d, 0xbf, 0xa5, 0xb6,
	0x77, 0x53, 0x4d, 0x63, 0x17, 0x6a, 0xa2, 0x8c, 0x8e, 0xaf, 0x6e, 0x7d, 0xe7, 0xfa, 0xd6, 0x77,
	0xfe, 0xdc, 0xfa, 0xce, 0xf7, 0x3b, 0xbf, 0x73, 0x7d, 0xe7, 0x77, 0x7e, 0xdd, 0xf9, 0x9d, 0xb3,
	0xb7, 0x8d, 0x17, 0x51, 0xf0, 0x89, 0x12, 0x4a, 0x73, 0x39, 0xe6, 0x47, 0x92, 0x87, 0x85, 0xb9,
	0xad, 0x3b, 0x92, 0x69, 0x71, 0xc9, 0xc3, 0x72, 0xfe, 0x60, 0xcc, 0x2b, 0x49, 0xba, 0xe6, 0xe9,
	0xbe, 0xf9, 0x17, 0x00, 0x00, 0xff, 0xff, 0xc9, 0xca, 0x9c, 0x8c, 0xb7, 0x04, 0x00, 0x00,
}

func (m *StakeAmount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StakeAmount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StakeAmount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DelegationCosmos) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DelegationCosmos) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DelegationCosmos) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Stake.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *IncomingTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IncomingTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IncomingTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.TxResponse.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *OutgoingTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OutgoingTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OutgoingTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.TxBody.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.OutgoingTxn.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.IncomingTxn) > 0 {
		for iNdEx := len(m.IncomingTxn) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.IncomingTxn[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.CosmosDelegations) > 0 {
		for iNdEx := len(m.CosmosDelegations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.CosmosDelegations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *StakeAmount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *DelegationCosmos) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = m.Stake.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *IncomingTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.TxResponse.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *OutgoingTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.TxBody.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.CosmosDelegations) > 0 {
		for _, e := range m.CosmosDelegations {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.IncomingTxn) > 0 {
		for _, e := range m.IncomingTxn {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.OutgoingTxn.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StakeAmount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: StakeAmount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StakeAmount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *DelegationCosmos) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: DelegationCosmos: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DelegationCosmos: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stake", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Stake.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *IncomingTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: IncomingTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IncomingTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxResponse", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TxResponse.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *OutgoingTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: OutgoingTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OutgoingTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBody", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TxBody.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CosmosDelegations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CosmosDelegations = append(m.CosmosDelegations, DelegationCosmos{})
			if err := m.CosmosDelegations[len(m.CosmosDelegations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncomingTxn", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IncomingTxn = append(m.IncomingTxn, IncomingTx{})
			if err := m.IncomingTxn[len(m.IncomingTxn)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutgoingTxn", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OutgoingTxn.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
