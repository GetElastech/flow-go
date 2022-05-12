// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// LedgerVerifier is an autogenerated mock type for the LedgerVerifier type
type LedgerVerifier struct {
	mock.Mock
}

// VerifyRegistersProof provides a mock function with given fields: registerIDs, stateCommitment, values, proof
func (_m *LedgerVerifier) VerifyRegistersProof(registerIDs []flow.RegisterID, stateCommitment flow.StateCommitment, values [][]byte, proof [][]byte) (bool, error) {
	ret := _m.Called(registerIDs, stateCommitment, values, proof)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]flow.RegisterID, flow.StateCommitment, [][]byte, [][]byte) bool); ok {
		r0 = rf(registerIDs, stateCommitment, values, proof)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]flow.RegisterID, flow.StateCommitment, [][]byte, [][]byte) error); ok {
		r1 = rf(registerIDs, stateCommitment, values, proof)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
