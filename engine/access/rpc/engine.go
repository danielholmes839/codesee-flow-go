package rpc

import (
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/execution"

	"github.com/dapperlabs/flow-go/engine"
	"github.com/dapperlabs/flow-go/engine/access/rpc/handler"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/state/protocol"
	"github.com/dapperlabs/flow-go/storage"
	grpcutils "github.com/dapperlabs/flow-go/utils/grpc"
)

// Config defines the configurable options for the gRPC grpcServer.
type Config struct {
	GRPCListenAddr string
	HTTPListenAddr string
	ExecutionAddr  string
	CollectionAddr string
	MaxMsgSize     int // In bytes
}

// Engine implements a gRPC grpcServer with a simplified version of the Observation API.
type Engine struct {
	unit       *engine.Unit
	log        zerolog.Logger
	handler    *handler.Handler // the gRPC service implementation
	grpcServer *grpc.Server     // the gRPC grpcServer
	httpServer *HTTPServer
	config     Config
}

// New returns a new RPC engine.
func New(log zerolog.Logger,
	state protocol.State,
	config Config,
	executionRPC execution.ExecutionAPIClient,
	collectionRPC access.AccessAPIClient,
	blocks storage.Blocks,
	headers storage.Headers,
	collections storage.Collections,
	transactions storage.Transactions,
	chainID flow.ChainID) *Engine {

	log = log.With().Str("engine", "rpc").Logger()

	if config.MaxMsgSize == 0 {
		config.MaxMsgSize = grpcutils.DefaultMaxMsgSize
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(config.MaxMsgSize),
		grpc.MaxSendMsgSize(config.MaxMsgSize),
	)

	httpServer := NewHTTPServer(grpcServer, 8080)

	eng := &Engine{
		log:        log,
		unit:       engine.NewUnit(),
		handler:    handler.NewHandler(log, state, executionRPC, collectionRPC, blocks, headers, collections, transactions, chainID),
		grpcServer: grpcServer,
		httpServer: httpServer,
		config:     config,
	}

	access.RegisterAccessAPIServer(eng.grpcServer, eng.handler)

	return eng
}

// Ready returns a ready channel that is closed once the engine has fully
// started. The RPC engine is ready when the gRPC grpcServer has successfully
// started.
func (e *Engine) Ready() <-chan struct{} {
	e.unit.Launch(e.serve)
	return e.unit.Ready()
}

// Done returns a done channel that is closed once the engine has fully stopped.
// It sends a signal to stop the gRPC grpcServer, then closes the channel.
func (e *Engine) Done() <-chan struct{} {
	return e.unit.Done(e.grpcServer.GracefulStop)
}

// serve starts the gRPC grpcServer and the http proxy server
// When this function returns, the grpcServer is considered ready.
func (e *Engine) serve() {
	e.log.Info().Msgf("starting grpc server on address %s", e.config.GRPCListenAddr)

	l, err := net.Listen("tcp", e.config.GRPCListenAddr)
	if err != nil {
		e.log.Err(err).Msg("failed to start grpcServer")
		return
	}

	err = e.grpcServer.Serve(l)
	if err != nil {
		e.log.Err(err).Msg("fatal error in grpcServer")
	}

	e.log.Info().Msgf("starting http server on address %s", e.config.HTTPListenAddr)
	e.httpServer.Start()
}
