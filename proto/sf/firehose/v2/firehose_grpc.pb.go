// copied the following files
// https://github.com/streamingfast/proto/blob/17ecce85bc8a6464a995933e68b172ba63646753/sf/firehose/v2/firehose.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: sf/firehose/v2/firehose.proto

package pbfirehose

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Stream_Blocks_FullMethodName = "/sf.firehose.v2.Stream/Blocks"
)

// StreamClient is the client API for Stream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamClient interface {
	Blocks(ctx context.Context, in *Request, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Response], error)
}

type streamClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamClient(cc grpc.ClientConnInterface) StreamClient {
	return &streamClient{cc}
}

func (c *streamClient) Blocks(ctx context.Context, in *Request, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Stream_ServiceDesc.Streams[0], Stream_Blocks_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Request, Response]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Stream_BlocksClient = grpc.ServerStreamingClient[Response]

// StreamServer is the server API for Stream service.
// All implementations must embed UnimplementedStreamServer
// for forward compatibility.
type StreamServer interface {
	Blocks(*Request, grpc.ServerStreamingServer[Response]) error
	mustEmbedUnimplementedStreamServer()
}

// UnimplementedStreamServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStreamServer struct{}

func (UnimplementedStreamServer) Blocks(*Request, grpc.ServerStreamingServer[Response]) error {
	return status.Errorf(codes.Unimplemented, "method Blocks not implemented")
}
func (UnimplementedStreamServer) mustEmbedUnimplementedStreamServer() {}
func (UnimplementedStreamServer) testEmbeddedByValue()                {}

// UnsafeStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamServer will
// result in compilation errors.
type UnsafeStreamServer interface {
	mustEmbedUnimplementedStreamServer()
}

func RegisterStreamServer(s grpc.ServiceRegistrar, srv StreamServer) {
	// If the following call pancis, it indicates UnimplementedStreamServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Stream_ServiceDesc, srv)
}

func _Stream_Blocks_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServer).Blocks(m, &grpc.GenericServerStream[Request, Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Stream_BlocksServer = grpc.ServerStreamingServer[Response]

// Stream_ServiceDesc is the grpc.ServiceDesc for Stream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sf.firehose.v2.Stream",
	HandlerType: (*StreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Blocks",
			Handler:       _Stream_Blocks_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "sf/firehose/v2/firehose.proto",
}

const (
	Fetch_Block_FullMethodName = "/sf.firehose.v2.Fetch/Block"
)

// FetchClient is the client API for Fetch service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FetchClient interface {
	Block(ctx context.Context, in *SingleBlockRequest, opts ...grpc.CallOption) (*SingleBlockResponse, error)
}

type fetchClient struct {
	cc grpc.ClientConnInterface
}

func NewFetchClient(cc grpc.ClientConnInterface) FetchClient {
	return &fetchClient{cc}
}

func (c *fetchClient) Block(ctx context.Context, in *SingleBlockRequest, opts ...grpc.CallOption) (*SingleBlockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SingleBlockResponse)
	err := c.cc.Invoke(ctx, Fetch_Block_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FetchServer is the server API for Fetch service.
// All implementations must embed UnimplementedFetchServer
// for forward compatibility.
type FetchServer interface {
	Block(context.Context, *SingleBlockRequest) (*SingleBlockResponse, error)
	mustEmbedUnimplementedFetchServer()
}

// UnimplementedFetchServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFetchServer struct{}

func (UnimplementedFetchServer) Block(context.Context, *SingleBlockRequest) (*SingleBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Block not implemented")
}
func (UnimplementedFetchServer) mustEmbedUnimplementedFetchServer() {}
func (UnimplementedFetchServer) testEmbeddedByValue()               {}

// UnsafeFetchServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FetchServer will
// result in compilation errors.
type UnsafeFetchServer interface {
	mustEmbedUnimplementedFetchServer()
}

func RegisterFetchServer(s grpc.ServiceRegistrar, srv FetchServer) {
	// If the following call pancis, it indicates UnimplementedFetchServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Fetch_ServiceDesc, srv)
}

func _Fetch_Block_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SingleBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FetchServer).Block(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Fetch_Block_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FetchServer).Block(ctx, req.(*SingleBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Fetch_ServiceDesc is the grpc.ServiceDesc for Fetch service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Fetch_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sf.firehose.v2.Fetch",
	HandlerType: (*FetchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Block",
			Handler:    _Fetch_Block_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sf/firehose/v2/firehose.proto",
}

const (
	EndpointInfo_Info_FullMethodName = "/sf.firehose.v2.EndpointInfo/Info"
)

// EndpointInfoClient is the client API for EndpointInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EndpointInfoClient interface {
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error)
}

type endpointInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewEndpointInfoClient(cc grpc.ClientConnInterface) EndpointInfoClient {
	return &endpointInfoClient{cc}
}

func (c *endpointInfoClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InfoResponse)
	err := c.cc.Invoke(ctx, EndpointInfo_Info_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointInfoServer is the server API for EndpointInfo service.
// All implementations must embed UnimplementedEndpointInfoServer
// for forward compatibility.
type EndpointInfoServer interface {
	Info(context.Context, *InfoRequest) (*InfoResponse, error)
	mustEmbedUnimplementedEndpointInfoServer()
}

// UnimplementedEndpointInfoServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEndpointInfoServer struct{}

func (UnimplementedEndpointInfoServer) Info(context.Context, *InfoRequest) (*InfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedEndpointInfoServer) mustEmbedUnimplementedEndpointInfoServer() {}
func (UnimplementedEndpointInfoServer) testEmbeddedByValue()                      {}

// UnsafeEndpointInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EndpointInfoServer will
// result in compilation errors.
type UnsafeEndpointInfoServer interface {
	mustEmbedUnimplementedEndpointInfoServer()
}

func RegisterEndpointInfoServer(s grpc.ServiceRegistrar, srv EndpointInfoServer) {
	// If the following call pancis, it indicates UnimplementedEndpointInfoServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EndpointInfo_ServiceDesc, srv)
}

func _EndpointInfo_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointInfoServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointInfo_Info_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointInfoServer).Info(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EndpointInfo_ServiceDesc is the grpc.ServiceDesc for EndpointInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EndpointInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sf.firehose.v2.EndpointInfo",
	HandlerType: (*EndpointInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Info",
			Handler:    _EndpointInfo_Info_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sf/firehose/v2/firehose.proto",
}
