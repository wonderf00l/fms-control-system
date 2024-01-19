package miller

import (
	"context"
	"math/rand"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

var _ Miller = (*mockMiller)(nil)

type mockMiller struct{}

func NewMiller() *mockMiller {
	return &mockMiller{}
}

func (m *mockMiller) IsReady(ctx context.Context) error {
	switch rand.Intn(5000) {
	case 0:
		return &errors.TimeoutExceededError{Service: errors.Miller}
	default:
		return nil
	}
}

func (m *mockMiller) HandleWorkpiece(ctx context.Context) error {
	return nil
}

func (m *mockMiller) Metrics(ctx context.Context) (*Metrics, error) {
	return &Metrics{
		Ready: m.IsReady(ctx) == nil,
	}, nil
}
