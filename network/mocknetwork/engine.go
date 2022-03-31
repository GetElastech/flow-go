// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocknetwork

import (
	mock "github.com/stretchr/testify/mock"

	flow "github.com/onflow/flow-go/model/flow"

	network "github.com/onflow/flow-go/network"
)

// Engine is an autogenerated mock type for the Engine type
type Engine struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *Engine) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Process provides a mock function with given fields: channel, originID, event
func (_m *Engine) Process(channel network.Channel, originID flow.Identifier, event interface{}) error {
	ret := _m.Called(channel, originID, event)

	var r0 error
	if rf, ok := ret.Get(0).(func(network.Channel, flow.Identifier, interface{}) error); ok {
		r0 = rf(channel, originID, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessLocal provides a mock function with given fields: event
func (_m *Engine) ProcessLocal(event interface{}) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *Engine) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Submit provides a mock function with given fields: channel, originID, event
func (_m *Engine) Submit(channel network.Channel, originID flow.Identifier, event interface{}) {
	_m.Called(channel, originID, event)
}

// SubmitLocal provides a mock function with given fields: event
func (_m *Engine) SubmitLocal(event interface{}) {
	_m.Called(event)
}
