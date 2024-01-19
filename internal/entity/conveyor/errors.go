package conveyor

import (
	"fmt"

	"github.com/wonderf00l/fms-control-system/internal/errors"
)

type invalidDistanceError struct {
	service  errors.ServiceName
	distance int
}

func (e *invalidDistanceError) Error() string {
	return fmt.Sprintf("Service %q: invalid distance %d", e.service, e.distance)
}

func (e *invalidDistanceError) Type() errors.ErrorType {
	return errors.ErrInvalidInput
}
