// Code generated by mockery v2.13.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	ledger "github.com/onflow/flow-go/ledger"
)

// Ledger is an autogenerated mock type for the Ledger type
type Ledger struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *Ledger) Done() <-chan struct{} {
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

// Get provides a mock function with given fields: query
func (_m *Ledger) Get(query *ledger.Query) ([]ledger.Value, error) {
	ret := _m.Called(query)

	var r0 []ledger.Value
	if rf, ok := ret.Get(0).(func(*ledger.Query) []ledger.Value); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ledger.Value)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ledger.Query) error); ok {
		r1 = rf(query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSingleValue provides a mock function with given fields: query
func (_m *Ledger) GetSingleValue(query *ledger.QuerySingleValue) (ledger.Value, error) {
	ret := _m.Called(query)

	var r0 ledger.Value
	if rf, ok := ret.Get(0).(func(*ledger.QuerySingleValue) ledger.Value); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.Value)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ledger.QuerySingleValue) error); ok {
		r1 = rf(query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitialState provides a mock function with given fields:
func (_m *Ledger) InitialState() ledger.State {
	ret := _m.Called()

	var r0 ledger.State
	if rf, ok := ret.Get(0).(func() ledger.State); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.State)
		}
	}

	return r0
}

// Prove provides a mock function with given fields: query
func (_m *Ledger) Prove(query *ledger.Query) (ledger.Proof, error) {
	ret := _m.Called(query)

	var r0 ledger.Proof
	if rf, ok := ret.Get(0).(func(*ledger.Query) ledger.Proof); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.Proof)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ledger.Query) error); ok {
		r1 = rf(query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ready provides a mock function with given fields:
func (_m *Ledger) Ready() <-chan struct{} {
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

// Set provides a mock function with given fields: update
func (_m *Ledger) Set(update *ledger.Update) (ledger.State, *ledger.TrieUpdate, error) {
	ret := _m.Called(update)

	var r0 ledger.State
	if rf, ok := ret.Get(0).(func(*ledger.Update) ledger.State); ok {
		r0 = rf(update)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.State)
		}
	}

	var r1 *ledger.TrieUpdate
	if rf, ok := ret.Get(1).(func(*ledger.Update) *ledger.TrieUpdate); ok {
		r1 = rf(update)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*ledger.TrieUpdate)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*ledger.Update) error); ok {
		r2 = rf(update)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type NewLedgerT interface {
	mock.TestingT
	Cleanup(func())
}

// NewLedger creates a new instance of Ledger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLedger(t NewLedgerT) *Ledger {
	mock := &Ledger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
