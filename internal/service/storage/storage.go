package storage

import (
	"context"

	"github.com/wonderf00l/fms-control-system/internal/entity"
	"github.com/wonderf00l/fms-control-system/internal/entity/conveyor"
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
	if _, err := s.pool.Conveyor.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Storage, Dependency: errors.Conveyor, Reason: err}
	}
	if err := s.pool.Storage.ProvideRawWorkpiece(ctx); err != nil {
		return err
	}
	if err := s.pool.Conveyor.SetWorkpieceLocation(ctx, conveyor.WorkpieceLocation{X: 0, Y: 1}, true); err != nil {
		return err
	}
	return nil
}

func (s *service) AcceptWorkpiece(ctx context.Context) error {
	if err := s.pool.Storage.IsReady(ctx); err != nil {
		return err
	}
	if _, err := s.pool.Conveyor.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Storage, Dependency: errors.Conveyor, Reason: err}
	}
	if err := s.pool.Storage.AcceptFinishedWorkpiece(ctx); err != nil {
		return err
	}
	if err := s.pool.Conveyor.SetWorkpieceLocation(ctx, conveyor.WorkpieceLocation{X: 0, Y: 0}, true); err != nil {
		return err
	}
	return nil
}

func (s *service) GetMetrics(ctx context.Context) (*storage.Metrics, error) {
	return s.pool.Storage.Metrics(ctx)
}
