package errors

import "fmt"

type ErrorType uint8

type ServiceName string

const (
	_ ErrorType = iota
	ErrNotFound
	ErrAlreadyExists
	ErrInvalidInput
	ErrNotImplemented
	ErrTimeout
	ErrServiceNotAvailable
	ErrServiceNotReady
	ErrDepsNotReady
	ErrInavlidStateForCMD
)

const (
	Storage     ServiceName = "storage"
	Conveyor    ServiceName = "conveyor"
	Recognition ServiceName = "recognition"
	Lathe       ServiceName = "lathe"
	Miller      ServiceName = "miller"
)

type DeclaredError interface {
	Type() ErrorType
}

type InternalError struct {
	Message string
	Layer   string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Internal error occured. Message: '%s'. Layer: %s", e.Message, e.Layer)
}

type TimeoutExceededError struct {
	Service ServiceName
}

func (e *TimeoutExceededError) Error() string {
	return fmt.Sprintf("service %s: timeout exceeded", e.Service)
}

func (e *TimeoutExceededError) Type() ErrorType {
	return ErrTimeout
}

type ServiceOfflineError struct {
	Service ServiceName
}

func (e *ServiceOfflineError) Error() string {
	return fmt.Sprintf("service %s is offline right now", e.Service)
}

func (e *ServiceOfflineError) Type() ErrorType {
	return ErrServiceNotAvailable
}

type ServiceNotReadyError struct {
	Service ServiceName
}

func (e *ServiceNotReadyError) Error() string {
	return fmt.Sprintf("service %s is not ready", e.Service)
}

func (e *ServiceNotReadyError) Type() ErrorType {
	return ErrServiceNotReady
}

type DepsNotReadyError struct {
	Service    ServiceName
	Dependency ServiceName
	Reason     error
}

func (e *DepsNotReadyError) Error() string {
	return fmt.Sprintf("%s: dependency %v is not ready -- %s", e.Service, e.Dependency, e.Reason.Error())
}

func (e *DepsNotReadyError) Type() ErrorType {
	return ErrDepsNotReady
}
