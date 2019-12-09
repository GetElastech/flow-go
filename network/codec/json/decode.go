// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package json

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/dapperlabs/flow-go/model/coldstuff"
	"github.com/dapperlabs/flow-go/model/collection"
	"github.com/dapperlabs/flow-go/model/trickle"
)

// decode will decode the envelope into an entity.
func decode(env Envelope) (interface{}, error) {

	// create the desired message
	var v interface{}
	switch env.Code {

	// trickle overlay network
	case CodePing:
		v = &trickle.Ping{}
	case CodePong:
		v = &trickle.Pong{}
	case CodeAuth:
		v = &trickle.Auth{}
	case CodeAnnounce:
		v = &trickle.Announce{}
	case CodeRequest:
		v = &trickle.Request{}
	case CodeResponse:
		v = &trickle.Response{}

	case CodeGuaranteedCollection:
		v = &collection.GuaranteedCollection{}

	case CodeBlockProposal:
		v = &coldstuff.BlockProposal{}
	case CodeBlockVote:
		v = &coldstuff.BlockVote{}
	case CodeBlockCommit:
		v = &coldstuff.BlockCommit{}

	default:
		return nil, errors.Errorf("invalid message code (%d)", env.Code)
	}

	// unmarshal the payload
	err := json.Unmarshal(env.Data, v)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode payload")
	}

	return v, nil
}
