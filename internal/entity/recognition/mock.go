package recognition

import (
	"context"
	"math/rand"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

var _ Recognition = (*mockRecognition)(nil)

type mockRecognition struct{}

func NewRecognition() *mockRecognition {
	return &mockRecognition{}
}

func (r *mockRecognition) IsReady(context.Context) error {
	switch rand.Intn(5000) {
	case 0:
		return &errors.TimeoutExceededError{Service: errors.Recognition}
	default:
		return nil
	}
}

func (r *mockRecognition) Recognize(ctx context.Context) (WorkpieceType, error) {
	return WorkpieceType(2), nil
}

func (r *mockRecognition) Metrics(ctx context.Context) (*Metrics, error) {
	return &Metrics{
		Ready: r.IsReady(ctx) == nil,
	}, nil
}
