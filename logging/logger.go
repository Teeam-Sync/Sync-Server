package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log   = mustInitialize()
	Debug = log.Debugln
	Info  = log.Infoln
	Error = log.Errorln
)

func mustInitialize() *zap.SugaredLogger {
	var logger *zap.Logger
	var err error

	env := os.Getenv("APP_ENV")
	switch env {
	case "prod":
		logger, err = newLogger(zapcore.InfoLevel)
	case "dev":
		fallthrough
	default:
		logger, err = newLogger(zapcore.DebugLevel)
	}
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

func newLogger(level zapcore.Level) (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
