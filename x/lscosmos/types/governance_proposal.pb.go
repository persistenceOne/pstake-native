// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pstake/lscosmos/v1beta1/governance_proposal.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type MinDepositAndFeeChangeProposal struct {
	Title               string                                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description         string                                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	MinDeposit          github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=min_deposit,json=minDeposit,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_deposit"`
	PstakeDepositFee    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=pstake_deposit_fee,json=pstakeDepositFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_deposit_fee"`
	PstakeRestakeFee    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=pstake_restake_fee,json=pstakeRestakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_restake_fee"`
	PstakeUnstakeFee    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=pstake_unstake_fee,json=pstakeUnstakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_unstake_fee"`
	PstakeRedemptionFee github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=pstake_redemption_fee,json=pstakeRedemptionFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_redemption_fee"`
}

func (m *MinDepositAndFeeChangeProposal) Reset()      { *m = MinDepositAndFeeChangeProposal{} }
func (*MinDepositAndFeeChangeProposal) ProtoMessage() {}
func (*MinDepositAndFeeChangeProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_47404a6acaa6ce8f, []int{0}
}
func (m *MinDepositAndFeeChangeProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MinDepositAndFeeChangeProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MinDepositAndFeeChangeProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MinDepositAndFeeChangeProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MinDepositAndFeeChangeProposal.Merge(m, src)
}
func (m *MinDepositAndFeeChangeProposal) XXX_Size() int {
	return m.Size()
}
func (m *MinDepositAndFeeChangeProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_MinDepositAndFeeChangeProposal.DiscardUnknown(m)
}

var xxx_messageInfo_MinDepositAndFeeChangeProposal proto.InternalMessageInfo

type PstakeFeeAddressChangeProposal struct {
	Title            string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description      string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	PstakeFeeAddress string `protobuf:"bytes,3,opt,name=pstake_fee_address,json=pstakeFeeAddress,proto3" json:"pstake_fee_address,omitempty"`
}

func (m *PstakeFeeAddressChangeProposal) Reset()      { *m = PstakeFeeAddressChangeProposal{} }
func (*PstakeFeeAddressChangeProposal) ProtoMessage() {}
func (*PstakeFeeAddressChangeProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_47404a6acaa6ce8f, []int{1}
}
func (m *PstakeFeeAddressChangeProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PstakeFeeAddressChangeProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PstakeFeeAddressChangeProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PstakeFeeAddressChangeProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PstakeFeeAddressChangeProposal.Merge(m, src)
}
func (m *PstakeFeeAddressChangeProposal) XXX_Size() int {
	return m.Size()
}
func (m *PstakeFeeAddressChangeProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_PstakeFeeAddressChangeProposal.DiscardUnknown(m)
}

var xxx_messageInfo_PstakeFeeAddressChangeProposal proto.InternalMessageInfo

type AllowListedValidatorSetChangeProposal struct {
	Title                 string                `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description           string                `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	AllowListedValidators AllowListedValidators `protobuf:"bytes,3,opt,name=allow_listed_validators,json=allowListedValidators,proto3" json:"allow_listed_validators"`
}

func (m *AllowListedValidatorSetChangeProposal) Reset()      { *m = AllowListedValidatorSetChangeProposal{} }
func (*AllowListedValidatorSetChangeProposal) ProtoMessage() {}
func (*AllowListedValidatorSetChangeProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_47404a6acaa6ce8f, []int{2}
}
func (m *AllowListedValidatorSetChangeProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AllowListedValidatorSetChangeProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AllowListedValidatorSetChangeProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AllowListedValidatorSetChangeProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllowListedValidatorSetChangeProposal.Merge(m, src)
}
func (m *AllowListedValidatorSetChangeProposal) XXX_Size() int {
	return m.Size()
}
func (m *AllowListedValidatorSetChangeProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_AllowListedValidatorSetChangeProposal.DiscardUnknown(m)
}

var xxx_messageInfo_AllowListedValidatorSetChangeProposal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MinDepositAndFeeChangeProposal)(nil), "pstake.lscosmos.v1beta1.MinDepositAndFeeChangeProposal")
	proto.RegisterType((*PstakeFeeAddressChangeProposal)(nil), "pstake.lscosmos.v1beta1.PstakeFeeAddressChangeProposal")
	proto.RegisterType((*AllowListedValidatorSetChangeProposal)(nil), "pstake.lscosmos.v1beta1.AllowListedValidatorSetChangeProposal")
}

func init() {
	proto.RegisterFile("pstake/lscosmos/v1beta1/governance_proposal.proto", fileDescriptor_47404a6acaa6ce8f)
}

var fileDescriptor_47404a6acaa6ce8f = []byte{
	// 544 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0x31, 0x6f, 0xd3, 0x40,
	0x14, 0xc7, 0x6d, 0x48, 0x0b, 0x5c, 0x16, 0x64, 0x52, 0x35, 0x74, 0xb8, 0x54, 0x95, 0xa8, 0x58,
	0x62, 0x2b, 0x65, 0x03, 0x96, 0x84, 0x28, 0x52, 0x25, 0x10, 0x55, 0x2a, 0x18, 0x90, 0x90, 0x75,
	0xb1, 0x5f, 0xdd, 0xa3, 0xf6, 0x9d, 0xe5, 0xbb, 0x1a, 0xd8, 0xf8, 0x08, 0x8c, 0x8c, 0x1d, 0x11,
	0x73, 0x3f, 0x44, 0xc6, 0xaa, 0x0b, 0x88, 0xa1, 0x02, 0x67, 0xe1, 0x1b, 0xb0, 0x22, 0xfb, 0xce,
	0x8e, 0x91, 0xda, 0xa1, 0x92, 0xa7, 0x4b, 0xfc, 0xde, 0xfb, 0xff, 0xfe, 0xef, 0xee, 0xe9, 0xa1,
	0x41, 0x2c, 0x24, 0x39, 0x02, 0x27, 0x14, 0x1e, 0x17, 0x11, 0x17, 0x4e, 0x3a, 0x98, 0x81, 0x24,
	0x03, 0x27, 0xe0, 0x29, 0x24, 0x8c, 0x30, 0x0f, 0xdc, 0x38, 0xe1, 0x31, 0x17, 0x24, 0xb4, 0xe3,
	0x84, 0x4b, 0x6e, 0xad, 0xab, 0x12, 0xbb, 0x2c, 0xb1, 0x75, 0xc9, 0x46, 0x27, 0xe0, 0x01, 0x2f,
	0x72, 0x9c, 0xfc, 0x97, 0x4a, 0xdf, 0xc0, 0x5a, 0x78, 0x46, 0x04, 0x54, 0xea, 0x1e, 0xa7, 0x4c,
	0xc7, 0xef, 0xab, 0xb8, 0xab, 0x0a, 0xb5, 0xa4, 0x0a, 0x6d, 0x5f, 0x65, 0xae, 0x42, 0x17, 0x79,
	0x5b, 0x7f, 0x5b, 0x08, 0xbf, 0xa0, 0x6c, 0x0c, 0x31, 0x17, 0x54, 0x0e, 0x99, 0x3f, 0x01, 0x78,
	0x76, 0x48, 0x58, 0x00, 0x7b, 0xda, 0xba, 0xd5, 0x41, 0x2b, 0x92, 0xca, 0x10, 0xba, 0xe6, 0xa6,
	0xf9, 0xf0, 0xce, 0x54, 0xfd, 0xb1, 0x36, 0x51, 0xdb, 0x07, 0xe1, 0x25, 0x34, 0x96, 0x94, 0xb3,
	0xee, 0x8d, 0x22, 0x56, 0xff, 0x64, 0xbd, 0x45, 0xed, 0x88, 0x32, 0xd7, 0x57, 0xd2, 0xdd, 0x9b,
	0x79, 0xc6, 0xe8, 0xe9, 0xfc, 0xa2, 0x67, 0xfc, 0xbc, 0xe8, 0x6d, 0x07, 0x54, 0x1e, 0x1e, 0xcf,
	0x6c, 0x8f, 0x47, 0xda, 0xb8, 0x3e, 0xfa, 0xc2, 0x3f, 0x72, 0xe4, 0xc7, 0x18, 0x84, 0xbd, 0xcb,
	0xe4, 0xf9, 0x69, 0x1f, 0x69, 0xbf, 0xbb, 0x4c, 0x4e, 0x51, 0x54, 0x59, 0xb5, 0xde, 0x21, 0x4b,
	0xf5, 0x58, 0x12, 0xdc, 0x03, 0x80, 0x6e, 0xeb, 0xda, 0x94, 0x31, 0x78, 0x35, 0xca, 0x18, 0xbc,
	0xe9, 0x5d, 0xa5, 0xab, 0x41, 0x13, 0x80, 0x1a, 0x2b, 0x01, 0x75, 0xe6, 0xac, 0x95, 0xe6, 0x58,
	0x53, 0x25, 0xfb, 0x3f, 0xeb, 0x98, 0x2d, 0x59, 0xab, 0xcd, 0xb1, 0x5e, 0xb1, 0x8a, 0x15, 0xa3,
	0xb5, 0xaa, 0x2f, 0x1f, 0xa2, 0xe2, 0xdd, 0x0a, 0xdc, 0xad, 0x06, 0x70, 0xf7, 0xca, 0xd6, 0x4a,
	0xe5, 0x09, 0xc0, 0xe3, 0xdb, 0x5f, 0x4e, 0x7a, 0xc6, 0x9f, 0x93, 0x9e, 0xb1, 0xf5, 0xcd, 0x44,
	0x78, 0xaf, 0x74, 0x32, 0xf4, 0xfd, 0x04, 0x84, 0x68, 0x68, 0xf2, 0x26, 0xd5, 0x15, 0x1e, 0x00,
	0xb8, 0x44, 0x69, 0xeb, 0x01, 0xec, 0x9e, 0x9f, 0xf6, 0x3b, 0xda, 0xa5, 0xa6, 0xee, 0xcb, 0x84,
	0xb2, 0xa0, 0xbc, 0x9e, 0xa5, 0x9b, 0x9a, 0xd9, 0xef, 0x26, 0x7a, 0x30, 0x0c, 0x43, 0xfe, 0xfe,
	0x39, 0x15, 0x12, 0xfc, 0xd7, 0x24, 0xa4, 0x3e, 0x91, 0x3c, 0xd9, 0x07, 0xd9, 0x90, 0xe7, 0x10,
	0xad, 0x93, 0x1c, 0xe0, 0x86, 0x05, 0xc1, 0x4d, 0x4b, 0x84, 0x32, 0xde, 0xde, 0xb1, 0xed, 0x2b,
	0x96, 0x87, 0x7d, 0x99, 0x31, 0x31, 0x6a, 0xe5, 0x8f, 0x37, 0x5d, 0x23, 0x97, 0x05, 0x97, 0x9d,
	0x8d, 0xc8, 0xfc, 0x37, 0x36, 0x3e, 0x65, 0xd8, 0xf8, 0x9a, 0x61, 0x73, 0x9e, 0x61, 0xf3, 0x2c,
	0xc3, 0xe6, 0xaf, 0x0c, 0x9b, 0x9f, 0x17, 0xd8, 0x38, 0x5b, 0x60, 0xe3, 0xc7, 0x02, 0x1b, 0x6f,
	0x9e, 0xd4, 0x26, 0x20, 0x86, 0x44, 0xe4, 0x7a, 0xcc, 0x83, 0x97, 0x0c, 0x1c, 0xe5, 0xaa, 0xcf,
	0x88, 0xa4, 0x29, 0x38, 0xe9, 0x8e, 0xf3, 0x61, 0xb9, 0x74, 0x8a, 0xd1, 0x98, 0xad, 0x16, 0xab,
	0xe6, 0xd1, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8f, 0x37, 0xa1, 0x50, 0x31, 0x05, 0x00, 0x00,
}

func (m *MinDepositAndFeeChangeProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MinDepositAndFeeChangeProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MinDepositAndFeeChangeProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.PstakeRedemptionFee.Size()
		i -= size
		if _, err := m.PstakeRedemptionFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.PstakeUnstakeFee.Size()
		i -= size
		if _, err := m.PstakeUnstakeFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.PstakeRestakeFee.Size()
		i -= size
		if _, err := m.PstakeRestakeFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.PstakeDepositFee.Size()
		i -= size
		if _, err := m.PstakeDepositFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.MinDeposit.Size()
		i -= size
		if _, err := m.MinDeposit.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PstakeFeeAddressChangeProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PstakeFeeAddressChangeProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PstakeFeeAddressChangeProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PstakeFeeAddress) > 0 {
		i -= len(m.PstakeFeeAddress)
		copy(dAtA[i:], m.PstakeFeeAddress)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.PstakeFeeAddress)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AllowListedValidatorSetChangeProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AllowListedValidatorSetChangeProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AllowListedValidatorSetChangeProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.AllowListedValidators.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGovernanceProposal(dAtA []byte, offset int, v uint64) int {
	offset -= sovGovernanceProposal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MinDepositAndFeeChangeProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = m.MinDeposit.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	l = m.PstakeDepositFee.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	l = m.PstakeRestakeFee.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	l = m.PstakeUnstakeFee.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	l = m.PstakeRedemptionFee.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	return n
}

func (m *PstakeFeeAddressChangeProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = len(m.PstakeFeeAddress)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	return n
}

func (m *AllowListedValidatorSetChangeProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = m.AllowListedValidators.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	return n
}

func sovGovernanceProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGovernanceProposal(x uint64) (n int) {
	return sovGovernanceProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MinDepositAndFeeChangeProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGovernanceProposal
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
			return fmt.Errorf("proto: MinDepositAndFeeChangeProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MinDepositAndFeeChangeProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinDeposit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinDeposit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PstakeDepositFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PstakeDepositFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PstakeRestakeFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PstakeRestakeFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PstakeUnstakeFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PstakeUnstakeFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PstakeRedemptionFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PstakeRedemptionFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGovernanceProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGovernanceProposal
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
func (m *PstakeFeeAddressChangeProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGovernanceProposal
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
			return fmt.Errorf("proto: PstakeFeeAddressChangeProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PstakeFeeAddressChangeProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PstakeFeeAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PstakeFeeAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGovernanceProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGovernanceProposal
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
func (m *AllowListedValidatorSetChangeProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGovernanceProposal
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
			return fmt.Errorf("proto: AllowListedValidatorSetChangeProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AllowListedValidatorSetChangeProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllowListedValidators", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
				return ErrInvalidLengthGovernanceProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGovernanceProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AllowListedValidators.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGovernanceProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGovernanceProposal
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
func skipGovernanceProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGovernanceProposal
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
					return 0, ErrIntOverflowGovernanceProposal
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
					return 0, ErrIntOverflowGovernanceProposal
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
				return 0, ErrInvalidLengthGovernanceProposal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGovernanceProposal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGovernanceProposal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGovernanceProposal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGovernanceProposal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGovernanceProposal = fmt.Errorf("proto: unexpected end of group")
)
