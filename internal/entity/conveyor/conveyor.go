package conveyor

import "context"

type Metrics struct {
	Ready             bool
	WorkpieceLocation any
}

type WorkpieceLocation struct {
	X int
	Y int
}

type Conveyor interface {
	IsReady(context.Context) (*WorkpieceLocation, error)
	SetWorkpieceLocation(context.Context, WorkpieceLocation, bool) error
	Metrics(context.Context) (*Metrics, error)
}
