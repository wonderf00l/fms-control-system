package mqtt

import (
	"context"
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessageWithCtx struct {
	mqtt.Message
	ctx context.Context
}

func (m *MessageWithCtx) Context() context.Context {
	return m.ctx
}

func GetContextFromMessage(msg mqtt.Message) (context.Context, error) {
	customMsg, ok := msg.(*MessageWithCtx)
	if !ok {
		return nil, errors.New("failed custom message conversion")
	}
	return customMsg.Context(), nil
}
