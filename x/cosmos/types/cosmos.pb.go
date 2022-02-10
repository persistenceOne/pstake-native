// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pstake/cosmos/v1beta1/cosmos.proto

package types

import (
	fmt "fmt"
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

type WeightedAddress struct {
	Address string                                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`
	Weight  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=weight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"weight" yaml:"weight"`
}

func (m *WeightedAddress) Reset()         { *m = WeightedAddress{} }
func (m *WeightedAddress) String() string { return proto.CompactTextString(m) }
func (*WeightedAddress) ProtoMessage()    {}
func (*WeightedAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_b56098799d2a1326, []int{0}
}
func (m *WeightedAddress) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WeightedAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WeightedAddress.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WeightedAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeightedAddress.Merge(m, src)
}
func (m *WeightedAddress) XXX_Size() int {
	return m.Size()
}
func (m *WeightedAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_WeightedAddress.DiscardUnknown(m)
}

var xxx_messageInfo_WeightedAddress proto.InternalMessageInfo

func (m *WeightedAddress) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type DistributionProportions struct {
	ValidatorRewards github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=validator_rewards,json=validatorRewards,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"validator_rewards" yaml:"validator_rewards"`
	DeveloperRewards github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=developer_rewards,json=developerRewards,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"developer_rewards" yaml:"developer_rewards"`
}

func (m *DistributionProportions) Reset()         { *m = DistributionProportions{} }
func (m *DistributionProportions) String() string { return proto.CompactTextString(m) }
func (*DistributionProportions) ProtoMessage()    {}
func (*DistributionProportions) Descriptor() ([]byte, []int) {
	return fileDescriptor_b56098799d2a1326, []int{1}
}
func (m *DistributionProportions) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DistributionProportions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DistributionProportions.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DistributionProportions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DistributionProportions.Merge(m, src)
}
func (m *DistributionProportions) XXX_Size() int {
	return m.Size()
}
func (m *DistributionProportions) XXX_DiscardUnknown() {
	xxx_messageInfo_DistributionProportions.DiscardUnknown(m)
}

var xxx_messageInfo_DistributionProportions proto.InternalMessageInfo

//Params hold the parameters for cosmos module
type Params struct {
	MinMintingAmount                  uint64                  `protobuf:"varint,1,opt,name=minMintingAmount,proto3" json:"minMintingAmount,omitempty" yaml:"minMintingAmount"`
	MaxMintingAmount                  uint64                  `protobuf:"varint,2,opt,name=maxMintingAmount,proto3" json:"maxMintingAmount,omitempty" yaml:"maxMintingAmount"`
	MinBurningAmount                  uint64                  `protobuf:"varint,3,opt,name=minBurningAmount,proto3" json:"minBurningAmount,omitempty" yaml:"minBurningAmount"`
	MaxBurningAmount                  uint64                  `protobuf:"varint,4,opt,name=maxBurningAmount,proto3" json:"maxBurningAmount,omitempty" yaml:"maxBurningAmount"`
	MaxValidatorToDelegate            uint64                  `protobuf:"varint,5,opt,name=maxValidatorToDelegate,proto3" json:"maxValidatorToDelegate,omitempty" yaml:"maxValidatorToDelegate"`
	ValidatorSetCosmosChain           []WeightedAddress       `protobuf:"bytes,6,rep,name=validatorSetCosmosChain,proto3" json:"validatorSetCosmosChain" yaml:"validatorSetCosmosChain"`
	ValidatorSetNativeChain           []WeightedAddress       `protobuf:"bytes,7,rep,name=validatorSetNativeChain,proto3" json:"validatorSetNativeChain" yaml:"validatorSetNativeChain"`
	WeightedDeveloperRewardsReceivers []WeightedAddress       `protobuf:"bytes,8,rep,name=weightedDeveloperRewardsReceivers,proto3" json:"weightedDeveloperRewardsReceivers" yaml:"weightedDeveloperRewardsReceivers"`
	DistributionProportion            DistributionProportions `protobuf:"bytes,9,opt,name=distributionProportion,proto3" json:"distributionProportion" yaml:"weightedDeveloperRewardsReceivers"`
	Epochs                            int64                   `protobuf:"varint,10,opt,name=epochs,proto3" json:"epochs,omitempty" yaml:"epochs"`
	MaxIncomingAndOutgoingTxns        int64                   `protobuf:"varint,11,opt,name=maxIncomingAndOutgoingTxns,proto3" json:"maxIncomingAndOutgoingTxns,omitempty" yaml:"maxIncomingAndOutgoingTxns"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_b56098799d2a1326, []int{2}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetMinMintingAmount() uint64 {
	if m != nil {
		return m.MinMintingAmount
	}
	return 0
}

func (m *Params) GetMaxMintingAmount() uint64 {
	if m != nil {
		return m.MaxMintingAmount
	}
	return 0
}

func (m *Params) GetMinBurningAmount() uint64 {
	if m != nil {
		return m.MinBurningAmount
	}
	return 0
}

func (m *Params) GetMaxBurningAmount() uint64 {
	if m != nil {
		return m.MaxBurningAmount
	}
	return 0
}

func (m *Params) GetMaxValidatorToDelegate() uint64 {
	if m != nil {
		return m.MaxValidatorToDelegate
	}
	return 0
}

func (m *Params) GetValidatorSetCosmosChain() []WeightedAddress {
	if m != nil {
		return m.ValidatorSetCosmosChain
	}
	return nil
}

func (m *Params) GetValidatorSetNativeChain() []WeightedAddress {
	if m != nil {
		return m.ValidatorSetNativeChain
	}
	return nil
}

func (m *Params) GetWeightedDeveloperRewardsReceivers() []WeightedAddress {
	if m != nil {
		return m.WeightedDeveloperRewardsReceivers
	}
	return nil
}

func (m *Params) GetDistributionProportion() DistributionProportions {
	if m != nil {
		return m.DistributionProportion
	}
	return DistributionProportions{}
}

func (m *Params) GetEpochs() int64 {
	if m != nil {
		return m.Epochs
	}
	return 0
}

func (m *Params) GetMaxIncomingAndOutgoingTxns() int64 {
	if m != nil {
		return m.MaxIncomingAndOutgoingTxns
	}
	return 0
}

func init() {
	proto.RegisterType((*WeightedAddress)(nil), "pstake.cosmos.v1beta1.WeightedAddress")
	proto.RegisterType((*DistributionProportions)(nil), "pstake.cosmos.v1beta1.DistributionProportions")
	proto.RegisterType((*Params)(nil), "pstake.cosmos.v1beta1.Params")
}

func init() {
	proto.RegisterFile("pstake/cosmos/v1beta1/cosmos.proto", fileDescriptor_b56098799d2a1326)
}

var fileDescriptor_b56098799d2a1326 = []byte{
	// 664 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x95, 0x4d, 0x4f, 0x13, 0x41,
	0x18, 0xc7, 0xbb, 0x14, 0x0b, 0x0c, 0x51, 0x61, 0xa3, 0xb0, 0xc1, 0xb8, 0x5b, 0x26, 0x91, 0xd4,
	0x44, 0xb6, 0x82, 0x89, 0x07, 0x2e, 0x86, 0xd2, 0xc4, 0x68, 0xa2, 0x90, 0x85, 0x48, 0xf4, 0x62,
	0xa6, 0xbb, 0x4f, 0xb6, 0x13, 0xba, 0x33, 0x9b, 0x9d, 0xe9, 0x0b, 0xdf, 0xc0, 0xc4, 0x8b, 0x47,
	0x8e, 0x9e, 0xfc, 0x2c, 0x1c, 0x39, 0x1a, 0x0f, 0x1b, 0x03, 0xdf, 0xa0, 0x37, 0x6f, 0x66, 0x5f,
	0xfa, 0xb6, 0x6d, 0x31, 0xa8, 0xa7, 0x4e, 0x9f, 0xfd, 0xcf, 0xef, 0xff, 0x6f, 0xe7, 0x99, 0x67,
	0x11, 0xf6, 0x85, 0x24, 0x27, 0x50, 0xb6, 0xb9, 0xf0, 0xb8, 0x28, 0xb7, 0xb6, 0x6a, 0x20, 0xc9,
	0x56, 0xfa, 0xd5, 0xf4, 0x03, 0x2e, 0xb9, 0x7a, 0x3f, 0xd1, 0x98, 0x69, 0x31, 0xd5, 0xac, 0xdd,
	0x73, 0xb9, 0xcb, 0x63, 0x45, 0x39, 0x5a, 0x25, 0x62, 0x7c, 0xa6, 0xa0, 0xbb, 0xc7, 0x40, 0xdd,
	0xba, 0x04, 0x67, 0xd7, 0x71, 0x02, 0x10, 0x42, 0x7d, 0x82, 0xe6, 0x48, 0xb2, 0xd4, 0x94, 0xa2,
	0x52, 0x5a, 0xa8, 0xa8, 0xdd, 0xd0, 0xb8, 0x73, 0x4a, 0xbc, 0xc6, 0x0e, 0x4e, 0x1f, 0x60, 0xab,
	0x27, 0x51, 0x8f, 0x51, 0xa1, 0x1d, 0x03, 0xb4, 0x99, 0x58, 0xfc, 0xe2, 0x3c, 0x34, 0x72, 0x3f,
	0x42, 0x63, 0xc3, 0xa5, 0xb2, 0xde, 0xac, 0x99, 0x36, 0xf7, 0x7a, 0x71, 0x93, 0x8f, 0x4d, 0xe1,
	0x9c, 0x94, 0xe5, 0xa9, 0x0f, 0xc2, 0xac, 0x82, 0xdd, 0x0d, 0x8d, 0xdb, 0x09, 0x3a, 0xa1, 0x60,
	0x2b, 0xc5, 0xe1, 0xcf, 0x33, 0x68, 0xb5, 0x4a, 0x85, 0x0c, 0x68, 0xad, 0x29, 0x29, 0x67, 0x07,
	0x01, 0xf7, 0x79, 0x10, 0xad, 0x84, 0xda, 0x46, 0xcb, 0x2d, 0xd2, 0xa0, 0x0e, 0x91, 0x3c, 0xf8,
	0x18, 0x40, 0x9b, 0x04, 0x4e, 0x2f, 0xec, 0xeb, 0x1b, 0xfb, 0x6b, 0x89, 0xff, 0x18, 0x10, 0x5b,
	0x4b, 0xfd, 0x9a, 0x95, 0x94, 0x22, 0x63, 0x07, 0x5a, 0xd0, 0xe0, 0x3e, 0x0c, 0x8c, 0xf3, 0xff,
	0x66, 0x3c, 0x06, 0xc4, 0xd6, 0x52, 0xbf, 0x96, 0x1a, 0xe3, 0x5f, 0xf3, 0xa8, 0x70, 0x40, 0x02,
	0xe2, 0x09, 0xf5, 0x25, 0x5a, 0xf2, 0x28, 0x7b, 0x43, 0x99, 0xa4, 0xcc, 0xdd, 0xf5, 0x78, 0x93,
	0xc9, 0xf8, 0xb7, 0xcf, 0x56, 0x1e, 0x74, 0x43, 0x63, 0x35, 0x81, 0x66, 0x15, 0xd8, 0x1a, 0xdb,
	0x14, 0x83, 0x48, 0x67, 0x14, 0x34, 0x33, 0x06, 0xca, 0x28, 0x22, 0x50, 0xa6, 0x94, 0x26, 0xaa,
	0x34, 0x03, 0x36, 0x00, 0xe5, 0x27, 0x24, 0x1a, 0x51, 0x24, 0x89, 0x46, 0x4a, 0x69, 0xa2, 0x51,
	0xd0, 0xec, 0x84, 0x44, 0x63, 0xa0, 0x4c, 0x49, 0x7d, 0x8f, 0x56, 0x3c, 0xd2, 0x79, 0xd7, 0x3b,
	0xbe, 0x23, 0x5e, 0x85, 0x06, 0xb8, 0x44, 0x82, 0x76, 0x2b, 0xc6, 0xad, 0x77, 0x43, 0xe3, 0x61,
	0x1f, 0x37, 0x41, 0x87, 0xad, 0x29, 0x00, 0xf5, 0x93, 0x82, 0x56, 0xfb, 0x7d, 0x71, 0x08, 0x72,
	0x2f, 0x3e, 0xdc, 0xbd, 0x3a, 0xa1, 0x4c, 0x2b, 0x14, 0xf3, 0xa5, 0xc5, 0xed, 0x0d, 0x73, 0xe2,
	0x15, 0x34, 0x33, 0x17, 0xad, 0xb2, 0x11, 0x75, 0x4c, 0x37, 0x34, 0xf4, 0x4c, 0x03, 0x8e, 0x42,
	0xb1, 0x35, 0xcd, 0x6e, 0x2c, 0xca, 0x5b, 0x22, 0x69, 0x0b, 0x92, 0x28, 0x73, 0xff, 0x2b, 0xca,
	0x10, 0x34, 0x13, 0x65, 0xe8, 0x89, 0xfa, 0x4d, 0x41, 0xeb, 0xed, 0x14, 0x5a, 0xcd, 0x34, 0xaf,
	0x05, 0x36, 0xd0, 0x16, 0x04, 0x42, 0x9b, 0xbf, 0x51, 0xa8, 0xa7, 0x69, 0xa8, 0xd2, 0xf0, 0x80,
	0xb8, 0x06, 0x8f, 0xad, 0x3f, 0x47, 0x50, 0xcf, 0x14, 0xb4, 0xe2, 0x4c, 0x1c, 0x2b, 0xda, 0x42,
	0x51, 0x29, 0x2d, 0x6e, 0x9b, 0x53, 0xd2, 0x4d, 0x99, 0x45, 0x7f, 0x91, 0x72, 0x8a, 0xbf, 0xfa,
	0x18, 0x15, 0xc0, 0xe7, 0x76, 0x5d, 0x68, 0xa8, 0xa8, 0x94, 0xf2, 0x95, 0xe5, 0xc1, 0x70, 0x4c,
	0xea, 0xd8, 0x4a, 0x05, 0x2a, 0xa0, 0x35, 0x8f, 0x74, 0x5e, 0x31, 0x9b, 0x7b, 0x51, 0xd3, 0x33,
	0x67, 0xbf, 0x29, 0x5d, 0x4e, 0x99, 0x7b, 0xd4, 0x61, 0x42, 0x5b, 0x8c, 0xb7, 0x3f, 0xea, 0x86,
	0xc6, 0x7a, 0xbf, 0xc7, 0xa7, 0x68, 0xb1, 0x75, 0x0d, 0x68, 0x67, 0xf6, 0xec, 0xab, 0x91, 0xab,
	0x1c, 0x9c, 0x5f, 0xea, 0xca, 0xc5, 0xa5, 0xae, 0xfc, 0xbc, 0xd4, 0x95, 0x2f, 0x57, 0x7a, 0xee,
	0xe2, 0x4a, 0xcf, 0x7d, 0xbf, 0xd2, 0x73, 0x1f, 0x9e, 0x0f, 0xcd, 0x3a, 0x1f, 0x02, 0x41, 0x85,
	0x04, 0x66, 0xc3, 0x3e, 0x83, 0xb2, 0x7f, 0x18, 0xfd, 0x89, 0x9b, 0x2c, 0xee, 0x91, 0x72, 0xa7,
	0x37, 0x0a, 0xe3, 0xf9, 0x57, 0x2b, 0xc4, 0x6f, 0x9f, 0x67, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff,
	0x75, 0xe2, 0x2b, 0x07, 0xd0, 0x06, 0x00, 0x00,
}

func (m *WeightedAddress) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WeightedAddress) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WeightedAddress) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Weight.Size()
		i -= size
		if _, err := m.Weight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintCosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintCosmos(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DistributionProportions) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DistributionProportions) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DistributionProportions) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.DeveloperRewards.Size()
		i -= size
		if _, err := m.DeveloperRewards.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintCosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.ValidatorRewards.Size()
		i -= size
		if _, err := m.ValidatorRewards.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintCosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MaxIncomingAndOutgoingTxns != 0 {
		i = encodeVarintCosmos(dAtA, i, uint64(m.MaxIncomingAndOutgoingTxns))
		i--
		dAtA[i] = 0x58
	}
	if m.Epochs != 0 {
		i = encodeVarintCosmos(dAtA, i, uint64(m.Epochs))
		i--
		dAtA[i] = 0x50
	}
	{
		size, err := m.DistributionProportion.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintCosmos(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	if len(m.WeightedDeveloperRewardsReceivers) > 0 {
		for iNdEx := len(m.WeightedDeveloperRewardsReceivers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.WeightedDeveloperRewardsReceivers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCosmos(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.ValidatorSetNativeChain) > 0 {
		for iNdEx := len(m.ValidatorSetNativeChain) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ValidatorSetNativeChain[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCosmos(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.ValidatorSetCosmosChain) > 0 {
		for iNdEx := len(m.ValidatorSetCosmosChain) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ValidatorSetCosmosChain[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCosmos(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if m.MaxValidatorToDelegate != 0 {
		i = encodeVarintCosmos(dAtA, i, uint64(m.MaxValidatorToDelegate))
		i--
		dAtA[i] = 0x28
	}
	if m.MaxBurningAmount != 0 {
		i = encodeVarintCosmos(dAtA, i, uint64(m.MaxBurningAmount))
		i--
		dAtA[i] = 0x20
	}
	if m.MinBurningAmount != 0 {
		i = encodeVarintCosmos(dAtA, i, uint64(m.MinBurningAmount))
		i--
		dAtA[i] = 0x18
	}
	if m.MaxMintingAmount != 0 {
		i = encodeVarintCosmos(dAtA, i, uint64(m.MaxMintingAmount))
		i--
		dAtA[i] = 0x10
	}
	if m.MinMintingAmount != 0 {
		i = encodeVarintCosmos(dAtA, i, uint64(m.MinMintingAmount))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintCosmos(dAtA []byte, offset int, v uint64) int {
	offset -= sovCosmos(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *WeightedAddress) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovCosmos(uint64(l))
	}
	l = m.Weight.Size()
	n += 1 + l + sovCosmos(uint64(l))
	return n
}

func (m *DistributionProportions) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ValidatorRewards.Size()
	n += 1 + l + sovCosmos(uint64(l))
	l = m.DeveloperRewards.Size()
	n += 1 + l + sovCosmos(uint64(l))
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MinMintingAmount != 0 {
		n += 1 + sovCosmos(uint64(m.MinMintingAmount))
	}
	if m.MaxMintingAmount != 0 {
		n += 1 + sovCosmos(uint64(m.MaxMintingAmount))
	}
	if m.MinBurningAmount != 0 {
		n += 1 + sovCosmos(uint64(m.MinBurningAmount))
	}
	if m.MaxBurningAmount != 0 {
		n += 1 + sovCosmos(uint64(m.MaxBurningAmount))
	}
	if m.MaxValidatorToDelegate != 0 {
		n += 1 + sovCosmos(uint64(m.MaxValidatorToDelegate))
	}
	if len(m.ValidatorSetCosmosChain) > 0 {
		for _, e := range m.ValidatorSetCosmosChain {
			l = e.Size()
			n += 1 + l + sovCosmos(uint64(l))
		}
	}
	if len(m.ValidatorSetNativeChain) > 0 {
		for _, e := range m.ValidatorSetNativeChain {
			l = e.Size()
			n += 1 + l + sovCosmos(uint64(l))
		}
	}
	if len(m.WeightedDeveloperRewardsReceivers) > 0 {
		for _, e := range m.WeightedDeveloperRewardsReceivers {
			l = e.Size()
			n += 1 + l + sovCosmos(uint64(l))
		}
	}
	l = m.DistributionProportion.Size()
	n += 1 + l + sovCosmos(uint64(l))
	if m.Epochs != 0 {
		n += 1 + sovCosmos(uint64(m.Epochs))
	}
	if m.MaxIncomingAndOutgoingTxns != 0 {
		n += 1 + sovCosmos(uint64(m.MaxIncomingAndOutgoingTxns))
	}
	return n
}

func sovCosmos(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCosmos(x uint64) (n int) {
	return sovCosmos(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *WeightedAddress) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCosmos
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
			return fmt.Errorf("proto: WeightedAddress: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WeightedAddress: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
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
				return ErrInvalidLengthCosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
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
				return ErrInvalidLengthCosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Weight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCosmos(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCosmos
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
func (m *DistributionProportions) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCosmos
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
			return fmt.Errorf("proto: DistributionProportions: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DistributionProportions: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorRewards", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
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
				return ErrInvalidLengthCosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ValidatorRewards.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeveloperRewards", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
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
				return ErrInvalidLengthCosmos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DeveloperRewards.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCosmos(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCosmos
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
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCosmos
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinMintingAmount", wireType)
			}
			m.MinMintingAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinMintingAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxMintingAmount", wireType)
			}
			m.MaxMintingAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxMintingAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinBurningAmount", wireType)
			}
			m.MinBurningAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinBurningAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxBurningAmount", wireType)
			}
			m.MaxBurningAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxBurningAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxValidatorToDelegate", wireType)
			}
			m.MaxValidatorToDelegate = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxValidatorToDelegate |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorSetCosmosChain", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
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
				return ErrInvalidLengthCosmos
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorSetCosmosChain = append(m.ValidatorSetCosmosChain, WeightedAddress{})
			if err := m.ValidatorSetCosmosChain[len(m.ValidatorSetCosmosChain)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorSetNativeChain", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
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
				return ErrInvalidLengthCosmos
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorSetNativeChain = append(m.ValidatorSetNativeChain, WeightedAddress{})
			if err := m.ValidatorSetNativeChain[len(m.ValidatorSetNativeChain)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WeightedDeveloperRewardsReceivers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
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
				return ErrInvalidLengthCosmos
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WeightedDeveloperRewardsReceivers = append(m.WeightedDeveloperRewardsReceivers, WeightedAddress{})
			if err := m.WeightedDeveloperRewardsReceivers[len(m.WeightedDeveloperRewardsReceivers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DistributionProportion", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
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
				return ErrInvalidLengthCosmos
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCosmos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DistributionProportion.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Epochs", wireType)
			}
			m.Epochs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Epochs |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxIncomingAndOutgoingTxns", wireType)
			}
			m.MaxIncomingAndOutgoingTxns = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCosmos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxIncomingAndOutgoingTxns |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCosmos(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCosmos
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
func skipCosmos(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCosmos
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
					return 0, ErrIntOverflowCosmos
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
					return 0, ErrIntOverflowCosmos
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
				return 0, ErrInvalidLengthCosmos
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCosmos
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCosmos
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCosmos        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCosmos          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCosmos = fmt.Errorf("proto: unexpected end of group")
)
