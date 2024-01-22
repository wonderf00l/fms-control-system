package configs

import (
	"crypto/tls"
	"fmt"

	"go.uber.org/config"
)

type ClientKey string

type Configs struct {
	BrokerConfig *BrokerConfig
	TLSConfigs   map[ClientKey]*tls.Config
	ClientIDs    map[ClientKey]string
}

type Parameters struct {
	EnvFile       string
	BrokerAddress string
	ClientIDs     map[ClientKey]string
	TLSParams     map[ClientKey]clientParamsTLS
}

func ExtractCfgValues(appCfg *config.YAML) (*Parameters, error) {
	var (
		tlsParams = map[ClientKey]clientParamsTLS{}
		clientIDs = map[ClientKey]string{}
	)

	envExtracted := ""
	if err := appCfg.Get(envKey).Populate(&envExtracted); err != nil {
		return nil, &configsError{inner: err}
	}

	brokerExtracted := brokerYAML{}
	if err := appCfg.Get(string(brokerKey)).Populate(&brokerExtracted); err != nil {
		return nil, &configsError{inner: err}
	}

	for _, clKey := range clientsCfgKeys {
		extracted := clientYAML{}
		if err := appCfg.Get(string(clKey)).Populate(&extracted); err != nil {
			return nil, &configsError{inner: fmt.Errorf("extract values: %w", err)}
		}
		tlsParams[clKey] = clientParamsTLS{
			brokerCA:         extracted.BrokerCA,
			clientCA:         extracted.ClientCA,
			clientCert:       extracted.Cert,
			clientPrivateKey: extracted.Key,
		}
		clientIDs[clKey] = extracted.ID
	}

	return &Parameters{
		EnvFile:       envExtracted,
		BrokerAddress: brokerAddress(&brokerExtracted),
		ClientIDs:     clientIDs,
		TLSParams:     tlsParams,
	}, nil
}

func NewConfigs(params *Parameters) (*Configs, error) {
	tlsConfigs := make(map[ClientKey]*tls.Config, len(params.TLSParams))
	for clKey, TLSparams := range params.TLSParams {
		tlsCfg, err := NewTLSConfig(TLSparams)
		if err != nil {
			return nil, &configsError{inner: fmt.Errorf("new configs: %w", err)}
		}
		tlsConfigs[clKey] = tlsCfg
	}

	brokerCfg, err := NewBrokerConfig(params.BrokerAddress)
	if err != nil {
		return nil, &configsError{inner: err}
	}

	return &Configs{
		BrokerConfig: brokerCfg,
		TLSConfigs:   tlsConfigs,
		ClientIDs:    params.ClientIDs,
	}, nil
}
