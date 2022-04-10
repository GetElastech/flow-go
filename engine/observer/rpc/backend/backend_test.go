package backend

import (
	"context"
	"math/rand"
	"testing"
	"time"

	accessproto "github.com/onflow/flow/protobuf/go/flow/access"
	entitiesproto "github.com/onflow/flow/protobuf/go/flow/entities"
	execproto "github.com/onflow/flow/protobuf/go/flow/execution"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	access "github.com/onflow/flow-go/engine/access/mock"
	backendmock "github.com/onflow/flow-go/engine/access/rpc/backend/mock"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/metrics"
	protocol "github.com/onflow/flow-go/state/protocol/mock"
	storagemock "github.com/onflow/flow-go/storage/mock"
	"github.com/onflow/flow-go/utils/unittest"
)

type Suite struct {
	suite.Suite

	state    *protocol.State
	snapshot *protocol.Snapshot
	log      zerolog.Logger

	blocks                 *storagemock.Blocks
	headers                *storagemock.Headers
	collections            *storagemock.Collections
	transactions           *storagemock.Transactions
	receipts               *storagemock.ExecutionReceipts
	results                *storagemock.ExecutionResults
	colClient              *access.AccessAPIClient
	execClient             *access.ExecutionAPIClient
	historicalAccessClient *access.AccessAPIClient
	connectionFactory      *backendmock.ConnectionFactory
	chainID                flow.ChainID
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	suite.log = zerolog.New(zerolog.NewConsoleWriter())
	suite.state = new(protocol.State)
	suite.snapshot = new(protocol.Snapshot)
	header := unittest.BlockHeaderFixture()
	params := new(protocol.Params)
	params.On("Root").Return(&header, nil)
	suite.state.On("Params").Return(params).Maybe()
	suite.blocks = new(storagemock.Blocks)
	suite.headers = new(storagemock.Headers)
	suite.transactions = new(storagemock.Transactions)
	suite.collections = new(storagemock.Collections)
	suite.receipts = new(storagemock.ExecutionReceipts)
	suite.results = new(storagemock.ExecutionResults)
	suite.colClient = new(access.AccessAPIClient)
	suite.execClient = new(access.ExecutionAPIClient)
	suite.chainID = flow.Testnet
	suite.historicalAccessClient = new(access.AccessAPIClient)
	suite.connectionFactory = new(backendmock.ConnectionFactory)
}

func (suite *Suite) TestPing() {
	suite.colClient.
		On("Ping", mock.Anything, &accessproto.PingRequest{}).
		Return(&accessproto.PingResponse{}, nil)

	backend := New(
		suite.colClient,
	)

	err := backend.Ping(context.Background())

	suite.Require().NoError(err)
}

func (suite *Suite) TestGetLatestSealedBlockHeader() {
	// setup the mocks
	suite.state.On("Sealed").Return(suite.snapshot, nil).Maybe()

	block := unittest.BlockHeaderFixture()
	suite.snapshot.On("Head").Return(&block, nil).Once()

	backend := &Backend{
		state:                suite.state,
		chainID:              suite.chainID,
		transactionMetrics:   metrics.NewNoopCollector(),
		maxHeightRange:       DefaultMaxHeightRange,
		log:                  suite.log,
		snapshot:             suite.snapshot,
		snapshotHistoryLimit: DefaultSnapshotHistoryLimit,
	}

	// query the handler for the latest sealed block
	header, err := backend.GetLatestBlockHeader(context.Background(), true)
	suite.checkResponse(header, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(block.ID(), header.ID())
	suite.Require().Equal(block.Height, header.Height)
	suite.Require().Equal(block.ParentID, header.ParentID)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockHeaderByID() {
	// setup the mocks
	suite.state.On("Sealed").Return(suite.snapshot, nil).Maybe()

	block := unittest.BlockHeaderFixture()
	suite.snapshot.On("Head").Return(&block, nil).Once()

	backend := &Backend{
		state:                suite.state,
		chainID:              suite.chainID,
		transactionMetrics:   metrics.NewNoopCollector(),
		maxHeightRange:       DefaultMaxHeightRange,
		log:                  suite.log,
		snapshot:             suite.snapshot,
		snapshotHistoryLimit: DefaultSnapshotHistoryLimit,
	}
	// query the handler for the latest sealed block
	header, err := backend.GetBlockHeaderByID(context.Background(), block.ID())
	suite.checkResponse(header, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(block.ID(), header.ID())
	suite.Require().Equal(block.Height, header.Height)
	suite.Require().Equal(block.ParentID, header.ParentID)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockByHeight() {
	// setup the mocks
	suite.state.On("Sealed").Return(suite.snapshot, nil).Maybe()

	block := unittest.BlockHeaderFixture()
	suite.snapshot.On("Head").Return(&block, nil).Once()

	backend := &Backend{
		state:                suite.state,
		chainID:              suite.chainID,
		transactionMetrics:   metrics.NewNoopCollector(),
		maxHeightRange:       DefaultMaxHeightRange,
		log:                  suite.log,
		snapshot:             suite.snapshot,
		snapshotHistoryLimit: DefaultSnapshotHistoryLimit,
	}

	// query the handler for the latest sealed block
	header, err := backend.GetBlockByHeight(context.Background(), block.Height)
	suite.checkResponse(header, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(block.ID(), header.ID())

	suite.assertAllExpectations()
}

func (suite *Suite) TestSendTransaction() {
	block := unittest.BlockHeaderFixture()
	transaction := unittest.TransactionFixture()
	transaction.SetReferenceBlockID(block.ID())

	suite.colClient.
		On("SendTransaction", mock.Anything, mock.Anything).
		Return(&accessproto.SendTransactionResponse{}, nil).
		Once()

	backend := &Backend{
		state:                suite.state,
		chainID:              suite.chainID,
		staticCollectionRPC:  suite.colClient,
		transactionMetrics:   metrics.NewNoopCollector(),
		maxHeightRange:       DefaultMaxHeightRange,
		transactions:         suite.transactions,
		log:                  suite.log,
		snapshot:             suite.snapshot,
		snapshotHistoryLimit: DefaultSnapshotHistoryLimit,
	}
	// Send transaction
	err := backend.SendTransaction(context.Background(), &transaction.TransactionBody)
	suite.Require().NoError(err)
}

func (suite *Suite) TestGetTransaction() {
	suite.state.On("Sealed").Return(suite.snapshot, nil).Maybe()

	transaction := unittest.TransactionFixture()
	expected := transaction.TransactionBody

	suite.transactions.
		On("ByID", transaction.ID()).
		Return(&expected, nil).
		Once()

	backend := &Backend{
		state:                suite.state,
		transactions:         suite.transactions,
		chainID:              suite.chainID,
		transactionMetrics:   metrics.NewNoopCollector(),
		maxHeightRange:       DefaultMaxHeightRange,
		log:                  suite.log,
		snapshotHistoryLimit: DefaultSnapshotHistoryLimit,
	}

	actual, err := backend.GetTransaction(context.Background(), transaction.ID())
	suite.checkResponse(actual, err)

	suite.Require().Equal(expected, *actual)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetAccount() {
	suite.state.On("Sealed").Return(suite.snapshot, nil).Maybe()
	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()

	address, err := suite.chainID.Chain().NewAddressGenerator().NextAddress()
	suite.Require().NoError(err)

	account := &entitiesproto.Account{
		Address: address.Bytes(),
	}
	ctx := context.Background()

	// setup the latest sealed block
	block := unittest.BlockFixture()
	header := block.Header          // create a mock header
	seal := unittest.Seal.Fixture() // create a mock seal
	seal.BlockID = header.ID()      // make the seal point to the header

	suite.snapshot.
		On("Head").
		Return(header, nil).
		Once()

	// create the expected execution API request
	blockID := header.ID()
	exeReq := &execproto.GetAccountAtBlockIDRequest{
		BlockId: blockID[:],
		Address: address.Bytes(),
	}

	// create the expected execution API response
	exeResp := &execproto.GetAccountAtBlockIDResponse{
		Account: account,
	}

	// setup the execution client mock
	suite.execClient.
		On("GetAccountAtBlockID", ctx, exeReq).
		Return(exeResp, nil).
		Once()

	receipts, ids := suite.setupReceipts(&block)

	suite.snapshot.On("Identities", mock.Anything).Return(ids, nil)
	// create a mock connection factory
	connFactory := new(backendmock.ConnectionFactory)
	connFactory.On("GetExecutionAPIClient", mock.Anything).Return(suite.execClient, &mockCloser{}, nil)

	// create the handler with the mock
	backend := &Backend{
		state:                suite.state,
		headers:              suite.headers,
		executionReceipts:    suite.receipts,
		executionResults:     suite.results,
		transactions:         suite.transactions,
		chainID:              suite.chainID,
		transactionMetrics:   metrics.NewNoopCollector(),
		connFactory:          connFactory,
		maxHeightRange:       DefaultMaxHeightRange,
		log:                  suite.log,
		snapshotHistoryLimit: DefaultSnapshotHistoryLimit,
	}
	preferredENIdentifiers = flow.IdentifierList{receipts[0].ExecutorID}

	suite.Run("happy path - valid request and valid response", func() {
		account, err := backend.GetAccountAtLatestBlock(ctx, address)
		suite.checkResponse(account, err)

		suite.Require().Equal(address, account.Address)

		suite.assertAllExpectations()
	})
}

func (suite *Suite) TestGetAccountAtBlockHeight() {
	suite.state.On("Sealed").Return(suite.snapshot, nil).Maybe()
	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()

	height := uint64(5)
	address := unittest.AddressFixture()
	account := &entitiesproto.Account{
		Address: address.Bytes(),
	}
	ctx := context.Background()

	// create a mock block header
	b := unittest.BlockFixture()
	h := b.Header

	// setup headers storage to return the header when queried by height
	suite.headers.
		On("ByHeight", height).
		Return(h, nil).
		Once()

	receipts, ids := suite.setupReceipts(&b)
	suite.snapshot.On("Identities", mock.Anything).Return(ids, nil)

	// create a mock connection factory
	connFactory := new(backendmock.ConnectionFactory)
	connFactory.On("GetExecutionAPIClient", mock.Anything).Return(suite.execClient, &mockCloser{}, nil)

	// create the expected execution API request
	blockID := h.ID()
	exeReq := &execproto.GetAccountAtBlockIDRequest{
		BlockId: blockID[:],
		Address: address.Bytes(),
	}

	// create the expected execution API response
	exeResp := &execproto.GetAccountAtBlockIDResponse{
		Account: account,
	}

	// setup the execution client mock
	suite.execClient.
		On("GetAccountAtBlockID", ctx, exeReq).
		Return(exeResp, nil).
		Once()

	// create the handler with the mock
	backend := &Backend{
		state:                suite.state,
		headers:              suite.headers,
		executionReceipts:    suite.receipts,
		executionResults:     suite.results,
		transactions:         suite.transactions,
		chainID:              flow.Testnet,
		transactionMetrics:   metrics.NewNoopCollector(),
		connFactory:          connFactory,
		maxHeightRange:       DefaultMaxHeightRange,
		log:                  suite.log,
		snapshotHistoryLimit: DefaultSnapshotHistoryLimit,
	}

	preferredENIdentifiers = flow.IdentifierList{receipts[0].ExecutorID}

	suite.Run("happy path - valid request and valid response", func() {
		account, err := backend.GetAccountAtBlockHeight(ctx, address, height)
		suite.checkResponse(account, err)

		suite.Require().Equal(address, account.Address)

		suite.assertAllExpectations()
	})
}

// TestExecuteScriptOnExecutionNode tests the method backend.scripts.executeScriptOnExecutionNode for script execution
func (suite *Suite) TestExecuteScriptOnExecutionNode() {

	// create a mock connection factory
	connFactory := new(backendmock.ConnectionFactory)
	connFactory.On("GetExecutionAPIClient", mock.Anything).
		Return(suite.execClient, &mockCloser{}, nil)

	// create the handler with the mock
	backend := &Backend{
		state:                suite.state,
		headers:              suite.headers,
		executionReceipts:    suite.receipts,
		executionResults:     suite.results,
		transactions:         suite.transactions,
		chainID:              flow.Mainnet,
		transactionMetrics:   metrics.NewNoopCollector(),
		connFactory:          connFactory,
		maxHeightRange:       DefaultMaxHeightRange,
		log:                  suite.log,
		snapshotHistoryLimit: DefaultSnapshotHistoryLimit,
	}

	// mock parameters
	ctx := context.Background()
	block := unittest.BlockFixture()
	blockID := block.ID()
	script := []byte("dummy script")
	arguments := [][]byte(nil)
	executionNode := unittest.IdentityFixture(unittest.WithRole(flow.RoleExecution))
	execReq := execproto.ExecuteScriptAtBlockIDRequest{
		BlockId:   blockID[:],
		Script:    script,
		Arguments: arguments,
	}
	execRes := execproto.ExecuteScriptAtBlockIDResponse{
		Value: []byte{4, 5, 6},
	}

	suite.Run("happy path script execution success", func() {
		suite.execClient.On("ExecuteScriptAtBlockID", ctx, &execReq).Return(&execRes, nil).Once()
		res, err := backend.tryExecuteScript(ctx, executionNode, execReq)
		suite.execClient.AssertExpectations(suite.T())
		suite.checkResponse(res, err)
	})

	suite.Run("script execution failure returns status OK", func() {
		suite.execClient.On("ExecuteScriptAtBlockID", ctx, &execReq).
			Return(nil, status.Error(codes.InvalidArgument, "execution failure!")).Once()
		_, err := backend.tryExecuteScript(ctx, executionNode, execReq)
		suite.execClient.AssertExpectations(suite.T())
		suite.Require().Error(err)
		suite.Require().Equal(status.Code(err), codes.InvalidArgument)
	})

	suite.Run("execution node internal failure returns status code Internal", func() {
		suite.execClient.On("ExecuteScriptAtBlockID", ctx, &execReq).
			Return(nil, status.Error(codes.Internal, "execution node internal error!")).Once()
		_, err := backend.tryExecuteScript(ctx, executionNode, execReq)
		suite.execClient.AssertExpectations(suite.T())
		suite.Require().Error(err)
		suite.Require().Equal(status.Code(err), codes.Internal)
	})
}

type mockCloser struct{}

func (mc *mockCloser) Close() error { return nil }

func (suite *Suite) setupReceipts(block *flow.Block) ([]*flow.ExecutionReceipt, flow.IdentityList) {
	ids := unittest.IdentityListFixture(2, unittest.WithRole(flow.RoleExecution))
	receipt1 := unittest.ReceiptForBlockFixture(block)
	receipt1.ExecutorID = ids[0].NodeID
	receipt2 := unittest.ReceiptForBlockFixture(block)
	receipt2.ExecutorID = ids[1].NodeID
	receipt1.ExecutionResult = receipt2.ExecutionResult

	receipts := flow.ExecutionReceiptList{receipt1, receipt2}
	suite.receipts.
		On("ByBlockID", block.ID()).
		Return(receipts, nil)

	return receipts, ids
}

func (suite *Suite) checkResponse(resp interface{}, err error) {
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}

func (suite *Suite) assertAllExpectations() {
	suite.snapshot.AssertExpectations(suite.T())
	suite.state.AssertExpectations(suite.T())
	suite.blocks.AssertExpectations(suite.T())
	suite.headers.AssertExpectations(suite.T())
	suite.collections.AssertExpectations(suite.T())
	suite.transactions.AssertExpectations(suite.T())
	suite.execClient.AssertExpectations(suite.T())
}
