package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type HandlerMQTT struct {
	log *zap.SugaredLogger
}

func (h *HandlerMQTT) ProcessStartCommand(client mqtt.Client, msg mqtt.Message) {}

func (h *HandlerMQTT) ProcessStopCommand(client mqtt.Client, msg mqtt.Message) {}

func (h *HandlerMQTT) ProcessContinueCommand(client mqtt.Client, msg mqtt.Message) {}

func (h *HandlerMQTT) OnConnectionLost(client mqtt.Client, err error) {}

func (h *HandlerMQTT) OnConnect(client mqtt.Client) {}

func (h *HandlerMQTT) OnReconnect(client mqtt.Client, opts *mqtt.ClientOptions) {}
