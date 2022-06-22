// Code generated by mockery v2.13.0. DO NOT EDIT.

package mocknetwork

import (
	mock "github.com/stretchr/testify/mock"

	network "github.com/onflow/flow-go/network"
)

// BlobServiceOption is an autogenerated mock type for the BlobServiceOption type
type BlobServiceOption struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *BlobServiceOption) Execute(_a0 network.BlobService) {
	_m.Called(_a0)
}

type NewBlobServiceOptionT interface {
	mock.TestingT
	Cleanup(func())
}

// NewBlobServiceOption creates a new instance of BlobServiceOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBlobServiceOption(t NewBlobServiceOptionT) *BlobServiceOption {
	mock := &BlobServiceOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
