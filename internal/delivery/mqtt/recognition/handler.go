package recognition

import (
	"context"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	deliveryMQTT "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
	"github.com/wonderf00l/fms-control-system/internal/service/recognition"
	"go.uber.org/zap"
)

type HandlerMQTT struct {
	Service recognition.Service
	log     *zap.SugaredLogger
}

func NewHandlerMQTT(service recognition.Service, log *zap.SugaredLogger) *HandlerMQTT {
	return &HandlerMQTT{
		Service: service,
		log:     log,
	}
}

func (h *HandlerMQTT) WorkpieceRecognition(client mqtt.Client, msg mqtt.Message) {
	ctx, err := deliveryMQTT.GetContextFromMessage(msg)
	reqTopic, respTopic := msg.Topic(), requestResponseTopics[msg.Topic()]

	if err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
		return
	}

	if workpieceType, err := h.Service.RecognizeWorkpiece(ctx); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	} else if err = deliveryMQTT.ResponseOk(h.log, client, reqTopic, respTopic, 1, false, fmt.Sprint("recognized workpiece: ", workpieceType), nil); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	}
}

func (h *HandlerMQTT) PushMetrics(ctx context.Context, client mqtt.Client) {
	if metrics, err := h.Service.GetMetrics(ctx); err != nil {
		deliveryMQTT.ResponseError(h.log, client, *new(string), pushMetrics, 1, false, err)
	} else if err = deliveryMQTT.ResponseOk(h.log, client, *new(string), pushMetrics, 1, true, "got recognition metrics sucessfully", metrics); err != nil {
		deliveryMQTT.ResponseError(h.log, client, *new(string), pushMetrics, 1, false, err)
	}
}
