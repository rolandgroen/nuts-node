// Code generated by MockGen. DO NOT EDIT.
// Source: network/transport/v2/gossip/manager.go

// Package gossip is a generated GoMock package.
package gossip

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hash "github.com/nuts-foundation/nuts-node/crypto/hash"
	transport "github.com/nuts-foundation/nuts-node/network/transport"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// GossipReceived mocks base method.
func (m *MockManager) GossipReceived(id transport.PeerID, refs ...hash.SHA256Hash) {
	m.ctrl.T.Helper()
	varargs := []interface{}{id}
	for _, a := range refs {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "GossipReceived", varargs...)
}

// GossipReceived indicates an expected call of GossipReceived.
func (mr *MockManagerMockRecorder) GossipReceived(id interface{}, refs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{id}, refs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GossipReceived", reflect.TypeOf((*MockManager)(nil).GossipReceived), varargs...)
}

// PeerConnected mocks base method.
func (m *MockManager) PeerConnected(peer transport.Peer, xor hash.SHA256Hash, clock uint32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PeerConnected", peer, xor, clock)
}

// PeerConnected indicates an expected call of PeerConnected.
func (mr *MockManagerMockRecorder) PeerConnected(peer, xor, clock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeerConnected", reflect.TypeOf((*MockManager)(nil).PeerConnected), peer, xor, clock)
}

// PeerDisconnected mocks base method.
func (m *MockManager) PeerDisconnected(peer transport.Peer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PeerDisconnected", peer)
}

// PeerDisconnected indicates an expected call of PeerDisconnected.
func (mr *MockManagerMockRecorder) PeerDisconnected(peer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeerDisconnected", reflect.TypeOf((*MockManager)(nil).PeerDisconnected), peer)
}

// RegisterSender mocks base method.
func (m *MockManager) RegisterSender(arg0 SenderFunc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterSender", arg0)
}

// RegisterSender indicates an expected call of RegisterSender.
func (mr *MockManagerMockRecorder) RegisterSender(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterSender", reflect.TypeOf((*MockManager)(nil).RegisterSender), arg0)
}

// TransactionRegistered mocks base method.
func (m *MockManager) TransactionRegistered(transaction, xor hash.SHA256Hash, clock uint32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "TransactionRegistered", transaction, xor, clock)
}

// TransactionRegistered indicates an expected call of TransactionRegistered.
func (mr *MockManagerMockRecorder) TransactionRegistered(transaction, xor, clock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionRegistered", reflect.TypeOf((*MockManager)(nil).TransactionRegistered), transaction, xor, clock)
}
