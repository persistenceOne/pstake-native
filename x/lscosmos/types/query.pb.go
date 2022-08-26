// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lscosmos/v1beta1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// QueryParamsRequest is request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_86447f92490bdee4, []int{0}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

// QueryParamsResponse is response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params holds all the parameters of this module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_86447f92490bdee4, []int{1}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

// QueryCosmosIBCParamsRequest is request for the Ouery/CosmosIBCParams methods.
type QueryCosmosIBCParamsRequest struct {
}

func (m *QueryCosmosIBCParamsRequest) Reset()         { *m = QueryCosmosIBCParamsRequest{} }
func (m *QueryCosmosIBCParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryCosmosIBCParamsRequest) ProtoMessage()    {}
func (*QueryCosmosIBCParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_86447f92490bdee4, []int{2}
}
func (m *QueryCosmosIBCParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCosmosIBCParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCosmosIBCParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCosmosIBCParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCosmosIBCParamsRequest.Merge(m, src)
}
func (m *QueryCosmosIBCParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryCosmosIBCParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCosmosIBCParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCosmosIBCParamsRequest proto.InternalMessageInfo

// QueryCosmosIBCParamsResponse is response for the Ouery/CosmosIBCParams methods.
type QueryCosmosIBCParamsResponse struct {
	CosmosIBCParams CosmosIBCParams `protobuf:"bytes,1,opt,name=cosmos_i_b_c_params,json=cosmosIBCParams,proto3" json:"cosmos_i_b_c_params"`
}

func (m *QueryCosmosIBCParamsResponse) Reset()         { *m = QueryCosmosIBCParamsResponse{} }
func (m *QueryCosmosIBCParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryCosmosIBCParamsResponse) ProtoMessage()    {}
func (*QueryCosmosIBCParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_86447f92490bdee4, []int{3}
}
func (m *QueryCosmosIBCParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCosmosIBCParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCosmosIBCParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCosmosIBCParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCosmosIBCParamsResponse.Merge(m, src)
}
func (m *QueryCosmosIBCParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryCosmosIBCParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCosmosIBCParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCosmosIBCParamsResponse proto.InternalMessageInfo

func (m *QueryCosmosIBCParamsResponse) GetCosmosIBCParams() CosmosIBCParams {
	if m != nil {
		return m.CosmosIBCParams
	}
	return CosmosIBCParams{}
}

// QueryDelegationStateRequest is request for the Ouery/DelegationState methods.
type QueryDelegationStateRequest struct {
}

func (m *QueryDelegationStateRequest) Reset()         { *m = QueryDelegationStateRequest{} }
func (m *QueryDelegationStateRequest) String() string { return proto.CompactTextString(m) }
func (*QueryDelegationStateRequest) ProtoMessage()    {}
func (*QueryDelegationStateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_86447f92490bdee4, []int{4}
}
func (m *QueryDelegationStateRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDelegationStateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDelegationStateRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDelegationStateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDelegationStateRequest.Merge(m, src)
}
func (m *QueryDelegationStateRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryDelegationStateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDelegationStateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDelegationStateRequest proto.InternalMessageInfo

// QueryDelegationStateResponse is response for the Ouery/DelegationState methods.
type QueryDelegationStateResponse struct {
	DelegationState DelegationState `protobuf:"bytes,1,opt,name=delegation_state,json=delegationState,proto3" json:"delegation_state"`
}

func (m *QueryDelegationStateResponse) Reset()         { *m = QueryDelegationStateResponse{} }
func (m *QueryDelegationStateResponse) String() string { return proto.CompactTextString(m) }
func (*QueryDelegationStateResponse) ProtoMessage()    {}
func (*QueryDelegationStateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_86447f92490bdee4, []int{5}
}
func (m *QueryDelegationStateResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDelegationStateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDelegationStateResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDelegationStateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDelegationStateResponse.Merge(m, src)
}
func (m *QueryDelegationStateResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryDelegationStateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDelegationStateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDelegationStateResponse proto.InternalMessageInfo

func (m *QueryDelegationStateResponse) GetDelegationState() DelegationState {
	if m != nil {
		return m.DelegationState
	}
	return DelegationState{}
}

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "lscosmos.v1beta1.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "lscosmos.v1beta1.QueryParamsResponse")
	proto.RegisterType((*QueryCosmosIBCParamsRequest)(nil), "lscosmos.v1beta1.QueryCosmosIBCParamsRequest")
	proto.RegisterType((*QueryCosmosIBCParamsResponse)(nil), "lscosmos.v1beta1.QueryCosmosIBCParamsResponse")
	proto.RegisterType((*QueryDelegationStateRequest)(nil), "lscosmos.v1beta1.QueryDelegationStateRequest")
	proto.RegisterType((*QueryDelegationStateResponse)(nil), "lscosmos.v1beta1.QueryDelegationStateResponse")
}

func init() { proto.RegisterFile("lscosmos/v1beta1/query.proto", fileDescriptor_86447f92490bdee4) }

var fileDescriptor_86447f92490bdee4 = []byte{
	// 486 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xbf, 0x8f, 0xd3, 0x30,
	0x1c, 0xc5, 0x6b, 0x7e, 0x74, 0x30, 0x43, 0x4f, 0xbe, 0x1b, 0x4e, 0xa5, 0x97, 0x83, 0x08, 0x24,
	0x84, 0xd4, 0x58, 0x77, 0x08, 0x38, 0x58, 0x40, 0x3d, 0x16, 0x06, 0x04, 0x14, 0xb1, 0xb0, 0x44,
	0x4e, 0xee, 0x2b, 0x13, 0xd1, 0xda, 0xbe, 0xd8, 0xad, 0xb8, 0x95, 0x99, 0x01, 0x89, 0xbf, 0x04,
	0xfe, 0x8a, 0x1b, 0x2b, 0xb1, 0x30, 0x21, 0xd4, 0xf2, 0x7f, 0x80, 0xe2, 0x38, 0x41, 0x49, 0x1a,
	0xe8, 0x56, 0xf9, 0xbd, 0xbe, 0xf7, 0xa9, 0xfd, 0x8a, 0x07, 0x13, 0x1d, 0x4b, 0x3d, 0x95, 0x9a,
	0xce, 0x0f, 0x22, 0x30, 0xec, 0x80, 0x9e, 0xce, 0x20, 0x3d, 0x0b, 0x54, 0x2a, 0x8d, 0x24, 0x5b,
	0x85, 0x1a, 0x38, 0xb5, 0xbf, 0xc3, 0x25, 0x97, 0x56, 0xa4, 0xd9, 0xa7, 0xdc, 0xd7, 0x1f, 0x70,
	0x29, 0xf9, 0x04, 0x28, 0x53, 0x09, 0x65, 0x42, 0x48, 0xc3, 0x4c, 0x22, 0x85, 0x76, 0xea, 0x6d,
	0xd7, 0x10, 0x31, 0x0d, 0x79, 0x7c, 0x59, 0xa6, 0x18, 0x4f, 0x84, 0x35, 0x3b, 0xef, 0x5e, 0x83,
	0x47, 0xb1, 0x94, 0x4d, 0xcb, 0xa8, 0x86, 0xcc, 0xe5, 0x1c, 0x52, 0xc1, 0x44, 0x0c, 0xa1, 0x4a,
	0xa5, 0x92, 0x9a, 0x4d, 0x9c, 0x77, 0xbf, 0xe1, 0x2d, 0x7f, 0x8d, 0x35, 0xf8, 0x3b, 0x98, 0xbc,
	0xcc, 0x68, 0x5e, 0xd8, 0x86, 0x31, 0x9c, 0xce, 0x40, 0x1b, 0xff, 0x19, 0xde, 0xae, 0x9c, 0x6a,
	0x25, 0x85, 0x06, 0x72, 0x0f, 0x77, 0x73, 0x92, 0x5d, 0x74, 0x0d, 0xdd, 0xba, 0x72, 0xb8, 0x1b,
	0xd4, 0xef, 0x26, 0xc8, 0xbf, 0x31, 0xba, 0x74, 0xfe, 0x63, 0xbf, 0x33, 0x76, 0x6e, 0x7f, 0x0f,
	0x5f, 0xb5, 0x71, 0xc7, 0xd6, 0xfb, 0x74, 0x74, 0x5c, 0x6d, 0x9b, 0xe1, 0xc1, 0x7a, 0xd9, 0xd5,
	0xbe, 0xc6, 0xdb, 0x79, 0x4b, 0x98, 0x84, 0x51, 0x18, 0x87, 0x15, 0x86, 0xeb, 0x4d, 0x86, 0x5a,
	0x8e, 0x83, 0xe9, 0xc5, 0xd5, 0xe3, 0x92, 0xea, 0x09, 0x4c, 0x80, 0xdb, 0xfb, 0x7f, 0x65, 0x98,
	0x81, 0x82, 0x2a, 0x75, 0x54, 0x0d, 0xd9, 0x51, 0x8d, 0xf1, 0xd6, 0x49, 0x29, 0x85, 0x3a, 0xd3,
	0xda, 0x91, 0x6a, 0x21, 0x05, 0xd2, 0x49, 0xf5, 0xf8, 0xf0, 0xf7, 0x45, 0x7c, 0xd9, 0x96, 0x92,
	0x8f, 0x08, 0x77, 0x73, 0x4e, 0x72, 0xa3, 0x19, 0xd7, 0x7c, 0xb2, 0xfe, 0xcd, 0xff, 0xb8, 0x72,
	0x6a, 0xff, 0xee, 0x87, 0x6f, 0xbf, 0x3e, 0x5f, 0xa0, 0x64, 0x48, 0x15, 0xa4, 0x3a, 0xd1, 0x06,
	0x44, 0x0c, 0xcf, 0x05, 0x50, 0xa5, 0x0d, 0x7b, 0x07, 0xc3, 0x6c, 0x88, 0x73, 0x28, 0x57, 0xe2,
	0x96, 0x47, 0xbe, 0x22, 0xdc, 0xab, 0x5d, 0x2b, 0x19, 0xb6, 0x34, 0xae, 0x7f, 0xe5, 0x7e, 0xb0,
	0xa9, 0xdd, 0x91, 0x3e, 0xb6, 0xa4, 0x0f, 0xc9, 0xd1, 0x86, 0xa4, 0xc5, 0x44, 0xa2, 0x62, 0x20,
	0xe4, 0x0b, 0xc2, 0xbd, 0xda, 0xc5, 0xb7, 0x42, 0xaf, 0x1f, 0x41, 0x2b, 0x74, 0xcb, 0x28, 0xfc,
	0x47, 0x16, 0xfa, 0x01, 0xb9, 0xbf, 0x21, 0x74, 0x7d, 0x41, 0xa3, 0xf1, 0xf9, 0xd2, 0x43, 0x8b,
	0xa5, 0x87, 0x7e, 0x2e, 0x3d, 0xf4, 0x69, 0xe5, 0x75, 0x16, 0x2b, 0xaf, 0xf3, 0x7d, 0xe5, 0x75,
	0xde, 0x1c, 0xf1, 0xc4, 0xbc, 0x9d, 0x45, 0x41, 0x2c, 0xa7, 0xff, 0x0e, 0x7f, 0xff, 0x37, 0xde,
	0x9c, 0x29, 0xd0, 0x51, 0xd7, 0xfe, 0xd5, 0xef, 0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0xc2, 0x01,
	0xf5, 0x34, 0xe8, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	CosmosIBCParams(ctx context.Context, in *QueryCosmosIBCParamsRequest, opts ...grpc.CallOption) (*QueryCosmosIBCParamsResponse, error)
	DelegationState(ctx context.Context, in *QueryDelegationStateRequest, opts ...grpc.CallOption) (*QueryDelegationStateResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/lscosmos.v1beta1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CosmosIBCParams(ctx context.Context, in *QueryCosmosIBCParamsRequest, opts ...grpc.CallOption) (*QueryCosmosIBCParamsResponse, error) {
	out := new(QueryCosmosIBCParamsResponse)
	err := c.cc.Invoke(ctx, "/lscosmos.v1beta1.Query/CosmosIBCParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DelegationState(ctx context.Context, in *QueryDelegationStateRequest, opts ...grpc.CallOption) (*QueryDelegationStateResponse, error) {
	out := new(QueryDelegationStateResponse)
	err := c.cc.Invoke(ctx, "/lscosmos.v1beta1.Query/DelegationState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	CosmosIBCParams(context.Context, *QueryCosmosIBCParamsRequest) (*QueryCosmosIBCParamsResponse, error)
	DelegationState(context.Context, *QueryDelegationStateRequest) (*QueryDelegationStateResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) CosmosIBCParams(ctx context.Context, req *QueryCosmosIBCParamsRequest) (*QueryCosmosIBCParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CosmosIBCParams not implemented")
}
func (*UnimplementedQueryServer) DelegationState(ctx context.Context, req *QueryDelegationStateRequest) (*QueryDelegationStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelegationState not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lscosmos.v1beta1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CosmosIBCParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCosmosIBCParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CosmosIBCParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lscosmos.v1beta1.Query/CosmosIBCParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CosmosIBCParams(ctx, req.(*QueryCosmosIBCParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DelegationState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDelegationStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DelegationState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lscosmos.v1beta1.Query/DelegationState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DelegationState(ctx, req.(*QueryDelegationStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "lscosmos.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "CosmosIBCParams",
			Handler:    _Query_CosmosIBCParams_Handler,
		},
		{
			MethodName: "DelegationState",
			Handler:    _Query_DelegationState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lscosmos/v1beta1/query.proto",
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryCosmosIBCParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCosmosIBCParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCosmosIBCParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryCosmosIBCParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCosmosIBCParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCosmosIBCParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.CosmosIBCParams.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryDelegationStateRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDelegationStateRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDelegationStateRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryDelegationStateResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDelegationStateResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDelegationStateResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.DelegationState.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryCosmosIBCParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryCosmosIBCParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.CosmosIBCParams.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryDelegationStateRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryDelegationStateResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.DelegationState.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryCosmosIBCParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryCosmosIBCParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCosmosIBCParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryCosmosIBCParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryCosmosIBCParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCosmosIBCParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CosmosIBCParams", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CosmosIBCParams.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryDelegationStateRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryDelegationStateRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDelegationStateRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryDelegationStateResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryDelegationStateResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDelegationStateResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegationState", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DelegationState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
