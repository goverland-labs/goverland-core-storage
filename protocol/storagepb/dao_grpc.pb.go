// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: storagepb/dao.proto

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
	Dao_GetByID_FullMethodName            = "/storagepb.Dao/GetByID"
	Dao_GetByFilter_FullMethodName        = "/storagepb.Dao/GetByFilter"
	Dao_GetTopByCategories_FullMethodName = "/storagepb.Dao/GetTopByCategories"
)

// DaoClient is the client API for Dao service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DaoClient interface {
	GetByID(ctx context.Context, in *DaoByIDRequest, opts ...grpc.CallOption) (*DaoByIDResponse, error)
	GetByFilter(ctx context.Context, in *DaoByFilterRequest, opts ...grpc.CallOption) (*DaoByFilterResponse, error)
	GetTopByCategories(ctx context.Context, in *TopByCategoriesRequest, opts ...grpc.CallOption) (*TopByCategoriesResponse, error)
}

type daoClient struct {
	cc grpc.ClientConnInterface
}

func NewDaoClient(cc grpc.ClientConnInterface) DaoClient {
	return &daoClient{cc}
}

func (c *daoClient) GetByID(ctx context.Context, in *DaoByIDRequest, opts ...grpc.CallOption) (*DaoByIDResponse, error) {
	out := new(DaoByIDResponse)
	err := c.cc.Invoke(ctx, Dao_GetByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daoClient) GetByFilter(ctx context.Context, in *DaoByFilterRequest, opts ...grpc.CallOption) (*DaoByFilterResponse, error) {
	out := new(DaoByFilterResponse)
	err := c.cc.Invoke(ctx, Dao_GetByFilter_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daoClient) GetTopByCategories(ctx context.Context, in *TopByCategoriesRequest, opts ...grpc.CallOption) (*TopByCategoriesResponse, error) {
	out := new(TopByCategoriesResponse)
	err := c.cc.Invoke(ctx, Dao_GetTopByCategories_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DaoServer is the server API for Dao service.
// All implementations must embed UnimplementedDaoServer
// for forward compatibility
type DaoServer interface {
	GetByID(context.Context, *DaoByIDRequest) (*DaoByIDResponse, error)
	GetByFilter(context.Context, *DaoByFilterRequest) (*DaoByFilterResponse, error)
	GetTopByCategories(context.Context, *TopByCategoriesRequest) (*TopByCategoriesResponse, error)
	mustEmbedUnimplementedDaoServer()
}

// UnimplementedDaoServer must be embedded to have forward compatible implementations.
type UnimplementedDaoServer struct {
}

func (UnimplementedDaoServer) GetByID(context.Context, *DaoByIDRequest) (*DaoByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedDaoServer) GetByFilter(context.Context, *DaoByFilterRequest) (*DaoByFilterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByFilter not implemented")
}
func (UnimplementedDaoServer) GetTopByCategories(context.Context, *TopByCategoriesRequest) (*TopByCategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopByCategories not implemented")
}
func (UnimplementedDaoServer) mustEmbedUnimplementedDaoServer() {}

// UnsafeDaoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DaoServer will
// result in compilation errors.
type UnsafeDaoServer interface {
	mustEmbedUnimplementedDaoServer()
}

func RegisterDaoServer(s grpc.ServiceRegistrar, srv DaoServer) {
	s.RegisterService(&Dao_ServiceDesc, srv)
}

func _Dao_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DaoByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaoServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Dao_GetByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaoServer).GetByID(ctx, req.(*DaoByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dao_GetByFilter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DaoByFilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaoServer).GetByFilter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Dao_GetByFilter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaoServer).GetByFilter(ctx, req.(*DaoByFilterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dao_GetTopByCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopByCategoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaoServer).GetTopByCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Dao_GetTopByCategories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaoServer).GetTopByCategories(ctx, req.(*TopByCategoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Dao_ServiceDesc is the grpc.ServiceDesc for Dao service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Dao_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storagepb.Dao",
	HandlerType: (*DaoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByID",
			Handler:    _Dao_GetByID_Handler,
		},
		{
			MethodName: "GetByFilter",
			Handler:    _Dao_GetByFilter_Handler,
		},
		{
			MethodName: "GetTopByCategories",
			Handler:    _Dao_GetTopByCategories_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storagepb/dao.proto",
}
