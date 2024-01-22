package storage

import (
	"context"
)

type Metrics struct {
	Ready             bool
	Workpieces        int
	RawWorkpieces     int
	HandledWorkpieces int
}

type Storage interface {
	IsReady(context.Context) error
	ProvideRawWorkpiece(context.Context) error
	AcceptFinishedWorkpiece(context.Context) error
	Metrics(context.Context) (*Metrics, error)
}
