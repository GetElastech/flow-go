package fvm

import (
	"github.com/onflow/cadence/runtime"
	"github.com/onflow/cadence/runtime/interpreter"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/fvm/programs"
	"github.com/onflow/flow-go/fvm/state"
	"github.com/onflow/flow-go/model/flow"
)

// An Procedure is an operation (or set of operations) that reads or writes ledger state.
type Procedure interface {
	Run(vm *VirtualMachine, ctx Context, sth *state.StateHolder, programs *programs.Programs) error
}

// A VirtualMachine augments the Cadence runtime with Flow host functionality.
type VirtualMachine struct {
	Runtime runtime.Runtime
}

// New creates a new virtual machine instance with the provided runtime.
func New(rt runtime.Runtime) *VirtualMachine {
	return &VirtualMachine{
		Runtime: rt,
	}
}

// Run runs a procedure against a ledger in the given context.
func (vm *VirtualMachine) Run(ctx Context, proc Procedure, v state.View, programs *programs.Programs) (err error) {

	st := state.NewState(v,
		state.WithMaxKeySizeAllowed(ctx.MaxStateKeySize),
		state.WithMaxValueSizeAllowed(ctx.MaxStateValueSize),
		state.WithMaxInteractionSizeAllowed(ctx.MaxStateInteractionSize))
	sth := state.NewStateHolder(st)

	defer func() {
		if r := recover(); r != nil {

			// Cadence may fail to encode certain values.
			// Return an error for now, which will cause transactions to revert.
			//
			if encodingErr, ok := r.(interpreter.EncodingUnsupportedValueError); ok {
				err = &EncodingUnsupportedValueError{
					Path:  encodingErr.Path,
					Value: encodingErr.Value,
				}
				return
			}

			panic(r)
		}
	}()

	err = proc.Run(vm, ctx, sth, programs)
	if err != nil {
		return err
	}

	return nil
}

// GetAccount returns an account by address or an error if none exists.
func (vm *VirtualMachine) GetAccount(ctx Context, address flow.Address, v state.View, programs *programs.Programs) (*flow.Account, error) {
	st := state.NewState(v,
		state.WithMaxKeySizeAllowed(ctx.MaxStateKeySize),
		state.WithMaxValueSizeAllowed(ctx.MaxStateValueSize),
		state.WithMaxInteractionSizeAllowed(ctx.MaxStateInteractionSize))

	sth := state.NewStateHolder(st)
	account, err := getAccount(vm, ctx, sth, programs, address)
	if err != nil {
		// TODO: wrap error
		return nil, err
	}
	return account, nil
}

// invokeMetaTransaction invokes a meta transaction inside the context of an outer transaction.
//
// Errors that occur in a meta transaction are propagated as a single error that can be
// captured by the Cadence runtime and eventually disambiguated by the parent context.
func (vm *VirtualMachine) invokeMetaTransaction(ctx Context, tx *TransactionProcedure, sth *state.StateHolder, programs *programs.Programs) error {
	invocator := NewTransactionInvocator(zerolog.Nop())
	err := invocator.Process(vm, ctx, tx, sth, programs)
	if err != nil {
		return err
	}

	if tx.Err != nil {
		return tx.Err
	}

	return nil
}
