// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
)

// FollowerLogic is an autogenerated mock type for the FollowerLogic type
type FollowerLogic struct {
	mock.Mock
}

// AddBlock provides a mock function with given fields: proposal
func (_m *FollowerLogic) AddBlock(proposal *model.Proposal) error {
	ret := _m.Called(proposal)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Proposal) error); ok {
		r0 = rf(proposal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FinalizedBlock provides a mock function with given fields:
func (_m *FollowerLogic) FinalizedBlock() *model.Block {
	ret := _m.Called()

	var r0 *model.Block
	if rf, ok := ret.Get(0).(func() *model.Block); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Block)
		}
	}

	return r0
}
