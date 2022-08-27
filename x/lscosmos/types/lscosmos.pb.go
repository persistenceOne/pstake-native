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

type AllowListedValidators struct {
	AllowListedValidators []AllowListedValidator `protobuf:"bytes,1,rep,name=allow_listed_validators,json=allowListedValidators,proto3" json:"allow_listed_validators" yaml:"allow_lised_validators"`
}

func (m *AllowListedValidators) Reset()         { *m = AllowListedValidators{} }
func (m *AllowListedValidators) String() string { return proto.CompactTextString(m) }
func (*AllowListedValidators) ProtoMessage()    {}
func (*AllowListedValidators) Descriptor() ([]byte, []int) {
	return fileDescriptor_1043ccbf14211c19, []int{0}
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
	return fileDescriptor_1043ccbf14211c19, []int{1}
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

// CosmosParams go into the DB
type CosmosParams struct {
	ConnectionID     string                                 `protobuf:"bytes,1,opt,name=connectionID,proto3" json:"connectionID,omitempty"`
	TransferChannel  string                                 `protobuf:"bytes,2,opt,name=transfer_channel,json=transferChannel,proto3" json:"transfer_channel,omitempty"`
	TransferPort     string                                 `protobuf:"bytes,3,opt,name=transfer_port,json=transferPort,proto3" json:"transfer_port,omitempty"`
	BaseDenom        string                                 `protobuf:"bytes,4,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
	MintDenom        string                                 `protobuf:"bytes,5,opt,name=mint_denom,json=mintDenom,proto3" json:"mint_denom,omitempty"`
	MinDeposit       github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=min_deposit,json=minDeposit,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_deposit"`
	PstakeDepositFee github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=pstake_deposit_fee,json=pstakeDepositFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_deposit_fee"`
	PstakeRestakeFee github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=pstake_restake_fee,json=pstakeRestakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_restake_fee"`
	PstakeUnstakeFee github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,9,opt,name=pstake_unstake_fee,json=pstakeUnstakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_unstake_fee"`
}

func (m *CosmosParams) Reset()         { *m = CosmosParams{} }
func (m *CosmosParams) String() string { return proto.CompactTextString(m) }
func (*CosmosParams) ProtoMessage()    {}
func (*CosmosParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_1043ccbf14211c19, []int{2}
}
func (m *CosmosParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CosmosParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CosmosParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CosmosParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CosmosParams.Merge(m, src)
}
func (m *CosmosParams) XXX_Size() int {
	return m.Size()
}
func (m *CosmosParams) XXX_DiscardUnknown() {
	xxx_messageInfo_CosmosParams.DiscardUnknown(m)
}

var xxx_messageInfo_CosmosParams proto.InternalMessageInfo

// DelegationState stores module account balance, ica account balance, delegation state, undelegation state
type DelegationState struct {
	//This field is necessary as the address of not blocked for send coins,
	// we only should care about funds that have come via proper channels.
	HostDelegationAccountBalance github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=host_delegation_account_balance,json=hostDelegationAccountBalance,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"host_delegation_account_balance"`
	HostChainDelegationAddress   string                                   `protobuf:"bytes,2,opt,name=host_chain_delegation_address,json=hostChainDelegationAddress,proto3" json:"host_chain_delegation_address,omitempty"`
}

func (m *DelegationState) Reset()         { *m = DelegationState{} }
func (m *DelegationState) String() string { return proto.CompactTextString(m) }
func (*DelegationState) ProtoMessage()    {}
func (*DelegationState) Descriptor() ([]byte, []int) {
	return fileDescriptor_1043ccbf14211c19, []int{3}
}
func (m *DelegationState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DelegationState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DelegationState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DelegationState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DelegationState.Merge(m, src)
}
func (m *DelegationState) XXX_Size() int {
	return m.Size()
}
func (m *DelegationState) XXX_DiscardUnknown() {
	xxx_messageInfo_DelegationState.DiscardUnknown(m)
}

var xxx_messageInfo_DelegationState proto.InternalMessageInfo

func init() {
	proto.RegisterType((*AllowListedValidators)(nil), "lscosmos.v1beta1.AllowListedValidators")
	proto.RegisterType((*AllowListedValidator)(nil), "lscosmos.v1beta1.AllowListedValidator")
	proto.RegisterType((*CosmosParams)(nil), "lscosmos.v1beta1.CosmosParams")
	proto.RegisterType((*DelegationState)(nil), "lscosmos.v1beta1.DelegationState")
}

func init() { proto.RegisterFile("lscosmos/v1beta1/lscosmos.proto", fileDescriptor_1043ccbf14211c19) }

var fileDescriptor_1043ccbf14211c19 = []byte{
	// 658 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xcf, 0x4f, 0xd4, 0x4e,
	0x18, 0xc6, 0xb7, 0x5f, 0xf8, 0xa2, 0x0c, 0x10, 0xd6, 0x06, 0x62, 0x25, 0xd0, 0x9a, 0x1a, 0x09,
	0x1e, 0x68, 0x45, 0x2f, 0x86, 0xdb, 0xfe, 0x08, 0x09, 0x89, 0x09, 0xa4, 0x46, 0x4d, 0x8c, 0xb1,
	0x99, 0x9d, 0xbe, 0xec, 0x36, 0xb4, 0x33, 0x9b, 0xce, 0x00, 0x72, 0xf3, 0xe2, 0xdd, 0x78, 0xf2,
	0xa6, 0x47, 0xe3, 0xcd, 0xff, 0x82, 0x93, 0xe1, 0x68, 0x3c, 0xac, 0xba, 0x9c, 0xbc, 0xf2, 0x17,
	0x98, 0xe9, 0xdb, 0x2d, 0xac, 0xee, 0x41, 0xc2, 0xa9, 0x9d, 0xe7, 0x79, 0xe7, 0xf3, 0x3e, 0x9d,
	0xce, 0x0c, 0x71, 0x12, 0xc9, 0x84, 0x4c, 0x85, 0xf4, 0xf7, 0xd7, 0x5a, 0xa0, 0xe8, 0x9a, 0x3f,
	0x10, 0xbc, 0x6e, 0x26, 0x94, 0x30, 0xab, 0xe5, 0xb8, 0x28, 0x58, 0x98, 0x6b, 0x8b, 0xb6, 0xc8,
	0x4d, 0x5f, 0xbf, 0x61, 0xdd, 0x82, 0x5d, 0x60, 0x5a, 0x54, 0x42, 0xc9, 0x62, 0x22, 0xe6, 0xe8,
	0xbb, 0xef, 0x0d, 0x32, 0x5f, 0x4b, 0x12, 0x71, 0xf0, 0x30, 0x96, 0x0a, 0xa2, 0x27, 0x34, 0x89,
	0x23, 0xaa, 0x44, 0x26, 0xcd, 0xd7, 0x06, 0xb9, 0x4e, 0xb5, 0x13, 0x26, 0xb9, 0x15, 0xee, 0x97,
	0x9e, 0x65, 0xdc, 0x1c, 0x5b, 0x99, 0xba, 0xb7, 0xec, 0xfd, 0x19, 0xc2, 0x1b, 0x85, 0xaa, 0xdf,
	0x3e, 0xea, 0x39, 0x95, 0xd3, 0x9e, 0xb3, 0x74, 0x48, 0xd3, 0x64, 0xdd, 0x2d, 0xa1, 0x43, 0x4c,
	0x37, 0x98, 0xa7, 0xa3, 0x72, 0xb8, 0x5f, 0x0c, 0x32, 0x37, 0x0a, 0x6b, 0x6e, 0x92, 0x6b, 0xe5,
	0xf4, 0x90, 0x46, 0x51, 0x06, 0x52, 0x27, 0x33, 0x56, 0x26, 0xeb, 0x8b, 0xa7, 0x3d, 0xc7, 0xc2,
	0x6e, 0x7f, 0x95, 0xb8, 0x41, 0xb5, 0xd4, 0x6a, 0x28, 0x99, 0xbb, 0x64, 0x46, 0xd1, 0xac, 0x0d,
	0x2a, 0x3c, 0x80, 0xb8, 0xdd, 0x51, 0xd6, 0x7f, 0x39, 0x66, 0x43, 0x07, 0xff, 0xd6, 0x73, 0x96,
	0xdb, 0xb1, 0xea, 0xec, 0xb5, 0x3c, 0x26, 0x52, 0xbf, 0x58, 0x4f, 0x7c, 0xac, 0xca, 0x68, 0xd7,
	0x57, 0x87, 0x5d, 0x90, 0x5e, 0x13, 0xd8, 0x69, 0xcf, 0x99, 0xc3, 0xa6, 0x43, 0x30, 0x37, 0x98,
	0xc6, 0xf1, 0x53, 0x1c, 0x7e, 0x1e, 0x27, 0xd3, 0x8d, 0x7c, 0xfa, 0x36, 0xcd, 0x68, 0x2a, 0x4d,
	0x97, 0x4c, 0x33, 0xc1, 0x39, 0x30, 0x15, 0x0b, 0xbe, 0xd9, 0xc4, 0x6f, 0x08, 0x86, 0x34, 0xf3,
	0x0e, 0xa9, 0xaa, 0x8c, 0x72, 0xb9, 0x03, 0x59, 0xc8, 0x3a, 0x94, 0x73, 0x48, 0x30, 0x64, 0x30,
	0x3b, 0xd0, 0x1b, 0x28, 0x9b, 0xb7, 0xc8, 0x4c, 0x59, 0xda, 0x15, 0x99, 0xb2, 0xc6, 0x90, 0x37,
	0x10, 0xb7, 0x45, 0xa6, 0xcc, 0x25, 0x42, 0xf4, 0x96, 0x08, 0x23, 0xe0, 0x22, 0xb5, 0xc6, 0xf3,
	0x8a, 0x49, 0xad, 0x34, 0xb5, 0xa0, 0xed, 0x34, 0xe6, 0xaa, 0xb0, 0xff, 0x47, 0x5b, 0x2b, 0x68,
	0x6f, 0x91, 0xa9, 0x34, 0xe6, 0x61, 0x04, 0x5d, 0x21, 0x63, 0x65, 0x4d, 0xe4, 0xab, 0xe5, 0x5d,
	0x60, 0xb5, 0x36, 0xb9, 0x0a, 0x74, 0x87, 0x26, 0x12, 0xcc, 0xe7, 0xc4, 0xec, 0x4a, 0x45, 0x77,
	0x61, 0xc0, 0x0c, 0x77, 0x00, 0xac, 0x2b, 0x17, 0xe6, 0x36, 0x81, 0x05, 0x55, 0x24, 0x15, 0xe8,
	0x0d, 0x80, 0x73, 0xf4, 0x0c, 0xf0, 0xa9, 0xe9, 0x57, 0x2f, 0x43, 0x0f, 0x10, 0x34, 0x4c, 0xdf,
	0xe3, 0x67, 0xf4, 0xc9, 0xcb, 0xd0, 0x1f, 0xf3, 0x01, 0x7d, 0x7d, 0xfc, 0xdd, 0x07, 0xc7, 0x70,
	0x7f, 0x19, 0x64, 0xb6, 0x09, 0x09, 0xb4, 0xa9, 0xde, 0x0f, 0x8f, 0x14, 0x55, 0x60, 0xbe, 0x35,
	0x88, 0xd3, 0x11, 0x52, 0xff, 0xa4, 0x81, 0x11, 0x52, 0xc6, 0xc4, 0x1e, 0x57, 0x61, 0x8b, 0x26,
	0x94, 0x33, 0x28, 0x0e, 0xea, 0x0d, 0xaf, 0x38, 0xa6, 0xfa, 0x07, 0x97, 0x67, 0xb5, 0x21, 0x62,
	0x5e, 0xbf, 0xab, 0x03, 0x7e, 0xfa, 0xee, 0xac, 0xfc, 0x43, 0x40, 0x3d, 0x41, 0x06, 0x8b, 0xba,
	0xe7, 0x59, 0x96, 0x1a, 0x76, 0xac, 0x63, 0x43, 0xb3, 0x46, 0x96, 0xf2, 0x4c, 0xac, 0x43, 0xf3,
	0x0d, 0x72, 0x96, 0xac, 0x38, 0xa0, 0xb8, 0x69, 0x17, 0x74, 0x51, 0x43, 0xd7, 0x9c, 0x23, 0x61,
	0x45, 0xfd, 0xc5, 0xd1, 0x4f, 0xbb, 0xf2, 0xaa, 0x6f, 0x57, 0x3e, 0xf6, 0x6d, 0xe3, 0xa8, 0x6f,
	0x1b, 0xc7, 0x7d, 0xdb, 0xf8, 0xd1, 0xb7, 0x8d, 0x37, 0x27, 0x76, 0xe5, 0xf8, 0xc4, 0xae, 0x7c,
	0x3d, 0xb1, 0x2b, 0xcf, 0x1e, 0x9c, 0x0b, 0xdc, 0x85, 0x4c, 0xea, 0xcb, 0x81, 0x33, 0xd8, 0xe2,
	0xe0, 0xe3, 0x1a, 0xae, 0x72, 0xaa, 0xe2, 0x7d, 0xf0, 0x5f, 0x96, 0x57, 0x27, 0x7e, 0x46, 0x6b,
	0x22, 0xbf, 0xf9, 0xee, 0xff, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xd9, 0x35, 0x43, 0x05, 0x64, 0x05,
	0x00, 0x00,
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
func (this *CosmosParams) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*CosmosParams)
	if !ok {
		that2, ok := that.(CosmosParams)
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
	if this.ConnectionID != that1.ConnectionID {
		return false
	}
	if this.TransferChannel != that1.TransferChannel {
		return false
	}
	if this.TransferPort != that1.TransferPort {
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
	if !this.PstakeDepositFee.Equal(that1.PstakeDepositFee) {
		return false
	}
	if !this.PstakeRestakeFee.Equal(that1.PstakeRestakeFee) {
		return false
	}
	if !this.PstakeUnstakeFee.Equal(that1.PstakeUnstakeFee) {
		return false
	}
	return true
}
func (this *DelegationState) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*DelegationState)
	if !ok {
		that2, ok := that.(DelegationState)
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
	if len(this.HostDelegationAccountBalance) != len(that1.HostDelegationAccountBalance) {
		return false
	}
	for i := range this.HostDelegationAccountBalance {
		if !this.HostDelegationAccountBalance[i].Equal(&that1.HostDelegationAccountBalance[i]) {
			return false
		}
	}
	if this.HostChainDelegationAddress != that1.HostChainDelegationAddress {
		return false
	}
	return true
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

func (m *CosmosParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CosmosParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CosmosParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.PstakeUnstakeFee.Size()
		i -= size
		if _, err := m.PstakeUnstakeFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLscosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.PstakeRestakeFee.Size()
		i -= size
		if _, err := m.PstakeRestakeFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLscosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.PstakeDepositFee.Size()
		i -= size
		if _, err := m.PstakeDepositFee.MarshalTo(dAtA[i:]); err != nil {
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
	if len(m.TransferPort) > 0 {
		i -= len(m.TransferPort)
		copy(dAtA[i:], m.TransferPort)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.TransferPort)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.TransferChannel) > 0 {
		i -= len(m.TransferChannel)
		copy(dAtA[i:], m.TransferChannel)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.TransferChannel)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ConnectionID) > 0 {
		i -= len(m.ConnectionID)
		copy(dAtA[i:], m.ConnectionID)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.ConnectionID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DelegationState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DelegationState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DelegationState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.HostChainDelegationAddress) > 0 {
		i -= len(m.HostChainDelegationAddress)
		copy(dAtA[i:], m.HostChainDelegationAddress)
		i = encodeVarintLscosmos(dAtA, i, uint64(len(m.HostChainDelegationAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.HostDelegationAccountBalance) > 0 {
		for iNdEx := len(m.HostDelegationAccountBalance) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.HostDelegationAccountBalance[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *CosmosParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ConnectionID)
	if l > 0 {
		n += 1 + l + sovLscosmos(uint64(l))
	}
	l = len(m.TransferChannel)
	if l > 0 {
		n += 1 + l + sovLscosmos(uint64(l))
	}
	l = len(m.TransferPort)
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
	l = m.PstakeDepositFee.Size()
	n += 1 + l + sovLscosmos(uint64(l))
	l = m.PstakeRestakeFee.Size()
	n += 1 + l + sovLscosmos(uint64(l))
	l = m.PstakeUnstakeFee.Size()
	n += 1 + l + sovLscosmos(uint64(l))
	return n
}

func (m *DelegationState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.HostDelegationAccountBalance) > 0 {
		for _, e := range m.HostDelegationAccountBalance {
			l = e.Size()
			n += 1 + l + sovLscosmos(uint64(l))
		}
	}
	l = len(m.HostChainDelegationAddress)
	if l > 0 {
		n += 1 + l + sovLscosmos(uint64(l))
	}
	return n
}

func sovLscosmos(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLscosmos(x uint64) (n int) {
	return sovLscosmos(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
func (m *CosmosParams) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: CosmosParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CosmosParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionID", wireType)
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
			m.ConnectionID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferChannel", wireType)
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
			m.TransferChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferPort", wireType)
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
			m.TransferPort = string(dAtA[iNdEx:postIndex])
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
				return fmt.Errorf("proto: wrong wireType = %d for field PstakeDepositFee", wireType)
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
			if err := m.PstakeDepositFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PstakeRestakeFee", wireType)
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
			if err := m.PstakeRestakeFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PstakeUnstakeFee", wireType)
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
			if err := m.PstakeUnstakeFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *DelegationState) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: DelegationState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DelegationState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HostDelegationAccountBalance", wireType)
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
			m.HostDelegationAccountBalance = append(m.HostDelegationAccountBalance, types.Coin{})
			if err := m.HostDelegationAccountBalance[len(m.HostDelegationAccountBalance)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HostChainDelegationAddress", wireType)
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
			m.HostChainDelegationAddress = string(dAtA[iNdEx:postIndex])
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
