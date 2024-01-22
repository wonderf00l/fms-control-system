package recognition

import "context"

type WorkpieceType uint8

const (
	Unknown WorkpieceType = iota
	King
	Queen
	Bishop
	Knight
	Castle
	Pawn
)

type Metrics struct {
	Ready bool
}

type Recognition interface {
	IsReady(context.Context) error
	Recognize(context.Context) (WorkpieceType, error)
	Metrics(context.Context) (*Metrics, error)
}
