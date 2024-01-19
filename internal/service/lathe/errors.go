package lathe

import "github.com/wonderf00l/fms-control-system/internal/errors"

type noWorkpieceForLatheToHandleError struct{}

func (e *noWorkpieceForLatheToHandleError) Error() string {
	return "There is no workpiece for lathe to handle"
}

func (e *noWorkpieceForLatheToHandleError) Type() errors.ErrorType {
	return errors.ErrInavlidStateForCMD
}
