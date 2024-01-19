package storage

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	deliveryMQTT "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
	"github.com/wonderf00l/fms-control-system/internal/service/storage"
	"go.uber.org/zap"
)

type HandlerMQTT struct {
	Service storage.Service
	log     *zap.SugaredLogger
}

func NewHandlerMQTT(service storage.Service, log *zap.SugaredLogger) *HandlerMQTT {
	return &HandlerMQTT{
		Service: service,
		log:     log,
	}
}

func (h *HandlerMQTT) ProvideWorkpiece(client mqtt.Client, msg mqtt.Message) {
	ctx, err := deliveryMQTT.GetContextFromMessage(msg)
	reqTopic, respTopic := msg.Topic(), requestResponseTopics[msg.Topic()]

	if err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
		return
	}

	if err := h.Service.ProvideWorkpiece(ctx); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	} else if err := deliveryMQTT.ResponseOk(h.log, client, reqTopic, respTopic, 1, false, "storage has provided workpiece successfully", nil); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	}
}

func (h *HandlerMQTT) AcceptWorkpiece(client mqtt.Client, msg mqtt.Message) {
	ctx, err := deliveryMQTT.GetContextFromMessage(msg)
	reqTopic, respTopic := msg.Topic(), requestResponseTopics[msg.Topic()]

	if err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
		return
	}

	if err := h.Service.AcceptWorkpiece(ctx); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	} else if err := deliveryMQTT.ResponseOk(h.log, client, reqTopic, respTopic, 1, false, "storage has accepted workpiece successfully", nil); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	}
}

func (h *HandlerMQTT) PushMetrics(ctx context.Context, client mqtt.Client) {
	if metrics, err := h.Service.GetMetrics(ctx); err != nil {
		deliveryMQTT.ResponseError(h.log, client, *new(string), pushMetrics, 1, false, err)
	} else if err := deliveryMQTT.ResponseOk(h.log, client, *new(string), pushMetrics, 1, true, "got storage metrics sucessfully", metrics); err != nil {
		deliveryMQTT.ResponseError(h.log, client, *new(string), pushMetrics, 1, false, err)
	}
}
