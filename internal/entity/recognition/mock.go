package recognition

import (
	"context"
	"math/rand"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

var _ Recognition = (*mockRecognition)(nil)

type mockRecognition struct {
	workpieceType WorkpieceType
}

func NewRecognition() *mockRecognition {
	return &mockRecognition{}
}

func (r *mockRecognition) IsReady(context.Context) error {
	switch rand.Intn(5000) {
	case 0:
		return &errors.ServiceOfflineError{Service: errors.Recognition}
	case 1:
		return &errors.ServiceNotReadyError{Service: errors.Recognition}
	case 2:
		return &errors.TimeoutExceededError{Service: errors.Recognition}
	default:
		return nil
	}
}

func (r *mockRecognition) Recognize(ctx context.Context) (WorkpieceType, error) {
	return WorkpieceType(2), nil
}

func (r *mockRecognition) Metrics(ctx context.Context) (*Metrics, error) {
	err := r.IsReady(ctx)
	return &Metrics{
		Ready:         err == nil,
		WorkpieceType: r.workpieceType,
	}, nil
}
