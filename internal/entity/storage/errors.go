package storage

import (
	"fmt"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

var (
	errPreifx = "Storage"
)

type allWorkpiecesHandledError struct{}

func (e *allWorkpiecesHandledError) Error() string {
	return fmt.Sprintf("%s: all workpieces are already handled", errPreifx)
}

func (e *allWorkpiecesHandledError) Type() errors.ErrorType {
	return errors.ErrInavlidStateForCMD
}

type workpieceAlreadyProcessingError struct{}

func (e *workpieceAlreadyProcessingError) Error() string {
	return fmt.Sprintf("%s: workpiece is already processing", errPreifx)
}

func (e *workpieceAlreadyProcessingError) Type() errors.ErrorType {
	return errors.ErrInavlidStateForCMD
}

type noWorkpieceProcessingError struct{}

func (e *noWorkpieceProcessingError) Error() string {
	return fmt.Sprintf("%s: no workpieces processing", errPreifx)
}

func (e *noWorkpieceProcessingError) Type() errors.ErrorType {
	return errors.ErrInavlidStateForCMD
}
