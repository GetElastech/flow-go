package observer

import (
	"context"
	"github.com/onflow/flow-go/engine/access/rpc/backend"
	"github.com/onflow/flow-go/utils/grpcutils"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

func NewFlowProxy(accessNodeAddressAndPort string, observerListenerAddress string, timeout time.Duration) (*FlowProxy, error) {

	collectionRPCConn, err := grpc.Dial(
		accessNodeAddressAndPort,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcutils.DefaultMaxMsgSize)),
		grpc.WithInsecure(), //nolint:staticcheck
		backend.WithClientUnaryInterceptor(timeout))
	if err != nil {
		return nil, err
	}

	ret := &FlowProxy{}
	ret.upstream = access.NewAccessAPIClient(collectionRPCConn)

	server := grpc.NewServer()
	ret.downstream = server
	access.RegisterAccessAPIServer(server, ret)

	listener, err := net.Listen("tcp4", observerListenerAddress)
	if err != nil {
		return nil, err
	}
	go func() {
		_ = server.Serve(listener)
	}()

	return ret, nil
}

func NewFlowProxyStandalone(accessNodeAddressAndPort string, fallback *access.AccessAPIServer, timeout time.Duration) (access.AccessAPIServer, error) {
	collectionRPCConn, err := grpc.Dial(
		accessNodeAddressAndPort,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcutils.DefaultMaxMsgSize)),
		grpc.WithInsecure(), //nolint:staticcheck
		backend.WithClientUnaryInterceptor(timeout))
	if err != nil {
		return nil, err
	}

	ret := &FlowProxy{}
	ret.upstream = access.NewAccessAPIClient(collectionRPCConn)
	ret.fallback = fallback
	return ret, nil
}

type FlowProxy struct {
	access.UnimplementedAccessAPIServer
	upstream   access.AccessAPIClient
	fallback   *access.AccessAPIServer
	downstream *grpc.Server
}

func (h FlowProxy) Ping(context context.Context, req *access.PingRequest) (*access.PingResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
	}
	return h.Ping(context, req)
}

func (h FlowProxy) GetLatestBlockHeader(context context.Context, req *access.GetLatestBlockHeaderRequest) (*access.BlockHeaderResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetLatestBlockHeader not implemented")
	}
	return h.GetLatestBlockHeader(context, req)
}

func (h FlowProxy) GetBlockHeaderByID(context context.Context, req *access.GetBlockHeaderByIDRequest) (*access.BlockHeaderResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetBlockHeaderByID not implemented")
	}
	return h.GetBlockHeaderByID(context, req)
}

func (h FlowProxy) GetBlockHeaderByHeight(context context.Context, req *access.GetBlockHeaderByHeightRequest) (*access.BlockHeaderResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetBlockHeaderByHeight not implemented")
	}
	return h.GetBlockHeaderByHeight(context, req)
}

func (h FlowProxy) GetLatestBlock(context context.Context, req *access.GetLatestBlockRequest) (*access.BlockResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetLatestBlock not implemented")
	}
	return h.GetLatestBlock(context, req)
}

func (h FlowProxy) GetBlockByID(context context.Context, req *access.GetBlockByIDRequest) (*access.BlockResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetBlockByID not implemented")
	}
	return h.GetBlockByID(context, req)
}

func (h FlowProxy) GetBlockByHeight(context context.Context, req *access.GetBlockByHeightRequest) (*access.BlockResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetBlockByHeight not implemented")
	}
	return h.GetBlockByHeight(context, req)
}

func (h FlowProxy) GetCollectionByID(context context.Context, req *access.GetCollectionByIDRequest) (*access.CollectionResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetCollectionByID not implemented")
	}
	return h.GetCollectionByID(context, req)
}

func (h FlowProxy) SendTransaction(context context.Context, req *access.SendTransactionRequest) (*access.SendTransactionResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method SendTransaction not implemented")
	}
	return h.SendTransaction(context, req)
}

func (h FlowProxy) GetTransaction(context context.Context, req *access.GetTransactionRequest) (*access.TransactionResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetTransaction not implemented")
	}
	return h.GetTransaction(context, req)
}

func (h FlowProxy) GetTransactionResult(context context.Context, req *access.GetTransactionRequest) (*access.TransactionResultResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetTransactionResult not implemented")
	}
	return h.GetTransactionResult(context, req)
}

func (h FlowProxy) GetTransactionResultByIndex(context context.Context, req *access.GetTransactionByIndexRequest) (*access.TransactionResultResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetTransactionResultByIndex not implemented")
	}
	return h.GetTransactionResultByIndex(context, req)
}

func (h FlowProxy) GetAccount(context context.Context, req *access.GetAccountRequest) (*access.GetAccountResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
	}
	return h.GetAccount(context, req)
}

func (h FlowProxy) GetAccountAtLatestBlock(context context.Context, req *access.GetAccountAtLatestBlockRequest) (*access.AccountResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetAccountAtLatestBlock not implemented")
	}
	return h.GetAccountAtLatestBlock(context, req)
}

func (h FlowProxy) GetAccountAtBlockHeight(context context.Context, req *access.GetAccountAtBlockHeightRequest) (*access.AccountResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetAccountAtBlockHeight not implemented")
	}
	return h.GetAccountAtBlockHeight(context, req)
}

func (h FlowProxy) ExecuteScriptAtLatestBlock(context context.Context, req *access.ExecuteScriptAtLatestBlockRequest) (*access.ExecuteScriptResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method ExecuteScriptAtLatestBlock not implemented")
	}
	return h.ExecuteScriptAtLatestBlock(context, req)
}

func (h FlowProxy) ExecuteScriptAtBlockID(context context.Context, req *access.ExecuteScriptAtBlockIDRequest) (*access.ExecuteScriptResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method ExecuteScriptAtBlockID not implemented")
	}
	return h.ExecuteScriptAtBlockID(context, req)
}

func (h FlowProxy) ExecuteScriptAtBlockHeight(context context.Context, req *access.ExecuteScriptAtBlockHeightRequest) (*access.ExecuteScriptResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method ExecuteScriptAtBlockHeight not implemented")
	}
	return h.ExecuteScriptAtBlockHeight(context, req)
}

func (h FlowProxy) GetEventsForHeightRange(context context.Context, req *access.GetEventsForHeightRangeRequest) (*access.EventsResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetEventsForHeightRange not implemented")
	}
	return h.GetEventsForHeightRange(context, req)
}

func (h FlowProxy) GetEventsForBlockIDs(context context.Context, req *access.GetEventsForBlockIDsRequest) (*access.EventsResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetEventsForBlockIDs not implemented")
	}
	return h.GetEventsForBlockIDs(context, req)
}

func (h FlowProxy) GetNetworkParameters(context context.Context, req *access.GetNetworkParametersRequest) (*access.GetNetworkParametersResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetNetworkParameters not implemented")
	}
	return h.GetNetworkParameters(context, req)
}

func (h FlowProxy) GetLatestProtocolStateSnapshot(context context.Context, req *access.GetLatestProtocolStateSnapshotRequest) (*access.ProtocolStateSnapshotResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetLatestProtocolStateSnapshot not implemented")
	}
	return h.GetLatestProtocolStateSnapshot(context, req)
}

func (h FlowProxy) GetExecutionResultForBlockID(context context.Context, req *access.GetExecutionResultForBlockIDRequest) (*access.ExecutionResultForBlockIDResponse, error) {
	if h.upstream == nil {
		return nil, status.Errorf(codes.Unimplemented, "method GetExecutionResultForBlockID not implemented")
	}
	return h.GetExecutionResultForBlockID(context, req)
}
