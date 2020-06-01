// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kacejot/resize-service/pkg/storage (interfaces: Storage)

// Package storage is a generated GoMock package.
package storage

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/kacejot/resize-service/pkg/api/graph/model"
	reflect "reflect"
)

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetRecordByID mocks base method
func (m *MockStorage) GetRecordByID(arg0 string) (model.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecordByID", arg0)
	ret0, _ := ret[0].(model.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecordByID indicates an expected call of GetRecordByID
func (mr *MockStorageMockRecorder) GetRecordByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecordByID", reflect.TypeOf((*MockStorage)(nil).GetRecordByID), arg0)
}

// ListUserRecords mocks base method
func (m *MockStorage) ListUserRecords(arg0 string) []model.ResizeResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserRecords", arg0)
	ret0, _ := ret[0].([]model.ResizeResult)
	return ret0
}

// ListUserRecords indicates an expected call of ListUserRecords
func (mr *MockStorageMockRecorder) ListUserRecords(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserRecords", reflect.TypeOf((*MockStorage)(nil).ListUserRecords), arg0)
}

// RecordResizeResult mocks base method
func (m *MockStorage) RecordResizeResult(arg0 string, arg1, arg2 model.Image) (model.ResizeResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordResizeResult", arg0, arg1, arg2)
	ret0, _ := ret[0].(model.ResizeResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RecordResizeResult indicates an expected call of RecordResizeResult
func (mr *MockStorageMockRecorder) RecordResizeResult(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordResizeResult", reflect.TypeOf((*MockStorage)(nil).RecordResizeResult), arg0, arg1, arg2)
}