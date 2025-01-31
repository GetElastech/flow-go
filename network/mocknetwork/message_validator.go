// Code generated by mockery v2.13.0. DO NOT EDIT.

package mocknetwork

import (
	message "github.com/onflow/flow-go/network/message"
	mock "github.com/stretchr/testify/mock"
)

// MessageValidator is an autogenerated mock type for the MessageValidator type
type MessageValidator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: msg
func (_m *MessageValidator) Validate(msg message.Message) bool {
	ret := _m.Called(msg)

	var r0 bool
	if rf, ok := ret.Get(0).(func(message.Message) bool); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type NewMessageValidatorT interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageValidator creates a new instance of MessageValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageValidator(t NewMessageValidatorT) *MessageValidator {
	mock := &MessageValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
