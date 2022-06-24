// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	badger "github.com/dgraph-io/badger/v2"
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// EpochStates is an autogenerated mock type for the EpochStates type
type EpochStates struct {
	mock.Mock
}

// ByBlockID provides a mock function with given fields: _a0
func (_m *EpochStates) ByBlockID(_a0 flow.Identifier) (*flow.EpochStatus, error) {
	ret := _m.Called(_a0)

	var r0 *flow.EpochStatus
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.EpochStatus); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.EpochStatus)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreTx provides a mock function with given fields: blockID, state
func (_m *EpochStates) StoreTx(blockID flow.Identifier, state *flow.EpochStatus) func(*badger.Txn) error {
	ret := _m.Called(blockID, state)

	var r0 func(*badger.Txn) error
	if rf, ok := ret.Get(0).(func(flow.Identifier, *flow.EpochStatus) func(*badger.Txn) error); ok {
		r0 = rf(blockID, state)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func(*badger.Txn) error)
		}
	}

	return r0
}
