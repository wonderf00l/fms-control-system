package lathe

import (
	"context"

	"github.com/wonderf00l/fms-control-system/internal/entity"
	"github.com/wonderf00l/fms-control-system/internal/entity/lathe"
	"go.uber.org/zap"
)

var _ Service = (*service)(nil)

type Service interface {
	HandleWorkpiece(context.Context) error
	GetMetrics(context.Context) (*lathe.Metrics, error)
}

type service struct {
	log  *zap.SugaredLogger
	pool *entity.Pool
}

func NewService(log *zap.SugaredLogger, pool *entity.Pool) *service {
	return &service{log: log, pool: pool}
}

func (m *service) HandleWorkpiece(ctx context.Context) error {
	if err := m.pool.Lathe.IsReady(ctx); err != nil {
		return err
	}

	location, err := m.pool.Conveyor.GetWorkpieceLocation(ctx)
	if err != nil {
		return err
	}
	if location.X != locationForLatheX || location.Y != locationForLatheY {
		return &noWorkpieceForLatheToHandleError{}
	}

	return m.pool.Lathe.HandleWorkpiece(ctx)
}

func (m *service) GetMetrics(ctx context.Context) (*lathe.Metrics, error) {
	return m.pool.Lathe.Metrics(ctx)
}
