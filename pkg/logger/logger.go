package logger

import "go.uber.org/zap"

type Logger interface {
	Info(name string, fields ...zap.Field)
	Debug(name string, fields ...zap.Field)
	Error(name string, fields ...zap.Field)
	Panic(name string, fields ...zap.Field)
}

func newZapLogger() Logger {
	logger, _ := zap.NewProduction()
	return logger
}

var DefaultLogger = newZapLogger()

func SetDefaultLogger(logger Logger) {
	DefaultLogger = logger
}

func Info(name string, fields ...zap.Field) {
	DefaultLogger.Info(name, fields...)
}

func Debug(name string, fields ...zap.Field) {
	DefaultLogger.Debug(name, fields...)
}

func Error(name string, fields ...zap.Field) {
	DefaultLogger.Error(name, fields...)
}

func Panic(name string, fields ...zap.Field) {
	DefaultLogger.Panic(name, fields...)
}
