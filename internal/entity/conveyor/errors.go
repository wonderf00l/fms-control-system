package conveyor

import (
	"fmt"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

type invalidLocationError struct {
	location WorkpieceLocation
}

func (e *invalidLocationError) Error() string {
	return fmt.Sprintf("Conveyor: invalid location: X - %d, Y - %d", e.location.X, e.location.Y)
}

func (e *invalidLocationError) Type() errors.ErrorType {
	return errors.ErrInvalidInput
}

type workpieceNotOnConveyorError struct{}

func (e *workpieceNotOnConveyorError) Error() string {
	return "workpiece is not on the conveyor"
}

func (e *workpieceNotOnConveyorError) Type() errors.ErrorType {
	return errors.ErrInavlidStateForCMD
}
