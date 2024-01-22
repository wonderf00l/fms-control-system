package app

import (
	"context"
	"errors"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/wonderf00l/fms-control-system/internal/configs"
	delivery "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
	"github.com/wonderf00l/fms-control-system/internal/entity"
	"go.uber.org/zap"

	entityConveyor "github.com/wonderf00l/fms-control-system/internal/entity/conveyor"
	entityLathe "github.com/wonderf00l/fms-control-system/internal/entity/lathe"
	entityMiller "github.com/wonderf00l/fms-control-system/internal/entity/miller"
	entityRecognition "github.com/wonderf00l/fms-control-system/internal/entity/recognition"
	entityStorage "github.com/wonderf00l/fms-control-system/internal/entity/storage"

	deliveryConveyor "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt/conveyor"
	deliveryLathe "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt/lathe"
	deliveryMiller "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt/miller"
	deliveryRecognition "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt/recognition"
	deliveryStorage "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt/storage"

	serviceConveyor "github.com/wonderf00l/fms-control-system/internal/service/conveyor"
	serviceLathe "github.com/wonderf00l/fms-control-system/internal/service/lathe"
	serviceMiller "github.com/wonderf00l/fms-control-system/internal/service/miller"
	serviceRecognition "github.com/wonderf00l/fms-control-system/internal/service/recognition"
	serviceStorage "github.com/wonderf00l/fms-control-system/internal/service/storage"
)

var (
	timeForDisconnect uint = 100
)

func Run(ctx context.Context, logger *zap.SugaredLogger, cfgs *configs.Configs) error {
	wg := &sync.WaitGroup{}
	tokenErrCh := make(chan struct{}, 1024)

	entityPool := entity.Pool{
		Storage:     entityStorage.NewStorage(),
		Conveyor:    entityConveyor.NewConveyor(),
		Recognition: entityRecognition.NewRecognition(),
		Lathe:       entityLathe.NewLathe(),
		Miller:      entityMiller.NewMiller(),
	}

	storageHandler := deliveryStorage.NewHandlerMQTT(serviceStorage.NewService(logger, &entityPool), logger)
	conveyorHandler := deliveryConveyor.NewHandlerMQTT(serviceConveyor.NewService(logger, &entityPool), logger)
	recognitionHandler := deliveryRecognition.NewHandlerMQTT(serviceRecognition.NewService(logger, &entityPool), logger)
	latheHandler := deliveryLathe.NewHandlerMQTT(serviceLathe.NewService(logger, &entityPool), logger)
	millerHandler := deliveryMiller.NewHandlerMQTT(serviceMiller.NewService(logger, &entityPool), logger)

	storageClient := delivery.NewClientMQTT(delivery.NewClientOptions(
		cfgs.BrokerConfig,
		cfgs.TLSConfigs[configs.StorageKey],
		delivery.DefaultHandlers{
			OnConnect:         deliveryStorage.OnConnect,
			OnConnectionLost:  deliveryStorage.OnConnectionLost,
			OnReconnect:       deliveryStorage.OnReconnect,
			ConnectionAttempt: deliveryStorage.ConnAtempt,
		}, cfgs.ClientIDs[configs.StorageKey]), logger)
	conveyorClient := delivery.NewClientMQTT(delivery.NewClientOptions(
		cfgs.BrokerConfig,
		cfgs.TLSConfigs[configs.ConveyorKey],
		delivery.DefaultHandlers{
			OnConnect:         deliveryConveyor.OnConnect,
			OnConnectionLost:  deliveryConveyor.OnConnectionLost,
			OnReconnect:       deliveryConveyor.OnReconnect,
			ConnectionAttempt: deliveryConveyor.ConnAtempt,
		}, cfgs.ClientIDs[configs.ConveyorKey]), logger)
	recognitionClient := delivery.NewClientMQTT(delivery.NewClientOptions(
		cfgs.BrokerConfig,
		cfgs.TLSConfigs[configs.RecognitionKey],
		delivery.DefaultHandlers{
			OnConnect:         deliveryRecognition.OnConnect,
			OnConnectionLost:  deliveryRecognition.OnConnectionLost,
			OnReconnect:       deliveryRecognition.OnReconnect,
			ConnectionAttempt: deliveryRecognition.ConnAtempt,
		}, cfgs.ClientIDs[configs.RecognitionKey]), logger)
	latheClient := delivery.NewClientMQTT(delivery.NewClientOptions(
		cfgs.BrokerConfig,
		cfgs.TLSConfigs[configs.LatheKey],
		delivery.DefaultHandlers{
			OnConnect:         deliveryLathe.OnConnect,
			OnConnectionLost:  deliveryLathe.OnConnectionLost,
			OnReconnect:       deliveryLathe.OnReconnect,
			ConnectionAttempt: deliveryLathe.ConnAtempt,
		}, cfgs.ClientIDs[configs.LatheKey]), logger)
	millerClient := delivery.NewClientMQTT(delivery.NewClientOptions(
		cfgs.BrokerConfig,
		cfgs.TLSConfigs[configs.MillerKey],
		delivery.DefaultHandlers{
			OnConnect:         deliveryMiller.OnConnect,
			OnConnectionLost:  deliveryMiller.OnConnectionLost,
			OnReconnect:       deliveryMiller.OnReconnect,
			ConnectionAttempt: deliveryMiller.ConnAtempt,
		}, cfgs.ClientIDs[configs.MillerKey]), logger)

	wg.Add(5)
	go delivery.CheckMQTTToken(ctx, storageClient.Connect(), "storage client connect", logger, wg, tokenErrCh)
	go delivery.CheckMQTTToken(ctx, conveyorClient.Connect(), "conveyor client connect", logger, wg, tokenErrCh)
	go delivery.CheckMQTTToken(ctx, recognitionClient.Connect(), "recognition client connect", logger, wg, tokenErrCh)
	go delivery.CheckMQTTToken(ctx, latheClient.Connect(), "lathe client connect", logger, wg, tokenErrCh)
	go delivery.CheckMQTTToken(ctx, millerClient.Connect(), "miller client connect", logger, wg, tokenErrCh)
	wg.Wait() // waiting is obligatory

	if len(tokenErrCh) > 0 {
		return &appRunError{inner: errors.New("one or more clients couldn't establish the connection")}
	}

	defer storageClient.Disconnect(timeForDisconnect)
	defer recognitionClient.Disconnect(timeForDisconnect)
	defer conveyorClient.Disconnect(timeForDisconnect)
	defer latheClient.Disconnect(timeForDisconnect)
	defer millerClient.Disconnect(timeForDisconnect)

	subscribeTokens := make([]mqtt.Token, 0)
	subscribeTokens = append(subscribeTokens, deliveryStorage.SetSubscribeRouter(ctx, storageClient, storageHandler)...)
	subscribeTokens = append(subscribeTokens, deliveryConveyor.SetSubscribeRouter(ctx, conveyorClient, conveyorHandler)...)
	subscribeTokens = append(subscribeTokens, deliveryRecognition.SetSubscribeRouter(ctx, recognitionClient, recognitionHandler)...)
	subscribeTokens = append(subscribeTokens, deliveryLathe.SetSubscribeRouter(ctx, latheClient, latheHandler)...)
	subscribeTokens = append(subscribeTokens, deliveryMiller.SetSubscribeRouter(ctx, millerClient, millerHandler)...)

	wg.Add(len(subscribeTokens))
	for _, token := range subscribeTokens {
		go delivery.CheckMQTTToken(ctx, token, "making subscription", logger, wg, tokenErrCh)
	}
	wg.Wait()

	if len(tokenErrCh) > 0 {
		return &appRunError{inner: errors.New("one or more clients couldn't make a subscription")}
	}

	wg.Add(5)
	go delivery.RunMetricsWorker(ctx, storageClient, storageHandler, "storage", logger, wg)
	go delivery.RunMetricsWorker(ctx, recognitionClient, conveyorHandler, "conveyor", logger, wg)
	go delivery.RunMetricsWorker(ctx, recognitionClient, recognitionHandler, "recognition", logger, wg)
	go delivery.RunMetricsWorker(ctx, latheClient, latheHandler, "lathe", logger, wg)
	go delivery.RunMetricsWorker(ctx, millerClient, millerHandler, "miller", logger, wg)

	<-ctx.Done()
	wg.Wait()
	logger.Infoln("Shutting down gracefully")

	return nil
}
