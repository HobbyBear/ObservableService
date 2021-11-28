package zaplogger

import (
	"ObservableService/pkg/logger"
	"go.uber.org/zap"
)

func NewLogger(config zap.Config, opts ...zap.Option) logger.Logger {
	logger, err := config.Build(opts...)
	if err != nil {
		panic(err)
	}
	return logger
}
