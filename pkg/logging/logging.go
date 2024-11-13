package logging

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger     *zap.Logger
	loggerOnce sync.Once
)

func Logger() *zap.Logger {
	loggerOnce.Do(InitLogger)
	return logger
}

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-readable time
	logger, _ = config.Build()
}

func SyncLogger() {
	if logger != nil {
		logger.Sync() // Flushes buffer, if any
	}
}
