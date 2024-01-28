package logger

import (
	"go.uber.org/zap"
)

var (
	log   = initialize()
	Debug = log.Debugln
	Info  = log.Infoln
	Error = log.Errorln
)

func initialize() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
