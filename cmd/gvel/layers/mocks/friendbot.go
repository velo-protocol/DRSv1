// Code generated by MockGen. DO NOT EDIT.
// Source: ./layers/repositories/friendbot/interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockFriendBotRepo is a mock of Repository interface
type MockFriendBotRepo struct {
	ctrl     *gomock.Controller
	recorder *MockFriendBotRepoMockRecorder
}

// MockFriendBotRepoMockRecorder is the mock recorder for MockFriendBotRepo
type MockFriendBotRepoMockRecorder struct {
	mock *MockFriendBotRepo
}

// NewMockFriendBotRepo creates a new mock instance
func NewMockFriendBotRepo(ctrl *gomock.Controller) *MockFriendBotRepo {
	mock := &MockFriendBotRepo{ctrl: ctrl}
	mock.recorder = &MockFriendBotRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFriendBotRepo) EXPECT() *MockFriendBotRepoMockRecorder {
	return m.recorder
}

// GetFreeLumens mocks base method
func (m *MockFriendBotRepo) GetFreeLumens(stellarAddress string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFreeLumens", stellarAddress)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetFreeLumens indicates an expected call of GetFreeLumens
func (mr *MockFriendBotRepoMockRecorder) GetFreeLumens(stellarAddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFreeLumens", reflect.TypeOf((*MockFriendBotRepo)(nil).GetFreeLumens), stellarAddress)
}