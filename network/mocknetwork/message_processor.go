// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocknetwork

import (
	mock "github.com/stretchr/testify/mock"

	flow "github.com/onflow/flow-go/model/flow"

	network "github.com/onflow/flow-go/network"
)

// MessageProcessor is an autogenerated mock type for the MessageProcessor type
type MessageProcessor struct {
	mock.Mock
}

// Process provides a mock function with given fields: channel, originID, message
func (_m *MessageProcessor) Process(channel network.Channel, originID flow.Identifier, message interface{}) error {
	ret := _m.Called(channel, originID, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(network.Channel, flow.Identifier, interface{}) error); ok {
		r0 = rf(channel, originID, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
