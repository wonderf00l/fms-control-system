package recognition

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	deliveryMQTT "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
	"github.com/wonderf00l/fms-control-system/internal/service/lathe"
	"go.uber.org/zap"
)

type HandlerMQTT struct {
	Service lathe.Service
	log     *zap.SugaredLogger
}

func NewHandlerMQTT(service lathe.Service, log *zap.SugaredLogger) *HandlerMQTT {
	return &HandlerMQTT{
		Service: service,
		log:     log,
	}
}

func (h *HandlerMQTT) HandleWorkpiece(client mqtt.Client, msg mqtt.Message) {
	ctx, err := deliveryMQTT.GetContextFromMessage(msg)
	reqTopic, respTopic := msg.Topic(), requestResponseTopics[msg.Topic()]

	if err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
		return
	}

	if err = h.Service.HandleWorkpiece(ctx); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	} else if err = deliveryMQTT.ResponseOk(h.log, client, reqTopic, respTopic, 1, false, "lathe handled workpiece sucessfully", nil); err != nil {
		deliveryMQTT.ResponseError(h.log, client, reqTopic, respTopic, 1, false, err)
	}
}

func (h *HandlerMQTT) PushMetrics(ctx context.Context, client mqtt.Client) {
	metrics, err := h.Service.GetMetrics(ctx)
	if err != nil {
		deliveryMQTT.ResponseError(h.log, client, *new(string), pushMetrics, 1, false, err)
		return
	}
	if err = deliveryMQTT.ResponseOk(h.log, client, *new(string), pushMetrics, 1, true, "got lathe metrics sucessfully", metrics); err != nil {
		deliveryMQTT.ResponseError(h.log, client, *new(string), pushMetrics, 1, false, err)
		return
	}
}
