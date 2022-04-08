package backend

import (
	"context"
	"math/rand"
	"testing"
	"time"

	accessproto "github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

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
		log:                  suite.log,
		snapshot:             suite.snapshot,
		snapshotHistoryLimit: DefaultSnapshotHistoryLimit,
	}
	// Send transaction
	err := backend.SendTransaction(context.Background(), &transaction.TransactionBody)
	suite.Require().NoError(err)
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
