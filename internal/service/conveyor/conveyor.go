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
	GetWorkpieceLocation(context.Context) (any, error)
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
	if err := c.pool.Conveyor.IsReady(ctx); err != nil {
		return err
	}
	if err := c.pool.Recognition.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Conveyor, Dependency: errors.Recognition, Reason: err}
	}
	return c.pool.Conveyor.MoveWorkpieceVertical(ctx, toRecognitionVerticalDistance)
}

func (c *service) MoveWorkpieceToLathe(ctx context.Context) error {
	if err := c.pool.Conveyor.IsReady(ctx); err != nil {
		return err
	}
	if err := c.pool.Lathe.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Conveyor, Dependency: errors.Lathe, Reason: err}
	}
	if err := c.pool.Conveyor.MoveWorkpieceVertical(ctx, toLatheVerticalDistance); err != nil {
		return err
	}
	if err := c.pool.Conveyor.MoveWorkpieceHorisontal(ctx, toLatheHorizontalDistance); err != nil {
		return err
	}
	return nil
}

func (c *service) MoveWorkpieceToMiller(ctx context.Context) error {
	if err := c.pool.Conveyor.IsReady(ctx); err != nil {
		return err
	}
	if err := c.pool.Miller.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Conveyor, Dependency: errors.Miller, Reason: err}
	}
	if err := c.pool.Conveyor.MoveWorkpieceHorisontal(ctx, toMillerHorizontalDistance); err != nil {
		return err
	}
	return nil
}

func (c *service) MoveWorkpieceToStorage(ctx context.Context) error {
	if err := c.pool.Conveyor.IsReady(ctx); err != nil {
		return err
	}
	if err := c.pool.Storage.IsReady(ctx); err != nil {
		return &errors.DepsNotReadyError{Service: errors.Conveyor, Dependency: errors.Storage, Reason: err}
	}
	if err := c.pool.Conveyor.MoveWorkpieceVertical(ctx, toStorageVerticalDistance); err != nil {
		return err
	}
	if err := c.pool.Conveyor.MoveWorkpieceHorisontal(ctx, toStorageHotizontalDistance); err != nil {
		return err
	}
	return nil
}

func (c *service) GetWorkpieceLocation(ctx context.Context) (any, error) {
	return c.pool.Conveyor.GetWorkpieceLocation(ctx)
}

func (c *service) GetMetrics(ctx context.Context) (*conveyor.Metrics, error) {
	return c.pool.Conveyor.Metrics(ctx)
}
