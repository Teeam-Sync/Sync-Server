package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	log    *zap.SugaredLogger
	Debug  func(args ...interface{})
	Info   func(args ...interface{})
	Error  func(args ...interface{})
)

func init() {
	mustInitialize()
}

func mustInitialize() {
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

	log = logger.Sugar()
	Debug = log.Debug
	Info = log.Info
	Error = log.Error
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
