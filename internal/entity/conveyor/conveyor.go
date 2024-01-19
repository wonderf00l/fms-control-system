package conveyor

import "context"

type Metrics struct {
	Ready             bool
	WorkpieceLocation any
}

type WorkpieceLocation struct {
	X int
	Y int
	Z int
}

type Conveyor interface {
	IsReady(context.Context) error
	MoveWorkpieceHorisontal(context.Context, int) error
	MoveWorkpieceVertical(context.Context, int) error
	GetWorkpieceLocation(context.Context) (*WorkpieceLocation, error)
	Metrics(context.Context) (*Metrics, error)
}
