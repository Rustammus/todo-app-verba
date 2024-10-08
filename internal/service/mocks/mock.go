// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mocks\mock.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	dto "ToDoVerba/internal/dto"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockITaskService is a mock of ITaskService interface.
type MockITaskService struct {
	ctrl     *gomock.Controller
	recorder *MockITaskServiceMockRecorder
}

// MockITaskServiceMockRecorder is the mock recorder for MockITaskService.
type MockITaskServiceMockRecorder struct {
	mock *MockITaskService
}

// NewMockITaskService creates a new mock instance.
func NewMockITaskService(ctrl *gomock.Controller) *MockITaskService {
	mock := &MockITaskService{ctrl: ctrl}
	mock.recorder = &MockITaskServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITaskService) EXPECT() *MockITaskServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockITaskService) Create(cTask *dto.TaskCreate) (*dto.TaskRead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", cTask)
	ret0, _ := ret[0].(*dto.TaskRead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockITaskServiceMockRecorder) Create(cTask any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockITaskService)(nil).Create), cTask)
}

// DeleteById mocks base method.
func (m *MockITaskService) DeleteById(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockITaskServiceMockRecorder) DeleteById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockITaskService)(nil).DeleteById), id)
}

// FindByID mocks base method.
func (m *MockITaskService) FindByID(id int) (*dto.TaskRead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*dto.TaskRead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockITaskServiceMockRecorder) FindByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockITaskService)(nil).FindByID), id)
}

// List mocks base method.
func (m *MockITaskService) List() ([]dto.TaskRead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]dto.TaskRead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockITaskServiceMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockITaskService)(nil).List))
}

// UpdateById mocks base method.
func (m *MockITaskService) UpdateById(id int, update *dto.TaskUpdate) (*dto.TaskRead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", id, update)
	ret0, _ := ret[0].(*dto.TaskRead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockITaskServiceMockRecorder) UpdateById(id, update any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockITaskService)(nil).UpdateById), id, update)
}
