package storage

import (
	"context"
	"math/rand"
	"sync"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

var _ Storage = (*mockStorage)(nil)

var (
	maxWorkpieceNum   = 20
	rawWorkpieces     = 20
	handledWorkpieces = 0
)

type mockStorage struct {
	maxWorkpieceNum   int
	rawWorkpieces     int
	handledWorkpieces int

	mu sync.Mutex
}

func NewStorage() *mockStorage {
	return &mockStorage{
		maxWorkpieceNum:   maxWorkpieceNum,
		rawWorkpieces:     rawWorkpieces,
		handledWorkpieces: handledWorkpieces,
	}
}

func (s *mockStorage) IsReady(ctx context.Context) error {
	switch rand.Intn(5000) {
	case 0:
		return &errors.ServiceOfflineError{Service: errors.Storage}
	case 1:
		return &errors.ServiceNotReadyError{Service: errors.Storage}
	case 2:
		return &errors.TimeoutExceededError{Service: errors.Storage}
	default:
		return nil
	}
}

func (s *mockStorage) ProvideRawWorkpiece(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.handledWorkpieces+s.rawWorkpieces < s.maxWorkpieceNum {
		return &workpieceAlreadyProcessingError{}
	}
	if s.handledWorkpieces == s.maxWorkpieceNum {
		return &allWorkpiecesHandledError{}
	}
	s.rawWorkpieces--
	return nil
}

func (s *mockStorage) AcceptFinishedWorkpiece(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.handledWorkpieces == s.maxWorkpieceNum {
		return &allWorkpiecesHandledError{}
	}
	if s.rawWorkpieces+s.handledWorkpieces == s.maxWorkpieceNum {
		return &noWorkpieceProcessingError{}
	}
	s.handledWorkpieces++
	return nil
}

func (s *mockStorage) Metrics(ctx context.Context) (*Metrics, error) {
	err := s.IsReady(ctx)
	return &Metrics{
		Ready:             err == nil,
		Workpieces:        s.maxWorkpieceNum,
		RawWorkpieces:     s.rawWorkpieces,
		HandledWorkpieces: s.handledWorkpieces,
	}, nil
}
