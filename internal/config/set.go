package config

import (
	"crypto/tls"
	"fmt"

	"github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
	"go.uber.org/config"
)

type CfgFiles struct {
	AppConfig    *config.YAML
	BrokerConfig *mqtt.BrokerConfig
	TLSConfig    *tls.Config
}

type CfgParameters struct {
	AppCfgFilename string
	TLSCfgParams   TlsCfgParameters
}

type TlsCfgParameters struct {
	BrokerCertFile   string
	ClientCertFile   string
	ClientPrivateKey string
	ClientCA         string
}

type ConfigFilesError struct {
	inner error
}

func (e *ConfigFilesError) Error() string {
	return fmt.Sprintf("Config files: %s", e.inner.Error())
}

func NewConfigFiles(params *CfgParameters) (*CfgFiles, error) {
	appCfg, err := NewYAML(params.AppCfgFilename)
	if err != nil {
		return nil, &ConfigFilesError{err}
	}

	tlsCfg, err := NewTLSConfig(params.TLSCfgParams)
	if err != nil {
		return nil, &ConfigFilesError{err}
	}

	brokerCfg, err := mqtt.NewBrokerConfig(mqtt.BrokerAddress(appCfg.Get(BrokerAddressName)))
	if err != nil {
		return nil, &ConfigFilesError{err}
	}

	return &CfgFiles{
		AppConfig:    appCfg,
		BrokerConfig: brokerCfg,
		TLSConfig:    tlsCfg,
	}, nil
}
