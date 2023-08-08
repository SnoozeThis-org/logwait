// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ObserverServiceClient is the client API for ObserverService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ObserverServiceClient interface {
	Communicate(ctx context.Context, opts ...grpc.CallOption) (ObserverService_CommunicateClient, error)
}

type observerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewObserverServiceClient(cc grpc.ClientConnInterface) ObserverServiceClient {
	return &observerServiceClient{cc}
}

func (c *observerServiceClient) Communicate(ctx context.Context, opts ...grpc.CallOption) (ObserverService_CommunicateClient, error) {
	stream, err := c.cc.NewStream(ctx, &ObserverService_ServiceDesc.Streams[0], "/com.snoozethis.logwait.ObserverService/Communicate", opts...)
	if err != nil {
		return nil, err
	}
	x := &observerServiceCommunicateClient{stream}
	return x, nil
}

type ObserverService_CommunicateClient interface {
	Send(*ScannerToObserver) error
	Recv() (*ObserverToScanner, error)
	grpc.ClientStream
}

type observerServiceCommunicateClient struct {
	grpc.ClientStream
}

func (x *observerServiceCommunicateClient) Send(m *ScannerToObserver) error {
	return x.ClientStream.SendMsg(m)
}

func (x *observerServiceCommunicateClient) Recv() (*ObserverToScanner, error) {
	m := new(ObserverToScanner)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ObserverServiceServer is the server API for ObserverService service.
// All implementations should embed UnimplementedObserverServiceServer
// for forward compatibility
type ObserverServiceServer interface {
	Communicate(ObserverService_CommunicateServer) error
}

// UnimplementedObserverServiceServer should be embedded to have forward compatible implementations.
type UnimplementedObserverServiceServer struct {
}

func (UnimplementedObserverServiceServer) Communicate(ObserverService_CommunicateServer) error {
	return status.Errorf(codes.Unimplemented, "method Communicate not implemented")
}

// UnsafeObserverServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ObserverServiceServer will
// result in compilation errors.
type UnsafeObserverServiceServer interface {
	mustEmbedUnimplementedObserverServiceServer()
}

func RegisterObserverServiceServer(s grpc.ServiceRegistrar, srv ObserverServiceServer) {
	s.RegisterService(&ObserverService_ServiceDesc, srv)
}

func _ObserverService_Communicate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ObserverServiceServer).Communicate(&observerServiceCommunicateServer{stream})
}

type ObserverService_CommunicateServer interface {
	Send(*ObserverToScanner) error
	Recv() (*ScannerToObserver, error)
	grpc.ServerStream
}

type observerServiceCommunicateServer struct {
	grpc.ServerStream
}

func (x *observerServiceCommunicateServer) Send(m *ObserverToScanner) error {
	return x.ServerStream.SendMsg(m)
}

func (x *observerServiceCommunicateServer) Recv() (*ScannerToObserver, error) {
	m := new(ScannerToObserver)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ObserverService_ServiceDesc is the grpc.ServiceDesc for ObserverService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ObserverService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.snoozethis.logwait.ObserverService",
	HandlerType: (*ObserverServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Communicate",
			Handler:       _ObserverService_Communicate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "scanner.proto",
}
