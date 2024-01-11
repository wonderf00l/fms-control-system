package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ClientMQTT struct {
	mqtt.Client
	Log *zap.SugaredLogger
}

func NewClientMQTT(opts *mqtt.ClientOptions, handler *HandlerMQTT, log *zap.SugaredLogger) *ClientMQTT {
	client := mqtt.NewClient(opts)
	SetClientRouter(client, handler)

	return &ClientMQTT{Client: client, Log: log}
}

func (c *ClientMQTT) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	token := c.Client.Publish(topic, qos, retained, payload)
	pubID, err := uuid.NewRandom()
	if err != nil {
		c.Log.Warnln("Generate PublishID: ", err)
	}
	c.Log.Infow(
		"Publish message:",
		zap.String("ID", pubID.String()),
		zap.String("topic", topic),
		zap.Int8("QOS", int8(qos)),
		zap.Bool("retained", retained),
	)
	return token
}
