package zap

import (
	"fmt"
	log "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Logger *log.Logger
var LogWriter *os.File

func Init() error {
	filePath := fmt.Sprintf("logs/%v", time.Now().Format(time.RFC3339))
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll("logs", 0777)
			if err != nil {
				return err
			}
			logFile, err = os.Create(filePath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	config := log.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	Logger = log.New(core, log.AddCaller(), log.AddStacktrace(zapcore.ErrorLevel))
	LogWriter = logFile
	//logFile.Close()
	return nil
}

func PrintLog(app, client, level string, v ...interface{}) {
	printLog(app, client, level, v...)
}

func printLog(app, client, level string, v ...interface{}) {
	msg := fmt.Sprint(v...)
	switch level {
	case "info":
		Logger.WithOptions(log.AddCallerSkip(2)).Info("", log.String("app", app), log.String("client", client), log.String("message:", msg))
	case "warn":
		Logger.WithOptions(log.AddCallerSkip(2)).Warn("", log.String("app", app), log.String("client", client), log.String("message:", msg))
	case "debug":
		Logger.WithOptions(log.AddCallerSkip(2)).Debug("", log.String("app", app), log.String("client", client), log.String("message:", msg))
	case "error":
		Logger.WithOptions(log.AddCallerSkip(2)).Error("", log.String("app", app), log.String("client", client), log.String("message:", msg))
	default:
		Logger.WithOptions(log.AddCallerSkip(2)).Info("", log.String("app", app), log.String("client", client), log.String("message:", msg))
	}
}
