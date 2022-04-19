// Code generated by MockGen. DO NOT EDIT.
// Source: network/transport/v2/senders.go

// Package v2 is a generated GoMock package.
package v2

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hash "github.com/nuts-foundation/nuts-node/crypto/hash"
	tree "github.com/nuts-foundation/nuts-node/network/dag/tree"
	transport "github.com/nuts-foundation/nuts-node/network/transport"
)

// MockmessageSender is a mock of messageSender interface.
type MockmessageSender struct {
	ctrl     *gomock.Controller
	recorder *MockmessageSenderMockRecorder
}

// MockmessageSenderMockRecorder is the mock recorder for MockmessageSender.
type MockmessageSenderMockRecorder struct {
	mock *MockmessageSender
}

// NewMockmessageSender creates a new mock instance.
func NewMockmessageSender(ctrl *gomock.Controller) *MockmessageSender {
	mock := &MockmessageSender{ctrl: ctrl}
	mock.recorder = &MockmessageSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmessageSender) EXPECT() *MockmessageSenderMockRecorder {
	return m.recorder
}

// sendGossipMsg mocks base method.
func (m *MockmessageSender) sendGossipMsg(id transport.PeerID, refs []hash.SHA256Hash, xor hash.SHA256Hash, clock uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "sendGossipMsg", id, refs, xor, clock)
	ret0, _ := ret[0].(error)
	return ret0
}

// sendGossipMsg indicates an expected call of sendGossipMsg.
func (mr *MockmessageSenderMockRecorder) sendGossipMsg(id, refs, xor, clock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "sendGossipMsg", reflect.TypeOf((*MockmessageSender)(nil).sendGossipMsg), id, refs, xor, clock)
}

// sendState mocks base method.
func (m *MockmessageSender) sendState(id transport.PeerID, xor hash.SHA256Hash, clock uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "sendState", id, xor, clock)
	ret0, _ := ret[0].(error)
	return ret0
}

// sendState indicates an expected call of sendState.
func (mr *MockmessageSenderMockRecorder) sendState(id, xor, clock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "sendState", reflect.TypeOf((*MockmessageSender)(nil).sendState), id, xor, clock)
}

// sendTransactionList mocks base method.
func (m *MockmessageSender) sendTransactionList(peerID transport.PeerID, conversationID conversationID, transactions []*Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "sendTransactionList", peerID, conversationID, transactions)
	ret0, _ := ret[0].(error)
	return ret0
}

// sendTransactionList indicates an expected call of sendTransactionList.
func (mr *MockmessageSenderMockRecorder) sendTransactionList(peerID, conversationID, transactions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "sendTransactionList", reflect.TypeOf((*MockmessageSender)(nil).sendTransactionList), peerID, conversationID, transactions)
}

// sendTransactionListQuery mocks base method.
func (m *MockmessageSender) sendTransactionListQuery(id transport.PeerID, refs []hash.SHA256Hash) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "sendTransactionListQuery", id, refs)
	ret0, _ := ret[0].(error)
	return ret0
}

// sendTransactionListQuery indicates an expected call of sendTransactionListQuery.
func (mr *MockmessageSenderMockRecorder) sendTransactionListQuery(id, refs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "sendTransactionListQuery", reflect.TypeOf((*MockmessageSender)(nil).sendTransactionListQuery), id, refs)
}

// sendTransactionRangeQuery mocks base method.
func (m *MockmessageSender) sendTransactionRangeQuery(id transport.PeerID, lcStart, lcEnd uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "sendTransactionRangeQuery", id, lcStart, lcEnd)
	ret0, _ := ret[0].(error)
	return ret0
}

// sendTransactionRangeQuery indicates an expected call of sendTransactionRangeQuery.
func (mr *MockmessageSenderMockRecorder) sendTransactionRangeQuery(id, lcStart, lcEnd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "sendTransactionRangeQuery", reflect.TypeOf((*MockmessageSender)(nil).sendTransactionRangeQuery), id, lcStart, lcEnd)
}

// sendTransactionSet mocks base method.
func (m *MockmessageSender) sendTransactionSet(id transport.PeerID, conversationID conversationID, LCReq, LC uint32, iblt tree.Iblt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "sendTransactionSet", id, conversationID, LCReq, LC, iblt)
	ret0, _ := ret[0].(error)
	return ret0
}

// sendTransactionSet indicates an expected call of sendTransactionSet.
func (mr *MockmessageSenderMockRecorder) sendTransactionSet(id, conversationID, LCReq, LC, iblt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "sendTransactionSet", reflect.TypeOf((*MockmessageSender)(nil).sendTransactionSet), id, conversationID, LCReq, LC, iblt)
}
