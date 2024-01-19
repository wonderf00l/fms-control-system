package mqtt

import (
	"crypto/tls"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/wonderf00l/fms-control-system/internal/configs"
)

var (
	_connectTimeout = 5 * time.Second
	_writeTimeout   = 5 * time.Second
	_keepAlive      = 50 * time.Second
	_pingTimeout    = 5 * time.Second
)

type DefaultHandlers struct {
	OnConnect         mqtt.OnConnectHandler
	OnConnectionLost  mqtt.ConnectionLostHandler
	OnReconnect       mqtt.ReconnectHandler
	ConnectionAttempt mqtt.ConnectionAttemptHandler
}

func setDefaultHandlers(opts *mqtt.ClientOptions, defaultHandlers DefaultHandlers) {
	opts.OnConnectionLost = defaultHandlers.OnConnectionLost
	opts.OnConnect = defaultHandlers.OnConnect
	opts.OnReconnecting = defaultHandlers.OnReconnect
}

func NewClientOptions(brokerCfg *configs.BrokerConfig, tlsCfg *tls.Config, defaultHandlers DefaultHandlers, clientID string) *mqtt.ClientOptions {
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

	setDefaultHandlers(opts, defaultHandlers)
	return opts
}
