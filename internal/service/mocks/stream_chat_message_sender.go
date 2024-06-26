// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/mistandok/chat-server/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// StreamChatMessageSender is an autogenerated mock type for the StreamChatMessageSender type
type StreamChatMessageSender struct {
	mock.Mock
}

// Context provides a mock function with given fields:
func (_m *StreamChatMessageSender) Context() context.Context {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Context")
	}

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// Send provides a mock function with given fields: _a0
func (_m *StreamChatMessageSender) Send(_a0 *model.Message) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Message) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStreamChatMessageSender creates a new instance of StreamChatMessageSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStreamChatMessageSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *StreamChatMessageSender {
	mock := &StreamChatMessageSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
