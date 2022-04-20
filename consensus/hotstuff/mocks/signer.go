// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
)

// Signer is an autogenerated mock type for the Signer type
type Signer struct {
	mock.Mock
}

// CreateProposal provides a mock function with given fields: block
func (_m *Signer) CreateProposal(block *model.Block) (*model.Proposal, error) {
	ret := _m.Called(block)

	var r0 *model.Proposal
	if rf, ok := ret.Get(0).(func(*model.Block) *model.Proposal); ok {
		r0 = rf(block)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Proposal)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Block) error); ok {
		r1 = rf(block)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateVote provides a mock function with given fields: block
func (_m *Signer) CreateVote(block *model.Block) (*model.Vote, error) {
	ret := _m.Called(block)

	var r0 *model.Vote
	if rf, ok := ret.Get(0).(func(*model.Block) *model.Vote); ok {
		r0 = rf(block)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Vote)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Block) error); ok {
		r1 = rf(block)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
