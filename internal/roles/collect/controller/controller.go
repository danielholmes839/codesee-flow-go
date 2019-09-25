package controller

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dapperlabs/flow-go/internal/roles/collect/storage"
	"github.com/dapperlabs/flow-go/internal/roles/collect/txpool"
	svc "github.com/dapperlabs/flow-go/pkg/grpc/services/collect"
)

type Controller struct {
	storage storage.Storage
	txPool  *txpool.TxPool
	log     *logrus.Entry
}

const msgInternalError = "internal error"

func New(storage storage.Storage, txPool *txpool.TxPool, log *logrus.Logger) *Controller {
	return &Controller{
		storage: storage,
		txPool:  txPool,
		log:     logrus.NewEntry(log),
	}
}

func (c *Controller) Ping(context.Context, *svc.PingRequest) (*svc.PingResponse, error) {
	return &svc.PingResponse{
		Address: []byte("pong!"),
	}, nil
}

func (c *Controller) SubmitCollection(context.Context, *svc.SubmitCollectionRequest) (*svc.SubmitCollectionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (c *Controller) GetTransaction(context.Context, *svc.GetTransactionRequest) (*svc.GetTransactionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (c *Controller) GetCollection(context.Context, *svc.GetCollectionRequest) (*svc.GetCollectionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
