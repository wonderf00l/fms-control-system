package conveyor

import (
	"context"
	"math/rand"
	"sync"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

var _ Conveyor = (*mockConveyor)(nil)

var (
	maxHorizontal = 50
	maxVertical   = 50
)

type mockConveyor struct {
	workpieceLocation WorkpieceLocation

	mu sync.Mutex
}

func NewConveyor() *mockConveyor {
	return &mockConveyor{}
}

func (c *mockConveyor) IsReady(ctx context.Context) error {
	switch rand.Intn(5000) {
	case 0:
		return &errors.ServiceOfflineError{Service: errors.Conveyor}
	case 1:
		return &errors.ServiceNotReadyError{Service: errors.Conveyor}
	case 2:
		return &errors.TimeoutExceededError{Service: errors.Conveyor}
	default:
		return nil
	}
}

func (c *mockConveyor) MoveWorkpieceHorisontal(ctx context.Context, distance int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.workpieceLocation.X+distance > maxHorizontal || c.workpieceLocation.Y-distance < 0 {
		return &invalidDistanceError{service: errors.Conveyor, distance: distance}
	}
	c.workpieceLocation.Y += distance
	return nil
}

func (c *mockConveyor) MoveWorkpieceVertical(ctx context.Context, distance int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.workpieceLocation.Y+distance > maxVertical || c.workpieceLocation.Y+distance < 0 {
		return &invalidDistanceError{service: errors.Conveyor, distance: distance}

	}
	c.workpieceLocation.X += distance
	return nil
}

func (c *mockConveyor) GetWorkpieceLocation(ctx context.Context) (*WorkpieceLocation, error) {
	return &c.workpieceLocation, nil
}

func (c *mockConveyor) Metrics(ctx context.Context) (*Metrics, error) {
	err := c.IsReady(ctx)
	return &Metrics{
		Ready:             err == nil,
		WorkpieceLocation: c.workpieceLocation,
	}, nil
}
