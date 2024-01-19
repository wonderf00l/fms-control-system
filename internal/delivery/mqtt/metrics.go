package mqtt

import (
	"context"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

var (
	metricsPeroid = 3 * time.Second
)

type MetricsProvider interface {
	PushMetrics(context.Context, mqtt.Client)
}

func RunMetricsWorker(ctx context.Context, client mqtt.Client, provider MetricsProvider, service string, log *zap.SugaredLogger, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Infof("Service %s beginning to push it's metrics", service)
	for {
		select {
		case <-time.After(metricsPeroid):
			provider.PushMetrics(ctx, client)
		case <-ctx.Done():
			log.Infof("Service %s finishing pushing metrics", service)
			return
		}
	}
}
