package storage

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	deliveryMQTT "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
)

var (
	requestCmdTopicPrefix  = "/storage/cmd/request"
	responseCmdTopicPrefix = "/storage/cmd/response"
)

var (
	provideWorkpieceReq  = requestCmdTopicPrefix + "/provide"
	provideWorkpieceResp = responseCmdTopicPrefix + "/provide"

	acceptWorkpieceReq  = requestCmdTopicPrefix + "/accept"
	acceptWorkpieceResp = responseCmdTopicPrefix + "/accept"

	pushMetrics = "/storage/metrics"

	requestResponseTopics = map[string]string{
		provideWorkpieceReq: provideWorkpieceResp,
		acceptWorkpieceReq:  acceptWorkpieceResp,
	}
)

func SetSubscribeRouter(subCtx context.Context, client *deliveryMQTT.ClientMQTT, handler *HandlerMQTT) []mqtt.Token {
	tokens := make([]mqtt.Token, 0, 2)

	tokens = append(tokens,
		client.Subscribe(provideWorkpieceReq, 1, deliveryMQTT.ApplyMiddlewareStack(
			handler.ProvideWorkpiece,
			deliveryMQTT.LoggingMiddleware,
			deliveryMQTT.CheckMsgWithCtxMiddleware,
			deliveryMQTT.ReplaceMessageClientMiddleware(subCtx, client),
			deliveryMQTT.RecoverMiddleware),
		),
		client.Subscribe(acceptWorkpieceReq, 1, deliveryMQTT.ApplyMiddlewareStack(
			handler.AcceptWorkpiece,
			deliveryMQTT.LoggingMiddleware,
			deliveryMQTT.CheckMsgWithCtxMiddleware,
			deliveryMQTT.ReplaceMessageClientMiddleware(subCtx, client),
			deliveryMQTT.RecoverMiddleware,
		)))

	return tokens
}
