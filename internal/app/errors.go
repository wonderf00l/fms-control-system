package app

import "fmt"

type initPrereqError struct {
	inner error
}

func (e *initPrereqError) Error() string {
	return fmt.Sprintf("Init prerequisites: %s", e.inner.Error())
}

type appRunError struct {
	inner error
}

func (e *appRunError) Error() string {
	return fmt.Sprintf("run application: %s", e.inner.Error())
}
