// Code generated by MockGen. DO NOT EDIT.
// Source: ./utils/config/interface.go

// Package mockutils is a generated GoMock package.
package mockutils

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockConfiguration is a mock of Configuration interface
type MockConfiguration struct {
	ctrl     *gomock.Controller
	recorder *MockConfigurationMockRecorder
}

// MockConfigurationMockRecorder is the mock recorder for MockConfiguration
type MockConfigurationMockRecorder struct {
	mock *MockConfiguration
}

// NewMockConfiguration creates a new mock instance
func NewMockConfiguration(ctrl *gomock.Controller) *MockConfiguration {
	mock := &MockConfiguration{ctrl: ctrl}
	mock.recorder = &MockConfigurationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfiguration) EXPECT() *MockConfigurationMockRecorder {
	return m.recorder
}

// InitConfigFile mocks base method
func (m *MockConfiguration) InitConfigFile(configFilePath string) error {
	ret := m.ctrl.Call(m, "InitConfigFile", configFilePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitConfigFile indicates an expected call of InitConfigFile
func (mr *MockConfigurationMockRecorder) InitConfigFile(configFilePath interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitConfigFile", reflect.TypeOf((*MockConfiguration)(nil).InitConfigFile), configFilePath)
}

// Exists mocks base method
func (m *MockConfiguration) Exists() bool {
	ret := m.ctrl.Call(m, "Exists")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Exists indicates an expected call of Exists
func (mr *MockConfigurationMockRecorder) Exists() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockConfiguration)(nil).Exists))
}

// SetDefaultAccount mocks base method
func (m *MockConfiguration) SetDefaultAccount(account string) error {
	ret := m.ctrl.Call(m, "SetDefaultAccount", account)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDefaultAccount indicates an expected call of SetDefaultAccount
func (mr *MockConfigurationMockRecorder) SetDefaultAccount(account interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDefaultAccount", reflect.TypeOf((*MockConfiguration)(nil).SetDefaultAccount), account)
}

// GetDefaultAccount mocks base method
func (m *MockConfiguration) GetDefaultAccount() string {
	ret := m.ctrl.Call(m, "GetDefaultAccount")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDefaultAccount indicates an expected call of GetDefaultAccount
func (mr *MockConfigurationMockRecorder) GetDefaultAccount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultAccount", reflect.TypeOf((*MockConfiguration)(nil).GetDefaultAccount))
}

// GetAccountDbPath mocks base method
func (m *MockConfiguration) GetAccountDbPath() string {
	ret := m.ctrl.Call(m, "GetAccountDbPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAccountDbPath indicates an expected call of GetAccountDbPath
func (mr *MockConfigurationMockRecorder) GetAccountDbPath() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountDbPath", reflect.TypeOf((*MockConfiguration)(nil).GetAccountDbPath))
}

// GetFriendBotUrl mocks base method
func (m *MockConfiguration) GetFriendBotUrl() string {
	ret := m.ctrl.Call(m, "GetFriendBotUrl")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFriendBotUrl indicates an expected call of GetFriendBotUrl
func (mr *MockConfigurationMockRecorder) GetFriendBotUrl() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriendBotUrl", reflect.TypeOf((*MockConfiguration)(nil).GetFriendBotUrl))
}

// GetHorizonUrl mocks base method
func (m *MockConfiguration) GetHorizonUrl() string {
	ret := m.ctrl.Call(m, "GetHorizonUrl")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetHorizonUrl indicates an expected call of GetHorizonUrl
func (mr *MockConfigurationMockRecorder) GetHorizonUrl() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHorizonUrl", reflect.TypeOf((*MockConfiguration)(nil).GetHorizonUrl))
}

// GetVeloNodeUrl mocks base method
func (m *MockConfiguration) GetVeloNodeUrl() string {
	ret := m.ctrl.Call(m, "GetVeloNodeUrl")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetVeloNodeUrl indicates an expected call of GetVeloNodeUrl
func (mr *MockConfigurationMockRecorder) GetVeloNodeUrl() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVeloNodeUrl", reflect.TypeOf((*MockConfiguration)(nil).GetVeloNodeUrl))
}

// GetNetworkPassphrase mocks base method
func (m *MockConfiguration) GetNetworkPassphrase() string {
	ret := m.ctrl.Call(m, "GetNetworkPassphrase")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNetworkPassphrase indicates an expected call of GetNetworkPassphrase
func (mr *MockConfigurationMockRecorder) GetNetworkPassphrase() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkPassphrase", reflect.TypeOf((*MockConfiguration)(nil).GetNetworkPassphrase))
}
