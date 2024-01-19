package mqtt

import (
	"context"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

var (
	topicForPanicMiddleware    = "/aux/accidents/panic"
	retainedForPanicMiddleware = false
	qosForPanicMiddleware      = 1

	topicForCheckMsgMiddleware    = "/aux/accidents/internals"
	retainedForCheckMsgMiddleware = false
	qosForCheckMsgMiddleware      = 1
)

type Middleware func(mqtt.MessageHandler) mqtt.MessageHandler

func ApplyMiddlewareStack(initialHandler mqtt.MessageHandler, middlewares ...Middleware) mqtt.MessageHandler {
	for _, middleware := range middlewares {
		initialHandler = middleware(initialHandler)
	}
	return initialHandler
}

func ReplaceMessageClientMiddleware(customClient *ClientMQTT, ctx context.Context) Middleware {
	return func(next mqtt.MessageHandler) mqtt.MessageHandler {
		return func(client mqtt.Client, msg mqtt.Message) {
			customMsg := &MessageWithCtx{Message: msg, ctx: ctx}
			next(customClient, customMsg)
		}
	}
}

func CheckMsgWithCtxMiddleware(next mqtt.MessageHandler) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		_, gotClient := client.(*ClientMQTT)
		if !gotClient {
			log.Println("Check msg and client middleware: can't extract custom client")
		}

		_, gotMsg := msg.(*MessageWithCtx)
		if !gotMsg {
			log.Println("Check msg and client middleware: can't extract custom message")
		}

		if !gotClient || !gotMsg {
			ResponseError(
				struct{}{},
				client, msg.Topic(),
				topicForCheckMsgMiddleware,
				byte(qosForCheckMsgMiddleware),
				retainedForCheckMsgMiddleware,
				&customMsgClientNotFoundError{msg: msg},
			)
			return
		}
		next(client, msg)
	}
}

func LoggingMiddleware(next mqtt.MessageHandler) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		customClient, ok := client.(*ClientMQTT)

		if !ok {
			log.Println("Logging middleware: can't extract client from the interface")
			next(client, msg)
			return
		}

		start := time.Now()
		next(client, msg)

		customClient.Log.Infow(
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
				ResponseError(
					struct{}{},
					client, msg.Topic(),
					topicForPanicMiddleware,
					byte(qosForPanicMiddleware),
					retainedForPanicMiddleware,
					&ServiceGotPanicError{msg: msg},
				)
			}
		}()

		next(client, msg)
	}
}
