package protocol

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	protocol "github.com/onflow/flow-go/state/protocol/mock"
	"github.com/onflow/flow-go/storage"
	storagemock "github.com/onflow/flow-go/storage/mock"
	"github.com/onflow/flow-go/utils/unittest"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var CodesNotFoundErr = status.Errorf(codes.NotFound, "not found")
var StorageNotFoundErr = status.Errorf(codes.NotFound, "not found: %v", storage.ErrNotFound)
var InternalErr = status.Errorf(codes.Internal, "internal")
var OutputInternalErr = status.Errorf(codes.Internal, "failed to find: %v", status.Errorf(codes.Internal, "internal"))

type Suite struct {
	suite.Suite

	state            *protocol.State
	blocks           *storagemock.Blocks
	headers          *storagemock.Headers
	executionResults *storagemock.ExecutionResults
	snapshot         *protocol.Snapshot
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	suite.snapshot = new(protocol.Snapshot)

	suite.state = new(protocol.State)
	suite.blocks = new(storagemock.Blocks)
	suite.headers = new(storagemock.Headers)
	suite.executionResults = new(storagemock.ExecutionResults)
}

func (suite *Suite) TestGetLatestFinalizedBlock_Success() {
	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()

	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.snapshot.
		On("Head").
		Return(header, nil).
		Once()

	suite.blocks.
		On("ByID", header.ID()).
		Return(&block, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized header
	respBlock, err := backend.GetLatestBlock(context.Background(), false)
	suite.checkResponse(respBlock, err)

	// make sure we got the latest header
	suite.Require().Equal(block, *respBlock)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetLatestSealedBlock_Success() {
	suite.state.On("Sealed").Return(suite.snapshot, nil).Maybe()

	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.snapshot.
		On("Head").
		Return(header, nil).
		Once()

	suite.blocks.
		On("ByID", header.ID()).
		Return(&block, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized header
	respBlock, err := backend.GetLatestBlock(context.Background(), true)
	suite.checkResponse(respBlock, err)

	// make sure we got the latest header
	suite.Require().Equal(block, *respBlock)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetLatestBlock_StorageNotFoundFailure() {
	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()

	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.snapshot.
		On("Head").
		Return(header, storage.ErrNotFound).
		Once()

	suite.blocks.
		On("ByID", header.ID()).
		Return(&block, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized header
	_, err := backend.GetLatestBlock(context.Background(), false)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, StorageNotFoundErr)
}

func (suite *Suite) TestGetLatestBlock_CodesNotFoundFailure() {
	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()

	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.snapshot.
		On("Head").
		Return(header, CodesNotFoundErr).
		Once()

	suite.blocks.
		On("ByID", header.ID()).
		Return(&block, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized header
	_, err := backend.GetLatestBlock(context.Background(), false)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, CodesNotFoundErr)
}

func (suite *Suite) TestGetLatestBlock_InternalFailure() {
	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()

	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.snapshot.
		On("Head").
		Return(header, InternalErr).
		Once()

	suite.blocks.
		On("ByID", header.ID()).
		Return(&block, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized header
	_, err := backend.GetLatestBlock(context.Background(), false)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, OutputInternalErr)
}

func (suite *Suite) TestGetBlockById_Success() {
	// setup the mocks
	block := unittest.BlockFixture()

	suite.blocks.
		On("ByID", block.ID()).
		Return(&block, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	respBlock, err := backend.GetBlockByID(context.Background(), block.ID())
	suite.checkResponse(respBlock, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(block.ID(), respBlock.ID())

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockById_StorageNotFoundFailure() {
	// setup the mocks
	block := unittest.BlockFixture()

	suite.blocks.
		On("ByID", block.ID()).
		Return(&block, storage.ErrNotFound).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockByID(context.Background(), block.ID())
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, StorageNotFoundErr)
}

func (suite *Suite) TestGetBlockById_CodesNotFoundFailure() {
	// setup the mocks
	block := unittest.BlockFixture()

	suite.blocks.
		On("ByID", block.ID()).
		Return(&block, CodesNotFoundErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockByID(context.Background(), block.ID())
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, CodesNotFoundErr)
}

func (suite *Suite) TestGetBlockById_InternalFailure() {
	// setup the mocks
	block := unittest.BlockFixture()

	suite.blocks.
		On("ByID", block.ID()).
		Return(&block, InternalErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockByID(context.Background(), block.ID())
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, OutputInternalErr)
}

func (suite *Suite) TestGetBlockByHeight_Success() {
	// setup the mocks
	block := unittest.BlockFixture()
	height := block.Header.Height

	suite.blocks.
		On("ByHeight", height).
		Return(&block, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	respBlock, err := backend.GetBlockByHeight(context.Background(), height)
	suite.checkResponse(respBlock, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(block.ID(), respBlock.ID())

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockByHeight_StorageNotFoundFailure() {
	// setup the mocks
	block := unittest.BlockFixture()
	height := block.Header.Height

	suite.blocks.
		On("ByHeight", height).
		Return(&block, StorageNotFoundErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockByHeight(context.Background(), height)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, StorageNotFoundErr)
}

func (suite *Suite) TestGetBlockByHeight_CodesNotFoundFailure() {
	// setup the mocks
	block := unittest.BlockFixture()
	height := block.Header.Height

	suite.blocks.
		On("ByHeight", height).
		Return(&block, CodesNotFoundErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockByHeight(context.Background(), height)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, CodesNotFoundErr)
}

func (suite *Suite) TestGetBlockByHeight_InternalFailure() {
	// setup the mocks
	block := unittest.BlockFixture()
	height := block.Header.Height

	suite.blocks.
		On("ByHeight", height).
		Return(&block, InternalErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockByHeight(context.Background(), height)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, OutputInternalErr)
}

func (suite *Suite) TestGetLatestFinalizedBlockHeader_Success() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()

	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()
	suite.snapshot.On("Head").Return(&blockHeader, nil).Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized block
	respHeader, err := backend.GetLatestBlockHeader(context.Background(), false)
	suite.checkResponse(respHeader, err)

	// make sure we got the latest block
	suite.Require().Equal(blockHeader.ID(), respHeader.ID())
	suite.Require().Equal(blockHeader.Height, respHeader.Height)
	suite.Require().Equal(blockHeader.ParentID, respHeader.ParentID)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetLatestSealedBlockHeader_Success() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()

	suite.state.On("Sealed").Return(suite.snapshot, nil).Maybe()
	suite.snapshot.On("Head").Return(&blockHeader, nil).Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized block
	respHeader, err := backend.GetLatestBlockHeader(context.Background(), true)
	suite.checkResponse(respHeader, err)

	// make sure we got the latest block
	suite.Require().Equal(blockHeader.ID(), respHeader.ID())
	suite.Require().Equal(blockHeader.Height, respHeader.Height)
	suite.Require().Equal(blockHeader.ParentID, respHeader.ParentID)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetLatestBlockHeader_StorageNotFoundFailure() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()

	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()
	suite.snapshot.On("Head").Return(&blockHeader, storage.ErrNotFound).Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized block
	_, err := backend.GetLatestBlockHeader(context.Background(), false)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, StorageNotFoundErr)
}

func (suite *Suite) TestGetLatestBlockHeader_CodesNotFoundFailure() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()

	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()
	suite.snapshot.On("Head").Return(&blockHeader, CodesNotFoundErr).Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized block
	_, err := backend.GetLatestBlockHeader(context.Background(), false)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, CodesNotFoundErr)
}

func (suite *Suite) TestGetLatestBlockHeader_InternalFailure() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()

	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()
	suite.snapshot.On("Head").Return(&blockHeader, InternalErr).Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest finalized block
	_, err := backend.GetLatestBlockHeader(context.Background(), false)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, OutputInternalErr)
}

func (suite *Suite) TestGetBlockHeaderByID_Success() {
	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.headers.
		On("ByBlockID", block.ID()).
		Return(header, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	respHeader, err := backend.GetBlockHeaderByID(context.Background(), block.ID())

	suite.checkResponse(respHeader, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(header.ID(), respHeader.ID())
	suite.Require().Equal(header.Height, respHeader.Height)
	suite.Require().Equal(header.ParentID, respHeader.ParentID)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockHeaderByID_StorageNotFoundFailure() {
	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.headers.
		On("ByBlockID", block.ID()).
		Return(header, StorageNotFoundErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockHeaderByID(context.Background(), block.ID())
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, StorageNotFoundErr)
}

func (suite *Suite) TestGetBlockHeaderByID_CodesNotFoundFailure() {
	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.headers.
		On("ByBlockID", block.ID()).
		Return(header, CodesNotFoundErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockHeaderByID(context.Background(), block.ID())
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, CodesNotFoundErr)
}

func (suite *Suite) TestGetBlockHeaderByID_InternalFailure() {
	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.headers.
		On("ByBlockID", block.ID()).
		Return(header, InternalErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockHeaderByID(context.Background(), block.ID())
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, OutputInternalErr)
}

func (suite *Suite) TestGetBlockHeaderByHeight_Success() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()
	headerHeight := blockHeader.Height

	suite.headers.
		On("ByHeight", headerHeight).
		Return(&blockHeader, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	respHeader, err := backend.GetBlockHeaderByHeight(context.Background(), headerHeight)

	suite.checkResponse(respHeader, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(blockHeader.Height, respHeader.Height)
	suite.Require().Equal(blockHeader.ParentID, respHeader.ParentID)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockHeaderByHeight_StorageNotFoundFailure() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()
	headerHeight := blockHeader.Height

	suite.headers.
		On("ByHeight", headerHeight).
		Return(&blockHeader, StorageNotFoundErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockHeaderByHeight(context.Background(), headerHeight)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, StorageNotFoundErr)
}

func (suite *Suite) TestGetBlockHeaderByHeight_CodesNotFoundFailure() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()
	headerHeight := blockHeader.Height

	suite.headers.
		On("ByHeight", headerHeight).
		Return(&blockHeader, CodesNotFoundErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockHeaderByHeight(context.Background(), headerHeight)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, CodesNotFoundErr)
}

func (suite *Suite) TestGetBlockHeaderByHeight_InternalFailure() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()
	headerHeight := blockHeader.Height

	suite.headers.
		On("ByHeight", headerHeight).
		Return(&blockHeader, InternalErr).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

	// query the handler for the latest sealed block
	_, err := backend.GetBlockHeaderByHeight(context.Background(), headerHeight)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, OutputInternalErr)
}

func (suite *Suite) TestGetExecutionResultsByBlockID() {

	nonexistingBlockID := unittest.IdentifierFixture()
	blockID := unittest.IdentifierFixture()

	executionResult := unittest.ExecutionResultFixture(
		unittest.WithExecutionResultBlockID(blockID))

	ctx := context.Background()

	suite.executionResults.
		On("ByBlockID", blockID).
		Return(executionResult, nil)

	suite.executionResults.
		On("ByBlockID", nonexistingBlockID).
		Return(nil, storage.ErrNotFound)

	suite.Run("success retrieving existing execution results", func() {

		// create the handler
		backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

		responseExecutionResult, err := backend.GetExecutionResultByBlockID(ctx, blockID)
		suite.checkResponse(responseExecutionResult, err)

		// make sure we got the correct execution results
		suite.Require().Equal(executionResult, responseExecutionResult)

	})

	suite.Run("failure retreiving nonexisting execution results", func() {

		// create the handler
		backend := New(suite.state, suite.blocks, suite.headers, suite.executionResults)

		_, err := backend.GetExecutionResultByBlockID(ctx, nonexistingBlockID)
		suite.Require().Error(err)
		suite.Require().ErrorIs(err, StorageNotFoundErr)
	})

	suite.assertAllExpectations()
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
}
