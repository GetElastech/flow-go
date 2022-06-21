package rpc

import (
	accessproto "github.com/onflow/flow/protobuf/go/flow/access"
	legacyaccessproto "github.com/onflow/flow/protobuf/go/flow/legacy/access"

	"github.com/onflow/flow-go/access"

	legacyaccess "github.com/onflow/flow-go/access/legacy"
	"github.com/onflow/flow-go/apiproxy"
)

// NewRPCEngineBuilder helps to build a new RPC engine.
func NewRPCEngineBuilder() *RPCEngineBuilder {
	builder := &RPCEngineBuilder{}
	return builder
}

type RPCEngineBuilder struct {
	*Engine
	// Use the parent interface instead of implementation, so that we can assign it to proxy.
	localAPIServer accessproto.AccessAPIServer
}

func (builder *RPCEngineBuilder) WithBase(engine *Engine) {
	builder.Engine = engine
	builder.localAPIServer = access.NewHandler(builder.backend, builder.chain)
}

func (builder *RPCEngineBuilder) WithRouting(router *apiproxy.FlowAccessAPIRouter) {
	router.SetLocalAPI(builder.localAPIServer)
	builder.localAPIServer = router
}

func (builder *RPCEngineBuilder) WithLegacy() {
	// Register legacy gRPC handlers for backwards compatibility, to be removed at a later date
	legacyaccessproto.RegisterAccessAPIServer(
		builder.unsecureGrpcServer,
		legacyaccess.NewHandler(builder.backend, builder.chain),
	)
	legacyaccessproto.RegisterAccessAPIServer(
		builder.secureGrpcServer,
		legacyaccess.NewHandler(builder.backend, builder.chain),
	)
}

func (builder *RPCEngineBuilder) WithRegister() {
	accessproto.RegisterAccessAPIServer(builder.unsecureGrpcServer, builder.localAPIServer)
	accessproto.RegisterAccessAPIServer(builder.secureGrpcServer, builder.localAPIServer)
}

func (builder *RPCEngineBuilder) Build() *Engine {
	return builder.Engine
}
