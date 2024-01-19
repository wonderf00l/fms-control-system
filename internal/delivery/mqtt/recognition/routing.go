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

	pushMetrics = "/storage/metrics"

	requestResponseTopics = map[string]string{
		recognizeWorkpieceReq: recognizeWorkpieceResp,
	}
)

func SetSubscribeRouter(client *deliveryMQTT.ClientMQTT, handler *HandlerMQTT, subCtx context.Context) []mqtt.Token {
	tokens := make([]mqtt.Token, 0, 1)

	tokens = append(tokens, client.Subscribe(recognizeWorkpieceReq, 1, deliveryMQTT.ApplyMiddlewareStack(
		handler.WorkpieceRecognition,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(client, subCtx),
		deliveryMQTT.RecoverMiddleware,
	)))

	return tokens
}
