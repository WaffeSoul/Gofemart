package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger(isDebug bool) *zap.Logger {
	var cfg zap.Config

	if isDebug {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	var err error
	logger, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger")
	}

	logger.Info("Logger initialized", zap.Bool("isDebug", isDebug))
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
