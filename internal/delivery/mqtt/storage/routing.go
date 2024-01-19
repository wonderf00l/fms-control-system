package storage

import (
	"context"

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

func AddRoutes(client *deliveryMQTT.ClientMQTT, handler *HandlerMQTT, subCtx context.Context) {
	client.AddRoute(provideWorkpieceReq, deliveryMQTT.ApplyMiddlewareStack(
		handler.ProvideWorkpiece,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(client, subCtx),
		deliveryMQTT.RecoverMiddleware))
	client.AddRoute(acceptWorkpieceReq, deliveryMQTT.ApplyMiddlewareStack(
		handler.AcceptWorkpiece,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(client, subCtx),
		deliveryMQTT.RecoverMiddleware))
}

func SetSubscribeRouter(client *deliveryMQTT.ClientMQTT, handler *HandlerMQTT, subCtx context.Context) error {
	if tok := client.Subscribe(provideWorkpieceReq, 1, deliveryMQTT.ApplyMiddlewareStack(
		handler.ProvideWorkpiece,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(client, subCtx),
		deliveryMQTT.RecoverMiddleware,
	)); tok.Wait() && tok.Error() != nil {
		return tok.Error()
	}

	if tok := client.Subscribe(acceptWorkpieceReq, 1, deliveryMQTT.ApplyMiddlewareStack(
		handler.AcceptWorkpiece,
		deliveryMQTT.LoggingMiddleware,
		deliveryMQTT.CheckMsgWithCtxMiddleware,
		deliveryMQTT.ReplaceMessageClientMiddleware(client, subCtx),
		deliveryMQTT.RecoverMiddleware,
	)); tok.Wait() && tok.Error() != nil {
		return tok.Error()
	}

	return nil
}
