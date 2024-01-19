package miller

import "github.com/wonderf00l/fms-control-system/internal/errors"

type noWorkpieceForMillerToHandleError struct{}

func (e *noWorkpieceForMillerToHandleError) Error() string {
	return "There is no workpiece for miller to handle"
}

func (e *noWorkpieceForMillerToHandleError) Type() errors.ErrorType {
	return errors.ErrInavlidStateForCMD
}
