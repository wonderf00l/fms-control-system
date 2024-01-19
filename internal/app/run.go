package app

import (
	"context"
	"errors"
	"sync"

	"github.com/wonderf00l/fms-control-system/internal/configs"
	"github.com/wonderf00l/fms-control-system/internal/entity"
	entityConveyor "github.com/wonderf00l/fms-control-system/internal/entity/conveyor"
	entityLathe "github.com/wonderf00l/fms-control-system/internal/entity/lathe"
	entityMiller "github.com/wonderf00l/fms-control-system/internal/entity/miller"
	entityRecognition "github.com/wonderf00l/fms-control-system/internal/entity/recognition"
	entityStorage "github.com/wonderf00l/fms-control-system/internal/entity/storage"

	delivery "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt"
	deliveryStorage "github.com/wonderf00l/fms-control-system/internal/delivery/mqtt/storage"
	serviceStorage "github.com/wonderf00l/fms-control-system/internal/service/storage"

	"go.uber.org/zap"
)

var (
	timeForDisconnect uint = 100
)

func Run(ctx context.Context, logger *zap.SugaredLogger, cfgs *configs.Configs) error {
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	wg := &sync.WaitGroup{}
	connectErrCh := make(chan struct{}, 5)

	entityPool := entity.Pool{
		Storage:     entityStorage.NewStorage(),
		Conveyor:    entityConveyor.NewConveyor(),
		Recognition: entityRecognition.NewRecognition(),
		Lathe:       entityLathe.NewLathe(),
		Miller:      entityMiller.NewMiller(),
	}

	storageClient := delivery.NewClientMQTT(delivery.NewClientOptions(
		cfgs.BrokerConfig,
		cfgs.TLSConfigs[configs.StorageKey],
		delivery.DefaultHandlers{
			OnConnect:         deliveryStorage.OnConnect,
			OnConnectionLost:  deliveryStorage.OnConnectionLost,
			OnReconnect:       deliveryStorage.OnReconnect,
			ConnectionAttempt: deliveryStorage.ConnAtempt,
		}, cfgs.ClientIDs[configs.StorageKey]), logger)

	storageHandler := deliveryStorage.NewHandlerMQTT(serviceStorage.NewService(logger, &entityPool), logger)

	subCtx := context.Background()

	deliveryStorage.AddRoutes(storageClient, storageHandler, subCtx)

	wg.Add(1)
	go delivery.CheckConnect(ctx, storageClient.Connect(), "storage", logger, wg, connectErrCh)
	// ...
	wg.Wait() // waiting is obligatory

	if len(connectErrCh) > 0 {
		return &appRunError{inner: errors.New("one or more clients couldn't establish the connection")}
	}

	defer storageClient.Disconnect(timeForDisconnect)

	if err := deliveryStorage.SetSubscribeRouter(storageClient, storageHandler, subCtx); err != nil {
		return &appRunError{inner: err}
	}

	wg.Add(1)
	go delivery.RunMetricsWorker(ctx, storageClient, storageHandler, "storage", logger, wg)
	//...
	<-ctx.Done()
	wg.Wait()
	logger.Infoln("Shutting down gracefully")

	return nil
}
