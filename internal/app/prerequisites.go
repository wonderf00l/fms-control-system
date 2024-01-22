package app

import (
	"context"
	"flag"
	"strings"

	"github.com/joho/godotenv"
	"github.com/wonderf00l/fms-control-system/internal/configs"
	"github.com/wonderf00l/fms-control-system/internal/pkg/logger"
	"go.uber.org/config"
	"go.uber.org/zap"
)

var (
	appConfigPath = "configs/app.yml"
)

var (
	timeKey             = "timestamp"
	logEncoding         = flag.String("logenc", "json", "preferred logging format")
	logOutputPaths      = flag.String("log", "stdout,access.log", "file paths to write logging output to")
	logErrorOutputPaths = flag.String("logerror", "stderr,error.log", "path to write internal logger errors to.")
)

func InitPrerequisites(_ context.Context) (*zap.SugaredLogger, *configs.Configs, error) {
	appCfg, err := config.NewYAML(config.File(appConfigPath))
	if err != nil {
		return nil, nil, &initPrereqError{err}
	}

	cfgParameters, err := configs.ExtractCfgValues(appCfg)
	if err != nil {
		return nil, nil, &initPrereqError{err}
	}

	err = godotenv.Load(cfgParameters.EnvFile)
	if err != nil {
		return nil, nil, &initPrereqError{err}
	}

	configs, err := configs.NewConfigs(cfgParameters)
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

	return log, configs, nil
}
