// Code generated by protoc-gen-go. DO NOT EDIT.
// source: system.proto

package magpiepb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type RpcRequest struct {
	Params               map[string]string `protobuf:"bytes,1,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Data                 string            `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RpcRequest) Reset()         { *m = RpcRequest{} }
func (m *RpcRequest) String() string { return proto.CompactTextString(m) }
func (*RpcRequest) ProtoMessage()    {}
func (*RpcRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_86a7260ebdc12f47, []int{0}
}

func (m *RpcRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RpcRequest.Unmarshal(m, b)
}
func (m *RpcRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RpcRequest.Marshal(b, m, deterministic)
}
func (m *RpcRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RpcRequest.Merge(m, src)
}
func (m *RpcRequest) XXX_Size() int {
	return xxx_messageInfo_RpcRequest.Size(m)
}
func (m *RpcRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RpcRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RpcRequest proto.InternalMessageInfo

func (m *RpcRequest) GetParams() map[string]string {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *RpcRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type Entry struct {
	Index                uint64   `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Data                 string   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Team                 string   `protobuf:"bytes,3,opt,name=team,proto3" json:"team,omitempty"`
	Address              string   `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Port                 int64    `protobuf:"varint,5,opt,name=port,proto3" json:"port,omitempty"`
	Timestamp            string   `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}
func (*Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_86a7260ebdc12f47, []int{1}
}

func (m *Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entry.Unmarshal(m, b)
}
func (m *Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entry.Marshal(b, m, deterministic)
}
func (m *Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entry.Merge(m, src)
}
func (m *Entry) XXX_Size() int {
	return xxx_messageInfo_Entry.Size(m)
}
func (m *Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Entry proto.InternalMessageInfo

func (m *Entry) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Entry) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *Entry) GetTeam() string {
	if m != nil {
		return m.Team
	}
	return ""
}

func (m *Entry) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Entry) GetPort() int64 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Entry) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

type RpcResponse struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data                 string   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RpcResponse) Reset()         { *m = RpcResponse{} }
func (m *RpcResponse) String() string { return proto.CompactTextString(m) }
func (*RpcResponse) ProtoMessage()    {}
func (*RpcResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_86a7260ebdc12f47, []int{2}
}

func (m *RpcResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RpcResponse.Unmarshal(m, b)
}
func (m *RpcResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RpcResponse.Marshal(b, m, deterministic)
}
func (m *RpcResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RpcResponse.Merge(m, src)
}
func (m *RpcResponse) XXX_Size() int {
	return xxx_messageInfo_RpcResponse.Size(m)
}
func (m *RpcResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RpcResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RpcResponse proto.InternalMessageInfo

func (m *RpcResponse) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *RpcResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *RpcResponse) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

//stream响应结构
type StreamResponse struct {
	Data                 string   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamResponse) Reset()         { *m = StreamResponse{} }
func (m *StreamResponse) String() string { return proto.CompactTextString(m) }
func (*StreamResponse) ProtoMessage()    {}
func (*StreamResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_86a7260ebdc12f47, []int{3}
}

func (m *StreamResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamResponse.Unmarshal(m, b)
}
func (m *StreamResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamResponse.Marshal(b, m, deterministic)
}
func (m *StreamResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamResponse.Merge(m, src)
}
func (m *StreamResponse) XXX_Size() int {
	return xxx_messageInfo_StreamResponse.Size(m)
}
func (m *StreamResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StreamResponse proto.InternalMessageInfo

func (m *StreamResponse) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*RpcRequest)(nil), "RpcRequest")
	proto.RegisterMapType((map[string]string)(nil), "RpcRequest.ParamsEntry")
	proto.RegisterType((*Entry)(nil), "Entry")
	proto.RegisterType((*RpcResponse)(nil), "RpcResponse")
	proto.RegisterType((*StreamResponse)(nil), "StreamResponse")
}

func init() { proto.RegisterFile("system.proto", fileDescriptor_86a7260ebdc12f47) }

var fileDescriptor_86a7260ebdc12f47 = []byte{
	// 429 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xbb, 0x75, 0x9c, 0x90, 0x71, 0xf9, 0xa3, 0x55, 0x25, 0xac, 0x0a, 0x89, 0xc8, 0xaa,
	0x50, 0xa8, 0x2a, 0x17, 0x85, 0x0b, 0x70, 0x24, 0x41, 0x5c, 0x5a, 0x81, 0x1c, 0x4e, 0xdc, 0xc6,
	0xf6, 0x28, 0xb2, 0xf0, 0xda, 0xcb, 0xee, 0x04, 0xe1, 0x77, 0xe0, 0xc2, 0xd3, 0xf0, 0x7a, 0x68,
	0x37, 0x0d, 0x49, 0x45, 0x65, 0xe5, 0xf6, 0xcd, 0xec, 0x6f, 0xbe, 0xfd, 0x76, 0xed, 0x85, 0x13,
	0xdb, 0x59, 0x26, 0x95, 0x6a, 0xd3, 0x72, 0x9b, 0xfc, 0x12, 0x00, 0x99, 0x2e, 0x32, 0xfa, 0xbe,
	0x26, 0xcb, 0xf2, 0x0a, 0x86, 0x1a, 0x0d, 0x2a, 0x1b, 0x8b, 0x49, 0x30, 0x8d, 0x66, 0x4f, 0xd3,
	0xdd, 0x62, 0xfa, 0xd9, 0xaf, 0x7c, 0x68, 0xd8, 0x74, 0xd9, 0x2d, 0x26, 0x25, 0x0c, 0x4a, 0x64,
	0x8c, 0x8f, 0x27, 0x62, 0x3a, 0xce, 0xbc, 0x3e, 0x7b, 0x0b, 0xd1, 0x1e, 0x2a, 0x9f, 0x40, 0xf0,
	0x8d, 0xba, 0x58, 0x78, 0xc2, 0x49, 0x79, 0x0a, 0xe1, 0x0f, 0xac, 0xd7, 0x74, 0x3b, 0xb5, 0x29,
	0xde, 0x1d, 0xbf, 0x11, 0xc9, 0x6f, 0x01, 0xe1, 0x66, 0xea, 0x14, 0xc2, 0xaa, 0x29, 0xe9, 0xa7,
	0x9f, 0x1b, 0x64, 0x9b, 0xe2, 0xbe, 0xed, 0x5c, 0x8f, 0x09, 0x55, 0x1c, 0x6c, 0x7a, 0x4e, 0xcb,
	0x18, 0x46, 0x58, 0x96, 0x86, 0xac, 0x8d, 0x07, 0xbe, 0xbd, 0x2d, 0x1d, 0xad, 0x5b, 0xc3, 0x71,
	0x38, 0x11, 0xd3, 0x20, 0xf3, 0x5a, 0x3e, 0x83, 0x31, 0x57, 0x8a, 0x2c, 0xa3, 0xd2, 0xf1, 0xd0,
	0xf3, 0xbb, 0x46, 0xf2, 0x09, 0x22, 0x7f, 0x09, 0x56, 0xb7, 0x8d, 0x25, 0x67, 0x50, 0xb4, 0x25,
	0xf9, 0x5c, 0x61, 0xe6, 0xb5, 0xdb, 0x4e, 0x91, 0xb5, 0xb8, 0xda, 0x1e, 0x69, 0x5b, 0xfe, 0x0b,
	0x1c, 0xec, 0x02, 0x27, 0xe7, 0xf0, 0x68, 0xc9, 0x86, 0x50, 0xed, 0x7b, 0x7a, 0x4a, 0xec, 0xa8,
	0xd9, 0x12, 0x46, 0x4b, 0x6e, 0x8d, 0x33, 0x49, 0x20, 0xf8, 0x48, 0x2c, 0xa3, 0xbd, 0x8f, 0x71,
	0x76, 0x92, 0xee, 0x85, 0x4a, 0x8e, 0xe4, 0x39, 0x84, 0xf3, 0x76, 0xdd, 0xf4, 0x53, 0xb3, 0x3f,
	0x02, 0x06, 0x37, 0xc4, 0x28, 0x2f, 0x20, 0x9a, 0x1b, 0x42, 0xa6, 0x2f, 0x98, 0xd7, 0xd4, 0x6f,
	0x7d, 0x01, 0xd1, 0x82, 0x6a, 0x3a, 0x88, 0xbd, 0x84, 0x87, 0x0b, 0xb2, 0x85, 0xa9, 0xf2, 0x43,
	0xe8, 0x97, 0x00, 0xd7, 0x95, 0x65, 0x4f, 0xda, 0xfe, 0xe4, 0x2f, 0x20, 0xb8, 0x6e, 0x57, 0xf2,
	0x39, 0x84, 0xa8, 0x75, 0xdd, 0xc9, 0x61, 0xea, 0xff, 0x93, 0xff, 0xb8, 0x02, 0x46, 0xf3, 0x7a,
	0x6d, 0x99, 0x8c, 0xbc, 0x84, 0x07, 0x0b, 0x64, 0x5c, 0x76, 0x4d, 0x71, 0xd7, 0xfb, 0x71, 0x7a,
	0xf7, 0xfe, 0x93, 0xa3, 0x57, 0xc2, 0x9d, 0xd2, 0x65, 0xb9, 0x21, 0x95, 0x93, 0xe9, 0x0f, 0xf3,
	0x3e, 0xfa, 0x3a, 0xbe, 0x52, 0xb8, 0xd2, 0x15, 0xe9, 0x3c, 0x1f, 0xfa, 0x97, 0xf4, 0xfa, 0x6f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x54, 0xc8, 0xb6, 0x98, 0x59, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StorageClient interface {
	Get(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
	Count(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
}

type storageClient struct {
	cc *grpc.ClientConn
}

func NewStorageClient(cc *grpc.ClientConn) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) Get(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Storage/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) Count(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Storage/Count", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageServer is the server API for Storage service.
type StorageServer interface {
	Get(context.Context, *RpcRequest) (*RpcResponse, error)
	Count(context.Context, *RpcRequest) (*RpcResponse, error)
}

// UnimplementedStorageServer can be embedded to have forward compatible implementations.
type UnimplementedStorageServer struct {
}

func (*UnimplementedStorageServer) Get(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedStorageServer) Count(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Count not implemented")
}

func RegisterStorageServer(s *grpc.Server, srv StorageServer) {
	s.RegisterService(&_Storage_serviceDesc, srv)
}

func _Storage_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Storage/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Get(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_Count_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Count(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Storage/Count",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Count(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Storage_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Storage",
	HandlerType: (*StorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Storage_Get_Handler,
		},
		{
			MethodName: "Count",
			Handler:    _Storage_Count_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "system.proto",
}

// MetaClient is the client API for Meta service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MetaClient interface {
	CreateTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
	DeleteTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
	DescribeTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
	ListTables(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
}

type metaClient struct {
	cc *grpc.ClientConn
}

func NewMetaClient(cc *grpc.ClientConn) MetaClient {
	return &metaClient{cc}
}

func (c *metaClient) CreateTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Meta/CreateTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DeleteTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Meta/DeleteTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DescribeTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Meta/DescribeTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) ListTables(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Meta/ListTables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetaServer is the server API for Meta service.
type MetaServer interface {
	CreateTable(context.Context, *RpcRequest) (*RpcResponse, error)
	DeleteTable(context.Context, *RpcRequest) (*RpcResponse, error)
	DescribeTable(context.Context, *RpcRequest) (*RpcResponse, error)
	ListTables(context.Context, *RpcRequest) (*RpcResponse, error)
}

// UnimplementedMetaServer can be embedded to have forward compatible implementations.
type UnimplementedMetaServer struct {
}

func (*UnimplementedMetaServer) CreateTable(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTable not implemented")
}
func (*UnimplementedMetaServer) DeleteTable(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTable not implemented")
}
func (*UnimplementedMetaServer) DescribeTable(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeTable not implemented")
}
func (*UnimplementedMetaServer) ListTables(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTables not implemented")
}

func RegisterMetaServer(s *grpc.Server, srv MetaServer) {
	s.RegisterService(&_Meta_serviceDesc, srv)
}

func _Meta_CreateTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).CreateTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Meta/CreateTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).CreateTable(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_DeleteTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).DeleteTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Meta/DeleteTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).DeleteTable(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_DescribeTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).DescribeTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Meta/DescribeTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).DescribeTable(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_ListTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).ListTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Meta/ListTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).ListTables(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Meta_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Meta",
	HandlerType: (*MetaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTable",
			Handler:    _Meta_CreateTable_Handler,
		},
		{
			MethodName: "DeleteTable",
			Handler:    _Meta_DeleteTable_Handler,
		},
		{
			MethodName: "DescribeTable",
			Handler:    _Meta_DescribeTable_Handler,
		},
		{
			MethodName: "ListTables",
			Handler:    _Meta_ListTables_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "system.proto",
}

// LogClient is the client API for Log service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LogClient interface {
	Apply(ctx context.Context, in *Entry, opts ...grpc.CallOption) (*RpcResponse, error)
}

type logClient struct {
	cc *grpc.ClientConn
}

func NewLogClient(cc *grpc.ClientConn) LogClient {
	return &logClient{cc}
}

func (c *logClient) Apply(ctx context.Context, in *Entry, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Log/apply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogServer is the server API for Log service.
type LogServer interface {
	Apply(context.Context, *Entry) (*RpcResponse, error)
}

// UnimplementedLogServer can be embedded to have forward compatible implementations.
type UnimplementedLogServer struct {
}

func (*UnimplementedLogServer) Apply(ctx context.Context, req *Entry) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Apply not implemented")
}

func RegisterLogServer(s *grpc.Server, srv LogServer) {
	s.RegisterService(&_Log_serviceDesc, srv)
}

func _Log_Apply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Entry)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServer).Apply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Log/Apply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServer).Apply(ctx, req.(*Entry))
	}
	return interceptor(ctx, in, info, handler)
}

var _Log_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Log",
	HandlerType: (*LogServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "apply",
			Handler:    _Log_Apply_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "system.proto",
}

// ClusterClient is the client API for Cluster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClusterClient interface {
	DataSync(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (Cluster_DataSyncClient, error)
	ListMembers(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
}

type clusterClient struct {
	cc *grpc.ClientConn
}

func NewClusterClient(cc *grpc.ClientConn) ClusterClient {
	return &clusterClient{cc}
}

func (c *clusterClient) DataSync(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (Cluster_DataSyncClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Cluster_serviceDesc.Streams[0], "/Cluster/DataSync", opts...)
	if err != nil {
		return nil, err
	}
	x := &clusterDataSyncClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Cluster_DataSyncClient interface {
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type clusterDataSyncClient struct {
	grpc.ClientStream
}

func (x *clusterDataSyncClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *clusterClient) ListMembers(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Cluster/ListMembers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterServer is the server API for Cluster service.
type ClusterServer interface {
	DataSync(*RpcRequest, Cluster_DataSyncServer) error
	ListMembers(context.Context, *RpcRequest) (*RpcResponse, error)
}

// UnimplementedClusterServer can be embedded to have forward compatible implementations.
type UnimplementedClusterServer struct {
}

func (*UnimplementedClusterServer) DataSync(req *RpcRequest, srv Cluster_DataSyncServer) error {
	return status.Errorf(codes.Unimplemented, "method DataSync not implemented")
}
func (*UnimplementedClusterServer) ListMembers(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMembers not implemented")
}

func RegisterClusterServer(s *grpc.Server, srv ClusterServer) {
	s.RegisterService(&_Cluster_serviceDesc, srv)
}

func _Cluster_DataSync_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RpcRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ClusterServer).DataSync(m, &clusterDataSyncServer{stream})
}

type Cluster_DataSyncServer interface {
	Send(*StreamResponse) error
	grpc.ServerStream
}

type clusterDataSyncServer struct {
	grpc.ServerStream
}

func (x *clusterDataSyncServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Cluster_ListMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).ListMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Cluster/ListMembers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).ListMembers(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Cluster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Cluster",
	HandlerType: (*ClusterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListMembers",
			Handler:    _Cluster_ListMembers_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DataSync",
			Handler:       _Cluster_DataSync_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "system.proto",
}
