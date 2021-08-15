// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package apiV1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CacosClient is the client API for Cacos service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CacosClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
	AuthLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
	NamespaceList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NamespaceListReply, error)
	AppList(ctx context.Context, in *AppListReq, opts ...grpc.CallOption) (*AppListReply, error)
}

type cacosClient struct {
	cc grpc.ClientConnInterface
}

func NewCacosClient(cc grpc.ClientConnInterface) CacosClient {
	return &cacosClient{cc}
}

func (c *cacosClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/apiV1.Cacos/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cacosClient) AuthLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/apiV1.Cacos/AuthLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cacosClient) NamespaceList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NamespaceListReply, error) {
	out := new(NamespaceListReply)
	err := c.cc.Invoke(ctx, "/apiV1.Cacos/NamespaceList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cacosClient) AppList(ctx context.Context, in *AppListReq, opts ...grpc.CallOption) (*AppListReply, error) {
	out := new(AppListReply)
	err := c.cc.Invoke(ctx, "/apiV1.Cacos/AppList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CacosServer is the server API for Cacos service.
// All implementations must embed UnimplementedCacosServer
// for forward compatibility
type CacosServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	AuthLogin(context.Context, *LoginRequest) (*LoginReply, error)
	NamespaceList(context.Context, *emptypb.Empty) (*NamespaceListReply, error)
	AppList(context.Context, *AppListReq) (*AppListReply, error)
	mustEmbedUnimplementedCacosServer()
}

// UnimplementedCacosServer must be embedded to have forward compatible implementations.
type UnimplementedCacosServer struct {
}

func (UnimplementedCacosServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedCacosServer) AuthLogin(context.Context, *LoginRequest) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthLogin not implemented")
}
func (UnimplementedCacosServer) NamespaceList(context.Context, *emptypb.Empty) (*NamespaceListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NamespaceList not implemented")
}
func (UnimplementedCacosServer) AppList(context.Context, *AppListReq) (*AppListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppList not implemented")
}
func (UnimplementedCacosServer) mustEmbedUnimplementedCacosServer() {}

// UnsafeCacosServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CacosServer will
// result in compilation errors.
type UnsafeCacosServer interface {
	mustEmbedUnimplementedCacosServer()
}

func RegisterCacosServer(s grpc.ServiceRegistrar, srv CacosServer) {
	s.RegisterService(&Cacos_ServiceDesc, srv)
}

func _Cacos_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacosServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apiV1.Cacos/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacosServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cacos_AuthLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacosServer).AuthLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apiV1.Cacos/AuthLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacosServer).AuthLogin(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cacos_NamespaceList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacosServer).NamespaceList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apiV1.Cacos/NamespaceList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacosServer).NamespaceList(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cacos_AppList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacosServer).AppList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apiV1.Cacos/AppList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacosServer).AppList(ctx, req.(*AppListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Cacos_ServiceDesc is the grpc.ServiceDesc for Cacos service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cacos_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "apiV1.Cacos",
	HandlerType: (*CacosServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Cacos_SayHello_Handler,
		},
		{
			MethodName: "AuthLogin",
			Handler:    _Cacos_AuthLogin_Handler,
		},
		{
			MethodName: "NamespaceList",
			Handler:    _Cacos_NamespaceList_Handler,
		},
		{
			MethodName: "AppList",
			Handler:    _Cacos_AppList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cacos.proto",
}
