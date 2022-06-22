// Code generated by mockery v2.13.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	flow "github.com/onflow/flow-go/model/flow"
)

// Params is an autogenerated mock type for the Params type
type Params struct {
	mock.Mock
}

// ChainID provides a mock function with given fields:
func (_m *Params) ChainID() (flow.ChainID, error) {
	ret := _m.Called()

	var r0 flow.ChainID
	if rf, ok := ret.Get(0).(func() flow.ChainID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(flow.ChainID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type NewParamsT interface {
	mock.TestingT
	Cleanup(func())
}

// NewParams creates a new instance of Params. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewParams(t NewParamsT) *Params {
	mock := &Params{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
