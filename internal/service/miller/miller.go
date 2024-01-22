package miller

import (
	"context"

	"github.com/wonderf00l/fms-control-system/internal/entity"
	"github.com/wonderf00l/fms-control-system/internal/entity/miller"
	"go.uber.org/zap"
)

var _ Service = (*service)(nil)

type Service interface {
	HandleWorkpiece(context.Context) error
	GetMetrics(context.Context) (*miller.Metrics, error)
}

type service struct {
	log  *zap.SugaredLogger
	pool *entity.Pool
}

func NewService(log *zap.SugaredLogger, pool *entity.Pool) *service {
	return &service{log: log, pool: pool}
}

func (m *service) HandleWorkpiece(ctx context.Context) error {
	if err := m.pool.Miller.IsReady(ctx); err != nil {
		return err
	}

	location, err := m.pool.Conveyor.IsReady(ctx)
	if err != nil {
		return err
	}
	if location.X != locationForMillerX || location.Y != locationForMillerY {
		return &noWorkpieceForMillerToHandleError{}
	}

	return m.pool.Miller.HandleWorkpiece(ctx)
}

func (m *service) GetMetrics(ctx context.Context) (*miller.Metrics, error) {
	return m.pool.Miller.Metrics(ctx)
}
