// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nitrictech/go-sdk/api/keyvalue (interfaces: KeyValue,Store)

// Package mockapi is a generated GoMock package.
package mockapi

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	keyvalue "github.com/nitrictech/go-sdk/api/keyvalue"
	kvstorepb "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
)

// MockKeyValue is a mock of KeyValue interface.
type MockKeyValue struct {
	ctrl     *gomock.Controller
	recorder *MockKeyValueMockRecorder
}

// MockKeyValueMockRecorder is the mock recorder for MockKeyValue.
type MockKeyValueMockRecorder struct {
	mock *MockKeyValue
}

// NewMockKeyValue creates a new mock instance.
func NewMockKeyValue(ctrl *gomock.Controller) *MockKeyValue {
	mock := &MockKeyValue{ctrl: ctrl}
	mock.recorder = &MockKeyValueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyValue) EXPECT() *MockKeyValueMockRecorder {
	return m.recorder
}

// Store mocks base method.
func (m *MockKeyValue) Store(arg0 string) keyvalue.Store {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0)
	ret0, _ := ret[0].(keyvalue.Store)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockKeyValueMockRecorder) Store(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockKeyValue)(nil).Store), arg0)
}

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockStore) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStore)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockStore) Get(arg0 context.Context, arg1 string) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStoreMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), arg0, arg1)
}

// Keys mocks base method.
func (m *MockStore) Keys(arg0 context.Context, arg1 ...func(*kvstorepb.KvStoreScanKeysRequest)) (*keyvalue.KeyStream, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Keys", varargs...)
	ret0, _ := ret[0].(*keyvalue.KeyStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Keys indicates an expected call of Keys.
func (mr *MockStoreMockRecorder) Keys(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Keys", reflect.TypeOf((*MockStore)(nil).Keys), varargs...)
}

// Name mocks base method.
func (m *MockStore) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockStoreMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockStore)(nil).Name))
}

// Set mocks base method.
func (m *MockStore) Set(arg0 context.Context, arg1 string, arg2 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockStoreMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStore)(nil).Set), arg0, arg1, arg2)
}