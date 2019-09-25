// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dapperlabs/flow-go/sdk/client (interfaces: RPCClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	observe "github.com/dapperlabs/flow-go/pkg/grpc/services/observe"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockRPCClient is a mock of RPCClient interface
type MockRPCClient struct {
	ctrl     *gomock.Controller
	recorder *MockRPCClientMockRecorder
}

// MockRPCClientMockRecorder is the mock recorder for MockRPCClient
type MockRPCClientMockRecorder struct {
	mock *MockRPCClient
}

// NewMockRPCClient creates a new mock instance
func NewMockRPCClient(ctrl *gomock.Controller) *MockRPCClient {
	mock := &MockRPCClient{ctrl: ctrl}
	mock.recorder = &MockRPCClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRPCClient) EXPECT() *MockRPCClientMockRecorder {
	return m.recorder
}

// CallScript mocks base method
func (m *MockRPCClient) CallScript(arg0 context.Context, arg1 *observe.CallScriptRequest, arg2 ...grpc.CallOption) (*observe.CallScriptResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CallScript", varargs...)
	ret0, _ := ret[0].(*observe.CallScriptResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CallScript indicates an expected call of CallScript
func (mr *MockRPCClientMockRecorder) CallScript(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CallScript", reflect.TypeOf((*MockRPCClient)(nil).CallScript), varargs...)
}

// GetAccount mocks base method
func (m *MockRPCClient) GetAccount(arg0 context.Context, arg1 *observe.GetAccountRequest, arg2 ...grpc.CallOption) (*observe.GetAccountResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccount", varargs...)
	ret0, _ := ret[0].(*observe.GetAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockRPCClientMockRecorder) GetAccount(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockRPCClient)(nil).GetAccount), varargs...)
}

// GetBlockByHash mocks base method
func (m *MockRPCClient) GetBlockByHash(arg0 context.Context, arg1 *observe.GetBlockByHashRequest, arg2 ...grpc.CallOption) (*observe.GetBlockByHashResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetBlockByHash", varargs...)
	ret0, _ := ret[0].(*observe.GetBlockByHashResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByHash indicates an expected call of GetBlockByHash
func (mr *MockRPCClientMockRecorder) GetBlockByHash(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByHash", reflect.TypeOf((*MockRPCClient)(nil).GetBlockByHash), varargs...)
}

// GetBlockByNumber mocks base method
func (m *MockRPCClient) GetBlockByNumber(arg0 context.Context, arg1 *observe.GetBlockByNumberRequest, arg2 ...grpc.CallOption) (*observe.GetBlockByNumberResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetBlockByNumber", varargs...)
	ret0, _ := ret[0].(*observe.GetBlockByNumberResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByNumber indicates an expected call of GetBlockByNumber
func (mr *MockRPCClientMockRecorder) GetBlockByNumber(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByNumber", reflect.TypeOf((*MockRPCClient)(nil).GetBlockByNumber), varargs...)
}

// GetLatestBlock mocks base method
func (m *MockRPCClient) GetLatestBlock(arg0 context.Context, arg1 *observe.GetLatestBlockRequest, arg2 ...grpc.CallOption) (*observe.GetLatestBlockResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLatestBlock", varargs...)
	ret0, _ := ret[0].(*observe.GetLatestBlockResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestBlock indicates an expected call of GetLatestBlock
func (mr *MockRPCClientMockRecorder) GetLatestBlock(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestBlock", reflect.TypeOf((*MockRPCClient)(nil).GetLatestBlock), varargs...)
}

// GetTransaction mocks base method
func (m *MockRPCClient) GetTransaction(arg0 context.Context, arg1 *observe.GetTransactionRequest, arg2 ...grpc.CallOption) (*observe.GetTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTransaction", varargs...)
	ret0, _ := ret[0].(*observe.GetTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction
func (mr *MockRPCClientMockRecorder) GetTransaction(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockRPCClient)(nil).GetTransaction), varargs...)
}

// Ping mocks base method
func (m *MockRPCClient) Ping(arg0 context.Context, arg1 *observe.PingRequest, arg2 ...grpc.CallOption) (*observe.PingResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Ping", varargs...)
	ret0, _ := ret[0].(*observe.PingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ping indicates an expected call of Ping
func (mr *MockRPCClientMockRecorder) Ping(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRPCClient)(nil).Ping), varargs...)
}

// SendTransaction mocks base method
func (m *MockRPCClient) SendTransaction(arg0 context.Context, arg1 *observe.SendTransactionRequest, arg2 ...grpc.CallOption) (*observe.SendTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendTransaction", varargs...)
	ret0, _ := ret[0].(*observe.SendTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendTransaction indicates an expected call of SendTransaction
func (mr *MockRPCClientMockRecorder) SendTransaction(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTransaction", reflect.TypeOf((*MockRPCClient)(nil).SendTransaction), varargs...)
}
