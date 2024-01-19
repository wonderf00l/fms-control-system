package recognition

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	deliveryMQTT "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
)

var (
	requestCmdTopicPrefix  = "/conveyor/cmd/request"
	responseCmdTopicPrefix = "/conveyor/cmd/response"
)

var (
	provideToRecReq  = requestCmdTopicPrefix + "/recognition"
	provideToRecResp = responseCmdTopicPrefix + "/recognition"

	pushMetrics = "/storage/metrics"

	requestResponseTopics = map[string]string{
		provideToRecReq: provideToRecResp,
	}
)

func SetSubscribeRouter(client *deliveryMQTT.ClientMQTT, handler *HandlerMQTT, subCtx context.Context) []mqtt.Token {
	tokens := make([]mqtt.Token, 0, 1)

	tokens = append(tokens, client.Subscribe(provideToRecReq, 1, deliveryMQTT.ApplyMiddlewareStack(
		handler.ProvideWorkpieceToRecognition,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(client, subCtx),
		deliveryMQTT.RecoverMiddleware,
	)))

	return tokens
}
