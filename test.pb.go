// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

/*
Package registry is a generated protocol buffer package.

It is generated from these files:
	test.proto

It has these top-level messages:
	Request
	Response
*/
package registry

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type Response struct {
	NumOf int32  `protobuf:"varint,1,opt,name=num_of,json=numOf" json:"num_of,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Error string `protobuf:"bytes,3,opt,name=error" json:"error,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Response) GetNumOf() int32 {
	if m != nil {
		return m.NumOf
	}
	return 0
}

func (m *Response) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Response) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "registry.Request")
	proto.RegisterType((*Response)(nil), "registry.Response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Handler service

type HandlerClient interface {
	Process(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type handlerClient struct {
	cc *grpc.ClientConn
}

func NewHandlerClient(cc *grpc.ClientConn) HandlerClient {
	return &handlerClient{cc}
}

func (c *handlerClient) Process(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/registry.Handler/Process", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Handler service

type HandlerServer interface {
	Process(context.Context, *Request) (*Response, error)
}

func RegisterHandlerServer(s *grpc.Server, srv HandlerServer) {
	s.RegisterService(&_Handler_serviceDesc, srv)
}

func _Handler_Process_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HandlerServer).Process(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registry.Handler/Process",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HandlerServer).Process(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Handler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "registry.Handler",
	HandlerType: (*HandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Process",
			Handler:    _Handler_Process_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test.proto",
}

func init() { proto.RegisterFile("test.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 174 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8e, 0xbd, 0x0a, 0xc2, 0x50,
	0x0c, 0x85, 0xad, 0xda, 0x1f, 0xb3, 0x19, 0x14, 0x8a, 0x20, 0x48, 0x27, 0xa7, 0x0e, 0x75, 0x76,
	0x17, 0x1c, 0x94, 0xfb, 0x02, 0x52, 0x35, 0x15, 0xc1, 0xde, 0x5b, 0x93, 0x74, 0xe8, 0xdb, 0x4b,
	0x6f, 0x15, 0xdc, 0x4e, 0x0e, 0xf9, 0x38, 0x1f, 0x80, 0x92, 0x68, 0xde, 0xb0, 0x53, 0x87, 0x09,
	0xd3, 0xe3, 0x29, 0xca, 0x5d, 0xb6, 0x86, 0xd8, 0xd0, 0xbb, 0x25, 0x51, 0x44, 0x98, 0x6a, 0xd7,
	0x50, 0x1a, 0x6c, 0x82, 0xed, 0xcc, 0xf8, 0x9c, 0x1d, 0x21, 0x31, 0x24, 0x8d, 0xb3, 0x42, 0xb8,
	0x84, 0xc8, 0xb6, 0xf5, 0xc5, 0x55, 0xfe, 0x23, 0x34, 0xa1, 0x6d, 0xeb, 0x53, 0xd5, 0x63, 0xb6,
	0xac, 0x29, 0x1d, 0x0f, 0x58, 0x9f, 0x71, 0x01, 0x21, 0x31, 0x3b, 0x4e, 0x27, 0xbe, 0x1c, 0x8e,
	0x62, 0x0f, 0xf1, 0xa1, 0xb4, 0xf7, 0x17, 0x31, 0x16, 0x10, 0x9f, 0xd9, 0xdd, 0x48, 0x04, 0xe7,
	0xf9, 0x4f, 0x26, 0xff, 0x9a, 0xac, 0xf0, 0xbf, 0x1a, 0xd6, 0xb3, 0xd1, 0x35, 0xf2, 0xee, 0xbb,
	0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xaa, 0x15, 0x7c, 0x27, 0xc9, 0x00, 0x00, 0x00,
}
