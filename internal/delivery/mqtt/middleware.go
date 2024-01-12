package mqtt

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type Middleware func(mqtt.MessageHandler) mqtt.MessageHandler

func ApplyMiddlewareStack(initialHandler func(mqtt.Client, mqtt.Message), middlewares ...Middleware) func(mqtt.Client, mqtt.Message) {
	for _, middleware := range middlewares {
		initialHandler = middleware(initialHandler)
	}
	return initialHandler
}

func LoggingMiddleware(next mqtt.MessageHandler) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		clientVal, ok := client.(*ClientMQTT)

		if !ok {
			log.Println("Can't extract client from the interface")
			next(client, msg)
			return
		}

		start := time.Now()
		next(client, msg)

		clientVal.Log.Infow(
			"Got message",
			zap.String("topic", msg.Topic()),
			zap.Int8("QOS", int8(msg.Qos())),
			zap.Bool("retained", msg.Retained()),
			zap.Uint16("MessageID", uint16(msg.MessageID())),
			zap.Duration("processing_time", time.Since(start)),
		)
	}
}

func RecoverMiddleware(next mqtt.MessageHandler) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		defer func() {
			if err := recover(); err != nil {
				if clientVal, ok := client.(*ClientMQTT); ok {
					clientVal.Log.Errorln("RECOVER MIDDLEWARE GOT PANIC")
				} else {
					log.Println("RECOVER MIDDLEWARE GOT PANIC")
				}
			}
		}()

		next(client, msg)
	}
}
