// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	flow "github.com/onflow/flow-go/model/flow"

	storage "github.com/onflow/flow-go/storage"
)

// TransactionResults is an autogenerated mock type for the TransactionResults type
type TransactionResults struct {
	mock.Mock
}

// BatchStore provides a mock function with given fields: blockID, transactionResults, batch
func (_m *TransactionResults) BatchStore(blockID flow.Identifier, transactionResults []flow.TransactionResult, batch storage.BatchStorage) error {
	ret := _m.Called(blockID, transactionResults, batch)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, []flow.TransactionResult, storage.BatchStorage) error); ok {
		r0 = rf(blockID, transactionResults, batch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ByBlockID provides a mock function with given fields: id
func (_m *TransactionResults) ByBlockID(id flow.Identifier) ([]flow.TransactionResult, error) {
	ret := _m.Called(id)

	var r0 []flow.TransactionResult
	if rf, ok := ret.Get(0).(func(flow.Identifier) []flow.TransactionResult); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.TransactionResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByBlockIDTransactionID provides a mock function with given fields: blockID, transactionID
func (_m *TransactionResults) ByBlockIDTransactionID(blockID flow.Identifier, transactionID flow.Identifier) (*flow.TransactionResult, error) {
	ret := _m.Called(blockID, transactionID)

	var r0 *flow.TransactionResult
	if rf, ok := ret.Get(0).(func(flow.Identifier, flow.Identifier) *flow.TransactionResult); ok {
		r0 = rf(blockID, transactionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.TransactionResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier, flow.Identifier) error); ok {
		r1 = rf(blockID, transactionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByBlockIDTransactionIndex provides a mock function with given fields: blockID, txIndex
func (_m *TransactionResults) ByBlockIDTransactionIndex(blockID flow.Identifier, txIndex uint32) (*flow.TransactionResult, error) {
	ret := _m.Called(blockID, txIndex)

	var r0 *flow.TransactionResult
	if rf, ok := ret.Get(0).(func(flow.Identifier, uint32) *flow.TransactionResult); ok {
		r0 = rf(blockID, txIndex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.TransactionResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier, uint32) error); ok {
		r1 = rf(blockID, txIndex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
