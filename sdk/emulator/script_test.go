package emulator_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/sdk/abi/values"
	"github.com/dapperlabs/flow-go/sdk/emulator"
	"github.com/dapperlabs/flow-go/sdk/keys"
)

func TestExecuteScript(t *testing.T) {
	b := emulator.NewEmulatedBlockchain(emulator.DefaultOptions)

	accountAddress := b.RootAccountAddress()

	tx := flow.Transaction{
		Script:             []byte(addTwoScript),
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       accountAddress,
		ScriptAccounts:     []flow.Address{accountAddress},
	}

	sig, err := keys.SignTransaction(tx, b.RootKey())
	assert.Nil(t, err)

	tx.AddSignature(accountAddress, sig)

	callScript := fmt.Sprintf(sampleCall, accountAddress)

	// Sample call (value is 0)
	value, err := b.ExecuteScript([]byte(callScript))
	assert.Nil(t, err)
	assert.Equal(t, values.Int(0), value)

	// Submit tx1 (script adds 2)
	err = b.SubmitTransaction(tx)
	assert.Nil(t, err)

	// Sample call (value is 2)
	value, err = b.ExecuteScript([]byte(callScript))
	assert.Nil(t, err)
	assert.Equal(t, values.Int(2), value)
}
