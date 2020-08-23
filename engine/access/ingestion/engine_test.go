package ingestion

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/dapperlabs/flow-go/consensus/hotstuff/model"
	"github.com/dapperlabs/flow-go/engine"
	"github.com/dapperlabs/flow-go/engine/access/rpc"
	"github.com/dapperlabs/flow-go/model/flow"

	"github.com/dapperlabs/flow-go/module/mempool/stdmap"
	"github.com/dapperlabs/flow-go/module/metrics"
	module "github.com/dapperlabs/flow-go/module/mock"
	network "github.com/dapperlabs/flow-go/network/mock"
	protocol "github.com/dapperlabs/flow-go/state/protocol/mock"
	storerr "github.com/dapperlabs/flow-go/storage"
	storage "github.com/dapperlabs/flow-go/storage/mock"
	"github.com/dapperlabs/flow-go/utils/unittest"
)

type Suite struct {
	suite.Suite

	// protocol state
	proto struct {
		state    *protocol.State
		snapshot *protocol.Snapshot
		mutator  *protocol.Mutator
	}

	me            *module.Local
	request       *module.Requester
	provider      *network.Engine
	blocks        *storage.Blocks
	headers       *storage.Headers
	collections   *storage.Collections
	transactions  *storage.Transactions
	rootBlkHeight uint64

	eng *Engine
}

func TestIngestEngine(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) SetupTest() {
	log := zerolog.New(os.Stderr)

	obsIdentity := unittest.IdentityFixture(unittest.WithRole(flow.RoleAccess))

	// mock out protocol state
	suite.proto.state = new(protocol.State)
	suite.proto.snapshot = new(protocol.Snapshot)
	suite.proto.state.On("Identity").Return(obsIdentity, nil)
	suite.proto.state.On("Final").Return(suite.proto.snapshot, nil)

	suite.me = new(module.Local)
	suite.me.On("NodeID").Return(obsIdentity.NodeID)

	net := new(module.Network)
	conduit := new(network.Conduit)
	net.On("Register", uint8(engine.ReceiveReceipts), mock.Anything).
		Return(conduit, nil).
		Once()
	suite.request = new(module.Requester)

	suite.provider = new(network.Engine)
	suite.blocks = new(storage.Blocks)
	suite.headers = new(storage.Headers)
	suite.collections = new(storage.Collections)
	suite.transactions = new(storage.Transactions)
	collectionsToMarkFinalized, err := stdmap.NewTimes(100)
	require.NoError(suite.T(), err)
	collectionsToMarkExecuted, err := stdmap.NewTimes(100)
	require.NoError(suite.T(), err)
	blocksToMarkExecuted, err := stdmap.NewTimes(100)
	require.NoError(suite.T(), err)
	suite.rootBlkHeight = 1000

	rpcEng := rpc.New(log, suite.proto.state, rpc.Config{}, nil, nil, suite.blocks, suite.headers, suite.collections,
		suite.transactions, flow.Testnet, metrics.NewNoopCollector(), 0)

	eng, err := New(log, net, suite.proto.state, suite.me, suite.request, suite.blocks, suite.headers, suite.collections,
		suite.transactions, metrics.NewNoopCollector(), collectionsToMarkFinalized, collectionsToMarkExecuted,
		blocksToMarkExecuted, rpcEng, suite.rootBlkHeight)
	require.NoError(suite.T(), err)

	suite.eng = eng

}

// TestOnFinalizedBlock checks that when a block is received, a request for each individual collection is made
func (suite *Suite) TestOnFinalizedBlock() {

	block := unittest.BlockFixture()
	modelBlock := model.Block{
		BlockID: block.ID(),
	}

	// we should query the block once and index the guarantee payload once
	suite.blocks.On("ByID", block.ID()).Return(&block, nil).Twice()
	for _, g := range block.Payload.Guarantees {
		collection := unittest.CollectionFixture(1)
		light := collection.Light()
		suite.collections.On("LightByID", g.CollectionID).Return(&light, nil).Twice()
	}

	// expect that the block storage is indexed with each of the collection guarantee
	suite.blocks.On("IndexBlockForCollections", block.ID(), flow.GetIDs(block.Payload.Guarantees)).Return(nil).Once()

	// for each of the guarantees, we should request the corresponding collection once
	needed := make(map[flow.Identifier]struct{})
	for _, guarantee := range block.Payload.Guarantees {
		needed[guarantee.ID()] = struct{}{}
	}
	suite.request.On("EntityByID", mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			collID := args.Get(0).(flow.Identifier)
			_, pending := needed[collID]
			suite.Assert().True(pending, "collection should be pending (%x)", collID)
			delete(needed, collID)
		},
	)

	// process the block through the finalized callback
	suite.eng.OnFinalizedBlock(&modelBlock)

	// wait for engine shutdown
	done := suite.eng.unit.Done()
	assert.Eventually(suite.T(), func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}, time.Second, 20*time.Millisecond)

	// assert that the block was retrieved and all collections were requested
	suite.headers.AssertExpectations(suite.T())
	suite.request.AssertNumberOfCalls(suite.T(), "EntityByID", len(block.Payload.Guarantees))
}

// TestOnCollection checks that when a Collection is received, it is persisted
func (suite *Suite) TestOnCollection() {

	originID := unittest.IdentifierFixture()
	collection := unittest.CollectionFixture(5)
	light := collection.Light()

	// we should store the light collection and index its transactions
	suite.collections.On("StoreLightAndIndexByTransaction", &light).Return(nil).Once()

	// for each transaction in the collection, we should store it
	needed := make(map[flow.Identifier]struct{})
	for _, txID := range light.Transactions {
		needed[txID] = struct{}{}
	}
	suite.transactions.On("Store", mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			tx := args.Get(0).(*flow.TransactionBody)
			_, pending := needed[tx.ID()]
			suite.Assert().True(pending, "tx not pending (%x)", tx.ID())
		},
	)

	// process the block through the collection callback
	suite.eng.OnCollection(originID, &collection)

	// wait for engine to be done processing
	done := suite.eng.unit.Done()
	assert.Eventually(suite.T(), func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}, time.Second, 20*time.Millisecond)

	// check that the collection was stored and indexed, and we stored all transactions
	suite.collections.AssertExpectations(suite.T())
	suite.transactions.AssertNumberOfCalls(suite.T(), "Store", len(collection.Transactions))
}

// TestOnCollection checks that when a duplicate collection is received, the node doesn't
// crash but just ignores its transactions.
func (suite *Suite) TestOnCollectionDuplicate() {

	originID := unittest.IdentifierFixture()
	collection := unittest.CollectionFixture(5)
	light := collection.Light()

	// we should store the light collection and index its transactions
	suite.collections.On("StoreLightAndIndexByTransaction", &light).Return(storerr.ErrAlreadyExists).Once()

	// for each transaction in the collection, we should store it
	needed := make(map[flow.Identifier]struct{})
	for _, txID := range light.Transactions {
		needed[txID] = struct{}{}
	}
	suite.transactions.On("Store", mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			tx := args.Get(0).(*flow.TransactionBody)
			_, pending := needed[tx.ID()]
			suite.Assert().True(pending, "tx not pending (%x)", tx.ID())
		},
	)

	// process the block through the collection callback
	suite.eng.OnCollection(originID, &collection)

	// wait for engine to be done processing
	done := suite.eng.unit.Done()
	assert.Eventually(suite.T(), func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}, time.Second, 20*time.Millisecond)

	// check that the collection was stored and indexed, and we stored all transactions
	suite.collections.AssertExpectations(suite.T())
	suite.transactions.AssertNotCalled(suite.T(), "Store", "should not store any transactions")
}

func (suite *Suite) TestRequestMissingCollections() {

	blkCnt := 3
	startHeight := suite.rootBlkHeight
	blocks := make([]flow.Block, blkCnt)
	heightMap := make(map[uint64]*flow.Block, blkCnt)

	var collIDs []flow.Identifier
	gap := 2
	for i := 0; i < blkCnt; i++ {
		block := unittest.BlockFixture()
		height := startHeight + uint64(i) + uint64(gap)
		block.Header.Height = height
		blocks[i] = block
		heightMap[height] = &block
		for _, c := range block.Payload.Guarantees {
			collIDs = append(collIDs, c.CollectionID)
		}
	}

	// each block should be queried by height
	suite.blocks.On("ByHeight", mock.IsType(uint64(0))).Return(
		func(h uint64) *flow.Block {
			// simulate a db lookup
			return heightMap[h]
		},
		func(h uint64) error {
			if _, ok := heightMap[h]; ok {
				return nil
			} else {
				return storerr.ErrNotFound
			}
		})

	// for the first lookup call for each collection, it will be reported as missing from db
	// for the second call, it will be reported as present
	rcvdColl := make(map[flow.Identifier]struct{})
	suite.collections.On("LightByID", mock.Anything).Return(
		func(cID flow.Identifier) *flow.LightCollection {
			return nil // the actual collection object return is never really read
		},
		func(cID flow.Identifier) error {
			if _, ok := rcvdColl[cID]; ok {
				return nil
			} else {
				 if rand.Float32() >= 0.8 {
					 rcvdColl[cID] = struct{}{}
				 }
				return storerr.ErrNotFound
			}
		}).
		// simulate some db i/o contention
		After(time.Millisecond * time.Duration(rand.Intn(5)))

	// consider the last test block as the head
	suite.proto.snapshot.On("Head").Return(blocks[blkCnt-1].Header, nil).Once()

	// setup the requester engine mock
	//// entityByID should be called once per collection
	for _, c := range collIDs {
		suite.request.On("EntityByID", c, mock.Anything).Return().Once()
	}
	//// force should be called once
	suite.request.On("Force").Return().Once()

	// finally, call the function under test
	err := suite.eng.requestMissingCollections(context.Background())
	require.NoError(suite.T(), err)
	require.Len(suite.T(), rcvdColl, len(collIDs))

	suite.request.AssertExpectations(suite.T())
	suite.collections.AssertExpectations(suite.T())
	suite.proto.snapshot.AssertExpectations(suite.T())
	suite.blocks.AssertExpectations(suite.T())
}
