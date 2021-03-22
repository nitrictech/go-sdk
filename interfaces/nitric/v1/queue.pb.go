// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: queue/v1/queue.proto

package v1

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Request to push a single event to a queue
type QueueSendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The Nitric name for the queue
	//  this will automatically be resolved to the provider specific queue identifier.
	Queue string `protobuf:"bytes,1,opt,name=queue,proto3" json:"queue,omitempty"`
	// The event to push to the queue
	Event *NitricEvent `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
}

func (x *QueueSendRequest) Reset() {
	*x = QueueSendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueSendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueSendRequest) ProtoMessage() {}

func (x *QueueSendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueSendRequest.ProtoReflect.Descriptor instead.
func (*QueueSendRequest) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{0}
}

func (x *QueueSendRequest) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

func (x *QueueSendRequest) GetEvent() *NitricEvent {
	if x != nil {
		return x.Event
	}
	return nil
}

// Result of pushing an event to a queue
type QueueSendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *QueueSendResponse) Reset() {
	*x = QueueSendResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueSendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueSendResponse) ProtoMessage() {}

func (x *QueueSendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueSendResponse.ProtoReflect.Descriptor instead.
func (*QueueSendResponse) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{1}
}

type QueueSendBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The Nitric name for the queue
	//  this will automatically be resolved to the provider specific queue identifier.
	Queue string `protobuf:"bytes,1,opt,name=queue,proto3" json:"queue,omitempty"`
	// Array of events to push to the queue
	Events []*NitricEvent `protobuf:"bytes,2,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *QueueSendBatchRequest) Reset() {
	*x = QueueSendBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueSendBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueSendBatchRequest) ProtoMessage() {}

func (x *QueueSendBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueSendBatchRequest.ProtoReflect.Descriptor instead.
func (*QueueSendBatchRequest) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{2}
}

func (x *QueueSendBatchRequest) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

func (x *QueueSendBatchRequest) GetEvents() []*NitricEvent {
	if x != nil {
		return x.Events
	}
	return nil
}

// An ordered array of booleans
// matching the same order as the events given
// in the original request, each one will mark if the
// Event was successful pushed
type QueueSendBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FailedEvents []*FailedEvent `protobuf:"bytes,1,rep,name=failedEvents,proto3" json:"failedEvents,omitempty"`
}

func (x *QueueSendBatchResponse) Reset() {
	*x = QueueSendBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueSendBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueSendBatchResponse) ProtoMessage() {}

func (x *QueueSendBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueSendBatchResponse.ProtoReflect.Descriptor instead.
func (*QueueSendBatchResponse) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{3}
}

func (x *QueueSendBatchResponse) GetFailedEvents() []*FailedEvent {
	if x != nil {
		return x.FailedEvents
	}
	return nil
}

type QueueReceiveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The nitric name for the queue
	//  this will automatically be resolved to the provider specific queue identifier.
	Queue string `protobuf:"bytes,1,opt,name=queue,proto3" json:"queue,omitempty"`
	// The max number of items to pop off the queue, may be capped by provider specific limitations
	Depth int32 `protobuf:"varint,2,opt,name=depth,proto3" json:"depth,omitempty"`
}

func (x *QueueReceiveRequest) Reset() {
	*x = QueueReceiveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueReceiveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueReceiveRequest) ProtoMessage() {}

func (x *QueueReceiveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueReceiveRequest.ProtoReflect.Descriptor instead.
func (*QueueReceiveRequest) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{4}
}

func (x *QueueReceiveRequest) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

func (x *QueueReceiveRequest) GetDepth() int32 {
	if x != nil {
		return x.Depth
	}
	return 0
}

type QueueReceiveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Array of items popped off the queue
	Items []*NitricQueueItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *QueueReceiveResponse) Reset() {
	*x = QueueReceiveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueReceiveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueReceiveResponse) ProtoMessage() {}

func (x *QueueReceiveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueReceiveResponse.ProtoReflect.Descriptor instead.
func (*QueueReceiveResponse) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{5}
}

func (x *QueueReceiveResponse) GetItems() []*NitricQueueItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type QueueCompleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The nitric name for the queue
	//  this will automatically be resolved to the provider specific queue identifier.
	Queue string `protobuf:"bytes,1,opt,name=queue,proto3" json:"queue,omitempty"`
	// Lease id of the event to be completed
	LeaseId string `protobuf:"bytes,2,opt,name=leaseId,proto3" json:"leaseId,omitempty"`
}

func (x *QueueCompleteRequest) Reset() {
	*x = QueueCompleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueCompleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueCompleteRequest) ProtoMessage() {}

func (x *QueueCompleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueCompleteRequest.ProtoReflect.Descriptor instead.
func (*QueueCompleteRequest) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{6}
}

func (x *QueueCompleteRequest) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

func (x *QueueCompleteRequest) GetLeaseId() string {
	if x != nil {
		return x.LeaseId
	}
	return ""
}

type QueueCompleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *QueueCompleteResponse) Reset() {
	*x = QueueCompleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueCompleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueCompleteResponse) ProtoMessage() {}

func (x *QueueCompleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueCompleteResponse.ProtoReflect.Descriptor instead.
func (*QueueCompleteResponse) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{7}
}

type FailedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The event that failed to be pushed
	Event *NitricEvent `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
	// A message describing the failure
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *FailedEvent) Reset() {
	*x = FailedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FailedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailedEvent) ProtoMessage() {}

func (x *FailedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailedEvent.ProtoReflect.Descriptor instead.
func (*FailedEvent) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{8}
}

func (x *FailedEvent) GetEvent() *NitricEvent {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *FailedEvent) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// A leased event, which must be completed or returned to the queue
type NitricQueueItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The event popped from the queue
	Event *NitricEvent `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
	// The lease id unique to the pop request, this must be used to complete, extend the lease or release the event.
	LeaseId string `protobuf:"bytes,2,opt,name=leaseId,proto3" json:"leaseId,omitempty"`
}

func (x *NitricQueueItem) Reset() {
	*x = NitricQueueItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_v1_queue_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NitricQueueItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NitricQueueItem) ProtoMessage() {}

func (x *NitricQueueItem) ProtoReflect() protoreflect.Message {
	mi := &file_queue_v1_queue_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NitricQueueItem.ProtoReflect.Descriptor instead.
func (*NitricQueueItem) Descriptor() ([]byte, []int) {
	return file_queue_v1_queue_proto_rawDescGZIP(), []int{9}
}

func (x *NitricQueueItem) GetEvent() *NitricEvent {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *NitricQueueItem) GetLeaseId() string {
	if x != nil {
		return x.LeaseId
	}
	return ""
}

var File_queue_v1_queue_proto protoreflect.FileDescriptor

var file_queue_v1_queue_proto_rawDesc = []byte{
	0x0a, 0x14, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x16, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x5d, 0x0a, 0x10, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x33, 0x0a, 0x05, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69,
	0x63, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x69, 0x74, 0x72,
	0x69, 0x63, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x13,
	0x0a, 0x11, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x64, 0x0a, 0x15, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65, 0x6e, 0x64,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65,
	0x75, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x5a, 0x0a, 0x16, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x0c, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x69, 0x74, 0x72,
	0x69, 0x63, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x69, 0x6c,
	0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x0c, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x41, 0x0a, 0x13, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65,
	0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x22, 0x4e, 0x0a, 0x14, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x36, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x4e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x51, 0x75, 0x65, 0x75, 0x65, 0x49, 0x74, 0x65,
	0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x46, 0x0a, 0x14, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x49, 0x64,
	0x22, 0x17, 0x0a, 0x15, 0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x5c, 0x0a, 0x0b, 0x46, 0x61, 0x69,
	0x6c, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x33, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x69, 0x74, 0x72, 0x69,
	0x63, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x60, 0x0a, 0x0f, 0x4e, 0x69, 0x74, 0x72, 0x69,
	0x63, 0x51, 0x75, 0x65, 0x75, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x33, 0x0a, 0x05, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6e, 0x69, 0x74, 0x72,
	0x69, 0x63, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x69, 0x74,
	0x72, 0x69, 0x63, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x49, 0x64, 0x32, 0xe7, 0x02, 0x0a, 0x05, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x12, 0x4d, 0x0a, 0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x21, 0x2e, 0x6e, 0x69,
	0x74, 0x72, 0x69, 0x63, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22,
	0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x5c, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12,
	0x26, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53,
	0x65, 0x6e, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x56, 0x0a, 0x07, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x12, 0x24, 0x2e, 0x6e, 0x69,
	0x74, 0x72, 0x69, 0x63, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x25, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x59, 0x0a, 0x08, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x12, 0x25, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x6e, 0x69,
	0x74, 0x72, 0x69, 0x63, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x49, 0x0a, 0x18, 0x69, 0x6f, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x42,
	0x05, 0x51, 0x75, 0x65, 0x75, 0x65, 0x50, 0x01, 0x5a, 0x0c, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63,
	0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0xca, 0x02, 0x15, 0x4e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x5c,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x51, 0x75, 0x65, 0x75, 0x65, 0x5c, 0x56, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_queue_v1_queue_proto_rawDescOnce sync.Once
	file_queue_v1_queue_proto_rawDescData = file_queue_v1_queue_proto_rawDesc
)

func file_queue_v1_queue_proto_rawDescGZIP() []byte {
	file_queue_v1_queue_proto_rawDescOnce.Do(func() {
		file_queue_v1_queue_proto_rawDescData = protoimpl.X.CompressGZIP(file_queue_v1_queue_proto_rawDescData)
	})
	return file_queue_v1_queue_proto_rawDescData
}

var file_queue_v1_queue_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_queue_v1_queue_proto_goTypes = []interface{}{
	(*QueueSendRequest)(nil),       // 0: nitric.queue.v1.QueueSendRequest
	(*QueueSendResponse)(nil),      // 1: nitric.queue.v1.QueueSendResponse
	(*QueueSendBatchRequest)(nil),  // 2: nitric.queue.v1.QueueSendBatchRequest
	(*QueueSendBatchResponse)(nil), // 3: nitric.queue.v1.QueueSendBatchResponse
	(*QueueReceiveRequest)(nil),    // 4: nitric.queue.v1.QueueReceiveRequest
	(*QueueReceiveResponse)(nil),   // 5: nitric.queue.v1.QueueReceiveResponse
	(*QueueCompleteRequest)(nil),   // 6: nitric.queue.v1.QueueCompleteRequest
	(*QueueCompleteResponse)(nil),  // 7: nitric.queue.v1.QueueCompleteResponse
	(*FailedEvent)(nil),            // 8: nitric.queue.v1.FailedEvent
	(*NitricQueueItem)(nil),        // 9: nitric.queue.v1.NitricQueueItem
	(*NitricEvent)(nil),            // 10: nitric.common.v1.NitricEvent
}
var file_queue_v1_queue_proto_depIdxs = []int32{
	10, // 0: nitric.queue.v1.QueueSendRequest.event:type_name -> nitric.common.v1.NitricEvent
	10, // 1: nitric.queue.v1.QueueSendBatchRequest.events:type_name -> nitric.common.v1.NitricEvent
	8,  // 2: nitric.queue.v1.QueueSendBatchResponse.failedEvents:type_name -> nitric.queue.v1.FailedEvent
	9,  // 3: nitric.queue.v1.QueueReceiveResponse.items:type_name -> nitric.queue.v1.NitricQueueItem
	10, // 4: nitric.queue.v1.FailedEvent.event:type_name -> nitric.common.v1.NitricEvent
	10, // 5: nitric.queue.v1.NitricQueueItem.event:type_name -> nitric.common.v1.NitricEvent
	0,  // 6: nitric.queue.v1.Queue.Send:input_type -> nitric.queue.v1.QueueSendRequest
	2,  // 7: nitric.queue.v1.Queue.SendBatch:input_type -> nitric.queue.v1.QueueSendBatchRequest
	4,  // 8: nitric.queue.v1.Queue.Receive:input_type -> nitric.queue.v1.QueueReceiveRequest
	6,  // 9: nitric.queue.v1.Queue.Complete:input_type -> nitric.queue.v1.QueueCompleteRequest
	1,  // 10: nitric.queue.v1.Queue.Send:output_type -> nitric.queue.v1.QueueSendResponse
	3,  // 11: nitric.queue.v1.Queue.SendBatch:output_type -> nitric.queue.v1.QueueSendBatchResponse
	5,  // 12: nitric.queue.v1.Queue.Receive:output_type -> nitric.queue.v1.QueueReceiveResponse
	7,  // 13: nitric.queue.v1.Queue.Complete:output_type -> nitric.queue.v1.QueueCompleteResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_queue_v1_queue_proto_init() }
func file_queue_v1_queue_proto_init() {
	if File_queue_v1_queue_proto != nil {
		return
	}
	file_common_v1_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_queue_v1_queue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueSendRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_queue_v1_queue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueSendResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_queue_v1_queue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueSendBatchRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_queue_v1_queue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueSendBatchResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_queue_v1_queue_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueReceiveRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_queue_v1_queue_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueReceiveResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_queue_v1_queue_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueCompleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_queue_v1_queue_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueCompleteResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_queue_v1_queue_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FailedEvent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_queue_v1_queue_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NitricQueueItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_queue_v1_queue_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_queue_v1_queue_proto_goTypes,
		DependencyIndexes: file_queue_v1_queue_proto_depIdxs,
		MessageInfos:      file_queue_v1_queue_proto_msgTypes,
	}.Build()
	File_queue_v1_queue_proto = out.File
	file_queue_v1_queue_proto_rawDesc = nil
	file_queue_v1_queue_proto_goTypes = nil
	file_queue_v1_queue_proto_depIdxs = nil
}
