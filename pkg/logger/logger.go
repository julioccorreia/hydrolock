package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(mode string) (*zap.Logger, error) {
	var config zap.Config

	if mode == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = "ts"
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func DefaultLogger() *zap.Logger {
	return zap.NewExample()
}
