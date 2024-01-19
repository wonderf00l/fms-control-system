package configs

import "fmt"

type yamlConfigError struct {
	inner    error
	filename string
}

func (e *yamlConfigError) Error() string {
	return fmt.Sprintf("YAML config %q: %s", e.filename, e.inner.Error())
}

type configsError struct {
	inner error
}

func (e *configsError) Error() string {
	return fmt.Sprintf("Configs: %s", e.inner.Error())
}

type tlsConfigError struct {
	inner error
}

func (e *tlsConfigError) Error() string {
	return fmt.Sprintf("TLS config: %s", e.inner.Error())
}
