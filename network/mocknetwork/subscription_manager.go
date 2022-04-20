// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocknetwork

import (
	mock "github.com/stretchr/testify/mock"

	network "github.com/onflow/flow-go/network"
)

// SubscriptionManager is an autogenerated mock type for the SubscriptionManager type
type SubscriptionManager struct {
	mock.Mock
}

// Channels provides a mock function with given fields:
func (_m *SubscriptionManager) Channels() network.ChannelList {
	ret := _m.Called()

	var r0 network.ChannelList
	if rf, ok := ret.Get(0).(func() network.ChannelList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.ChannelList)
		}
	}

	return r0
}

// GetEngine provides a mock function with given fields: channel
func (_m *SubscriptionManager) GetEngine(channel network.Channel) (network.MessageProcessor, error) {
	ret := _m.Called(channel)

	var r0 network.MessageProcessor
	if rf, ok := ret.Get(0).(func(network.Channel) network.MessageProcessor); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.MessageProcessor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(network.Channel) error); ok {
		r1 = rf(channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: channel, engine
func (_m *SubscriptionManager) Register(channel network.Channel, engine network.MessageProcessor) error {
	ret := _m.Called(channel, engine)

	var r0 error
	if rf, ok := ret.Get(0).(func(network.Channel, network.MessageProcessor) error); ok {
		r0 = rf(channel, engine)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unregister provides a mock function with given fields: channel
func (_m *SubscriptionManager) Unregister(channel network.Channel) error {
	ret := _m.Called(channel)

	var r0 error
	if rf, ok := ret.Get(0).(func(network.Channel) error); ok {
		r0 = rf(channel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
