// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lscosmos/v1beta1/governance_proposal.proto

package types

import (
	fmt "fmt"
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

// RegisterCosmosChainProposal defines the details needed to register cosmos chain for
// liquid staking transactions through lscosmos module
type RegisterCosmosChainProposal struct {
	Title                 string                                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description           string                                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	ModuleEnabled         bool                                   `protobuf:"varint,3,opt,name=module_enabled,json=moduleEnabled,proto3" json:"module_enabled,omitempty"`
	ConnectionID          string                                 `protobuf:"bytes,4,opt,name=connectionID,proto3" json:"connectionID,omitempty"`
	TransferChannel       string                                 `protobuf:"bytes,5,opt,name=transfer_channel,json=transferChannel,proto3" json:"transfer_channel,omitempty"`
	TransferPort          string                                 `protobuf:"bytes,6,opt,name=transfer_port,json=transferPort,proto3" json:"transfer_port,omitempty"`
	BaseDenom             string                                 `protobuf:"bytes,7,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
	MintDenom             string                                 `protobuf:"bytes,8,opt,name=mint_denom,json=mintDenom,proto3" json:"mint_denom,omitempty"`
	MinDeposit            github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,9,opt,name=min_deposit,json=minDeposit,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_deposit"`
	AllowListedValidators AllowListedValidators                  `protobuf:"bytes,10,opt,name=allow_listed_validators,json=allowListedValidators,proto3" json:"allow_listed_validators"`
	PstakeDepositFee      github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,11,opt,name=pstake_deposit_fee,json=pstakeDepositFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_deposit_fee"`
	PstakeRestakeFee      github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,12,opt,name=pstake_restake_fee,json=pstakeRestakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_restake_fee"`
	PstakeUnstakeFee      github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,13,opt,name=pstake_unstake_fee,json=pstakeUnstakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pstake_unstake_fee"`
}

func (m *RegisterCosmosChainProposal) Reset()      { *m = RegisterCosmosChainProposal{} }
func (*RegisterCosmosChainProposal) ProtoMessage() {}
func (*RegisterCosmosChainProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_abbb79eadcf33bd7, []int{0}
}
func (m *RegisterCosmosChainProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RegisterCosmosChainProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RegisterCosmosChainProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RegisterCosmosChainProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterCosmosChainProposal.Merge(m, src)
}
func (m *RegisterCosmosChainProposal) XXX_Size() int {
	return m.Size()
}
func (m *RegisterCosmosChainProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterCosmosChainProposal.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterCosmosChainProposal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*RegisterCosmosChainProposal)(nil), "lscosmos.v1beta1.RegisterCosmosChainProposal")
}

func init() {
	proto.RegisterFile("lscosmos/v1beta1/governance_proposal.proto", fileDescriptor_abbb79eadcf33bd7)
}

var fileDescriptor_abbb79eadcf33bd7 = []byte{
	// 542 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xc1, 0x6f, 0xd3, 0x30,
	0x14, 0xc6, 0x13, 0xd8, 0x4a, 0xeb, 0xae, 0x50, 0x59, 0x43, 0x44, 0x43, 0xa4, 0xd5, 0x10, 0x50,
	0x90, 0x96, 0x68, 0x70, 0x41, 0xdc, 0x68, 0x0b, 0xd2, 0x24, 0xa4, 0x4d, 0x91, 0xe0, 0x80, 0x10,
	0x91, 0x9b, 0xbc, 0xb5, 0xd6, 0x1c, 0x3b, 0xb2, 0xdd, 0x02, 0x37, 0xce, 0x9c, 0x38, 0x72, 0xdc,
	0x91, 0x3f, 0xa5, 0xc7, 0x1d, 0x11, 0x87, 0x09, 0xda, 0x0b, 0x7f, 0x06, 0x72, 0x9c, 0x96, 0x0e,
	0xb8, 0x4c, 0x3b, 0x39, 0xf9, 0xbe, 0xcf, 0xbf, 0xf7, 0xf2, 0x62, 0xa3, 0x07, 0x4c, 0x25, 0x42,
	0x65, 0x42, 0x85, 0x93, 0xdd, 0x01, 0x68, 0xb2, 0x1b, 0x0e, 0xc5, 0x04, 0x24, 0x27, 0x3c, 0x81,
	0x38, 0x97, 0x22, 0x17, 0x8a, 0xb0, 0x20, 0x97, 0x42, 0x0b, 0xdc, 0x5c, 0x64, 0x83, 0x32, 0xbb,
	0xb5, 0x39, 0x14, 0x43, 0x51, 0x98, 0xa1, 0x79, 0xb2, 0xb9, 0x2d, 0xbf, 0x24, 0x0e, 0x88, 0x82,
	0x25, 0x36, 0x11, 0x94, 0x97, 0x7e, 0xeb, 0x9f, 0x9a, 0x4b, 0x70, 0x11, 0xd8, 0xfe, 0x54, 0x41,
	0x37, 0x23, 0x18, 0x52, 0xa5, 0x41, 0xf6, 0x0a, 0xa3, 0x37, 0x22, 0x94, 0x1f, 0x94, 0xed, 0xe0,
	0x4d, 0xb4, 0xae, 0xa9, 0x66, 0xe0, 0xb9, 0x6d, 0xb7, 0x53, 0x8b, 0xec, 0x0b, 0x6e, 0xa3, 0x7a,
	0x0a, 0x2a, 0x91, 0x34, 0xd7, 0x54, 0x70, 0xef, 0x52, 0xe1, 0xad, 0x4a, 0xf8, 0x0e, 0xba, 0x9a,
	0x89, 0x74, 0xcc, 0x20, 0x06, 0x4e, 0x06, 0x0c, 0x52, 0xef, 0x72, 0xdb, 0xed, 0x54, 0xa3, 0x86,
	0x55, 0x9f, 0x59, 0x11, 0x6f, 0xa3, 0x8d, 0x44, 0x70, 0x0e, 0x89, 0xd9, 0xb4, 0xd7, 0xf7, 0xd6,
	0x0a, 0xd2, 0x19, 0x0d, 0xdf, 0x47, 0x4d, 0x2d, 0x09, 0x57, 0x87, 0x20, 0xe3, 0x64, 0x44, 0x38,
	0x07, 0xe6, 0xad, 0x17, 0xb9, 0x6b, 0x0b, 0xbd, 0x67, 0x65, 0x7c, 0x1b, 0x35, 0x96, 0xd1, 0x5c,
	0x48, 0xed, 0x55, 0x2c, 0x6f, 0x21, 0x1e, 0x08, 0xa9, 0xf1, 0x2d, 0x84, 0xcc, 0xb8, 0xe2, 0x14,
	0xb8, 0xc8, 0xbc, 0x2b, 0x45, 0xa2, 0x66, 0x94, 0xbe, 0x11, 0x8c, 0x9d, 0x51, 0xae, 0x4b, 0xbb,
	0x6a, 0x6d, 0xa3, 0x58, 0x7b, 0x1f, 0xd5, 0x33, 0xca, 0xe3, 0x14, 0x72, 0xa1, 0xa8, 0xf6, 0x6a,
	0xc6, 0xef, 0x06, 0xd3, 0xd3, 0x96, 0xf3, 0xfd, 0xb4, 0x75, 0x77, 0x48, 0xf5, 0x68, 0x3c, 0x08,
	0x12, 0x91, 0x85, 0xe5, 0xdc, 0xed, 0xb2, 0xa3, 0xd2, 0xa3, 0x50, 0x7f, 0xc8, 0x41, 0x05, 0x7b,
	0x5c, 0x47, 0xa6, 0x42, 0xdf, 0x12, 0x30, 0xa0, 0x1b, 0x84, 0x31, 0xf1, 0x2e, 0x66, 0xe6, 0x27,
	0xa4, 0xf1, 0x84, 0x30, 0x9a, 0x12, 0x2d, 0xa4, 0xf2, 0x50, 0xdb, 0xed, 0xd4, 0x1f, 0xde, 0x0b,
	0xfe, 0x3e, 0x0c, 0xc1, 0x53, 0xb3, 0xe1, 0x45, 0x91, 0x7f, 0xb5, 0x8c, 0x77, 0xd7, 0x4c, 0x17,
	0xd1, 0x75, 0xf2, 0x3f, 0x13, 0xbf, 0x41, 0x38, 0x57, 0x9a, 0x1c, 0xc1, 0xa2, 0xf5, 0xf8, 0x10,
	0xc0, 0xab, 0x9f, 0xbb, 0xfd, 0x3e, 0x24, 0x51, 0xd3, 0x92, 0xca, 0x2f, 0x78, 0x0e, 0xb0, 0x42,
	0x97, 0x60, 0x57, 0x43, 0xdf, 0xb8, 0x08, 0x3d, 0xb2, 0xa0, 0xb3, 0xf4, 0x31, 0xff, 0x43, 0x6f,
	0x5c, 0x84, 0xfe, 0x92, 0x2f, 0xe8, 0x4f, 0xaa, 0x5f, 0x8e, 0x5b, 0xce, 0xaf, 0xe3, 0x96, 0xd3,
	0x7d, 0x3b, 0xfd, 0xe9, 0x3b, 0x1f, 0x67, 0xbe, 0xf3, 0x75, 0xe6, 0xbb, 0xd3, 0x99, 0xef, 0x9e,
	0xcc, 0x7c, 0xf7, 0xc7, 0xcc, 0x77, 0x3f, 0xcf, 0x7d, 0xe7, 0x64, 0xee, 0x3b, 0xdf, 0xe6, 0xbe,
	0xf3, 0xfa, 0xf1, 0x4a, 0xa5, 0x1c, 0xa4, 0x32, 0x93, 0xe6, 0x09, 0xec, 0x73, 0x08, 0x2d, 0x7b,
	0x87, 0x13, 0x4d, 0x27, 0x10, 0xbe, 0x5f, 0x5e, 0x36, 0x5b, 0x7f, 0x50, 0x29, 0xee, 0xdc, 0xa3,
	0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc8, 0x3f, 0x77, 0xa2, 0x0a, 0x04, 0x00, 0x00,
}

func (m *RegisterCosmosChainProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RegisterCosmosChainProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RegisterCosmosChainProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x6a
	{
		size := m.PstakeRestakeFee.Size()
		i -= size
		if _, err := m.PstakeRestakeFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	{
		size := m.PstakeDepositFee.Size()
		i -= size
		if _, err := m.PstakeDepositFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x5a
	{
		size, err := m.AllowListedValidators.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	{
		size := m.MinDeposit.Size()
		i -= size
		if _, err := m.MinDeposit.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	if len(m.MintDenom) > 0 {
		i -= len(m.MintDenom)
		copy(dAtA[i:], m.MintDenom)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.MintDenom)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.BaseDenom) > 0 {
		i -= len(m.BaseDenom)
		copy(dAtA[i:], m.BaseDenom)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.BaseDenom)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.TransferPort) > 0 {
		i -= len(m.TransferPort)
		copy(dAtA[i:], m.TransferPort)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.TransferPort)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.TransferChannel) > 0 {
		i -= len(m.TransferChannel)
		copy(dAtA[i:], m.TransferChannel)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.TransferChannel)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.ConnectionID) > 0 {
		i -= len(m.ConnectionID)
		copy(dAtA[i:], m.ConnectionID)
		i = encodeVarintGovernanceProposal(dAtA, i, uint64(len(m.ConnectionID)))
		i--
		dAtA[i] = 0x22
	}
	if m.ModuleEnabled {
		i--
		if m.ModuleEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
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
func (m *RegisterCosmosChainProposal) Size() (n int) {
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
	if m.ModuleEnabled {
		n += 2
	}
	l = len(m.ConnectionID)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = len(m.TransferChannel)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = len(m.TransferPort)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = len(m.BaseDenom)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = len(m.MintDenom)
	if l > 0 {
		n += 1 + l + sovGovernanceProposal(uint64(l))
	}
	l = m.MinDeposit.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	l = m.AllowListedValidators.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	l = m.PstakeDepositFee.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	l = m.PstakeRestakeFee.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	l = m.PstakeUnstakeFee.Size()
	n += 1 + l + sovGovernanceProposal(uint64(l))
	return n
}

func sovGovernanceProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGovernanceProposal(x uint64) (n int) {
	return sovGovernanceProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RegisterCosmosChainProposal) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: RegisterCosmosChainProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RegisterCosmosChainProposal: illegal tag %d (wire type %d)", fieldNum, wire)
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ModuleEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernanceProposal
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
			m.ModuleEnabled = bool(v != 0)
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionID", wireType)
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
			m.ConnectionID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferChannel", wireType)
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
			m.TransferChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferPort", wireType)
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
			m.TransferPort = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseDenom", wireType)
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
			m.BaseDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MintDenom", wireType)
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
			m.MintDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
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
		case 10:
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
		case 11:
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
		case 12:
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
		case 13:
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
