package example

import (
	"github.com/dapperlabs/flow-go/sdk/abi/encode"
	"github.com/dapperlabs/flow-go/sdk/abi/types"
	"github.com/dapperlabs/flow-go/sdk/abi/values"
)

type Person interface {
	FullName() string
	Value() values.Value
}

var PersonType types.Type = types.Composite{
	FieldTypes: []types.Type{types.String{}},
}

func EncodePerson(p Person) ([]byte, error) {
	return encode.Encode(p.Value())
}

func DecodePerson(b []byte) (Person, error) {
	v, err := encode.Decode(PersonType, b)
	if err != nil {
		return nil, err
	}

	return newPersonFromValue(v), nil
}

type person struct {
	value values.Composite
}

func newPersonFromValue(v values.Value) Person {
	value := v.(values.Composite)
	return person{value}
}

func (p person) FullName() string {
	return string(p.value.Fields[0].(values.String))
}

func (p person) Value() values.Value {
	return p.value
}
