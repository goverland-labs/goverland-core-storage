// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: storagepb/vote.proto

package storagepb

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
	Vote_GetVotes_FullMethodName       = "/storagepb.Vote/GetVotes"
	Vote_Validate_FullMethodName       = "/storagepb.Vote/Validate"
	Vote_Prepare_FullMethodName        = "/storagepb.Vote/Prepare"
	Vote_Vote_FullMethodName           = "/storagepb.Vote/Vote"
	Vote_GetDaosVotedIn_FullMethodName = "/storagepb.Vote/GetDaosVotedIn"
)

// VoteClient is the client API for Vote service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VoteClient interface {
	GetVotes(ctx context.Context, in *VotesFilterRequest, opts ...grpc.CallOption) (*VotesFilterResponse, error)
	Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
	Prepare(ctx context.Context, in *PrepareRequest, opts ...grpc.CallOption) (*PrepareResponse, error)
	Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error)
	GetDaosVotedIn(ctx context.Context, in *DaosVotedInRequest, opts ...grpc.CallOption) (*DaosVotedInResponse, error)
}

type voteClient struct {
	cc grpc.ClientConnInterface
}

func NewVoteClient(cc grpc.ClientConnInterface) VoteClient {
	return &voteClient{cc}
}

func (c *voteClient) GetVotes(ctx context.Context, in *VotesFilterRequest, opts ...grpc.CallOption) (*VotesFilterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VotesFilterResponse)
	err := c.cc.Invoke(ctx, Vote_GetVotes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *voteClient) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, Vote_Validate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *voteClient) Prepare(ctx context.Context, in *PrepareRequest, opts ...grpc.CallOption) (*PrepareResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PrepareResponse)
	err := c.cc.Invoke(ctx, Vote_Prepare_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *voteClient) Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VoteResponse)
	err := c.cc.Invoke(ctx, Vote_Vote_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *voteClient) GetDaosVotedIn(ctx context.Context, in *DaosVotedInRequest, opts ...grpc.CallOption) (*DaosVotedInResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DaosVotedInResponse)
	err := c.cc.Invoke(ctx, Vote_GetDaosVotedIn_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VoteServer is the server API for Vote service.
// All implementations must embed UnimplementedVoteServer
// for forward compatibility.
type VoteServer interface {
	GetVotes(context.Context, *VotesFilterRequest) (*VotesFilterResponse, error)
	Validate(context.Context, *ValidateRequest) (*ValidateResponse, error)
	Prepare(context.Context, *PrepareRequest) (*PrepareResponse, error)
	Vote(context.Context, *VoteRequest) (*VoteResponse, error)
	GetDaosVotedIn(context.Context, *DaosVotedInRequest) (*DaosVotedInResponse, error)
	mustEmbedUnimplementedVoteServer()
}

// UnimplementedVoteServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedVoteServer struct{}

func (UnimplementedVoteServer) GetVotes(context.Context, *VotesFilterRequest) (*VotesFilterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVotes not implemented")
}
func (UnimplementedVoteServer) Validate(context.Context, *ValidateRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}
func (UnimplementedVoteServer) Prepare(context.Context, *PrepareRequest) (*PrepareResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Prepare not implemented")
}
func (UnimplementedVoteServer) Vote(context.Context, *VoteRequest) (*VoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
}
func (UnimplementedVoteServer) GetDaosVotedIn(context.Context, *DaosVotedInRequest) (*DaosVotedInResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDaosVotedIn not implemented")
}
func (UnimplementedVoteServer) mustEmbedUnimplementedVoteServer() {}
func (UnimplementedVoteServer) testEmbeddedByValue()              {}

// UnsafeVoteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VoteServer will
// result in compilation errors.
type UnsafeVoteServer interface {
	mustEmbedUnimplementedVoteServer()
}

func RegisterVoteServer(s grpc.ServiceRegistrar, srv VoteServer) {
	// If the following call pancis, it indicates UnimplementedVoteServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Vote_ServiceDesc, srv)
}

func _Vote_GetVotes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VotesFilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoteServer).GetVotes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vote_GetVotes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoteServer).GetVotes(ctx, req.(*VotesFilterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vote_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoteServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vote_Validate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoteServer).Validate(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vote_Prepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrepareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoteServer).Prepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vote_Prepare_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoteServer).Prepare(ctx, req.(*PrepareRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vote_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoteServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vote_Vote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoteServer).Vote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vote_GetDaosVotedIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DaosVotedInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoteServer).GetDaosVotedIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Vote_GetDaosVotedIn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoteServer).GetDaosVotedIn(ctx, req.(*DaosVotedInRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Vote_ServiceDesc is the grpc.ServiceDesc for Vote service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Vote_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storagepb.Vote",
	HandlerType: (*VoteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVotes",
			Handler:    _Vote_GetVotes_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _Vote_Validate_Handler,
		},
		{
			MethodName: "Prepare",
			Handler:    _Vote_Prepare_Handler,
		},
		{
			MethodName: "Vote",
			Handler:    _Vote_Vote_Handler,
		},
		{
			MethodName: "GetDaosVotedIn",
			Handler:    _Vote_GetDaosVotedIn_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storagepb/vote.proto",
}
