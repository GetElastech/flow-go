// Code generated by mockery v2.13.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// BackendScriptsMetrics is an autogenerated mock type for the BackendScriptsMetrics type
type BackendScriptsMetrics struct {
	mock.Mock
}

// ScriptExecuted provides a mock function with given fields: dur, size
func (_m *BackendScriptsMetrics) ScriptExecuted(dur time.Duration, size int) {
	_m.Called(dur, size)
}

type NewBackendScriptsMetricsT interface {
	mock.TestingT
	Cleanup(func())
}

// NewBackendScriptsMetrics creates a new instance of BackendScriptsMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBackendScriptsMetrics(t NewBackendScriptsMetricsT) *BackendScriptsMetrics {
	mock := &BackendScriptsMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
