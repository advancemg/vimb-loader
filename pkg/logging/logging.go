package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func init() {
	LogFile := fmt.Sprintf("logs/%v", time.Now().Format(time.RFC3339))
	logFile, err := os.OpenFile(LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll("logs", 0777)
			if err != nil {
				log.Panic(err)
			}
			logFile, err = os.Create(LogFile)
			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Panic(err)
		}
	}
	log.SetOutput(logFile)
}

func printLog(app, client, level, format string, v ...interface{}) {
	setLogFormatter(format)
	setLogLevel(level)
	msg := map[string]interface{}{}
	msg["app"] = app
	msg["client"] = client
	msgs := fmt.Sprint(v...)
	switch level {
	case "info":
		log.WithFields(msg).Info(msgs)
	case "warn":
		log.WithFields(msg).Warn(msgs)
	case "debug":
		log.WithFields(msg).Debug(msgs)
	case "error":
		log.WithFields(msg).Error(msgs)
	default:
		log.WithFields(msg).Info(msgs)
	}
}

func PrintJsonLog(app, client, level string, v ...interface{}) {
	printLog(app, client, level, "json", v...)
}
func PrintLog(app, client, level string, v ...interface{}) {
	printLog(app, client, level, "text", v...)
}

type UTCFormatter struct {
	log.Formatter
}

func (u UTCFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}
func setLogFormatter(f string) {
	s := strings.ToLower(f)
	switch s {
	case "text":
		log.SetFormatter(UTCFormatter{&log.TextFormatter{
			ForceColors:               false,
			DisableColors:             true,
			ForceQuote:                false,
			DisableQuote:              false,
			EnvironmentOverrideColors: false,
			DisableTimestamp:          false,
			FullTimestamp:             true,
			TimestampFormat:           "2006-01-02T15:04:05Z07:00",
			DisableSorting:            true,
			SortingFunc:               nil,
			DisableLevelTruncation:    false,
			PadLevelText:              false,
			QuoteEmptyFields:          false,
			FieldMap:                  nil,
			CallerPrettyfier:          nil,
		}})
	case "json":
		log.SetFormatter(UTCFormatter{&log.JSONFormatter{}})
	default:
		log.SetFormatter(UTCFormatter{&log.JSONFormatter{}})
	}
}
func setLogLevel(s string) {
	log.SetLevel(getLogLevel(s))
}
func getLogLevel(str string) log.Level {
	s := strings.ToLower(str)
	switch s {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	default:
		return log.DebugLevel
	}
}

type BadgerLog struct {
}

func (l *BadgerLog) Errorf(f string, v ...interface{}) {
	PrintLog("", "", "BADGER ERROR: "+f, v)
}

func (l *BadgerLog) Warningf(f string, v ...interface{}) {
	PrintLog("", "", "BADGER WARNING: "+f, v)
}

func (l *BadgerLog) Infof(f string, v ...interface{}) {
	PrintLog("", "", "BADGER INFO: "+f, v)
}

func (l *BadgerLog) Debugf(f string, v ...interface{}) {
	PrintLog("", "", "BADGER DEBUG: "+f, v)
}
