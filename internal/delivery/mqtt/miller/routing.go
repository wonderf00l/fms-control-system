package recognition

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	deliveryMQTT "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
)

var (
	requestCmdTopicPrefix  = "/miller/cmd/request"
	responseCmdTopicPrefix = "/miller/cmd/response"
)

var (
	hanldeReq  = requestCmdTopicPrefix + "/handle"
	handleResp = responseCmdTopicPrefix + "/handle"

	pushMetrics = "/miller/metrics"

	requestResponseTopics = map[string]string{
		hanldeReq: handleResp,
	}
)

func SetSubscribeRouter(subCtx context.Context, client *deliveryMQTT.ClientMQTT, handler *HandlerMQTT) []mqtt.Token {
	tokens := make([]mqtt.Token, 0, 1)

	tokens = append(tokens, client.Subscribe(hanldeReq, 1, deliveryMQTT.ApplyMiddlewareStack(
		handler.HandleWorkpiece,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(subCtx, client),
		deliveryMQTT.RecoverMiddleware,
	)))

	return tokens
}
