package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wonderf00l/fms-control-system/internal/app"
)

func main() {
	serviceCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	serviceLogger, cfgFiles, err := app.InitPrerequisites(serviceCtx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = serviceLogger.Sync(); err != nil {
			log.Fatal("Sync service logger: ", err)
		}
	}()

	if err = app.Run(serviceCtx, serviceLogger, cfgFiles); err != nil {
		serviceLogger.Fatal(err)
	}
}
