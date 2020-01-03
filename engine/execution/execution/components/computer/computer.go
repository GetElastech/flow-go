package computer

import (
	"errors"
	"fmt"

	"github.com/dapperlabs/flow-go/crypto"
	"github.com/dapperlabs/flow-go/engine/execution/execution/modules/context"
	"github.com/dapperlabs/flow-go/language/runtime"
	"github.com/dapperlabs/flow-go/model/flow"
)

// A TransactionResult is the result of executing a transaction.
type TransactionResult struct {
	TxHash  crypto.Hash
	Error   error
	GasUsed uint64
}

// A Computer uses the Cadence runtime to compute transaction results.
type Computer struct {
	runtime         runtime.Runtime
	contextProvider context.Provider
}

// New initializes a new computer with a runtime and context provider.
func New(runtime runtime.Runtime, contextProvider context.Provider) *Computer {
	return &Computer{
		runtime:         runtime,
		contextProvider: contextProvider,
	}
}

// ExecuteTransaction computes the result of a transaction.
func (c *Computer) ExecuteTransaction(tx *flow.Transaction) (*TransactionResult, error) {
	ctx := c.contextProvider.NewTransactionContext(tx)

	location := runtime.TransactionLocation(tx.Hash())

	err := c.runtime.ExecuteTransaction(tx.Script, ctx, location)
	if err != nil {
		if errors.As(err, &runtime.Error{}) {
			// runtime errors occur when the execution reverts
			return &TransactionResult{
				TxHash: tx.Hash(),
				Error:  err,
			}, nil
		}

		// other errors are unexpected and should be treated as fatal
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}

	return &TransactionResult{
		TxHash: tx.Hash(),
		Error:  nil,
	}, nil
}
