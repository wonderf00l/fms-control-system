package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/wonderf00l/fms-control-system/internal/config"
	"github.com/wonderf00l/fms-control-system/internal/pkg/logger"
	"go.uber.org/zap"
)

var (
	envFilePath    = ".env"
	appConfigPath  = "configs/app.yml"
	brokerCertPath = "certs/rootCA.crt"
)

var (
	timeKey             = "timestamp"
	logEncoding         = flag.String("logenc", "json", "preferred logging format")
	logOutputPaths      = flag.String("log", "stdout,access.log", "file paths to write logging output to")
	logErrorOutputPaths = flag.String("logerror", "stderr,error.log", "path to write internal logger errors to.")
)

type initPrereqError struct {
	inner error
}

func (e *initPrereqError) Error() string {
	return fmt.Sprintf("Init prerequisites: %s", e.inner.Error())
}

func initPrerequisites() (*zap.SugaredLogger, *config.CfgFiles, error) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return nil, nil, &initPrereqError{err}
	}

	flag.Parse()

	log, err := logger.New(logger.NewConfig(
		logger.ConfigureTimeKey(timeKey),
		logger.ConfigureEncoding(*logEncoding),
		logger.ConfigureOutput(strings.Split(*logOutputPaths, ",")),
		logger.ConfigureErrorOutput(strings.Split(*logErrorOutputPaths, ",")),
	))
	if err != nil {
		return nil, nil, &initPrereqError{err}
	}

	cfgFiles, err := config.NewConfigFiles(&config.CfgParameters{
		AppCfgFilename:    appConfigPath,
		TLSBrokerCertFile: brokerCertPath,
	})
	if err != nil {
		return nil, nil, &initPrereqError{err}
	}

	return log, cfgFiles, nil
}
