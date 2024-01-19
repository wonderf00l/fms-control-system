package storage

import (
	"context"

	"github.com/wonderf00l/fms-control-system/internal/entity"
	storage "github.com/wonderf00l/fms-control-system/internal/entity/storage"
	"github.com/wonderf00l/fms-control-system/internal/errors"
	"go.uber.org/zap"
)

var _ Service = (*service)(nil)

type Service interface {
	ProvideWorkpiece(context.Context) error
	AcceptWorkpiece(context.Context) error
	GetMetrics(context.Context) (*storage.Metrics, error)
}

type service struct {
	log  *zap.SugaredLogger
	pool *entity.Pool
}

func NewService(log *zap.SugaredLogger, pool *entity.Pool) *service {
	return &service{log: log, pool: pool}
}

func (s *service) ProvideWorkpiece(ctx context.Context) error {
	if err := s.pool.Storage.IsReady(ctx); err != nil {
		return err
	}
	if err := s.pool.Conveyor.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Storage, Dependency: errors.Conveyor, Reason: err}
	}
	return s.pool.Storage.ProvideRawWorkpiece(ctx)
}

func (s *service) AcceptWorkpiece(ctx context.Context) error {
	if err := s.pool.Storage.IsReady(ctx); err != nil {
		return err
	}
	if err := s.pool.Conveyor.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Storage, Dependency: errors.Conveyor, Reason: err}
	}
	return s.pool.Storage.AcceptFinishedWorkpiece(ctx)
}

func (s *service) GetMetrics(ctx context.Context) (*storage.Metrics, error) {
	return s.pool.Storage.Metrics(ctx)
}
