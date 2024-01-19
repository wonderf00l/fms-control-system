package lathe

import "context"

type Metrics struct {
	Ready bool
}

type Lathe interface {
	IsReady(context.Context) error
	HandleWorkpiece(context.Context) error
	Metrics(context.Context) (*Metrics, error)
}
