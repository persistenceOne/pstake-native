// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lscosmos/v1beta1/lscosmos.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
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

type DepositAmount struct {
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount" yaml:"amount"`
}

func (m *DepositAmount) Reset()         { *m = DepositAmount{} }
func (m *DepositAmount) String() string { return proto.CompactTextString(m) }
func (*DepositAmount) ProtoMessage()    {}
func (*DepositAmount) Descriptor() ([]byte, []int) {
	return fileDescriptor_1043ccbf14211c19, []int{0}
}
func (m *DepositAmount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DepositAmount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DepositAmount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DepositAmount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DepositAmount.Merge(m, src)
}
func (m *DepositAmount) XXX_Size() int {
	return m.Size()
}
func (m *DepositAmount) XXX_DiscardUnknown() {
	xxx_messageInfo_DepositAmount.DiscardUnknown(m)
}

var xxx_messageInfo_DepositAmount proto.InternalMessageInfo

type AllowListedValidators struct {
	AllowListedValidators []AllowListedValidator `protobuf:"bytes,1,rep,name=allow_listed_validators,json=allowListedValidators,proto3" json:"allow_listed_validators" yaml:"allow_lised_validators"`
}

func (m *AllowListedValidators) Reset()         { *m = AllowListedValidators{} }
func (m *AllowListedValidators) String() string { return proto.CompactTextString(m) }
func (*AllowListedValidators) ProtoMessage()    {}
func (*AllowListedValidators) Descriptor() ([]byte, []int) {
	return fileDescriptor_1043ccbf14211c19, []int{1}
}
func (m *AllowListedValidators) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AllowListedValidators) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AllowListedValidators.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AllowListedValidators) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllowListedValidators.Merge(m, src)
}
func (m *AllowListedValidators) XXX_Size() int {
	return m.Size()
}
func (m *AllowListedValidators) XXX_DiscardUnknown() {
	xxx_messageInfo_AllowListedValidators.DiscardUnknown(m)
}

var xxx_messageInfo_AllowListedValidators proto.InternalMessageInfo

type AllowListedValidator struct {
	// validator_address defines the bech32-encoded address the allowlisted validator
	ValidatorAddress string `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty" yaml:"validator_address"`
	// target_weight specifies the target weight for liquid staking, unstaking amount, which is a value for calculating
	// the real weight to be derived according to the active status
	TargetWeight github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=target_weight,json=targetWeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"target_weight" yaml:"target_weight"`
}

func (m *AllowListedValidator) Reset()         { *m = AllowListedValidator{} }
func (m *AllowListedValidator) String() string { return proto.CompactTextString(m) }
func (*AllowListedValidator) ProtoMessage()    {}
func (*AllowListedValidator) Descriptor() ([]byte, []int) {
	return fileDescriptor_1043ccbf14211c19, []int{2}
}
func (m *AllowListedValidator) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AllowListedValidator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AllowListedValidator.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AllowListedValidator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllowListedValidator.Merge(m, src)
}
func (m *AllowListedValidator) XXX_Size() int {
	return m.Size()
}
func (m *AllowListedValidator) XXX_DiscardUnknown() {
	xxx_messageInfo_AllowListedValidator.DiscardUnknown(m)
}

var xxx_messageInfo_AllowListedValidator proto.InternalMessageInfo

// CosmosIBCParams go into the DB
type CosmosIBCParams struct {
	IBCConnection        string                                 `protobuf:"bytes,1,opt,name=i_b_c_connection,json=iBCConnection,proto3" json:"i_b_c_connection,omitempty"`
	TokenTransferChannel string                                 `protobuf:"bytes,2,opt,name=token_transfer_channel,json=tokenTransferChannel,proto3" json:"token_transfer_channel,omitempty"`
	TokenTransferPort    string                                 `protobuf:"bytes,3,opt,name=token_transfer_port,json=tokenTransferPort,proto3" json:"token_transfer_port,omitempty"`
	BaseDenom            string                                 `protobuf:"bytes,4,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
	MintDenom            string                                 `protobuf:"bytes,5,opt,name=mint_denom,json=mintDenom,proto3" json:"mint_denom,omitempty"`
	MinDeposit           github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=min_deposit,json=minDeposit,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_deposit"`
	PStakeDepositFee     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=p_stake_deposit_fee,json=pStakeDepositFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"p_stake_deposit_fee"`
	PStakeRestakeFee     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=p_stake_restake_fee,json=pStakeRestakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"p_stake_restake_fee"`
	PStakeUnstakeFee     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,9,opt,name=p_stake_unstake_fee,json=pStakeUnstakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"p_stake_unstake_fee"`
}

func (m *CosmosIBCParams) Reset()         { *m = CosmosIBCParams{} }
func (m *CosmosIBCParams) String() string { return proto.CompactTextString(m) }
func (*CosmosIBCParams) ProtoMessage()    {}
func (*CosmosIBCParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_1043ccbf14211c19, []int{3}
}
func (m *CosmosIBCParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CosmosIBCParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CosmosIBCParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CosmosIBCParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CosmosIBCParams.Merge(m, src)
}
func (m *CosmosIBCParams) XXX_Size() int {
	return m.Size()
}
func (m *CosmosIBCParams) XXX_DiscardUnknown() {
	xxx_messageInfo_CosmosIBCParams.DiscardUnknown(m)
}

var xxx_messageInfo_CosmosIBCParams proto.InternalMessageInfo

func init() {
	proto.RegisterType((*DepositAmount)(nil), "lscosmos.v1beta1.DepositAmount")
	proto.RegisterType((*AllowListedValidators)(nil), "lscosmos.v1beta1.AllowListedValidators")
	proto.RegisterType((*AllowListedValidator)(nil), "lscosmos.v1beta1.AllowListedValidator")
	proto.RegisterType((*CosmosIBCParams)(nil), "lscosmos.v1beta1.CosmosIBCParams")
}

func init() { proto.RegisterFile("lscosmos/v1beta1/lscosmos.proto", fileDescriptor_1043ccbf14211c19) }

var fileDescriptor_1043ccbf14211c19 = []byte{
	// 629 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0x4f, 0x4f, 0xd4, 0x4e,
	0x18, 0xc7, 0xb7, 0x3f, 0xf8, 0xa1, 0x0c, 0x6e, 0x5c, 0xca, 0xa2, 0x95, 0x48, 0x6b, 0x9a, 0x88,
	0x5c, 0x68, 0x83, 0x7a, 0x30, 0xdc, 0x76, 0xbb, 0x21, 0x21, 0x31, 0x81, 0xd4, 0x7f, 0x89, 0x89,
	0x4e, 0x66, 0xdb, 0x87, 0x65, 0xb2, 0xed, 0x4c, 0xd3, 0x19, 0x40, 0x6e, 0x5e, 0xb8, 0x7b, 0xf4,
	0xa6, 0x47, 0xe3, 0x2b, 0xe1, 0x64, 0x88, 0x27, 0xe3, 0x61, 0xd5, 0xe5, 0x1d, 0xf0, 0x0a, 0xcc,
	0x74, 0xba, 0x15, 0x90, 0x83, 0x84, 0x53, 0x3b, 0xcf, 0xf7, 0x79, 0x3e, 0xdf, 0x6f, 0x9f, 0xb4,
	0x45, 0x4e, 0x22, 0x22, 0x2e, 0x52, 0x2e, 0xfc, 0x9d, 0xe5, 0x2e, 0x48, 0xb2, 0xec, 0x8f, 0x0a,
	0x5e, 0x96, 0x73, 0xc9, 0xcd, 0x46, 0x75, 0x2e, 0x1b, 0xe6, 0x9a, 0x3d, 0xde, 0xe3, 0x85, 0xe8,
	0xab, 0x3b, 0xdd, 0x37, 0x67, 0x97, 0x98, 0x2e, 0x11, 0x50, 0xb1, 0x22, 0x4e, 0x99, 0xd6, 0xdd,
	0x7d, 0x03, 0xd5, 0x3b, 0x90, 0x71, 0x41, 0x65, 0x2b, 0xe5, 0xdb, 0x4c, 0x9a, 0x12, 0x4d, 0x90,
	0xe2, 0xce, 0x32, 0xee, 0x8c, 0x2d, 0x4e, 0xdd, 0xbf, 0xe5, 0x95, 0x46, 0x0a, 0x31, 0x72, 0xf3,
	0x02, 0x4e, 0x59, 0xbb, 0x75, 0x30, 0x70, 0x6a, 0xc7, 0x03, 0xa7, 0xbe, 0x47, 0xd2, 0x64, 0xc5,
	0xd5, 0x63, 0xee, 0xe7, 0x1f, 0xce, 0x62, 0x8f, 0xca, 0xad, 0xed, 0xae, 0x17, 0xf1, 0xd4, 0x2f,
	0x03, 0xe8, 0xcb, 0x92, 0x88, 0xfb, 0xbe, 0xdc, 0xcb, 0x40, 0x14, 0x04, 0x11, 0x96, 0x5e, 0xee,
	0x07, 0x03, 0xcd, 0xb6, 0x92, 0x84, 0xef, 0x3e, 0xa6, 0x42, 0x42, 0xfc, 0x9c, 0x24, 0x34, 0x26,
	0x92, 0xe7, 0xc2, 0xdc, 0x37, 0xd0, 0x4d, 0xa2, 0x14, 0x9c, 0x14, 0x12, 0xde, 0xa9, 0xb4, 0x32,
	0xe1, 0x82, 0x77, 0x76, 0x19, 0xde, 0x79, 0xa8, 0xf6, 0xdd, 0x32, 0xee, 0x7c, 0x19, 0x77, 0x04,
	0x3d, 0xc5, 0x74, 0xc3, 0x59, 0x72, 0x5e, 0x0e, 0xf7, 0x8b, 0x81, 0x9a, 0xe7, 0x61, 0xcd, 0x35,
	0x34, 0x5d, 0x8d, 0x63, 0x12, 0xc7, 0x39, 0x08, 0x95, 0xcc, 0x58, 0x9c, 0x6c, 0xdf, 0x3e, 0x1e,
	0x38, 0x96, 0x76, 0xfb, 0xab, 0xc5, 0x0d, 0x1b, 0x55, 0xad, 0xa5, 0x4b, 0x66, 0x1f, 0xd5, 0x25,
	0xc9, 0x7b, 0x20, 0xf1, 0x2e, 0xd0, 0xde, 0x96, 0xb4, 0xfe, 0x2b, 0x30, 0xab, 0x2a, 0xf8, 0xf7,
	0x81, 0xb3, 0xf0, 0x0f, 0x6b, 0xed, 0x40, 0x74, 0x3c, 0x70, 0x9a, 0xda, 0xf4, 0x14, 0xcc, 0x0d,
	0xaf, 0xe9, 0xf3, 0x0b, 0x7d, 0xfc, 0x3a, 0x8e, 0xae, 0x07, 0xc5, 0xf8, 0x5a, 0x3b, 0xd8, 0x20,
	0x39, 0x49, 0x85, 0x79, 0x0f, 0x35, 0x28, 0xee, 0xe2, 0x08, 0x47, 0x9c, 0x31, 0x88, 0x24, 0xe5,
	0x4c, 0x3f, 0x4a, 0x58, 0xa7, 0xed, 0x20, 0xa8, 0x8a, 0xe6, 0x43, 0x74, 0x43, 0xf2, 0x3e, 0x30,
	0x2c, 0x73, 0xc2, 0xc4, 0x26, 0xe4, 0x38, 0xda, 0x22, 0x8c, 0x41, 0xa2, 0x23, 0x87, 0xcd, 0x42,
	0x7d, 0x5a, 0x8a, 0x81, 0xd6, 0x4c, 0x0f, 0xcd, 0x9c, 0x99, 0xca, 0x78, 0x2e, 0xad, 0xb1, 0x62,
	0x64, 0xfa, 0xd4, 0xc8, 0x06, 0xcf, 0xa5, 0x39, 0x8f, 0x90, 0x7a, 0xeb, 0x70, 0x0c, 0x8c, 0xa7,
	0xd6, 0x78, 0xd1, 0x36, 0xa9, 0x2a, 0x1d, 0x55, 0x50, 0x72, 0x4a, 0x99, 0x2c, 0xe5, 0xff, 0xb5,
	0xac, 0x2a, 0x5a, 0x5e, 0x47, 0x53, 0x29, 0x65, 0x38, 0xd6, 0xaf, 0xb7, 0x35, 0x51, 0xec, 0xd2,
	0xbb, 0xc0, 0x2e, 0xd7, 0x98, 0x0c, 0x95, 0x43, 0xf9, 0x81, 0x98, 0xaf, 0xd0, 0x4c, 0x86, 0x85,
	0x24, 0x7d, 0x18, 0x41, 0xf1, 0x26, 0x80, 0x75, 0xe5, 0xc2, 0xe0, 0x0e, 0x44, 0x61, 0x23, 0x7b,
	0xa2, 0x48, 0x25, 0x7b, 0x15, 0xe0, 0x24, 0x3e, 0x07, 0x7d, 0x55, 0xf8, 0xab, 0x97, 0xc1, 0x87,
	0x1a, 0x74, 0x06, 0xbf, 0xcd, 0xfe, 0xe0, 0x27, 0x2f, 0x83, 0x7f, 0xc6, 0x46, 0xf8, 0x95, 0xf1,
	0xf7, 0x1f, 0x1d, 0xa3, 0xfd, 0xfa, 0xe0, 0x97, 0x5d, 0x7b, 0x3b, 0xb4, 0x6b, 0x9f, 0x86, 0xb6,
	0x71, 0x30, 0xb4, 0x8d, 0xc3, 0xa1, 0x6d, 0xfc, 0x1c, 0xda, 0xc6, 0xbb, 0x23, 0xbb, 0x76, 0x78,
	0x64, 0xd7, 0xbe, 0x1d, 0xd9, 0xb5, 0x97, 0x8f, 0x4e, 0xb8, 0x64, 0x90, 0x0b, 0xf5, 0x45, 0xb1,
	0x08, 0xd6, 0x19, 0xf8, 0x59, 0x01, 0x5c, 0x62, 0x44, 0xd2, 0x1d, 0xf0, 0xdf, 0x54, 0xff, 0x3d,
	0xed, 0xdd, 0x9d, 0x28, 0x7e, 0x5b, 0x0f, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0x22, 0x91, 0xea,
	0x16, 0x21, 0x05, 0x00, 0x00,
}

func (this *DepositAmount) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*DepositAmount)
	if !ok {
		that2, ok := that.(DepositAmount)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Amount) != len(that1.Amount) {
		return false
	}
	for i := range this.Amount {
		if !this.Amount[i].Equal(&that1.Amount[i]) {
			return false
		}
	}
	return true
}
func (this *AllowListedValidators) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AllowListedValidators)
	if !ok {
		that2, ok := that.(AllowListedValidators)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.AllowListedValidators) != len(that1.AllowListedValidators) {
		return false
	}
	for i := range this.AllowListedValidators {
		if !this.AllowListedValidators[i].Equal(&that1.AllowListedValidators[i]) {
			return false
		}
	}
	return true
}
func (this *AllowListedValidator) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AllowListedValidator)
	if !ok {
		that2, ok := that.(AllowListedValidator)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.ValidatorAddress != that1.ValidatorAddress {
		return false
	}
	if !this.TargetWeight.Equal(that1.TargetWeight) {
		return false
	}
	return true
}
func (this *CosmosIBCParams) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*CosmosIBCParams)
	if !ok {
		that2, ok := that.(CosmosIBCParams)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.IBCConnection != that1.IBCConnection {
		return false
	}
	if this.TokenTransferChannel != that1.TokenTransferChannel {
		return false
	}
	if this.TokenTransferPort != that1.TokenTransferPort {
		return false
	}
	if this.BaseDenom != that1.BaseDenom {
		return false
	}
	if this.MintDenom != that1.MintDenom {
		return false
	}
	if !this.MinDeposit.Equal(that1.MinDeposit) {
		return false
	}
	if !this.PStakeDepositFee.Equal(that1.PStakeDepositFee) {
		return false
	}
	if !this.PStakeRestakeFee.Equal(that1.PStakeRestakeFee) {
		return false
	}
	if !this.PStakeUnstakeFee.Equal(that1.PStakeUnstakeFee) {
		return false
	}
	return true
}
func (m *DepositAmount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DepositAmount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DepositAmount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
				i = encodeVarintLscosmos(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *AllowListedValidators) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AllowListedValidators) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AllowListedValidators) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AllowListedValidators) > 0 {
		for iNdEx := len(m.AllowListedValidators) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AllowListedValidators[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintLscosmos(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *AllowListedValidator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AllowListedValidator) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AllowListedValidator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.TargetWeight.Size()
		i -= size
		if _, err := m.TargetWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLscosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *CosmosIBCParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CosmosIBCParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CosmosIBCParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.PStakeUnstakeFee.Size()
		i -= size
		if _, err := m.PStakeUnstakeFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLscosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.PStakeRestakeFee.Size()
		i -= size
		if _, err := m.PStakeRestakeFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLscosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.PStakeDepositFee.Size()
		i -= size
		if _, err := m.PStakeDepositFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLscosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.MinDeposit.Size()
		i -= size
		if _, err := m.MinDeposit.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLscosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if len(m.MintDenom) > 0 {
		i -= len(m.MintDenom)
		copy(dAtA[i:], m.MintDenom)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.MintDenom)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.BaseDenom) > 0 {
		i -= len(m.BaseDenom)
		copy(dAtA[i:], m.BaseDenom)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.BaseDenom)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.TokenTransferPort) > 0 {
		i -= len(m.TokenTransferPort)
		copy(dAtA[i:], m.TokenTransferPort)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.TokenTransferPort)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.TokenTransferChannel) > 0 {
		i -= len(m.TokenTransferChannel)
		copy(dAtA[i:], m.TokenTransferChannel)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.TokenTransferChannel)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.IBCConnection) > 0 {
		i -= len(m.IBCConnection)
		copy(dAtA[i:], m.IBCConnection)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.IBCConnection)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintLscosmos(dAtA []byte, offset int, v uint64) int {
	offset -= sovLscosmos(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DepositAmount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovLscosmos(uint64(l))
		}
	}
	return n
}

func (m *AllowListedValidators) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.AllowListedValidators) > 0 {
		for _, e := range m.AllowListedValidators {
			l = e.Size()
			n += 1 + l + sovLscosmos(uint64(l))
		}
	}
	return n
}

func (m *AllowListedValidator) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovLscosmos(uint64(l))
	}
	l = m.TargetWeight.Size()
	n += 1 + l + sovLscosmos(uint64(l))
	return n
}

func (m *CosmosIBCParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.IBCConnection)
	if l > 0 {
		n += 1 + l + sovLscosmos(uint64(l))
	}
	l = len(m.TokenTransferChannel)
	if l > 0 {
		n += 1 + l + sovLscosmos(uint64(l))
	}
	l = len(m.TokenTransferPort)
	if l > 0 {
		n += 1 + l + sovLscosmos(uint64(l))
	}
	l = len(m.BaseDenom)
	if l > 0 {
		n += 1 + l + sovLscosmos(uint64(l))
	}
	l = len(m.MintDenom)
	if l > 0 {
		n += 1 + l + sovLscosmos(uint64(l))
	}
	l = m.MinDeposit.Size()
	n += 1 + l + sovLscosmos(uint64(l))
	l = m.PStakeDepositFee.Size()
	n += 1 + l + sovLscosmos(uint64(l))
	l = m.PStakeRestakeFee.Size()
	n += 1 + l + sovLscosmos(uint64(l))
	l = m.PStakeUnstakeFee.Size()
	n += 1 + l + sovLscosmos(uint64(l))
	return n
}

func sovLscosmos(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLscosmos(x uint64) (n int) {
	return sovLscosmos(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DepositAmount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLscosmos
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
			return fmt.Errorf("proto: DepositAmount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DepositAmount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
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
			skippy, err := skipLscosmos(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLscosmos
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
func (m *AllowListedValidators) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLscosmos
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
			return fmt.Errorf("proto: AllowListedValidators: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AllowListedValidators: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllowListedValidators", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AllowListedValidators = append(m.AllowListedValidators, AllowListedValidator{})
			if err := m.AllowListedValidators[len(m.AllowListedValidators)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLscosmos(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLscosmos
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
func (m *AllowListedValidator) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLscosmos
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
			return fmt.Errorf("proto: AllowListedValidator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AllowListedValidator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TargetWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TargetWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLscosmos(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLscosmos
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
func (m *CosmosIBCParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLscosmos
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
			return fmt.Errorf("proto: CosmosIBCParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CosmosIBCParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IBCConnection", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IBCConnection = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenTransferChannel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenTransferChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenTransferPort", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenTransferPort = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BaseDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MintDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MintDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinDeposit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinDeposit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PStakeDepositFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PStakeDepositFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PStakeRestakeFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PStakeRestakeFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PStakeUnstakeFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLscosmos
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
				return ErrInvalidLengthLscosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLscosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PStakeUnstakeFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLscosmos(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLscosmos
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
func skipLscosmos(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLscosmos
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
					return 0, ErrIntOverflowLscosmos
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
					return 0, ErrIntOverflowLscosmos
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
				return 0, ErrInvalidLengthLscosmos
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLscosmos
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLscosmos
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLscosmos        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLscosmos          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLscosmos = fmt.Errorf("proto: unexpected end of group")
)
