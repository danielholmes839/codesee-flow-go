// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package json

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/dapperlabs/flow-go/model/coldstuff"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/model/libp2p/message"
	"github.com/dapperlabs/flow-go/model/messages"
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

	// ColdStuff
	case CodeColdStuffBlockProposal:
		v = &coldstuff.BlockProposal{}
	case CodeColdStuffBlockCommit:
		v = &coldstuff.BlockCommit{}
	case CodeColdStuffBlockVote:
		v = &coldstuff.BlockVote{}

	case CodeCollectionGuarantee:
		v = &flow.CollectionGuarantee{}
	case CodeTransaction:
		v = &flow.Transaction{}

	case CodeBlock:
		v = &flow.Block{}

	case CodeCollectionRequest:
		v = &messages.CollectionRequest{}
	case CodeCollectionResponse:
		v = &messages.CollectionResponse{}

	case CodeEcho:
		v = &message.Echo{}

	case CodeExecutionRecipt:
		v = &flow.ExecutionReceipt{}

	case CodeExecutionStateRequest:
		v = &messages.ExecutionStateRequest{}

	case CodeExecutionStateResponse:
		v = &messages.ExecutionStateResponse{}

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
