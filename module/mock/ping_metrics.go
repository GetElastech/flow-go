// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	flow "github.com/onflow/flow-go/model/flow"

	time "time"
)

// PingMetrics is an autogenerated mock type for the PingMetrics type
type PingMetrics struct {
	mock.Mock
}

// NodeInfo provides a mock function with given fields: node, nodeInfo, version, sealedHeight, hotstuffCurView
func (_m *PingMetrics) NodeInfo(node *flow.Identity, nodeInfo string, version string, sealedHeight uint64, hotstuffCurView uint64) {
	_m.Called(node, nodeInfo, version, sealedHeight, hotstuffCurView)
}

// NodeReachable provides a mock function with given fields: node, nodeInfo, rtt
func (_m *PingMetrics) NodeReachable(node *flow.Identity, nodeInfo string, rtt time.Duration) {
	_m.Called(node, nodeInfo, rtt)
}
