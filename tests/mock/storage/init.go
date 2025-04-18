// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/storage/init.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	reflect "reflect"

	storage "github.com/GregChrisnaDev/Amartha-Sol-3/internal/storage"
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

// DownloadFile mocks base method.
func (m *MockClient) DownloadFile(path string) (storage.DownloadFileResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFile", path)
	ret0, _ := ret[0].(storage.DownloadFileResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadFile indicates an expected call of DownloadFile.
func (mr *MockClientMockRecorder) DownloadFile(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockClient)(nil).DownloadFile), path)
}

// GetMainPath mocks base method.
func (m *MockClient) GetMainPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMainPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetMainPath indicates an expected call of GetMainPath.
func (mr *MockClientMockRecorder) GetMainPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMainPath", reflect.TypeOf((*MockClient)(nil).GetMainPath))
}

// UploadImage mocks base method.
func (m *MockClient) UploadImage(fileData []byte, dest, filename string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", fileData, dest, filename)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockClientMockRecorder) UploadImage(fileData, dest, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockClient)(nil).UploadImage), fileData, dest, filename)
}
