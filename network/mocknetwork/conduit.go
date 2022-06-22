// Code generated by mockery v2.13.0. DO NOT EDIT.

package mocknetwork

import (
	mock "github.com/stretchr/testify/mock"

	flow "github.com/onflow/flow-go/model/flow"
)

// Conduit is an autogenerated mock type for the Conduit type
type Conduit struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Conduit) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Multicast provides a mock function with given fields: event, num, targetIDs
func (_m *Conduit) Multicast(event interface{}, num uint, targetIDs ...flow.Identifier) error {
	_va := make([]interface{}, len(targetIDs))
	for _i := range targetIDs {
		_va[_i] = targetIDs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, event, num)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, uint, ...flow.Identifier) error); ok {
		r0 = rf(event, num, targetIDs...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Publish provides a mock function with given fields: event, targetIDs
func (_m *Conduit) Publish(event interface{}, targetIDs ...flow.Identifier) error {
	_va := make([]interface{}, len(targetIDs))
	for _i := range targetIDs {
		_va[_i] = targetIDs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, event)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, ...flow.Identifier) error); ok {
		r0 = rf(event, targetIDs...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unicast provides a mock function with given fields: event, targetID
func (_m *Conduit) Unicast(event interface{}, targetID flow.Identifier) error {
	ret := _m.Called(event, targetID)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, flow.Identifier) error); ok {
		r0 = rf(event, targetID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewConduitT interface {
	mock.TestingT
	Cleanup(func())
}

// NewConduit creates a new instance of Conduit. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConduit(t NewConduitT) *Conduit {
	mock := &Conduit{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
