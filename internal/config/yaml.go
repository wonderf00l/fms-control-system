package config

import (
	"fmt"

	"go.uber.org/config"
)

var (
	BrokerAddressName = "app.broker"
)

type YAMLConfigError struct {
	inner    error
	filename string
}

func (e *YAMLConfigError) Error() string {
	return fmt.Sprintf("YAML config %q: %s", e.filename, e.inner.Error())
}

func NewYAML(filename string) (*config.YAML, error) {
	cfg, err := config.NewYAML(config.File(filename))
	if err != nil {
		return nil, &YAMLConfigError{filename: filename, inner: err}
	}
	return cfg, nil
}
