// Code generated by MockGen. DO NOT EDIT.
// Source: ./common/mail/init.go

// Package mock_mail is a generated GoMock package.
package mock_mail

import (
	reflect "reflect"

	mail "github.com/GregChrisnaDev/Amartha-Sol-3/common/mail"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// SendMail mocks base method.
func (m *MockClient) SendMail(templateName, emailDest string, data mail.AgreementMailReq) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMail", templateName, emailDest, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMail indicates an expected call of SendMail.
func (mr *MockClientMockRecorder) SendMail(templateName, emailDest, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMail", reflect.TypeOf((*MockClient)(nil).SendMail), templateName, emailDest, data)
}
