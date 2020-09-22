// Code generated by MockGen. DO NOT EDIT.
// Source: msg.go

// Package msgmock is a generated GoMock package.
package msgmock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	msg "test/msg"
)

// MockMessages is a mock of Messages interface
type MockMessages struct {
	ctrl     *gomock.Controller
	recorder *MockMessagesMockRecorder
}

// MockMessagesMockRecorder is the mock recorder for MockMessages
type MockMessagesMockRecorder struct {
	mock *MockMessages
}

// NewMockMessages creates a new mock instance
func NewMockMessages(ctrl *gomock.Controller) *MockMessages {
	mock := &MockMessages{ctrl: ctrl}
	mock.recorder = &MockMessagesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMessages) EXPECT() *MockMessagesMockRecorder {
	return m.recorder
}

// GetUnread mocks base method
func (m *MockMessages) GetUnread() ([]msg.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnread")
	ret0, _ := ret[0].([]msg.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnread indicates an expected call of GetUnread
func (mr *MockMessagesMockRecorder) GetUnread() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnread", reflect.TypeOf((*MockMessages)(nil).GetUnread))
}

// GetIndex mocks base method
func (m *MockMessages) GetIndex(id int) (*msg.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIndex", id)
	ret0, _ := ret[0].(*msg.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIndex indicates an expected call of GetIndex
func (mr *MockMessagesMockRecorder) GetIndex(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIndex", reflect.TypeOf((*MockMessages)(nil).GetIndex), id)
}
