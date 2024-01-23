package conveyor

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

	provideToLatheReq  = requestCmdTopicPrefix + "/lathe"
	provideToLatheResp = responseCmdTopicPrefix + "/lathe"

	provideToMillerReq  = requestCmdTopicPrefix + "/miller"
	provideToMillerResp = responseCmdTopicPrefix + "/miller"

	pushMetrics = "/conveyor/metrics"

	requestResponseTopics = map[string]string{
		provideToRecReq:    provideToRecResp,
		provideToLatheReq:  provideToLatheResp,
		provideToMillerReq: provideToMillerResp,
	}
)

func SetSubscribeRouter(subCtx context.Context, client *deliveryMQTT.ClientMQTT, handler *HandlerMQTT) []mqtt.Token {
	tokens := make([]mqtt.Token, 0, 3)

	tokens = append(tokens, client.Subscribe(provideToRecReq, 1, deliveryMQTT.ApplyMiddlewareStack(
		handler.ProvideWorkpieceToRecognition,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(subCtx, client),
		deliveryMQTT.RecoverMiddleware,
	)), client.Subscribe(provideToLatheReq, 1, deliveryMQTT.ApplyMiddlewareStack(
		handler.ProvideWorkpieceToLathe,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(subCtx, client),
		deliveryMQTT.RecoverMiddleware,
	)), client.Subscribe(provideToMillerReq, 1, deliveryMQTT.ApplyMiddlewareStack(
		handler.ProvideWorkpieceToMiller,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(subCtx, client),
		deliveryMQTT.RecoverMiddleware,
	)))

	return tokens
}
