package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wonderf00l/fms-control-system/internal/app"
)

func main() {
	serviceLogger, cfgFiles, err := initPrerequisites()
	if err != nil {
		log.Fatal(err)
	}
	defer serviceLogger.Sync()

	if err := app.Run(serviceLogger, cfgFiles); err != nil {
		serviceLogger.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	serviceLogger.Infoln("Shutting down gracefully")
}
