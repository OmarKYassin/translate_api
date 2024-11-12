package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-readable time
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger
}

func SyncLogger() {
	if Logger != nil {
		Logger.Sync() // Flushes buffer, if any
	}
}
