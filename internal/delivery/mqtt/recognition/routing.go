package recognition

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	deliveryMQTT "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
)

var (
	requestCmdTopicPrefix  = "/recognition/cmd/request"
	responseCmdTopicPrefix = "/recognition/cmd/response"
)

var (
	recognizeWorkpieceReq  = requestCmdTopicPrefix + "/recognize"
	recognizeWorkpieceResp = responseCmdTopicPrefix + "/recognize"

	pushMetrics = "/recognition/metrics"

	requestResponseTopics = map[string]string{
		recognizeWorkpieceReq: recognizeWorkpieceResp,
	}
)

func SetSubscribeRouter(subCtx context.Context, client *deliveryMQTT.ClientMQTT, handler *HandlerMQTT) []mqtt.Token {
	tokens := make([]mqtt.Token, 0, 1)

	tokens = append(tokens, client.Subscribe(recognizeWorkpieceReq, 1, deliveryMQTT.ApplyMiddlewareStack(
		handler.WorkpieceRecognition,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(subCtx, client),
		deliveryMQTT.RecoverMiddleware,
	)))

	return tokens
}
