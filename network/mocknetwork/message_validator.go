// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocknetwork

import (
	mock "github.com/stretchr/testify/mock"

	message "github.com/onflow/flow-go/network/message"
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
