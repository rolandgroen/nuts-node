// Code generated by MockGen. DO NOT EDIT.
// Source: vcr/signature/signature.go

// Package signature is a generated GoMock package.
package signature

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ssi "github.com/nuts-foundation/go-did"
	crypto "github.com/nuts-foundation/nuts-node/crypto"
)

// MockSuite is a mock of Suite interface.
type MockSuite struct {
	ctrl     *gomock.Controller
	recorder *MockSuiteMockRecorder
}

// MockSuiteMockRecorder is the mock recorder for MockSuite.
type MockSuiteMockRecorder struct {
	mock *MockSuite
}

// NewMockSuite creates a new mock instance.
func NewMockSuite(ctrl *gomock.Controller) *MockSuite {
	mock := &MockSuite{ctrl: ctrl}
	mock.recorder = &MockSuiteMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSuite) EXPECT() *MockSuiteMockRecorder {
	return m.recorder
}

// CalculateDigest mocks base method.
func (m *MockSuite) CalculateDigest(doc []byte) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateDigest", doc)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// CalculateDigest indicates an expected call of CalculateDigest.
func (mr *MockSuiteMockRecorder) CalculateDigest(doc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateDigest", reflect.TypeOf((*MockSuite)(nil).CalculateDigest), doc)
}

// CanonicalizeDocument mocks base method.
func (m *MockSuite) CanonicalizeDocument(doc interface{}) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanonicalizeDocument", doc)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CanonicalizeDocument indicates an expected call of CanonicalizeDocument.
func (mr *MockSuiteMockRecorder) CanonicalizeDocument(doc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanonicalizeDocument", reflect.TypeOf((*MockSuite)(nil).CanonicalizeDocument), doc)
}

// GetType mocks base method.
func (m *MockSuite) GetType() ssi.ProofType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetType")
	ret0, _ := ret[0].(ssi.ProofType)
	return ret0
}

// GetType indicates an expected call of GetType.
func (mr *MockSuiteMockRecorder) GetType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetType", reflect.TypeOf((*MockSuite)(nil).GetType))
}

// Sign mocks base method.
func (m *MockSuite) Sign(doc []byte, key crypto.Key) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", doc, key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockSuiteMockRecorder) Sign(doc, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockSuite)(nil).Sign), doc, key)
}
