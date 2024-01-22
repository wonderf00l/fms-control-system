package conveyor

import (
	"context"
	"math/rand"
	"sync"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

var _ Conveyor = (*mockConveyor)(nil)

var (
	maxX = 50
	maxY = 50
)

type mockConveyor struct {
	workpieceLocation WorkpieceLocation

	mu sync.Mutex
}

func NewConveyor() *mockConveyor {
	return &mockConveyor{}
}

func (c *mockConveyor) IsReady(ctx context.Context) (*WorkpieceLocation, error) {
	switch rand.Intn(5000) {
	case 0:
		return nil, &errors.TimeoutExceededError{Service: errors.Conveyor}
	default:
		return &c.workpieceLocation, nil
	}
}

func (c *mockConveyor) SetWorkpieceLocation(ctx context.Context, location WorkpieceLocation, fromStorage bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !fromStorage && c.workpieceLocation.X == 0 && c.workpieceLocation.Y == 0 {
		return &workpieceNotOnConveyorError{}
	}
	if location.X > maxX || location.X < 0 || location.Y > maxY || location.Y < 0 {
		return &invalidLocationError{location: location}
	}
	c.workpieceLocation.X = location.X
	c.workpieceLocation.Y = location.Y
	return nil
}

func (c *mockConveyor) Metrics(ctx context.Context) (*Metrics, error) {
	location, err := c.IsReady(ctx)
	return &Metrics{
		Ready:             err == nil,
		WorkpieceLocation: location,
	}, nil
}
