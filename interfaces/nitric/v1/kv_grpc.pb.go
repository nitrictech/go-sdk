// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// KeyValueClient is the client API for KeyValue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeyValueClient interface {
	// Get an existing key
	Get(ctx context.Context, in *KeyValueGetRequest, opts ...grpc.CallOption) (*KeyValueGetResponse, error)
	// Create a new or overwrite and existing key
	Put(ctx context.Context, in *KeyValuePutRequest, opts ...grpc.CallOption) (*KeyValuePutResponse, error)
	// Delete an existing
	Delete(ctx context.Context, in *KeyValueDeleteRequest, opts ...grpc.CallOption) (*KeyValueDeleteResponse, error)
}

type keyValueClient struct {
	cc grpc.ClientConnInterface
}

func NewKeyValueClient(cc grpc.ClientConnInterface) KeyValueClient {
	return &keyValueClient{cc}
}

func (c *keyValueClient) Get(ctx context.Context, in *KeyValueGetRequest, opts ...grpc.CallOption) (*KeyValueGetResponse, error) {
	out := new(KeyValueGetResponse)
	err := c.cc.Invoke(ctx, "/nitric.kv.v1.KeyValue/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueClient) Put(ctx context.Context, in *KeyValuePutRequest, opts ...grpc.CallOption) (*KeyValuePutResponse, error) {
	out := new(KeyValuePutResponse)
	err := c.cc.Invoke(ctx, "/nitric.kv.v1.KeyValue/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueClient) Delete(ctx context.Context, in *KeyValueDeleteRequest, opts ...grpc.CallOption) (*KeyValueDeleteResponse, error) {
	out := new(KeyValueDeleteResponse)
	err := c.cc.Invoke(ctx, "/nitric.kv.v1.KeyValue/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeyValueServer is the server API for KeyValue service.
// All implementations must embed UnimplementedKeyValueServer
// for forward compatibility
type KeyValueServer interface {
	// Get an existing key
	Get(context.Context, *KeyValueGetRequest) (*KeyValueGetResponse, error)
	// Create a new or overwrite and existing key
	Put(context.Context, *KeyValuePutRequest) (*KeyValuePutResponse, error)
	// Delete an existing
	Delete(context.Context, *KeyValueDeleteRequest) (*KeyValueDeleteResponse, error)
	mustEmbedUnimplementedKeyValueServer()
}

// UnimplementedKeyValueServer must be embedded to have forward compatible implementations.
type UnimplementedKeyValueServer struct {
}

func (UnimplementedKeyValueServer) Get(context.Context, *KeyValueGetRequest) (*KeyValueGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedKeyValueServer) Put(context.Context, *KeyValuePutRequest) (*KeyValuePutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedKeyValueServer) Delete(context.Context, *KeyValueDeleteRequest) (*KeyValueDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedKeyValueServer) mustEmbedUnimplementedKeyValueServer() {}

// UnsafeKeyValueServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeyValueServer will
// result in compilation errors.
type UnsafeKeyValueServer interface {
	mustEmbedUnimplementedKeyValueServer()
}

func RegisterKeyValueServer(s grpc.ServiceRegistrar, srv KeyValueServer) {
	s.RegisterService(&KeyValue_ServiceDesc, srv)
}

func _KeyValue_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValueGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.kv.v1.KeyValue/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServer).Get(ctx, req.(*KeyValueGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValue_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValuePutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.kv.v1.KeyValue/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServer).Put(ctx, req.(*KeyValuePutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValue_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValueDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.kv.v1.KeyValue/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServer).Delete(ctx, req.(*KeyValueDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeyValue_ServiceDesc is the grpc.ServiceDesc for KeyValue service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeyValue_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nitric.kv.v1.KeyValue",
	HandlerType: (*KeyValueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _KeyValue_Get_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _KeyValue_Put_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _KeyValue_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kv/v1/kv.proto",
}
