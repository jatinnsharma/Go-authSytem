
package logger

import (
	"os"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	var err error
	
	if os.Getenv("ENVIRONMENT") == "production" {
		Logger, err = zap.NewProduction()
	} else {
		Logger, err = zap.NewDevelopment()
	}
	
	if err != nil {
		panic(err)
	}
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}