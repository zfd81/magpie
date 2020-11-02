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
	return fileDescriptor_86a7260ebdc12f47, []int{1}
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

func init() {
	proto.RegisterType((*RpcRequest)(nil), "RpcRequest")
	proto.RegisterMapType((map[string]string)(nil), "RpcRequest.ParamsEntry")
	proto.RegisterType((*RpcResponse)(nil), "RpcResponse")
}

func init() { proto.RegisterFile("system.proto", fileDescriptor_86a7260ebdc12f47) }

var fileDescriptor_86a7260ebdc12f47 = []byte{
	// 285 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0xdd, 0xc6, 0xb4, 0x74, 0xb6, 0x82, 0x2c, 0x82, 0xa1, 0xa7, 0x12, 0x3c, 0xd4, 0x22,
	0x29, 0xd4, 0x8b, 0x7a, 0xb4, 0x15, 0x2f, 0x82, 0x92, 0x7a, 0xf2, 0xb6, 0x69, 0x87, 0x12, 0x4c,
	0xb2, 0xeb, 0xce, 0x46, 0xc8, 0x7f, 0xf0, 0x07, 0xf9, 0xf3, 0x24, 0xdb, 0xd6, 0xe6, 0x14, 0x7a,
	0x7b, 0x8f, 0xf7, 0xbd, 0x99, 0x59, 0x58, 0x18, 0x50, 0x45, 0x16, 0xf3, 0x48, 0x1b, 0x65, 0x55,
	0xf8, 0xc3, 0x00, 0x62, 0xbd, 0x8a, 0xf1, 0xab, 0x44, 0xb2, 0x62, 0x0a, 0x5d, 0x2d, 0x8d, 0xcc,
	0x29, 0x60, 0x23, 0x6f, 0xcc, 0x67, 0x97, 0xd1, 0x21, 0x8c, 0xde, 0x5c, 0xf2, 0x54, 0x58, 0x53,
	0xc5, 0x3b, 0x4c, 0x08, 0x38, 0x5d, 0x4b, 0x2b, 0x83, 0xce, 0x88, 0x8d, 0xfb, 0xb1, 0xd3, 0xc3,
	0x7b, 0xe0, 0x0d, 0x54, 0x9c, 0x83, 0xf7, 0x89, 0x55, 0xc0, 0x1c, 0x51, 0x4b, 0x71, 0x01, 0xfe,
	0xb7, 0xcc, 0x4a, 0xdc, 0xb5, 0xb6, 0xe6, 0xa1, 0x73, 0xc7, 0xc2, 0x57, 0xe0, 0x6e, 0x21, 0x69,
	0x55, 0x10, 0xd6, 0xd3, 0x57, 0x6a, 0x8d, 0xae, 0xeb, 0xc7, 0x4e, 0x8b, 0x00, 0x7a, 0x39, 0x12,
	0xc9, 0xcd, 0xbe, 0xbe, 0xb7, 0xff, 0xb7, 0x78, 0x87, 0x5b, 0x66, 0x4b, 0xe8, 0x2d, 0xad, 0x32,
	0x75, 0x1c, 0x82, 0xf7, 0x8c, 0x56, 0xf0, 0xc6, 0x93, 0x86, 0x83, 0xa8, 0xb1, 0x2e, 0x3c, 0x11,
	0x57, 0xe0, 0xcf, 0x55, 0x59, 0xb4, 0x53, 0xb3, 0x5f, 0x06, 0xfe, 0xbb, 0x4c, 0x32, 0x14, 0x13,
	0xe0, 0x73, 0x83, 0xd2, 0xe2, 0xd6, 0xb6, 0xce, 0x9e, 0x00, 0x5f, 0x60, 0x86, 0x47, 0xb1, 0x37,
	0x70, 0xb6, 0x40, 0x5a, 0x99, 0x34, 0x39, 0x86, 0xbe, 0x06, 0x78, 0x49, 0xc9, 0x3a, 0x92, 0x5a,
	0xd1, 0x47, 0xfe, 0xd1, 0x9f, 0xe6, 0x72, 0xa3, 0x53, 0xd4, 0x49, 0xd2, 0x75, 0x7f, 0xe0, 0xf6,
	0x2f, 0x00, 0x00, 0xff, 0xff, 0x92, 0xdd, 0x8f, 0xd9, 0x13, 0x02, 0x00, 0x00,
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

// TableClient is the client API for Table service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TableClient interface {
	CreateTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
	DeleteTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
	DescribeTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
	ListTables(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error)
}

type tableClient struct {
	cc *grpc.ClientConn
}

func NewTableClient(cc *grpc.ClientConn) TableClient {
	return &tableClient{cc}
}

func (c *tableClient) CreateTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Table/CreateTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tableClient) DeleteTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Table/DeleteTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tableClient) DescribeTable(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Table/DescribeTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tableClient) ListTables(ctx context.Context, in *RpcRequest, opts ...grpc.CallOption) (*RpcResponse, error) {
	out := new(RpcResponse)
	err := c.cc.Invoke(ctx, "/Table/ListTables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TableServer is the server API for Table service.
type TableServer interface {
	CreateTable(context.Context, *RpcRequest) (*RpcResponse, error)
	DeleteTable(context.Context, *RpcRequest) (*RpcResponse, error)
	DescribeTable(context.Context, *RpcRequest) (*RpcResponse, error)
	ListTables(context.Context, *RpcRequest) (*RpcResponse, error)
}

// UnimplementedTableServer can be embedded to have forward compatible implementations.
type UnimplementedTableServer struct {
}

func (*UnimplementedTableServer) CreateTable(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTable not implemented")
}
func (*UnimplementedTableServer) DeleteTable(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTable not implemented")
}
func (*UnimplementedTableServer) DescribeTable(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeTable not implemented")
}
func (*UnimplementedTableServer) ListTables(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTables not implemented")
}

func RegisterTableServer(s *grpc.Server, srv TableServer) {
	s.RegisterService(&_Table_serviceDesc, srv)
}

func _Table_CreateTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TableServer).CreateTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Table/CreateTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TableServer).CreateTable(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Table_DeleteTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TableServer).DeleteTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Table/DeleteTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TableServer).DeleteTable(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Table_DescribeTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TableServer).DescribeTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Table/DescribeTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TableServer).DescribeTable(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Table_ListTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TableServer).ListTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Table/ListTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TableServer).ListTables(ctx, req.(*RpcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Table_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Table",
	HandlerType: (*TableServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTable",
			Handler:    _Table_CreateTable_Handler,
		},
		{
			MethodName: "DeleteTable",
			Handler:    _Table_DeleteTable_Handler,
		},
		{
			MethodName: "DescribeTable",
			Handler:    _Table_DescribeTable_Handler,
		},
		{
			MethodName: "ListTables",
			Handler:    _Table_ListTables_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "system.proto",
}