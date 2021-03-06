// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/logdog/api/endpoints/coordinator/registration/v1/service.proto
// DO NOT EDIT!

/*
Package logdog is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/logdog/api/endpoints/coordinator/registration/v1/service.proto

It has these top-level messages:
	RegisterPrefixRequest
	RegisterPrefixResponse
*/
package logdog

import prpc "github.com/luci/luci-go/grpc/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/luci/luci-go/common/proto/google"

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

// RegisterPrefixRequest registers a new Prefix with the Coordinator.
type RegisterPrefixRequest struct {
	// The log stream's project.
	Project string `protobuf:"bytes,1,opt,name=project" json:"project,omitempty"`
	// The log stream prefix to register.
	Prefix string `protobuf:"bytes,2,opt,name=prefix" json:"prefix,omitempty"`
	// Optional information about the registering agent.
	SourceInfo []string `protobuf:"bytes,3,rep,name=source_info,json=sourceInfo" json:"source_info,omitempty"`
	// The prefix expiration time. If <= 0, the project's default prefix
	// expiration period will be applied.
	//
	// The prefix will be closed by the Coordinator after its expiration period.
	// Once closed, new stream registration requests will no longer be accepted.
	//
	// If supplied, this value should exceed the timeout of the local task, else
	// some of the task's streams may be dropped due to failing registration.
	Expiration *google_protobuf.Duration `protobuf:"bytes,10,opt,name=expiration" json:"expiration,omitempty"`
}

func (m *RegisterPrefixRequest) Reset()                    { *m = RegisterPrefixRequest{} }
func (m *RegisterPrefixRequest) String() string            { return proto.CompactTextString(m) }
func (*RegisterPrefixRequest) ProtoMessage()               {}
func (*RegisterPrefixRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RegisterPrefixRequest) GetExpiration() *google_protobuf.Duration {
	if m != nil {
		return m.Expiration
	}
	return nil
}

// The response message for the RegisterPrefix RPC.
type RegisterPrefixResponse struct {
	// Secret is the prefix's secret. This must be included verbatim in Butler
	// bundles to assert ownership of this prefix.
	Secret []byte `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
	// The name of the Pub/Sub topic to publish butlerproto-formatted Butler log
	// bundles to.
	LogBundleTopic string `protobuf:"bytes,2,opt,name=log_bundle_topic,json=logBundleTopic" json:"log_bundle_topic,omitempty"`
}

func (m *RegisterPrefixResponse) Reset()                    { *m = RegisterPrefixResponse{} }
func (m *RegisterPrefixResponse) String() string            { return proto.CompactTextString(m) }
func (*RegisterPrefixResponse) ProtoMessage()               {}
func (*RegisterPrefixResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*RegisterPrefixRequest)(nil), "logdog.RegisterPrefixRequest")
	proto.RegisterType((*RegisterPrefixResponse)(nil), "logdog.RegisterPrefixResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Registration service

type RegistrationClient interface {
	// RegisterStream allows a Butler instance to register a log stream with the
	// Coordinator. Upon success, the Coordinator will return registration
	// information and streaming parameters to the Butler.
	//
	// This should be called by a Butler instance to gain the ability to publish
	// to a prefix space. The caller must have WRITE access to its project's
	// stream space. If WRITE access is not present, this will fail with the
	// "PermissionDenied" gRPC code.
	//
	// A stream prefix may be registered at most once. Additional registration
	// requests will fail with the "AlreadyExists" gRPC code.
	RegisterPrefix(ctx context.Context, in *RegisterPrefixRequest, opts ...grpc.CallOption) (*RegisterPrefixResponse, error)
}
type registrationPRPCClient struct {
	client *prpc.Client
}

func NewRegistrationPRPCClient(client *prpc.Client) RegistrationClient {
	return &registrationPRPCClient{client}
}

func (c *registrationPRPCClient) RegisterPrefix(ctx context.Context, in *RegisterPrefixRequest, opts ...grpc.CallOption) (*RegisterPrefixResponse, error) {
	out := new(RegisterPrefixResponse)
	err := c.client.Call(ctx, "logdog.Registration", "RegisterPrefix", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type registrationClient struct {
	cc *grpc.ClientConn
}

func NewRegistrationClient(cc *grpc.ClientConn) RegistrationClient {
	return &registrationClient{cc}
}

func (c *registrationClient) RegisterPrefix(ctx context.Context, in *RegisterPrefixRequest, opts ...grpc.CallOption) (*RegisterPrefixResponse, error) {
	out := new(RegisterPrefixResponse)
	err := grpc.Invoke(ctx, "/logdog.Registration/RegisterPrefix", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Registration service

type RegistrationServer interface {
	// RegisterStream allows a Butler instance to register a log stream with the
	// Coordinator. Upon success, the Coordinator will return registration
	// information and streaming parameters to the Butler.
	//
	// This should be called by a Butler instance to gain the ability to publish
	// to a prefix space. The caller must have WRITE access to its project's
	// stream space. If WRITE access is not present, this will fail with the
	// "PermissionDenied" gRPC code.
	//
	// A stream prefix may be registered at most once. Additional registration
	// requests will fail with the "AlreadyExists" gRPC code.
	RegisterPrefix(context.Context, *RegisterPrefixRequest) (*RegisterPrefixResponse, error)
}

func RegisterRegistrationServer(s prpc.Registrar, srv RegistrationServer) {
	s.RegisterService(&_Registration_serviceDesc, srv)
}

func _Registration_RegisterPrefix_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterPrefixRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationServer).RegisterPrefix(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logdog.Registration/RegisterPrefix",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationServer).RegisterPrefix(ctx, req.(*RegisterPrefixRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Registration_serviceDesc = grpc.ServiceDesc{
	ServiceName: "logdog.Registration",
	HandlerType: (*RegistrationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterPrefix",
			Handler:    _Registration_RegisterPrefix_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/logdog/api/endpoints/coordinator/registration/v1/service.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 309 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x90, 0x4f, 0x4b, 0x03, 0x31,
	0x10, 0xc5, 0xa9, 0x85, 0x8a, 0x69, 0x29, 0x12, 0xb0, 0xac, 0x05, 0x6b, 0xe9, 0xa9, 0x17, 0x13,
	0xac, 0x27, 0xaf, 0xe2, 0xc5, 0x93, 0xb2, 0x88, 0x07, 0x2f, 0x4b, 0x37, 0x3b, 0x1b, 0x23, 0x6b,
	0x26, 0xe6, 0x4f, 0xe9, 0x17, 0xf2, 0x7b, 0xba, 0x4d, 0xb6, 0x50, 0x45, 0x2f, 0x81, 0x79, 0xf3,
	0x5e, 0xf2, 0xcb, 0x23, 0x2f, 0x52, 0xf9, 0xb7, 0x50, 0x32, 0x81, 0x1f, 0xbc, 0x09, 0x42, 0xc5,
	0xe3, 0x4a, 0x22, 0x6f, 0x50, 0x56, 0x28, 0xf9, 0xda, 0x28, 0x0e, 0xba, 0x32, 0xa8, 0xb4, 0x77,
	0x5c, 0x20, 0xda, 0x4a, 0xe9, 0xb5, 0x47, 0xcb, 0x2d, 0x48, 0xe5, 0xbc, 0x5d, 0x7b, 0x85, 0x9a,
	0x6f, 0xae, 0xb9, 0x03, 0xbb, 0x51, 0x02, 0x98, 0xb1, 0xe8, 0x91, 0x0e, 0x52, 0x7e, 0x3a, 0x93,
	0x88, 0xb2, 0x01, 0x1e, 0xd5, 0x32, 0xd4, 0xbc, 0x0a, 0x29, 0x92, 0x7c, 0x8b, 0xaf, 0x1e, 0x39,
	0xcb, 0xe3, 0x4d, 0x60, 0x9f, 0x2c, 0xd4, 0x6a, 0x9b, 0xc3, 0x67, 0x00, 0xe7, 0x69, 0x46, 0x8e,
	0x5b, 0xcb, 0x3b, 0x08, 0x9f, 0xf5, 0xe6, 0xbd, 0xe5, 0x49, 0xbe, 0x1f, 0xe9, 0x84, 0x0c, 0x4c,
	0xb4, 0x66, 0x47, 0x71, 0xd1, 0x4d, 0xf4, 0x92, 0x0c, 0x1d, 0x06, 0x2b, 0xa0, 0x50, 0xba, 0xc6,
	0xac, 0x3f, 0xef, 0xb7, 0x4b, 0x92, 0xa4, 0x87, 0x56, 0xa1, 0xb7, 0x84, 0xc0, 0xd6, 0xa8, 0x04,
	0x90, 0x91, 0x36, 0x3c, 0x5c, 0x9d, 0xb3, 0x44, 0xc8, 0xf6, 0x84, 0xec, 0xbe, 0x23, 0xcc, 0x0f,
	0xcc, 0x8b, 0x57, 0x32, 0xf9, 0x8d, 0xe9, 0x0c, 0x6a, 0x07, 0x3b, 0x1a, 0x07, 0xc2, 0x42, 0xc2,
	0x1c, 0xe5, 0xdd, 0x44, 0x97, 0xe4, 0xb4, 0xed, 0xa0, 0x28, 0x83, 0xae, 0x1a, 0x28, 0x3c, 0x1a,
	0x25, 0x3a, 0xde, 0x71, 0xab, 0xdf, 0x45, 0xf9, 0x79, 0xa7, 0xae, 0x0a, 0x32, 0xca, 0x0f, 0xca,
	0xa4, 0x8f, 0x64, 0xfc, 0xf3, 0x2d, 0x7a, 0xc1, 0x52, 0x9d, 0xec, 0xcf, 0xaa, 0xa6, 0xb3, 0xff,
	0xd6, 0x09, 0xb1, 0x1c, 0xc4, 0xbf, 0xdd, 0x7c, 0x07, 0x00, 0x00, 0xff, 0xff, 0x15, 0x5f, 0xdf,
	0x53, 0xed, 0x01, 0x00, 0x00,
}
