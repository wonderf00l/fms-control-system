package recognition

import "github.com/wonderf00l/fms-control-system/internal/errors"

type noWorkpieceToRecognizeError struct{}

func (e *noWorkpieceToRecognizeError) Error() string {
	return "There is no workpiece for recognition module"
}

func (e *noWorkpieceToRecognizeError) Type() errors.ErrorType {
	return errors.ErrInavlidStateForCMD
}
