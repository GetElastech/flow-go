// Code generated by mockery v2.13.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	chunks "github.com/onflow/flow-go/model/chunks"

	verification "github.com/onflow/flow-go/model/verification"
)

// ChunkVerifier is an autogenerated mock type for the ChunkVerifier type
type ChunkVerifier struct {
	mock.Mock
}

// SystemChunkVerify provides a mock function with given fields: ch
func (_m *ChunkVerifier) SystemChunkVerify(ch *verification.VerifiableChunkData) ([]byte, chunks.ChunkFault, error) {
	ret := _m.Called(ch)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*verification.VerifiableChunkData) []byte); ok {
		r0 = rf(ch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 chunks.ChunkFault
	if rf, ok := ret.Get(1).(func(*verification.VerifiableChunkData) chunks.ChunkFault); ok {
		r1 = rf(ch)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(chunks.ChunkFault)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*verification.VerifiableChunkData) error); ok {
		r2 = rf(ch)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Verify provides a mock function with given fields: ch
func (_m *ChunkVerifier) Verify(ch *verification.VerifiableChunkData) ([]byte, chunks.ChunkFault, error) {
	ret := _m.Called(ch)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*verification.VerifiableChunkData) []byte); ok {
		r0 = rf(ch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 chunks.ChunkFault
	if rf, ok := ret.Get(1).(func(*verification.VerifiableChunkData) chunks.ChunkFault); ok {
		r1 = rf(ch)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(chunks.ChunkFault)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*verification.VerifiableChunkData) error); ok {
		r2 = rf(ch)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type NewChunkVerifierT interface {
	mock.TestingT
	Cleanup(func())
}

// NewChunkVerifier creates a new instance of ChunkVerifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewChunkVerifier(t NewChunkVerifierT) *ChunkVerifier {
	mock := &ChunkVerifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
