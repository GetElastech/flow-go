// Code generated by mockery v2.13.0. DO NOT EDIT.

package mocknetwork

import (
	datastore "github.com/ipfs/go-datastore"
	mock "github.com/stretchr/testify/mock"

	irrecoverable "github.com/onflow/flow-go/module/irrecoverable"

	network "github.com/onflow/flow-go/network"

	protocol "github.com/libp2p/go-libp2p-core/protocol"
)

// Network is an autogenerated mock type for the Network type
type Network struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *Network) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *Network) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Register provides a mock function with given fields: channel, messageProcessor
func (_m *Network) Register(channel network.Channel, messageProcessor network.MessageProcessor) (network.Conduit, error) {
	ret := _m.Called(channel, messageProcessor)

	var r0 network.Conduit
	if rf, ok := ret.Get(0).(func(network.Channel, network.MessageProcessor) network.Conduit); ok {
		r0 = rf(channel, messageProcessor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.Conduit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(network.Channel, network.MessageProcessor) error); ok {
		r1 = rf(channel, messageProcessor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterBlobService provides a mock function with given fields: channel, store, opts
func (_m *Network) RegisterBlobService(channel network.Channel, store datastore.Batching, opts ...network.BlobServiceOption) (network.BlobService, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, channel, store)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 network.BlobService
	if rf, ok := ret.Get(0).(func(network.Channel, datastore.Batching, ...network.BlobServiceOption) network.BlobService); ok {
		r0 = rf(channel, store, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.BlobService)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(network.Channel, datastore.Batching, ...network.BlobServiceOption) error); ok {
		r1 = rf(channel, store, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterPingService provides a mock function with given fields: pingProtocolID, pingInfoProvider
func (_m *Network) RegisterPingService(pingProtocolID protocol.ID, pingInfoProvider network.PingInfoProvider) (network.PingService, error) {
	ret := _m.Called(pingProtocolID, pingInfoProvider)

	var r0 network.PingService
	if rf, ok := ret.Get(0).(func(protocol.ID, network.PingInfoProvider) network.PingService); ok {
		r0 = rf(pingProtocolID, pingInfoProvider)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.PingService)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(protocol.ID, network.PingInfoProvider) error); ok {
		r1 = rf(pingProtocolID, pingInfoProvider)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: _a0
func (_m *Network) Start(_a0 irrecoverable.SignalerContext) {
	_m.Called(_a0)
}

type NewNetworkT interface {
	mock.TestingT
	Cleanup(func())
}

// NewNetwork creates a new instance of Network. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNetwork(t NewNetworkT) *Network {
	mock := &Network{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
