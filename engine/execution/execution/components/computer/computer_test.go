package computer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/flow-go/engine/execution/execution/components/computer"
	context "github.com/dapperlabs/flow-go/engine/execution/execution/modules/context/mock"
	"github.com/dapperlabs/flow-go/engine/execution/execution/modules/ledger"
	"github.com/dapperlabs/flow-go/language/runtime"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/utils/unittest"
)

func TestExecuteTransaction(t *testing.T) {
	rt := runtime.NewInterpreterRuntime()
	prov := &context.Provider{}
	ctx := &context.TransactionContext{}

	c := computer.New(rt, prov)

	t.Run("transaction success", func(t *testing.T) {
		tx := flow.TransactionBody{
			Script: []byte(`
                transaction {
                  prepare(signer: Account) {}
                }
            `),
		}

		view := ledger.NewView(func(key string) ([]byte, error) { return nil, nil })

		prov.On("NewTransactionContext", tx, view).
			Return(ctx).
			Once()

		ctx.On("GetSigningAccounts").
			Return([]runtime.Address{runtime.Address(unittest.AddressFixture())}).
			Once()

		result, err := c.ExecuteTransaction(view, tx)
		assert.NoError(t, err)
		assert.NoError(t, result.Error)

		prov.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("transaction failure", func(t *testing.T) {
		tx := flow.TransactionBody{
			Script: []byte(`
                transaction {
                  var x: Int

                  prepare(signer: Account) {
                    self.x = 0
                  }

                  execute {
                    self.x = 1
                  }

                  post {
                    self.x == 2
                  }
                }
            `),
		}

		view := ledger.NewView(func(key string) ([]byte, error) { return nil, nil })

		prov.On("NewTransactionContext", tx, view).
			Return(ctx).
			Once()

		ctx.On("GetSigningAccounts").
			Return([]runtime.Address{runtime.Address(unittest.AddressFixture())}).
			Once()

		result, err := c.ExecuteTransaction(view, tx)
		assert.NoError(t, err)
		assert.Error(t, result.Error)

		prov.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})
}
