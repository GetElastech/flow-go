// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
)

// FinalizationConsumer is an autogenerated mock type for the FinalizationConsumer type
type FinalizationConsumer struct {
	mock.Mock
}

// OnBlockIncorporated provides a mock function with given fields: _a0
func (_m *FinalizationConsumer) OnBlockIncorporated(_a0 *model.Block) {
	_m.Called(_a0)
}

// OnDoubleProposeDetected provides a mock function with given fields: _a0, _a1
func (_m *FinalizationConsumer) OnDoubleProposeDetected(_a0 *model.Block, _a1 *model.Block) {
	_m.Called(_a0, _a1)
}

// OnFinalizedBlock provides a mock function with given fields: _a0
func (_m *FinalizationConsumer) OnFinalizedBlock(_a0 *model.Block) {
	_m.Called(_a0)
}
