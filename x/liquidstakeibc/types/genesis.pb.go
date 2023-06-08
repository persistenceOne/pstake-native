// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pstake/liquidstakeibc/v1beta1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// GenesisState defines the liquidstakeibc module's genesis state.
type GenesisState struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// initial host chain list
	HostChains []*HostChain `protobuf:"bytes,2,rep,name=host_chains,json=hostChains,proto3" json:"host_chains,omitempty"`
	// initial deposit list
	Deposits []*Deposit `protobuf:"bytes,3,rep,name=deposits,proto3" json:"deposits,omitempty"`
	// initial unbondings
	Unbondings []*Unbonding `protobuf:"bytes,4,rep,name=unbondings,proto3" json:"unbondings,omitempty"`
	// initial user unbondings
	UserUnbondings []*UserUnbonding `protobuf:"bytes,5,rep,name=user_unbondings,json=userUnbondings,proto3" json:"user_unbondings,omitempty"`
	// validator unbondings
	ValidatorUnbondings []*ValidatorUnbonding `protobuf:"bytes,6,rep,name=validator_unbondings,json=validatorUnbondings,proto3" json:"validator_unbondings,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d650226665335af, []int{0}
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

func (m *GenesisState) GetHostChains() []*HostChain {
	if m != nil {
		return m.HostChains
	}
	return nil
}

func (m *GenesisState) GetDeposits() []*Deposit {
	if m != nil {
		return m.Deposits
	}
	return nil
}

func (m *GenesisState) GetUnbondings() []*Unbonding {
	if m != nil {
		return m.Unbondings
	}
	return nil
}

func (m *GenesisState) GetUserUnbondings() []*UserUnbonding {
	if m != nil {
		return m.UserUnbondings
	}
	return nil
}

func (m *GenesisState) GetValidatorUnbondings() []*ValidatorUnbonding {
	if m != nil {
		return m.ValidatorUnbondings
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "pstake.liquidstakeibc.v1beta1.GenesisState")
}

func init() {
	proto.RegisterFile("pstake/liquidstakeibc/v1beta1/genesis.proto", fileDescriptor_1d650226665335af)
}

var fileDescriptor_1d650226665335af = []byte{
	// 373 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4d, 0x4b, 0xe3, 0x40,
	0x18, 0xc7, 0x93, 0x6d, 0xb7, 0x2c, 0xd3, 0x65, 0x17, 0xb2, 0x3d, 0x84, 0xc2, 0xc6, 0x22, 0x28,
	0xc5, 0x97, 0x0c, 0x8d, 0x9f, 0xc0, 0x56, 0xb0, 0x9e, 0x94, 0x48, 0x3d, 0xe8, 0xa1, 0x4c, 0x92,
	0x21, 0x19, 0x6c, 0x67, 0x62, 0x9e, 0x49, 0xd0, 0x6f, 0xe1, 0xc7, 0xea, 0xb1, 0x47, 0x4f, 0x22,
	0xed, 0xc9, 0x6f, 0x21, 0x9d, 0xb4, 0xf6, 0x45, 0x68, 0xbc, 0x3d, 0x33, 0xfc, 0x7f, 0xbf, 0xff,
	0x73, 0x78, 0xd0, 0x61, 0x0c, 0x92, 0xdc, 0x53, 0x3c, 0x60, 0x0f, 0x29, 0x0b, 0xd4, 0xcc, 0x3c,
	0x1f, 0x67, 0x2d, 0x8f, 0x4a, 0xd2, 0xc2, 0x21, 0xe5, 0x14, 0x18, 0xd8, 0x71, 0x22, 0xa4, 0x30,
	0xfe, 0xe7, 0x61, 0x7b, 0x3d, 0x6c, 0xcf, 0xc3, 0xf5, 0x5a, 0x28, 0x42, 0xa1, 0x92, 0x78, 0x36,
	0xe5, 0x50, 0xfd, 0x60, 0x7b, 0x43, 0x4c, 0x12, 0x32, 0x9c, 0x17, 0xd4, 0x9d, 0xed, 0xd9, 0x8d,
	0x5e, 0xc5, 0xec, 0xbe, 0x97, 0xd0, 0xef, 0xf3, 0x7c, 0xcd, 0x6b, 0x49, 0x24, 0x35, 0x3a, 0xa8,
	0x92, 0x4b, 0x4d, 0xbd, 0xa1, 0x37, 0xab, 0xce, 0x9e, 0xbd, 0x75, 0x6d, 0xfb, 0x4a, 0x85, 0xdb,
	0xe5, 0xd1, 0xeb, 0x8e, 0xe6, 0xce, 0x51, 0xe3, 0x02, 0x55, 0x23, 0x01, 0xb2, 0xef, 0x47, 0x84,
	0x71, 0x30, 0x7f, 0x34, 0x4a, 0xcd, 0xaa, 0xd3, 0x2c, 0x30, 0x75, 0x05, 0xc8, 0xce, 0x0c, 0x70,
	0x51, 0xb4, 0x18, 0xc1, 0x68, 0xa3, 0x5f, 0x01, 0x8d, 0x05, 0x30, 0x09, 0x66, 0x49, 0x79, 0xf6,
	0x0b, 0x3c, 0x67, 0x79, 0xdc, 0xfd, 0xe4, 0x8c, 0x2e, 0x42, 0x29, 0xf7, 0x04, 0x0f, 0x18, 0x0f,
	0xc1, 0x2c, 0x7f, 0x6b, 0x9b, 0xde, 0x02, 0x70, 0x57, 0x58, 0xa3, 0x87, 0xfe, 0xa6, 0x40, 0x93,
	0xfe, 0x8a, 0xee, 0xa7, 0xd2, 0x1d, 0x15, 0xe9, 0x80, 0x26, 0x4b, 0xe5, 0x9f, 0x74, 0xf5, 0x09,
	0x46, 0x80, 0x6a, 0x19, 0x19, 0xb0, 0x80, 0x48, 0xb1, 0xe6, 0xae, 0x28, 0x77, 0xab, 0xc0, 0x7d,
	0xb3, 0x40, 0x97, 0x05, 0xff, 0xb2, 0x2f, 0x7f, 0xd0, 0xbe, 0x1b, 0x4d, 0x2c, 0x7d, 0x3c, 0xb1,
	0xf4, 0xb7, 0x89, 0xa5, 0x3f, 0x4f, 0x2d, 0x6d, 0x3c, 0xb5, 0xb4, 0x97, 0xa9, 0xa5, 0xdd, 0x9e,
	0x86, 0x4c, 0x46, 0xa9, 0x67, 0xfb, 0x62, 0x88, 0x63, 0x9a, 0x00, 0x03, 0x49, 0xb9, 0x4f, 0x2f,
	0x39, 0xc5, 0x79, 0xf5, 0x31, 0x27, 0x92, 0x65, 0x14, 0x67, 0x0e, 0x7e, 0xdc, 0xbc, 0x2f, 0xf9,
	0x14, 0x53, 0xf0, 0x2a, 0xea, 0x9e, 0x4e, 0x3e, 0x02, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x9d, 0xaf,
	0x0f, 0x13, 0x03, 0x00, 0x00,
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
	if len(m.ValidatorUnbondings) > 0 {
		for iNdEx := len(m.ValidatorUnbondings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ValidatorUnbondings[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.UserUnbondings) > 0 {
		for iNdEx := len(m.UserUnbondings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.UserUnbondings[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Unbondings) > 0 {
		for iNdEx := len(m.Unbondings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Unbondings[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Deposits) > 0 {
		for iNdEx := len(m.Deposits) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Deposits[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.HostChains) > 0 {
		for iNdEx := len(m.HostChains) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.HostChains[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.HostChains) > 0 {
		for _, e := range m.HostChains {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Deposits) > 0 {
		for _, e := range m.Deposits {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Unbondings) > 0 {
		for _, e := range m.Unbondings {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.UserUnbondings) > 0 {
		for _, e := range m.UserUnbondings {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.ValidatorUnbondings) > 0 {
		for _, e := range m.ValidatorUnbondings {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
				return fmt.Errorf("proto: wrong wireType = %d for field HostChains", wireType)
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
			m.HostChains = append(m.HostChains, &HostChain{})
			if err := m.HostChains[len(m.HostChains)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deposits", wireType)
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
			m.Deposits = append(m.Deposits, &Deposit{})
			if err := m.Deposits[len(m.Deposits)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Unbondings", wireType)
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
			m.Unbondings = append(m.Unbondings, &Unbonding{})
			if err := m.Unbondings[len(m.Unbondings)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserUnbondings", wireType)
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
			m.UserUnbondings = append(m.UserUnbondings, &UserUnbonding{})
			if err := m.UserUnbondings[len(m.UserUnbondings)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorUnbondings", wireType)
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
			m.ValidatorUnbondings = append(m.ValidatorUnbondings, &ValidatorUnbonding{})
			if err := m.ValidatorUnbondings[len(m.ValidatorUnbondings)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
