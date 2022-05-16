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

var InternalErr = status.Errorf(codes.Internal, "failed to find: %v", status.Errorf(codes.Internal, "internal"))
var CodesNotFoundErr = status.Errorf(codes.NotFound, "not found: %v", storage.ErrNotFound)

type Suite struct {
	suite.Suite

	state   *protocol.State
	blocks  *storagemock.Blocks
	headers *storagemock.Headers

	snapshot *protocol.Snapshot
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
}

func (suite *Suite) TestGetLatestFinalizedBlockHeader_Success() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()

	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()
	suite.snapshot.On("Head").Return(&blockHeader, nil).Once()

	backend := New(suite.state, suite.blocks, suite.headers)

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

	backend := New(suite.state, suite.blocks, suite.headers)

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

	backend := New(suite.state, suite.blocks, suite.headers)

	// query the handler for the latest finalized block
	_, err := backend.GetLatestBlockHeader(context.Background(), false)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, CodesNotFoundErr)
}

func (suite *Suite) TestGetLatestBlockHeader_CodesNotFoundFailure() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()

	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()
	suite.snapshot.On("Head").Return(&blockHeader, status.Errorf(codes.NotFound, "not found")).Once()

	backend := New(suite.state, suite.blocks, suite.headers)

	// query the handler for the latest finalized block
	_, err := backend.GetLatestBlockHeader(context.Background(), false)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, status.Errorf(codes.NotFound, "not found"))
}

func (suite *Suite) TestGetLatestBlockHeader_InternalFailure() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()

	suite.state.On("Final").Return(suite.snapshot, nil).Maybe()
	suite.snapshot.On("Head").Return(&blockHeader, status.Errorf(codes.Internal, "internal")).Once()

	backend := New(suite.state, suite.blocks, suite.headers)

	// query the handler for the latest finalized block
	_, err := backend.GetLatestBlockHeader(context.Background(), false)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, InternalErr)
}

func (suite *Suite) TestGetLatestFinalizedBlock() {
	suite.state.On("Sealed").Return(suite.snapshot, nil).Maybe()
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

	backend := New(suite.state, suite.blocks, suite.headers)

	// query the handler for the latest finalized header
	respBlock, err := backend.GetLatestBlock(context.Background(), false)
	suite.checkResponse(respBlock, err)

	// make sure we got the latest header
	suite.Require().Equal(block, *respBlock)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockHeaderByID() {
	// setup the mocks
	block := unittest.BlockFixture()
	header := block.Header

	suite.headers.
		On("ByBlockID", block.ID()).
		Return(header, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers)

	// query the handler for the latest sealed block
	respHeader, err := backend.GetBlockHeaderByID(context.Background(), block.ID())

	suite.checkResponse(respHeader, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(header.ID(), respHeader.ID())
	suite.Require().Equal(header.Height, respHeader.Height)
	suite.Require().Equal(header.ParentID, respHeader.ParentID)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockHeaderByHeight() {
	// setup the mocks
	blockHeader := unittest.BlockHeaderFixture()
	headerHeight := blockHeader.Height

	suite.headers.
		On("ByHeight", headerHeight).
		Return(&blockHeader, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers)

	// query the handler for the latest sealed block
	respHeader, err := backend.GetBlockHeaderByHeight(context.Background(), headerHeight)

	suite.checkResponse(respHeader, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(blockHeader.Height, respHeader.Height)
	suite.Require().Equal(blockHeader.ParentID, respHeader.ParentID)

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockByHeight() {
	// setup the mocks
	block := unittest.BlockFixture()
	height := block.Header.Height

	suite.blocks.
		On("ByHeight", height).
		Return(&block, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers)

	// query the handler for the latest sealed block
	respBlock, err := backend.GetBlockByHeight(context.Background(), height)
	suite.checkResponse(respBlock, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(block.ID(), respBlock.ID())

	suite.assertAllExpectations()
}

func (suite *Suite) TestGetBlockById() {
	// setup the mocks
	block := unittest.BlockFixture()

	suite.blocks.
		On("ByID", block.ID()).
		Return(&block, nil).
		Once()

	backend := New(suite.state, suite.blocks, suite.headers)

	// query the handler for the latest sealed block
	respBlock, err := backend.GetBlockByID(context.Background(), block.ID())
	suite.checkResponse(respBlock, err)

	// make sure we got the latest sealed block
	suite.Require().Equal(block.ID(), respBlock.ID())

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
