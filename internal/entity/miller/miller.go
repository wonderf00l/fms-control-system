package miller

import "context"

type Metrics struct {
	Ready bool
}

type Miller interface {
	IsReady(context.Context) error
	HandleWorkpiece(context.Context) error
	Metrics(context.Context) (*Metrics, error)
}
