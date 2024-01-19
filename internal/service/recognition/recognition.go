package recognition

import (
	"context"

	"github.com/wonderf00l/fms-control-system/internal/entity"
	"github.com/wonderf00l/fms-control-system/internal/entity/recognition"
	"go.uber.org/zap"
)

var _ Service = (*service)(nil)

type Service interface {
	RecognizeWorkpiece(context.Context) (recognition.WorkpieceType, error)
	GetMetrics(context.Context) (*recognition.Metrics, error)
}

type service struct {
	log  *zap.SugaredLogger
	pool *entity.Pool
}

func NewService(log *zap.SugaredLogger, pool *entity.Pool) *service {
	return &service{log: log, pool: pool}
}

func (r *service) RecognizeWorkpiece(ctx context.Context) (recognition.WorkpieceType, error) {
	if err := r.pool.Recognition.IsReady(ctx); err != nil {
		return recognition.Unknown, err
	}

	location, err := r.pool.Conveyor.IsReady(ctx)
	if err != nil {
		return recognition.Unknown, err
	}

	if location.X != locationForRecognitionX || location.Y != locationForRecognitionY {
		return recognition.Unknown, &noWorkpieceToRecognizeError{}
	}

	return r.pool.Recognition.Recognize(ctx)
}

func (r *service) GetMetrics(ctx context.Context) (*recognition.Metrics, error) {
	return r.pool.Recognition.Metrics(ctx)
}
