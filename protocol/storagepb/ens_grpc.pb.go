// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: storagepb/ens.proto

package storagepb

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

const (
	Ens_GetEnsByAddresses_FullMethodName = "/storagepb.Ens/GetEnsByAddresses"
)

// EnsClient is the client API for Ens service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnsClient interface {
	GetEnsByAddresses(ctx context.Context, in *EnsByAddressesRequest, opts ...grpc.CallOption) (*EnsByAddressesResponse, error)
}

type ensClient struct {
	cc grpc.ClientConnInterface
}

func NewEnsClient(cc grpc.ClientConnInterface) EnsClient {
	return &ensClient{cc}
}

func (c *ensClient) GetEnsByAddresses(ctx context.Context, in *EnsByAddressesRequest, opts ...grpc.CallOption) (*EnsByAddressesResponse, error) {
	out := new(EnsByAddressesResponse)
	err := c.cc.Invoke(ctx, Ens_GetEnsByAddresses_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnsServer is the server API for Ens service.
// All implementations must embed UnimplementedEnsServer
// for forward compatibility
type EnsServer interface {
	GetEnsByAddresses(context.Context, *EnsByAddressesRequest) (*EnsByAddressesResponse, error)
	mustEmbedUnimplementedEnsServer()
}

// UnimplementedEnsServer must be embedded to have forward compatible implementations.
type UnimplementedEnsServer struct {
}

func (UnimplementedEnsServer) GetEnsByAddresses(context.Context, *EnsByAddressesRequest) (*EnsByAddressesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnsByAddresses not implemented")
}
func (UnimplementedEnsServer) mustEmbedUnimplementedEnsServer() {}

// UnsafeEnsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EnsServer will
// result in compilation errors.
type UnsafeEnsServer interface {
	mustEmbedUnimplementedEnsServer()
}

func RegisterEnsServer(s grpc.ServiceRegistrar, srv EnsServer) {
	s.RegisterService(&Ens_ServiceDesc, srv)
}

func _Ens_GetEnsByAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnsByAddressesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsServer).GetEnsByAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ens_GetEnsByAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsServer).GetEnsByAddresses(ctx, req.(*EnsByAddressesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Ens_ServiceDesc is the grpc.ServiceDesc for Ens service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ens_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storagepb.Ens",
	HandlerType: (*EnsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEnsByAddresses",
			Handler:    _Ens_GetEnsByAddresses_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storagepb/ens.proto",
}
