package configs

import (
	"go.uber.org/config"
)

var (
	envKey                   = "app.env_file"
	brokerKey      ClientKey = "app.broker"
	StorageKey     ClientKey = "app.clients.storage"
	ConveyorKey    ClientKey = "app.clients.conveyor"
	RecognitionKey ClientKey = "app.clients.recognition"
	LatheKey       ClientKey = "app.clients.lathe"
	MillerKey      ClientKey = "app.clients.miller"
	clientsCfgKeys           = [5]ClientKey{StorageKey, ConveyorKey, RecognitionKey, LatheKey, MillerKey}
)

type brokerYAML struct {
	Scheme  string `yaml:"scheme"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type clientYAML struct {
	ID       string `yaml:"ID"`
	BrokerCA string `yaml:"brokerCA"`
	ClientCA string `yaml:"clientCA"`
	Cert     string `yaml:"cert"`
	Key      string `yaml:"key"`
}

func NewYAML(filename string) (*config.YAML, error) {
	cfg, err := config.NewYAML(config.File(filename))
	if err != nil {
		return nil, &yamlConfigError{filename: filename, inner: err}
	}
	return cfg, nil
}
