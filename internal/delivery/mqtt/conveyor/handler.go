package recognition

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	deliveryMQTT "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
	"github.com/wonderf00l/fms-control-system/internal/service/conveyor"
	"go.uber.org/zap"
)

type HandlerMQTT struct {
	Service conveyor.Service
	log     *zap.SugaredLogger
}

func NewHandlerMQTT(service conveyor.Service, log *zap.SugaredLogger) *HandlerMQTT {
	return &HandlerMQTT{
		Service: service,
		log:     log,
	}
}

func (h *HandlerMQTT) ProvideWorkpieceToRecognition(client mqtt.Client, msg mqtt.Message) {
	ctx, err := deliveryMQTT.GetContextFromMessage(msg)
	reqTopic, respTopic := msg.Topic(), requestResponseTopics[msg.Topic()]

	if err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
		return
	}

	if err = h.Service.MoveWorkpieceToRecognition(ctx); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	} else if err = deliveryMQTT.ResponseOk(h.log, client, reqTopic, respTopic, 1, false, "conveyor provided workpiece to the recognition successfully", nil); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	}
}

func (h *HandlerMQTT) ProvideWorkpieceToLathe(client mqtt.Client, msg mqtt.Message) {
	ctx, err := deliveryMQTT.GetContextFromMessage(msg)
	reqTopic, respTopic := msg.Topic(), requestResponseTopics[msg.Topic()]

	if err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
		return
	}

	if err = h.Service.MoveWorkpieceToLathe(ctx); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	} else if err = deliveryMQTT.ResponseOk(h.log, client, reqTopic, respTopic, 1, false, "conveyor provided workpiece to the recognition successfully", nil); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	}
}

func (h *HandlerMQTT) ProvideWorkpieceToMiller(client mqtt.Client, msg mqtt.Message) {
	ctx, err := deliveryMQTT.GetContextFromMessage(msg)
	reqTopic, respTopic := msg.Topic(), requestResponseTopics[msg.Topic()]

	if err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
		return
	}

	if err = h.Service.MoveWorkpieceToMiller(ctx); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	} else if err = deliveryMQTT.ResponseOk(h.log, client, reqTopic, respTopic, 1, false, "conveyor provided workpiece to the recognition successfully", nil); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	}
}

func (h *HandlerMQTT) PushMetrics(ctx context.Context, client mqtt.Client) {
	metrics, err := h.Service.GetMetrics(ctx)
	if err != nil {
		deliveryMQTT.ResponseError(h.log, client, *new(string), pushMetrics, 1, false, err)
		return
	}
	if err = deliveryMQTT.ResponseOk(h.log, client, *new(string), pushMetrics, 1, true, "got conveyor metrics sucessfully", metrics); err != nil {
		deliveryMQTT.ResponseError(h.log, client, *new(string), pushMetrics, 1, false, err)
		return
	}
}
