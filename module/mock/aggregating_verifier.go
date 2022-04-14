// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	crypto "github.com/onflow/flow-go/crypto"
)

// AggregatingVerifier is an autogenerated mock type for the AggregatingVerifier type
type AggregatingVerifier struct {
	mock.Mock
}

// Verify provides a mock function with given fields: msg, sig, key
func (_m *AggregatingVerifier) Verify(msg []byte, sig crypto.Signature, key crypto.PublicKey) (bool, error) {
	ret := _m.Called(msg, sig, key)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte, crypto.Signature, crypto.PublicKey) bool); ok {
		r0 = rf(msg, sig, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, crypto.Signature, crypto.PublicKey) error); ok {
		r1 = rf(msg, sig, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyMany provides a mock function with given fields: msg, sig, keys
func (_m *AggregatingVerifier) VerifyMany(msg []byte, sig crypto.Signature, keys []crypto.PublicKey) (bool, error) {
	ret := _m.Called(msg, sig, keys)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte, crypto.Signature, []crypto.PublicKey) bool); ok {
		r0 = rf(msg, sig, keys)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, crypto.Signature, []crypto.PublicKey) error); ok {
		r1 = rf(msg, sig, keys)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
