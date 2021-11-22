package rest

import (
	"context"
	"fmt"
	"github.com/onflow/flow-go/model/flow"
	"github.com/rs/zerolog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onflow/flow-go/access"
	"github.com/onflow/flow-go/engine/access/rest/generated"
)

const ExpandableFieldPayload = "payload"
const ExpandableExecutionResult = "execution_result"

// getBlocksByID gets blocks by provided ID or collection of IDs.
func getBlocksByID(
	r *requestDecorator,
	backend access.API,
	linkGenerator LinkGenerator,
	logger zerolog.Logger,
) (interface{}, StatusError) {

	ids, err := r.ids()
	if err != nil {
		return nil, NewBadRequestError(err.Error(), err)
	}

	blocks := make([]*generated.Block, len(ids))
	for i, id := range ids {
		block, err := getBlockByID(r.Context(), id, r, backend, linkGenerator)
		if err != nil {
			return nil, err
		}
		blocks[i] = block
	}

	return blocks, nil
}

func getBlockByID(ctx context.Context, id flow.Identifier, req *requestDecorator, backend access.API, linkGenerator LinkGenerator) (*generated.Block, StatusError) {
	var responseBlock = new(generated.Block)
	if req.expands(ExpandableFieldPayload) {
		flowBlock, err := backend.GetBlockByID(ctx, id)
		if err != nil {
			return nil, blockLookupError(id, err)
		}
		responseBlock.Payload = blockPayloadResponse(flowBlock.Payload)
		responseBlock.Header = blockHeaderResponse(flowBlock.Header)
	} else {
		flowBlockHeader, err := backend.GetBlockHeaderByID(ctx, id)
		if err != nil {
			return nil, blockLookupError(id, err)
		}
		responseBlock.Payload = nil
		responseBlock.Header = blockHeaderResponse(flowBlockHeader)
	}
	//if req.expands(ExpandableExecutionResult) {
	//	// lookup ER here and add to response
	//}

	blockLink, err := linkGenerator.BlockLink(id)
	if err != nil {
		msg := fmt.Sprintf("failed to generate respose for block ID %s", id.String())
		return nil, NewRestError(http.StatusInternalServerError, msg, err)
	}

	responseBlock.Links = new(generated.Links)
	responseBlock.Links.Self = blockLink

	return responseBlock, nil
}

func blockLookupError(id flow.Identifier, err error) StatusError {
	msg := fmt.Sprintf("block with ID %s not found", id.String())
	// if error has GRPC code NotFound, then return HTTP NotFound error
	if status.Code(err) == codes.NotFound {
		return NewNotFoundError(msg, err)
	}

	return NewRestError(http.StatusInternalServerError, msg, err)
}
