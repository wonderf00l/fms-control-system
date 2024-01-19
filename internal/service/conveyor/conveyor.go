package conveyor

import (
	"context"

	"github.com/wonderf00l/fms-control-system/internal/entity"
	"github.com/wonderf00l/fms-control-system/internal/entity/conveyor"
	"github.com/wonderf00l/fms-control-system/internal/errors"
	"go.uber.org/zap"
)

var _ Service = (*service)(nil)

type Service interface {
	MoveWorkpieceToRecognition(context.Context) error
	MoveWorkpieceToLathe(context.Context) error
	MoveWorkpieceToMiller(context.Context) error
	MoveWorkpieceToStorage(context.Context) error
	GetMetrics(context.Context) (*conveyor.Metrics, error)
}

type service struct {
	log  *zap.SugaredLogger
	pool *entity.Pool
}

func NewService(log *zap.SugaredLogger, pool *entity.Pool) *service {
	return &service{log: log, pool: pool}
}

func (c *service) MoveWorkpieceToRecognition(ctx context.Context) error {
	if _, err := c.pool.Conveyor.IsReady(ctx); err != nil {
		return err
	}
	if err := c.pool.Recognition.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Conveyor, Dependency: errors.Recognition, Reason: err}
	}
	return c.pool.Conveyor.SetWorkpieceLocation(ctx, conveyor.WorkpieceLocation{
		X: locationForRecognitionX,
		Y: locationForRecognitionY,
	}, false)
}

func (c *service) MoveWorkpieceToLathe(ctx context.Context) error {
	if _, err := c.pool.Conveyor.IsReady(ctx); err != nil {
		return err
	}
	if err := c.pool.Lathe.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Conveyor, Dependency: errors.Lathe, Reason: err}
	}
	return c.pool.Conveyor.SetWorkpieceLocation(ctx, conveyor.WorkpieceLocation{
		X: locationForLatheX,
		Y: locationForLatheY,
	}, false)
}

func (c *service) MoveWorkpieceToMiller(ctx context.Context) error {
	if _, err := c.pool.Conveyor.IsReady(ctx); err != nil {
		return err
	}
	if err := c.pool.Miller.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Conveyor, Dependency: errors.Miller, Reason: err}
	}
	return c.pool.Conveyor.SetWorkpieceLocation(ctx, conveyor.WorkpieceLocation{
		X: locationForMillerX,
		Y: locationForMillerY,
	}, false)
}

func (c *service) MoveWorkpieceToStorage(ctx context.Context) error {
	if _, err := c.pool.Conveyor.IsReady(ctx); err != nil {
		return err
	}
	if err := c.pool.Storage.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Conveyor, Dependency: errors.Storage, Reason: err}
	}
	return c.pool.Conveyor.SetWorkpieceLocation(ctx, conveyor.WorkpieceLocation{
		X: locationForStorageX,
		Y: locationForStorageY,
	}, false)
}

func (c *service) GetMetrics(ctx context.Context) (*conveyor.Metrics, error) {
	return c.pool.Conveyor.Metrics(ctx)
}
