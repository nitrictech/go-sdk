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

// QueueClient is the client API for Queue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueueClient interface {
	// Send a single event to a queue
	Send(ctx context.Context, in *QueueSendRequest, opts ...grpc.CallOption) (*QueueSendResponse, error)
	// Send multiple events to a queue
	SendBatch(ctx context.Context, in *QueueSendBatchRequest, opts ...grpc.CallOption) (*QueueSendBatchResponse, error)
	// Receive event(s) off a queue
	Receive(ctx context.Context, in *QueueReceiveRequest, opts ...grpc.CallOption) (*QueueReceiveResponse, error)
	// Complete an event previously popped from a queue
	Complete(ctx context.Context, in *QueueCompleteRequest, opts ...grpc.CallOption) (*QueueCompleteResponse, error)
}

type queueClient struct {
	cc grpc.ClientConnInterface
}

func NewQueueClient(cc grpc.ClientConnInterface) QueueClient {
	return &queueClient{cc}
}

func (c *queueClient) Send(ctx context.Context, in *QueueSendRequest, opts ...grpc.CallOption) (*QueueSendResponse, error) {
	out := new(QueueSendResponse)
	err := c.cc.Invoke(ctx, "/nitric.queue.v1.Queue/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) SendBatch(ctx context.Context, in *QueueSendBatchRequest, opts ...grpc.CallOption) (*QueueSendBatchResponse, error) {
	out := new(QueueSendBatchResponse)
	err := c.cc.Invoke(ctx, "/nitric.queue.v1.Queue/SendBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) Receive(ctx context.Context, in *QueueReceiveRequest, opts ...grpc.CallOption) (*QueueReceiveResponse, error) {
	out := new(QueueReceiveResponse)
	err := c.cc.Invoke(ctx, "/nitric.queue.v1.Queue/Receive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) Complete(ctx context.Context, in *QueueCompleteRequest, opts ...grpc.CallOption) (*QueueCompleteResponse, error) {
	out := new(QueueCompleteResponse)
	err := c.cc.Invoke(ctx, "/nitric.queue.v1.Queue/Complete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueueServer is the server API for Queue service.
// All implementations must embed UnimplementedQueueServer
// for forward compatibility
type QueueServer interface {
	// Send a single event to a queue
	Send(context.Context, *QueueSendRequest) (*QueueSendResponse, error)
	// Send multiple events to a queue
	SendBatch(context.Context, *QueueSendBatchRequest) (*QueueSendBatchResponse, error)
	// Receive event(s) off a queue
	Receive(context.Context, *QueueReceiveRequest) (*QueueReceiveResponse, error)
	// Complete an event previously popped from a queue
	Complete(context.Context, *QueueCompleteRequest) (*QueueCompleteResponse, error)
	mustEmbedUnimplementedQueueServer()
}

// UnimplementedQueueServer must be embedded to have forward compatible implementations.
type UnimplementedQueueServer struct {
}

func (UnimplementedQueueServer) Send(context.Context, *QueueSendRequest) (*QueueSendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedQueueServer) SendBatch(context.Context, *QueueSendBatchRequest) (*QueueSendBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendBatch not implemented")
}
func (UnimplementedQueueServer) Receive(context.Context, *QueueReceiveRequest) (*QueueReceiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Receive not implemented")
}
func (UnimplementedQueueServer) Complete(context.Context, *QueueCompleteRequest) (*QueueCompleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Complete not implemented")
}
func (UnimplementedQueueServer) mustEmbedUnimplementedQueueServer() {}

// UnsafeQueueServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueueServer will
// result in compilation errors.
type UnsafeQueueServer interface {
	mustEmbedUnimplementedQueueServer()
}

func RegisterQueueServer(s grpc.ServiceRegistrar, srv QueueServer) {
	s.RegisterService(&Queue_ServiceDesc, srv)
}

func _Queue_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueSendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.queue.v1.Queue/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).Send(ctx, req.(*QueueSendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_SendBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueSendBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).SendBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.queue.v1.Queue/SendBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).SendBatch(ctx, req.(*QueueSendBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_Receive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueReceiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).Receive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.queue.v1.Queue/Receive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).Receive(ctx, req.(*QueueReceiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_Complete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueCompleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).Complete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.queue.v1.Queue/Complete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).Complete(ctx, req.(*QueueCompleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Queue_ServiceDesc is the grpc.ServiceDesc for Queue service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Queue_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nitric.queue.v1.Queue",
	HandlerType: (*QueueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _Queue_Send_Handler,
		},
		{
			MethodName: "SendBatch",
			Handler:    _Queue_SendBatch_Handler,
		},
		{
			MethodName: "Receive",
			Handler:    _Queue_Receive_Handler,
		},
		{
			MethodName: "Complete",
			Handler:    _Queue_Complete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "queue/v1/queue.proto",
}
