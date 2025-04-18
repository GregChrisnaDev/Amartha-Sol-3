// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/lend.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	model "github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockLendRepository is a mock of LendRepository interface.
type MockLendRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLendRepositoryMockRecorder
}

// MockLendRepositoryMockRecorder is the mock recorder for MockLendRepository.
type MockLendRepositoryMockRecorder struct {
	mock *MockLendRepository
}

// NewMockLendRepository creates a new mock instance.
func NewMockLendRepository(ctrl *gomock.Controller) *MockLendRepository {
	mock := &MockLendRepository{ctrl: ctrl}
	mock.recorder = &MockLendRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLendRepository) EXPECT() *MockLendRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockLendRepository) Add(ctx context.Context, params model.Lend) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockLendRepositoryMockRecorder) Add(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockLendRepository)(nil).Add), ctx, params)
}

// GetByLoanId mocks base method.
func (m *MockLendRepository) GetByLoanId(ctx context.Context, loanId uint64) ([]model.Lend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByLoanId", ctx, loanId)
	ret0, _ := ret[0].([]model.Lend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByLoanId indicates an expected call of GetByLoanId.
func (mr *MockLendRepositoryMockRecorder) GetByLoanId(ctx, loanId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByLoanId", reflect.TypeOf((*MockLendRepository)(nil).GetByLoanId), ctx, loanId)
}

// GetByUID mocks base method.
func (m *MockLendRepository) GetByUID(ctx context.Context, userId uint64) ([]model.Lend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUID", ctx, userId)
	ret0, _ := ret[0].([]model.Lend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUID indicates an expected call of GetByUID.
func (mr *MockLendRepositoryMockRecorder) GetByUID(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUID", reflect.TypeOf((*MockLendRepository)(nil).GetByUID), ctx, userId)
}

// GetByUidLoanId mocks base method.
func (m *MockLendRepository) GetByUidLoanId(ctx context.Context, loanId, userId uint64) (model.Lend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUidLoanId", ctx, loanId, userId)
	ret0, _ := ret[0].(model.Lend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUidLoanId indicates an expected call of GetByUidLoanId.
func (mr *MockLendRepositoryMockRecorder) GetByUidLoanId(ctx, loanId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUidLoanId", reflect.TypeOf((*MockLendRepository)(nil).GetByUidLoanId), ctx, loanId, userId)
}

// GetListLenderByLoanerID mocks base method.
func (m *MockLendRepository) GetListLenderByLoanerID(ctx context.Context, loanId, userId uint64) ([]model.Lend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListLenderByLoanerID", ctx, loanId, userId)
	ret0, _ := ret[0].([]model.Lend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListLenderByLoanerID indicates an expected call of GetListLenderByLoanerID.
func (mr *MockLendRepositoryMockRecorder) GetListLenderByLoanerID(ctx, loanId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListLenderByLoanerID", reflect.TypeOf((*MockLendRepository)(nil).GetListLenderByLoanerID), ctx, loanId, userId)
}

// Update mocks base method.
func (m *MockLendRepository) Update(ctx context.Context, params model.Lend) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockLendRepositoryMockRecorder) Update(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockLendRepository)(nil).Update), ctx, params)
}
