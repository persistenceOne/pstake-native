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

func init() {
	proto.RegisterType((*GenesisState)(nil), "pstake.liquidstakeibc.v1beta1.GenesisState")
}

func init() {
	proto.RegisterFile("pstake/liquidstakeibc/v1beta1/genesis.proto", fileDescriptor_1d650226665335af)
}

var fileDescriptor_1d650226665335af = []byte{
	// 343 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0x40, 0x93, 0xaf, 0xfd, 0x8a, 0x4c, 0x45, 0x21, 0xb8, 0x08, 0x05, 0x63, 0x11, 0x94, 0xe2,
	0x4f, 0x86, 0xc6, 0x27, 0xb0, 0x15, 0xac, 0x2b, 0x25, 0xd2, 0x8d, 0x2e, 0xca, 0x24, 0xbd, 0x24,
	0x83, 0x76, 0x26, 0xe6, 0x4e, 0x8a, 0xae, 0x7d, 0x01, 0x1f, 0xab, 0xcb, 0x2e, 0x5d, 0x89, 0xb4,
	0x2f, 0x22, 0x9d, 0xb4, 0x5a, 0xbb, 0x68, 0xdc, 0xdd, 0x81, 0x73, 0xce, 0xbd, 0x30, 0xe4, 0x38,
	0x41, 0xc5, 0x1e, 0x80, 0x3e, 0xf2, 0xa7, 0x8c, 0xf7, 0xf5, 0xcc, 0x83, 0x90, 0x0e, 0x9b, 0x01,
	0x28, 0xd6, 0xa4, 0x11, 0x08, 0x40, 0x8e, 0x6e, 0x92, 0x4a, 0x25, 0xad, 0xdd, 0x1c, 0x76, 0x7f,
	0xc3, 0xee, 0x1c, 0xae, 0xed, 0x44, 0x32, 0x92, 0x9a, 0xa4, 0xb3, 0x29, 0x97, 0x6a, 0x47, 0xeb,
	0x37, 0x24, 0x2c, 0x65, 0x83, 0xf9, 0x82, 0x9a, 0xb7, 0x9e, 0x5d, 0xd9, 0xab, 0x9d, 0xfd, 0xd7,
	0x12, 0xd9, 0xbc, 0xcc, 0xcf, 0xbc, 0x55, 0x4c, 0x81, 0xd5, 0x26, 0x95, 0x3c, 0x6a, 0x9b, 0x75,
	0xb3, 0x51, 0xf5, 0x0e, 0xdc, 0xb5, 0x67, 0xbb, 0x37, 0x1a, 0x6e, 0x95, 0x47, 0x1f, 0x7b, 0x86,
	0x3f, 0x57, 0xad, 0x2b, 0x52, 0x8d, 0x25, 0xaa, 0x5e, 0x18, 0x33, 0x2e, 0xd0, 0xfe, 0x57, 0x2f,
	0x35, 0xaa, 0x5e, 0xa3, 0xa0, 0xd4, 0x91, 0xa8, 0xda, 0x33, 0xc1, 0x27, 0xf1, 0x62, 0x44, 0xab,
	0x45, 0x36, 0xfa, 0x90, 0x48, 0xe4, 0x0a, 0xed, 0x92, 0xee, 0x1c, 0x16, 0x74, 0x2e, 0x72, 0xdc,
	0xff, 0xf6, 0xac, 0x0e, 0x21, 0x99, 0x08, 0xa4, 0xe8, 0x73, 0x11, 0xa1, 0x5d, 0xfe, 0xd3, 0x35,
	0xdd, 0x85, 0xe0, 0x2f, 0xb9, 0x56, 0x97, 0x6c, 0x67, 0x08, 0x69, 0x6f, 0x29, 0xf7, 0x5f, 0xe7,
	0x4e, 0x8a, 0x72, 0x08, 0xe9, 0x4f, 0x72, 0x2b, 0x5b, 0x7e, 0x62, 0xeb, 0x7e, 0x34, 0x71, 0xcc,
	0xf1, 0xc4, 0x31, 0x3f, 0x27, 0x8e, 0xf9, 0x36, 0x75, 0x8c, 0xf1, 0xd4, 0x31, 0xde, 0xa7, 0x8e,
	0x71, 0x77, 0x1e, 0x71, 0x15, 0x67, 0x81, 0x1b, 0xca, 0x01, 0x4d, 0x20, 0x45, 0x8e, 0x0a, 0x44,
	0x08, 0xd7, 0x02, 0x68, 0xbe, 0xf0, 0x54, 0x30, 0xc5, 0x87, 0x40, 0x87, 0x1e, 0x7d, 0x5e, 0xfd,
	0x79, 0xf5, 0x92, 0x00, 0x06, 0x15, 0xfd, 0xd3, 0x67, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xb9,
	0x29, 0x90, 0xe1, 0xad, 0x02, 0x00, 0x00,
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
