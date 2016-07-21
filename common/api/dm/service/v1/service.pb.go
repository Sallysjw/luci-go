// Code generated by protoc-gen-go.
// source: service.proto
// DO NOT EDIT!

package dm

import prpccommon "github.com/luci/luci-go/common/prpc"
import prpc "github.com/luci/luci-go/server/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/luci/luci-go/common/proto/google"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Deps service

type DepsClient interface {
	// allows you to add additional data to the current dependency graph.
	EnsureGraphData(ctx context.Context, in *EnsureGraphDataReq, opts ...grpc.CallOption) (*EnsureGraphDataRsp, error)
	// is called by Execution clients to activate themselves with DM.
	ActivateExecution(ctx context.Context, in *ActivateExecutionReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	// is called by Execution clients to indicate that an Attempt is finished.
	FinishAttempt(ctx context.Context, in *FinishAttemptReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	// runs queries, and walks along the dependency graph from the query results.
	WalkGraph(ctx context.Context, in *WalkGraphReq, opts ...grpc.CallOption) (*GraphData, error)
}
type depsPRPCClient struct {
	client *prpccommon.Client
}

func NewDepsPRPCClient(client *prpccommon.Client) DepsClient {
	return &depsPRPCClient{client}
}

func (c *depsPRPCClient) EnsureGraphData(ctx context.Context, in *EnsureGraphDataReq, opts ...grpc.CallOption) (*EnsureGraphDataRsp, error) {
	out := new(EnsureGraphDataRsp)
	err := c.client.Call(ctx, "dm.Deps", "EnsureGraphData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsPRPCClient) ActivateExecution(ctx context.Context, in *ActivateExecutionReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := c.client.Call(ctx, "dm.Deps", "ActivateExecution", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsPRPCClient) FinishAttempt(ctx context.Context, in *FinishAttemptReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := c.client.Call(ctx, "dm.Deps", "FinishAttempt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsPRPCClient) WalkGraph(ctx context.Context, in *WalkGraphReq, opts ...grpc.CallOption) (*GraphData, error) {
	out := new(GraphData)
	err := c.client.Call(ctx, "dm.Deps", "WalkGraph", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type depsClient struct {
	cc *grpc.ClientConn
}

func NewDepsClient(cc *grpc.ClientConn) DepsClient {
	return &depsClient{cc}
}

func (c *depsClient) EnsureGraphData(ctx context.Context, in *EnsureGraphDataReq, opts ...grpc.CallOption) (*EnsureGraphDataRsp, error) {
	out := new(EnsureGraphDataRsp)
	err := grpc.Invoke(ctx, "/dm.Deps/EnsureGraphData", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsClient) ActivateExecution(ctx context.Context, in *ActivateExecutionReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/dm.Deps/ActivateExecution", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsClient) FinishAttempt(ctx context.Context, in *FinishAttemptReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/dm.Deps/FinishAttempt", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsClient) WalkGraph(ctx context.Context, in *WalkGraphReq, opts ...grpc.CallOption) (*GraphData, error) {
	out := new(GraphData)
	err := grpc.Invoke(ctx, "/dm.Deps/WalkGraph", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Deps service

type DepsServer interface {
	// allows you to add additional data to the current dependency graph.
	EnsureGraphData(context.Context, *EnsureGraphDataReq) (*EnsureGraphDataRsp, error)
	// is called by Execution clients to activate themselves with DM.
	ActivateExecution(context.Context, *ActivateExecutionReq) (*google_protobuf1.Empty, error)
	// is called by Execution clients to indicate that an Attempt is finished.
	FinishAttempt(context.Context, *FinishAttemptReq) (*google_protobuf1.Empty, error)
	// runs queries, and walks along the dependency graph from the query results.
	WalkGraph(context.Context, *WalkGraphReq) (*GraphData, error)
}

func RegisterDepsServer(s prpc.Registrar, srv DepsServer) {
	s.RegisterService(&_Deps_serviceDesc, srv)
}

func _Deps_EnsureGraphData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnsureGraphDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepsServer).EnsureGraphData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dm.Deps/EnsureGraphData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepsServer).EnsureGraphData(ctx, req.(*EnsureGraphDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deps_ActivateExecution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateExecutionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepsServer).ActivateExecution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dm.Deps/ActivateExecution",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepsServer).ActivateExecution(ctx, req.(*ActivateExecutionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deps_FinishAttempt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinishAttemptReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepsServer).FinishAttempt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dm.Deps/FinishAttempt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepsServer).FinishAttempt(ctx, req.(*FinishAttemptReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deps_WalkGraph_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WalkGraphReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepsServer).WalkGraph(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dm.Deps/WalkGraph",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepsServer).WalkGraph(ctx, req.(*WalkGraphReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Deps_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dm.Deps",
	HandlerType: (*DepsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EnsureGraphData",
			Handler:    _Deps_EnsureGraphData_Handler,
		},
		{
			MethodName: "ActivateExecution",
			Handler:    _Deps_ActivateExecution_Handler,
		},
		{
			MethodName: "FinishAttempt",
			Handler:    _Deps_FinishAttempt_Handler,
		},
		{
			MethodName: "WalkGraph",
			Handler:    _Deps_WalkGraph_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor5,
}

func init() { proto.RegisterFile("service.proto", fileDescriptor5) }

var fileDescriptor5 = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x90, 0x41, 0x4b, 0x03, 0x31,
	0x10, 0x85, 0x51, 0x44, 0x30, 0xb0, 0x58, 0x87, 0xa2, 0x65, 0xfd, 0x0b, 0x92, 0x82, 0x9e, 0x3d,
	0x2c, 0x74, 0xf5, 0xee, 0xc5, 0x63, 0x48, 0x77, 0xa7, 0xdb, 0x60, 0x77, 0x13, 0x93, 0xd9, 0xaa,
	0x3f, 0x5e, 0xb0, 0x93, 0xec, 0x16, 0xb4, 0xd2, 0x63, 0xbe, 0x99, 0xf7, 0xe6, 0xe5, 0x89, 0x2c,
	0xa0, 0xdf, 0x9a, 0x0a, 0xa5, 0xf3, 0x96, 0x2c, 0x9c, 0xd6, 0x6d, 0x7e, 0xdb, 0x58, 0xdb, 0x6c,
	0x70, 0x1e, 0xc9, 0xb2, 0x5f, 0xcd, 0xb1, 0x75, 0xf4, 0x95, 0x16, 0xf2, 0x49, 0xe3, 0xb5, 0x5b,
	0xab, 0x5a, 0x93, 0x1e, 0xc8, 0x0d, 0x76, 0xa1, 0xf7, 0xa8, 0x0e, 0x06, 0x33, 0x5d, 0x91, 0xd9,
	0x6a, 0x42, 0x85, 0x9f, 0x58, 0xf5, 0x64, 0x6c, 0x37, 0x4c, 0xa6, 0x2b, 0xd3, 0x99, 0xb0, 0x56,
	0x9a, 0x88, 0xbd, 0x47, 0xeb, 0x0f, 0xbd, 0x79, 0x4b, 0x36, 0x89, 0xdc, 0x7f, 0x9f, 0x88, 0xb3,
	0x05, 0xba, 0x00, 0x85, 0xb8, 0x2c, 0xe3, 0x95, 0x67, 0x9e, 0x2e, 0x76, 0x37, 0xe0, 0x5a, 0xd6,
	0xad, 0xfc, 0x03, 0x5f, 0xf0, 0x3d, 0xff, 0x97, 0x07, 0x07, 0xa5, 0xb8, 0x2a, 0x86, 0x3c, 0xe5,
	0x18, 0x07, 0x66, 0xbc, 0x7c, 0x80, 0x93, 0x4d, 0x6a, 0x41, 0x8e, 0x2d, 0xc8, 0x92, 0x5b, 0x80,
	0x47, 0x91, 0x3d, 0xc5, 0xf0, 0x45, 0xca, 0x0e, 0x53, 0xb6, 0xf8, 0x85, 0x8e, 0xc9, 0xef, 0xc4,
	0xc5, 0xeb, 0xee, 0x97, 0x31, 0x19, 0x4c, 0x58, 0xba, 0x7f, 0xb2, 0x2c, 0x63, 0xb2, 0x8f, 0xbd,
	0x3c, 0x8f, 0xea, 0x87, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x21, 0x8f, 0xba, 0xcb, 0xa5, 0x01,
	0x00, 0x00,
}
