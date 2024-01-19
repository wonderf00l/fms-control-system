package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type customMsgClientNotFoundError struct {
	msg mqtt.Message
}

func (e *customMsgClientNotFoundError) Error() string {
	return fmt.Sprintf("SERVICE WON'T PROCESS MESSAGE BECAUSE OF INTERNAL ERROR: topic - %q, messageID - %d\n", e.msg.Topic(), e.msg.MessageID())
}

type ServiceGotPanicError struct {
	msg mqtt.Message
}

func (e *ServiceGotPanicError) Error() string {
	return fmt.Sprintf("SERVICE GOT PANIC: topic - %q, messageID - %q\n", e.msg.Topic(), e.msg.MessageID())
}

type connectionError struct {
	service string
	inner   error
}

func (e *connectionError) Error() string {
	return fmt.Sprintf("Service %s: connection error: %s", e.service, e.inner.Error())
}
