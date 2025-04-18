// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecase/user.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	model "github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	usecase "github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
	gomock "github.com/golang/mock/gomock"
)

// MockUserUsecase is a mock of UserUsecase interface.
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase.
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance.
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// GenerateUser mocks base method.
func (m *MockUserUsecase) GenerateUser(ctx context.Context, params usecase.UserGenerateReq) (usecase.UserResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateUser", ctx, params)
	ret0, _ := ret[0].(usecase.UserResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateUser indicates an expected call of GenerateUser.
func (mr *MockUserUsecaseMockRecorder) GenerateUser(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateUser", reflect.TypeOf((*MockUserUsecase)(nil).GenerateUser), ctx, params)
}

// GetAllUser mocks base method.
func (m *MockUserUsecase) GetAllUser(ctx context.Context) ([]usecase.UserResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUser", ctx)
	ret0, _ := ret[0].([]usecase.UserResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUser indicates an expected call of GetAllUser.
func (mr *MockUserUsecaseMockRecorder) GetAllUser(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUser", reflect.TypeOf((*MockUserUsecase)(nil).GetAllUser), ctx)
}

// ValidateUser mocks base method.
func (m *MockUserUsecase) ValidateUser(ctx context.Context, params usecase.ValidateUserReq) *model.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateUser", ctx, params)
	ret0, _ := ret[0].(*model.User)
	return ret0
}

// ValidateUser indicates an expected call of ValidateUser.
func (mr *MockUserUsecaseMockRecorder) ValidateUser(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateUser", reflect.TypeOf((*MockUserUsecase)(nil).ValidateUser), ctx, params)
}
