package lathe

import (
	"context"
	"math/rand"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

var _ Lathe = (*mockLathe)(nil)

type mockLathe struct{}

func NewLathe() *mockLathe {
	return &mockLathe{}
}

func (l *mockLathe) IsReady(ctx context.Context) error {
	switch rand.Intn(5000) {
	case 0:
		return &errors.ServiceOfflineError{Service: errors.Lathe}
	case 1:
		return &errors.ServiceNotReadyError{Service: errors.Lathe}
	case 2:
		return &errors.TimeoutExceededError{Service: errors.Lathe}
	default:
		return nil
	}
}

func (l *mockLathe) HandleWorkpiece(context.Context) error {
	return nil
}

func (l *mockLathe) Metrics(ctx context.Context) (*Metrics, error) {
	err := l.IsReady(ctx)
	return &Metrics{
		Ready: err == nil,
	}, nil
}
