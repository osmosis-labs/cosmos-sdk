// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cometbft/cometbft-db (interfaces: DB)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	db "github.com/cosmos/cosmos-db"
	gomock "github.com/golang/mock/gomock"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDB)(nil).Close))
}

// Delete mocks base method.
func (m *MockDB) Delete(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDBMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDB)(nil).Delete), arg0)
}

// DeleteSync mocks base method.
func (m *MockDB) DeleteSync(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSync", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSync indicates an expected call of DeleteSync.
func (mr *MockDBMockRecorder) DeleteSync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSync", reflect.TypeOf((*MockDB)(nil).DeleteSync), arg0)
}

// Get mocks base method.
func (m *MockDB) Get(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDBMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDB)(nil).Get), arg0)
}

// Has mocks base method.
func (m *MockDB) Has(arg0 []byte) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Has indicates an expected call of Has.
func (mr *MockDBMockRecorder) Has(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockDB)(nil).Has), arg0)
}

// Iterator mocks base method.
func (m *MockDB) Iterator(arg0, arg1 []byte) (db.Iterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Iterator", arg0, arg1)
	ret0, _ := ret[0].(db.Iterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Iterator indicates an expected call of Iterator.
func (mr *MockDBMockRecorder) Iterator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Iterator", reflect.TypeOf((*MockDB)(nil).Iterator), arg0, arg1)
}

// NewBatch mocks base method.
func (m *MockDB) NewBatch() db.Batch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewBatch")
	ret0, _ := ret[0].(db.Batch)
	return ret0
}

// NewBatch indicates an expected call of NewBatch.
func (mr *MockDBMockRecorder) NewBatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewBatch", reflect.TypeOf((*MockDB)(nil).NewBatch))
}

// Print mocks base method.
func (m *MockDB) Print() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Print")
	ret0, _ := ret[0].(error)
	return ret0
}

// Print indicates an expected call of Print.
func (mr *MockDBMockRecorder) Print() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Print", reflect.TypeOf((*MockDB)(nil).Print))
}

// ReverseIterator mocks base method.
func (m *MockDB) ReverseIterator(arg0, arg1 []byte) (db.Iterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReverseIterator", arg0, arg1)
	ret0, _ := ret[0].(db.Iterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReverseIterator indicates an expected call of ReverseIterator.
func (mr *MockDBMockRecorder) ReverseIterator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReverseIterator", reflect.TypeOf((*MockDB)(nil).ReverseIterator), arg0, arg1)
}

// Set mocks base method.
func (m *MockDB) Set(arg0, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockDBMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockDB)(nil).Set), arg0, arg1)
}

// SetSync mocks base method.
func (m *MockDB) SetSync(arg0, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSync", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSync indicates an expected call of SetSync.
func (mr *MockDBMockRecorder) SetSync(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSync", reflect.TypeOf((*MockDB)(nil).SetSync), arg0, arg1)
}

// Stats mocks base method.
func (m *MockDB) Stats() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// Stats indicates an expected call of Stats.
func (mr *MockDBMockRecorder) Stats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockDB)(nil).Stats))
}
