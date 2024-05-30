// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nitrictech/go-sdk/api/topics (interfaces: Topics,Topic)

// Package mockapi is a generated GoMock package.
package mockapi

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	topics "github.com/nitrictech/go-sdk/api/topics"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

// MockTopics is a mock of Topics interface.
type MockTopics struct {
	ctrl     *gomock.Controller
	recorder *MockTopicsMockRecorder
}

// MockTopicsMockRecorder is the mock recorder for MockTopics.
type MockTopicsMockRecorder struct {
	mock *MockTopics
}

// NewMockTopics creates a new mock instance.
func NewMockTopics(ctrl *gomock.Controller) *MockTopics {
	mock := &MockTopics{ctrl: ctrl}
	mock.recorder = &MockTopicsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopics) EXPECT() *MockTopicsMockRecorder {
	return m.recorder
}

// Topic mocks base method.
func (m *MockTopics) Topic(arg0 string) topics.Topic {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Topic", arg0)
	ret0, _ := ret[0].(topics.Topic)
	return ret0
}

// Topic indicates an expected call of Topic.
func (mr *MockTopicsMockRecorder) Topic(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Topic", reflect.TypeOf((*MockTopics)(nil).Topic), arg0)
}

// MockTopic is a mock of Topic interface.
type MockTopic struct {
	ctrl     *gomock.Controller
	recorder *MockTopicMockRecorder
}

// MockTopicMockRecorder is the mock recorder for MockTopic.
type MockTopicMockRecorder struct {
	mock *MockTopic
}

// NewMockTopic creates a new mock instance.
func NewMockTopic(ctrl *gomock.Controller) *MockTopic {
	mock := &MockTopic{ctrl: ctrl}
	mock.recorder = &MockTopicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopic) EXPECT() *MockTopicMockRecorder {
	return m.recorder
}

// Name mocks base method.
func (m *MockTopic) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockTopicMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockTopic)(nil).Name))
}

// Publish mocks base method.
func (m *MockTopic) Publish(arg0 context.Context, arg1 map[string]interface{}, arg2 ...func(*topicspb.TopicPublishRequest)) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Publish", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockTopicMockRecorder) Publish(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockTopic)(nil).Publish), varargs...)
}