// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nitrictech/apis/go/nitric/v1 (interfaces: DocumentServiceServer,EventServiceServer,TopicServiceServer,QueueServiceServer,StorageServiceServer,FaasServiceServer,FaasService_TriggerStreamServer,DocumentService_QueryStreamServer,SecretServiceServer)

// Package mock_v1 is a generated GoMock package.
package mock_v1

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/nitrictech/apis/go/nitric/v1"
	metadata "google.golang.org/grpc/metadata"
)

// MockDocumentServiceServer is a mock of DocumentServiceServer interface.
type MockDocumentServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockDocumentServiceServerMockRecorder
}

// MockDocumentServiceServerMockRecorder is the mock recorder for MockDocumentServiceServer.
type MockDocumentServiceServerMockRecorder struct {
	mock *MockDocumentServiceServer
}

// NewMockDocumentServiceServer creates a new mock instance.
func NewMockDocumentServiceServer(ctrl *gomock.Controller) *MockDocumentServiceServer {
	mock := &MockDocumentServiceServer{ctrl: ctrl}
	mock.recorder = &MockDocumentServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDocumentServiceServer) EXPECT() *MockDocumentServiceServerMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDocumentServiceServer) Delete(arg0 context.Context, arg1 *v1.DocumentDeleteRequest) (*v1.DocumentDeleteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*v1.DocumentDeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockDocumentServiceServerMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDocumentServiceServer)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockDocumentServiceServer) Get(arg0 context.Context, arg1 *v1.DocumentGetRequest) (*v1.DocumentGetResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*v1.DocumentGetResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDocumentServiceServerMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDocumentServiceServer)(nil).Get), arg0, arg1)
}

// Query mocks base method.
func (m *MockDocumentServiceServer) Query(arg0 context.Context, arg1 *v1.DocumentQueryRequest) (*v1.DocumentQueryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0, arg1)
	ret0, _ := ret[0].(*v1.DocumentQueryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockDocumentServiceServerMockRecorder) Query(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDocumentServiceServer)(nil).Query), arg0, arg1)
}

// QueryStream mocks base method.
func (m *MockDocumentServiceServer) QueryStream(arg0 *v1.DocumentQueryStreamRequest, arg1 v1.DocumentService_QueryStreamServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryStream", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueryStream indicates an expected call of QueryStream.
func (mr *MockDocumentServiceServerMockRecorder) QueryStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryStream", reflect.TypeOf((*MockDocumentServiceServer)(nil).QueryStream), arg0, arg1)
}

// Set mocks base method.
func (m *MockDocumentServiceServer) Set(arg0 context.Context, arg1 *v1.DocumentSetRequest) (*v1.DocumentSetResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1)
	ret0, _ := ret[0].(*v1.DocumentSetResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockDocumentServiceServerMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockDocumentServiceServer)(nil).Set), arg0, arg1)
}

// MockEventServiceServer is a mock of EventServiceServer interface.
type MockEventServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockEventServiceServerMockRecorder
}

// MockEventServiceServerMockRecorder is the mock recorder for MockEventServiceServer.
type MockEventServiceServerMockRecorder struct {
	mock *MockEventServiceServer
}

// NewMockEventServiceServer creates a new mock instance.
func NewMockEventServiceServer(ctrl *gomock.Controller) *MockEventServiceServer {
	mock := &MockEventServiceServer{ctrl: ctrl}
	mock.recorder = &MockEventServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventServiceServer) EXPECT() *MockEventServiceServerMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockEventServiceServer) Publish(arg0 context.Context, arg1 *v1.EventPublishRequest) (*v1.EventPublishResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(*v1.EventPublishResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockEventServiceServerMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockEventServiceServer)(nil).Publish), arg0, arg1)
}

// MockTopicServiceServer is a mock of TopicServiceServer interface.
type MockTopicServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockTopicServiceServerMockRecorder
}

// MockTopicServiceServerMockRecorder is the mock recorder for MockTopicServiceServer.
type MockTopicServiceServerMockRecorder struct {
	mock *MockTopicServiceServer
}

// NewMockTopicServiceServer creates a new mock instance.
func NewMockTopicServiceServer(ctrl *gomock.Controller) *MockTopicServiceServer {
	mock := &MockTopicServiceServer{ctrl: ctrl}
	mock.recorder = &MockTopicServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopicServiceServer) EXPECT() *MockTopicServiceServerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockTopicServiceServer) List(arg0 context.Context, arg1 *v1.TopicListRequest) (*v1.TopicListResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*v1.TopicListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTopicServiceServerMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTopicServiceServer)(nil).List), arg0, arg1)
}

// MockQueueServiceServer is a mock of QueueServiceServer interface.
type MockQueueServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockQueueServiceServerMockRecorder
}

// MockQueueServiceServerMockRecorder is the mock recorder for MockQueueServiceServer.
type MockQueueServiceServerMockRecorder struct {
	mock *MockQueueServiceServer
}

// NewMockQueueServiceServer creates a new mock instance.
func NewMockQueueServiceServer(ctrl *gomock.Controller) *MockQueueServiceServer {
	mock := &MockQueueServiceServer{ctrl: ctrl}
	mock.recorder = &MockQueueServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueueServiceServer) EXPECT() *MockQueueServiceServerMockRecorder {
	return m.recorder
}

// Complete mocks base method.
func (m *MockQueueServiceServer) Complete(arg0 context.Context, arg1 *v1.QueueCompleteRequest) (*v1.QueueCompleteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Complete", arg0, arg1)
	ret0, _ := ret[0].(*v1.QueueCompleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Complete indicates an expected call of Complete.
func (mr *MockQueueServiceServerMockRecorder) Complete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Complete", reflect.TypeOf((*MockQueueServiceServer)(nil).Complete), arg0, arg1)
}

// Receive mocks base method.
func (m *MockQueueServiceServer) Receive(arg0 context.Context, arg1 *v1.QueueReceiveRequest) (*v1.QueueReceiveResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Receive", arg0, arg1)
	ret0, _ := ret[0].(*v1.QueueReceiveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Receive indicates an expected call of Receive.
func (mr *MockQueueServiceServerMockRecorder) Receive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Receive", reflect.TypeOf((*MockQueueServiceServer)(nil).Receive), arg0, arg1)
}

// Send mocks base method.
func (m *MockQueueServiceServer) Send(arg0 context.Context, arg1 *v1.QueueSendRequest) (*v1.QueueSendResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1)
	ret0, _ := ret[0].(*v1.QueueSendResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockQueueServiceServerMockRecorder) Send(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockQueueServiceServer)(nil).Send), arg0, arg1)
}

// SendBatch mocks base method.
func (m *MockQueueServiceServer) SendBatch(arg0 context.Context, arg1 *v1.QueueSendBatchRequest) (*v1.QueueSendBatchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendBatch", arg0, arg1)
	ret0, _ := ret[0].(*v1.QueueSendBatchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendBatch indicates an expected call of SendBatch.
func (mr *MockQueueServiceServerMockRecorder) SendBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendBatch", reflect.TypeOf((*MockQueueServiceServer)(nil).SendBatch), arg0, arg1)
}

// MockStorageServiceServer is a mock of StorageServiceServer interface.
type MockStorageServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockStorageServiceServerMockRecorder
}

// MockStorageServiceServerMockRecorder is the mock recorder for MockStorageServiceServer.
type MockStorageServiceServerMockRecorder struct {
	mock *MockStorageServiceServer
}

// NewMockStorageServiceServer creates a new mock instance.
func NewMockStorageServiceServer(ctrl *gomock.Controller) *MockStorageServiceServer {
	mock := &MockStorageServiceServer{ctrl: ctrl}
	mock.recorder = &MockStorageServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageServiceServer) EXPECT() *MockStorageServiceServerMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockStorageServiceServer) Delete(arg0 context.Context, arg1 *v1.StorageDeleteRequest) (*v1.StorageDeleteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*v1.StorageDeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockStorageServiceServerMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorageServiceServer)(nil).Delete), arg0, arg1)
}

// PreSignUrl mocks base method.
func (m *MockStorageServiceServer) PreSignUrl(arg0 context.Context, arg1 *v1.StoragePreSignUrlRequest) (*v1.StoragePreSignUrlResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreSignUrl", arg0, arg1)
	ret0, _ := ret[0].(*v1.StoragePreSignUrlResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PreSignUrl indicates an expected call of PreSignUrl.
func (mr *MockStorageServiceServerMockRecorder) PreSignUrl(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreSignUrl", reflect.TypeOf((*MockStorageServiceServer)(nil).PreSignUrl), arg0, arg1)
}

// Read mocks base method.
func (m *MockStorageServiceServer) Read(arg0 context.Context, arg1 *v1.StorageReadRequest) (*v1.StorageReadResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0, arg1)
	ret0, _ := ret[0].(*v1.StorageReadResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockStorageServiceServerMockRecorder) Read(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockStorageServiceServer)(nil).Read), arg0, arg1)
}

// Write mocks base method.
func (m *MockStorageServiceServer) Write(arg0 context.Context, arg1 *v1.StorageWriteRequest) (*v1.StorageWriteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0, arg1)
	ret0, _ := ret[0].(*v1.StorageWriteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockStorageServiceServerMockRecorder) Write(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockStorageServiceServer)(nil).Write), arg0, arg1)
}

// MockFaasServiceServer is a mock of FaasServiceServer interface.
type MockFaasServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockFaasServiceServerMockRecorder
}

// MockFaasServiceServerMockRecorder is the mock recorder for MockFaasServiceServer.
type MockFaasServiceServerMockRecorder struct {
	mock *MockFaasServiceServer
}

// NewMockFaasServiceServer creates a new mock instance.
func NewMockFaasServiceServer(ctrl *gomock.Controller) *MockFaasServiceServer {
	mock := &MockFaasServiceServer{ctrl: ctrl}
	mock.recorder = &MockFaasServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFaasServiceServer) EXPECT() *MockFaasServiceServerMockRecorder {
	return m.recorder
}

// TriggerStream mocks base method.
func (m *MockFaasServiceServer) TriggerStream(arg0 v1.FaasService_TriggerStreamServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TriggerStream", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// TriggerStream indicates an expected call of TriggerStream.
func (mr *MockFaasServiceServerMockRecorder) TriggerStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TriggerStream", reflect.TypeOf((*MockFaasServiceServer)(nil).TriggerStream), arg0)
}

// MockFaasService_TriggerStreamServer is a mock of FaasService_TriggerStreamServer interface.
type MockFaasService_TriggerStreamServer struct {
	ctrl     *gomock.Controller
	recorder *MockFaasService_TriggerStreamServerMockRecorder
}

// MockFaasService_TriggerStreamServerMockRecorder is the mock recorder for MockFaasService_TriggerStreamServer.
type MockFaasService_TriggerStreamServerMockRecorder struct {
	mock *MockFaasService_TriggerStreamServer
}

// NewMockFaasService_TriggerStreamServer creates a new mock instance.
func NewMockFaasService_TriggerStreamServer(ctrl *gomock.Controller) *MockFaasService_TriggerStreamServer {
	mock := &MockFaasService_TriggerStreamServer{ctrl: ctrl}
	mock.recorder = &MockFaasService_TriggerStreamServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFaasService_TriggerStreamServer) EXPECT() *MockFaasService_TriggerStreamServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockFaasService_TriggerStreamServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockFaasService_TriggerStreamServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockFaasService_TriggerStreamServer)(nil).Context))
}

// Recv mocks base method.
func (m *MockFaasService_TriggerStreamServer) Recv() (*v1.ClientMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*v1.ClientMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockFaasService_TriggerStreamServerMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockFaasService_TriggerStreamServer)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockFaasService_TriggerStreamServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockFaasService_TriggerStreamServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockFaasService_TriggerStreamServer)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockFaasService_TriggerStreamServer) Send(arg0 *v1.ServerMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockFaasService_TriggerStreamServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockFaasService_TriggerStreamServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockFaasService_TriggerStreamServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockFaasService_TriggerStreamServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockFaasService_TriggerStreamServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockFaasService_TriggerStreamServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockFaasService_TriggerStreamServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockFaasService_TriggerStreamServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockFaasService_TriggerStreamServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockFaasService_TriggerStreamServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockFaasService_TriggerStreamServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockFaasService_TriggerStreamServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockFaasService_TriggerStreamServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockFaasService_TriggerStreamServer)(nil).SetTrailer), arg0)
}

// MockDocumentService_QueryStreamServer is a mock of DocumentService_QueryStreamServer interface.
type MockDocumentService_QueryStreamServer struct {
	ctrl     *gomock.Controller
	recorder *MockDocumentService_QueryStreamServerMockRecorder
}

// MockDocumentService_QueryStreamServerMockRecorder is the mock recorder for MockDocumentService_QueryStreamServer.
type MockDocumentService_QueryStreamServerMockRecorder struct {
	mock *MockDocumentService_QueryStreamServer
}

// NewMockDocumentService_QueryStreamServer creates a new mock instance.
func NewMockDocumentService_QueryStreamServer(ctrl *gomock.Controller) *MockDocumentService_QueryStreamServer {
	mock := &MockDocumentService_QueryStreamServer{ctrl: ctrl}
	mock.recorder = &MockDocumentService_QueryStreamServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDocumentService_QueryStreamServer) EXPECT() *MockDocumentService_QueryStreamServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockDocumentService_QueryStreamServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDocumentService_QueryStreamServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDocumentService_QueryStreamServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m *MockDocumentService_QueryStreamServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDocumentService_QueryStreamServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDocumentService_QueryStreamServer)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockDocumentService_QueryStreamServer) Send(arg0 *v1.DocumentQueryStreamResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockDocumentService_QueryStreamServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDocumentService_QueryStreamServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockDocumentService_QueryStreamServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockDocumentService_QueryStreamServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockDocumentService_QueryStreamServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockDocumentService_QueryStreamServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDocumentService_QueryStreamServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDocumentService_QueryStreamServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockDocumentService_QueryStreamServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockDocumentService_QueryStreamServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockDocumentService_QueryStreamServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockDocumentService_QueryStreamServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockDocumentService_QueryStreamServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockDocumentService_QueryStreamServer)(nil).SetTrailer), arg0)
}

// MockSecretServiceServer is a mock of SecretServiceServer interface.
type MockSecretServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockSecretServiceServerMockRecorder
}

// MockSecretServiceServerMockRecorder is the mock recorder for MockSecretServiceServer.
type MockSecretServiceServerMockRecorder struct {
	mock *MockSecretServiceServer
}

// NewMockSecretServiceServer creates a new mock instance.
func NewMockSecretServiceServer(ctrl *gomock.Controller) *MockSecretServiceServer {
	mock := &MockSecretServiceServer{ctrl: ctrl}
	mock.recorder = &MockSecretServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretServiceServer) EXPECT() *MockSecretServiceServerMockRecorder {
	return m.recorder
}

// Access mocks base method.
func (m *MockSecretServiceServer) Access(arg0 context.Context, arg1 *v1.SecretAccessRequest) (*v1.SecretAccessResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Access", arg0, arg1)
	ret0, _ := ret[0].(*v1.SecretAccessResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Access indicates an expected call of Access.
func (mr *MockSecretServiceServerMockRecorder) Access(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Access", reflect.TypeOf((*MockSecretServiceServer)(nil).Access), arg0, arg1)
}

// Put mocks base method.
func (m *MockSecretServiceServer) Put(arg0 context.Context, arg1 *v1.SecretPutRequest) (*v1.SecretPutResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", arg0, arg1)
	ret0, _ := ret[0].(*v1.SecretPutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put.
func (mr *MockSecretServiceServerMockRecorder) Put(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockSecretServiceServer)(nil).Put), arg0, arg1)
}