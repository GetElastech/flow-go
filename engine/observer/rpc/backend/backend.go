package backend

import (
	"context"
	"fmt"

	accessproto "github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/access"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/storage"
)

// maxExecutionNodesCnt is the max number of execution nodes that will be contacted to complete an execution api request
const maxExecutionNodesCnt = 3

// minExecutionNodesCnt is the minimum number of execution nodes expected to have sent the execution receipt for a block
const minExecutionNodesCnt = 2

// maxAttemptsForExecutionReceipt is the maximum number of attempts to find execution receipts for a given block ID
const maxAttemptsForExecutionReceipt = 3

// DefaultMaxHeightRange is the default maximum size of range requests.
const DefaultMaxHeightRange = 250

// DefaultSnapshotHistoryLimit the amount of blocks to look back in state
// when recursively searching for a valid snapshot
const DefaultSnapshotHistoryLimit = 50

var preferredENIdentifiers flow.IdentifierList

type Backend struct {
	access.API
	blocks               storage.Blocks
	headers              storage.Headers
	collections          storage.Collections
	transactions         storage.Transactions
	executionReceipts    storage.ExecutionReceipts
	executionResults     storage.ExecutionResults
	staticCollectionRPC  accessproto.AccessAPIClient
	state                protocol.State
	snapshot             protocol.Snapshot
	chainID              flow.ChainID
	log                  zerolog.Logger
	transactionMetrics   module.TransactionMetrics
	connFactory          ConnectionFactory
	maxHeightRange       uint
	snapshotHistoryLimit int
}

func New(
	collectionRPC accessproto.AccessAPIClient,
) *Backend {
	return &Backend{
		staticCollectionRPC: collectionRPC,
	}
}

// Ping responds to requests when the server is up.
func (b *Backend) Ping(ctx context.Context) error {
	// staticCollectionRPC is only set if a collection node address was provided at startup
	if b.staticCollectionRPC != nil {
		_, err := b.staticCollectionRPC.Ping(ctx, &accessproto.PingRequest{})
		if err != nil {
			return fmt.Errorf("could not ping collection node: %w", err)
		}
	}
	return nil
}
