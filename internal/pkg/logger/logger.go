package logger

import "go.uber.org/zap"

func New(config *zap.Config) (*zap.SugaredLogger, error) {
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
