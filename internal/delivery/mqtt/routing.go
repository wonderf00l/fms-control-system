package mqtt

import mqtt "github.com/eclipse/paho.mqtt.golang"

var (
	startFmsRoute    = "/fms/cmd/start"
	stopFmsRoute     = "/fms/cmd/stop"
	continueFmsRoute = "/fms/cmd/continue"
)

func SetClientRouter(client mqtt.Client, handler *HandlerMQTT) {
	client.AddRoute(startFmsRoute, ApplyMiddlewareStack(handler.ProcessStartCommand, LoggingMiddleware, RecoverMiddleware))
	client.AddRoute(stopFmsRoute, ApplyMiddlewareStack(handler.ProcessStartCommand, LoggingMiddleware, RecoverMiddleware))
	client.AddRoute(continueFmsRoute, ApplyMiddlewareStack(handler.ProcessStartCommand, LoggingMiddleware, RecoverMiddleware))
}
