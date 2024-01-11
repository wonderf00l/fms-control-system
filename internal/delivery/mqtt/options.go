package mqtt

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	yaml "go.uber.org/config"
)

var (
	_connectTimeout = 5 * time.Second
	_writeTimeout   = 5 * time.Second
	_keepAlive      = 50 * time.Second
	_pingTimeout    = 5 * time.Second
)

type BrokerConfig struct {
	Address  string
	Username string
	Password string
}

func BrokerAddress(cfgValue yaml.Value) string {
	return fmt.Sprintf("%s://%s:%s",
		cfgValue.Get("scheme").String(),
		cfgValue.Get("address").String(),
		cfgValue.Get("port").String(),
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

func setDefaultHandlers(opts *mqtt.ClientOptions, handler *HandlerMQTT) {
	opts.OnConnectionLost = handler.OnConnectionLost
	opts.OnConnect = handler.OnConnect
	opts.OnReconnecting = handler.OnReconnect
}

func NewClientOptions(brokerCfg *BrokerConfig, tlsCfg *tls.Config, handler *HandlerMQTT, clientID string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.SetTLSConfig(tlsCfg)

	opts.AddBroker(brokerCfg.Address)
	opts.SetUsername(brokerCfg.Username)
	opts.SetPassword(brokerCfg.Password)

	opts.SetClientID(clientID)
	opts.SetOrderMatters(true)
	opts.SetCleanSession(false)

	opts.SetConnectTimeout(_connectTimeout)
	opts.SetWriteTimeout(_writeTimeout)
	opts.SetKeepAlive(_keepAlive)
	opts.SetPingTimeout(_pingTimeout)

	opts.SetConnectRetry(true)
	opts.SetAutoReconnect(true)

	setDefaultHandlers(opts, handler)
	return opts
}
