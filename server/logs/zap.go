package logs

import (
	"go.uber.org/zap"
)

/**
This is the logger client that we are working with, initialises the logger with custom configurations so that we can add middlewares if we need into the logger.
*/

var logger *zap.SugaredLogger

func init() {
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"stdout"}
	prodLogger, _ := conf.Build(zap.AddCallerSkip(1))
	logger = prodLogger.Sugar()
}

func Info(log string, v ...any) {
	logger.Infof(log, v...)
}

func Warn(log string, v ...any) {
	logger.Warnf(log, v...)
}

func Error(log string, v ...any) {
	logger.Warnf(log, v...)
}

func Fatal(log string, v ...any) {
	logger.Fatalf(log, v...)
}
