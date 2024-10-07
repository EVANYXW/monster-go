// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.2.0
// source: gate.proto

package pb

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

// GateSerClient is the client API for GateSer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GateSerClient interface {
	SendToCli(ctx context.Context, in *SendToCli_Req, opts ...grpc.CallOption) (*SendToCli_Rsp, error)
}

type gateSerClient struct {
	cc grpc.ClientConnInterface
}

func NewGateSerClient(cc grpc.ClientConnInterface) GateSerClient {
	return &gateSerClient{cc}
}

func (c *gateSerClient) SendToCli(ctx context.Context, in *SendToCli_Req, opts ...grpc.CallOption) (*SendToCli_Rsp, error) {
	out := new(SendToCli_Rsp)
	err := c.cc.Invoke(ctx, "/proto.GateSer/SendToCli", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GateSerServer is the server API for GateSer service.
// All implementations must embed UnimplementedGateSerServer
// for forward compatibility
type GateSerServer interface {
	SendToCli(context.Context, *SendToCli_Req) (*SendToCli_Rsp, error)
	mustEmbedUnimplementedGateSerServer()
}

// UnimplementedGateSerServer must be embedded to have forward compatible implementations.
type UnimplementedGateSerServer struct {
}

func (UnimplementedGateSerServer) SendToCli(context.Context, *SendToCli_Req) (*SendToCli_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendToCli not implemented")
}
func (UnimplementedGateSerServer) mustEmbedUnimplementedGateSerServer() {}

// UnsafeGateSerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GateSerServer will
// result in compilation errors.
type UnsafeGateSerServer interface {
	mustEmbedUnimplementedGateSerServer()
}

func RegisterGateSerServer(s grpc.ServiceRegistrar, srv GateSerServer) {
	s.RegisterService(&GateSer_ServiceDesc, srv)
}

func _GateSer_SendToCli_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendToCli_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GateSerServer).SendToCli(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GateSer/SendToCli",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GateSerServer).SendToCli(ctx, req.(*SendToCli_Req))
	}
	return interceptor(ctx, in, info, handler)
}

// GateSer_ServiceDesc is the grpc.ServiceDesc for GateSer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GateSer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GateSer",
	HandlerType: (*GateSerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendToCli",
			Handler:    _GateSer_SendToCli_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gate.proto",
}