package mqtt

import (
	"encoding/json"
	"errors"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	errorsPkg "github.com/wonderf00l/fms-control-system/internal/errors"
	"go.uber.org/zap"
)

var (
	StatusOK            = "ok"
	StatusError         = "error"
	InternalErrorCode   = "internal_error"
	InternalErorMessage = "internal service error occured"
)

type JSONResponse struct {
	RespUUID     string `json:"responseID"`
	RequestTopic string `json:"request_topic"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	Payload      any    `json:"payload"`
}

type JSONErrResponse struct {
	RespUUID     string `json:"responseID"`
	RequestTopic string `json:"request_topic"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	Code         string `json:"code"`
}

func ConvertErrorMQTT(err error) (code string, msg string) {
	var declaredErr errorsPkg.DeclaredError
	if errors.As(err, &declaredErr) {
		switch declaredErr.Type() {
		case errorsPkg.ErrNotFound:
			return "not_found", err.Error()
		case errorsPkg.ErrAlreadyExists:
			return "already_exists", err.Error()
		case errorsPkg.ErrInvalidInput:
			return "invalid_input", err.Error()
		case errorsPkg.ErrTimeout:
			return "timeout", err.Error()
		case errorsPkg.ErrServiceNotReady:
			return "service_not_ready", err.Error()
		case errorsPkg.ErrServiceNotAvailable:
			return "not_available", err.Error()
		case errorsPkg.ErrDepsNotReady:
			return "deps_not_ready", err.Error()
		case errorsPkg.ErrInavlidStateForCMD:
			return "invalid_cmd_state", err.Error()
		}
	}
	return InternalErrorCode, InternalErorMessage
}

func ResponseOk(log *zap.SugaredLogger, client mqtt.Client, reqTopic, respTopic string, qos byte, retained bool, msg string, payload any) error {
	serialized, err := json.Marshal(JSONResponse{
		RespUUID:     uuid.NewString(),
		RequestTopic: reqTopic,
		Status:       "ok",
		Message:      msg,
		Payload:      payload,
	})
	if err != nil {
		return err
	}
	token := client.Publish(respTopic, qos, retained, serialized)
	go func() {
		<-token.Done()
		if err := token.Error(); err != nil {
			log.Errorf("publish ok response: topic - %s, error - %s", respTopic, err.Error())
		}
	}()
	return nil
}

func ResponseError(serviceLog any, client mqtt.Client, reqTopic, respTopic string, qos byte, retained bool, err error) {
	logger, ok := serviceLog.(*zap.SugaredLogger)

	code, msg := ConvertErrorMQTT(err)
	if code == InternalErrorCode {
		if ok {
			logger.Errorf("Internal service error occured: %s", err.Error())
		} else {
			log.Printf("Internal service error occured: %s", err.Error())
		}
	}

	serialized, _ := json.Marshal(JSONErrResponse{
		RespUUID:     uuid.NewString(),
		RequestTopic: reqTopic,
		Status:       "error",
		Message:      msg,
		Code:         code,
	})
	token := client.Publish(respTopic, qos, retained, serialized)
	go func() {
		<-token.Done()
		if err := token.Error(); err != nil {
			if ok {
				logger.Errorf("publish error response: topic - %s, error - %s", respTopic, err.Error())
			} else {
				log.Printf("publish error response: topic - %s, error - %s", respTopic, err.Error())
			}

		}
	}()
}
