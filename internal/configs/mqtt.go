package configs

import (
	"fmt"
	"os"
)

type BrokerConfig struct {
	Address  string
	Username string
	Password string
}

func brokerAddress(brokerParams *brokerYAML) string {
	return fmt.Sprintf("%s://%s:%d",
		brokerParams.Scheme,
		brokerParams.Address,
		brokerParams.Port,
	)
}

func NewBrokerConfig(address string) (*BrokerConfig, error) {
	username := os.Getenv("BROKER_USERNAME")
	password := os.Getenv("BROKER_PASSWORD")
	if len(username) == 0 || len(password) == 0 {
		return nil, fmt.Errorf("incorrect .env broker credentials")
	}

	return &BrokerConfig{
		Address:  address,
		Username: username,
		Password: password,
	}, nil
}
