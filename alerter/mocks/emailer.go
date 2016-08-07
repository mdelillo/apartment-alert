// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/mdelillo/apartment-alert/alerter (interfaces: Emailer)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Emailer interface
type MockEmailer struct {
	ctrl     *gomock.Controller
	recorder *_MockEmailerRecorder
}

// Recorder for MockEmailer (not exported)
type _MockEmailerRecorder struct {
	mock *MockEmailer
}

func NewMockEmailer(ctrl *gomock.Controller) *MockEmailer {
	mock := &MockEmailer{ctrl: ctrl}
	mock.recorder = &_MockEmailerRecorder{mock}
	return mock
}

func (_m *MockEmailer) EXPECT() *_MockEmailerRecorder {
	return _m.recorder
}

func (_m *MockEmailer) Send(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Send", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockEmailerRecorder) Send(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Send", arg0)
}