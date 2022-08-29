package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	GetLogger() *zap.Logger
}

type LoggerImpl struct {
	logger *zap.Logger
}

func NewLog(appName string, appEnvironment string, logPath string) Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	writer := zapcore.AddSync(logFile)

	core := ecszap.NewCore(encoderConfig, writer, zap.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	logger = logger.With(zap.String("app", appName)).With(zap.String("environment", appEnvironment))

	return &LoggerImpl{
		logger: logger,
	}
}

func (logger LoggerImpl) GetLogger() *zap.Logger {
	return logger.logger
}
