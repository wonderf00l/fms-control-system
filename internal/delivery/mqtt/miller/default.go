package miller

import (
	"crypto/tls"
	"log"
	"net/url"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	OnConnect mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Println("on connect miller: successfully")
	}
	OnConnectionLost mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Printf("on connection lost miller: %s\n", err.Error())
	}
	OnReconnect mqtt.ReconnectHandler = func(client mqtt.Client, _ *mqtt.ClientOptions) {
		log.Println("on reconnect miller: trying to reconnect")
	}
	ConnAtempt mqtt.ConnectionAttemptHandler = func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		log.Printf("conntAtempt miller: broker address - %q\n", broker.String())
		return nil
	}
)
