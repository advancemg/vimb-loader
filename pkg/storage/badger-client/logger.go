package badger_client

import log "github.com/advancemg/vimb-loader/pkg/logging/logrus"

type loggingLevel int

const (
	DEBUG loggingLevel = iota
	INFO
	WARNING
	ERROR
)

type badgerLog struct {
	level loggingLevel
}

func defaultLogger(level loggingLevel) *badgerLog {
	return &badgerLog{level: level}
}

func (l *badgerLog) Errorf(f string, v ...interface{}) {
	if l.level <= ERROR {
		log.PrintLog("", "badger", "error", f, v)
	}
}

func (l *badgerLog) Warningf(f string, v ...interface{}) {
	if l.level <= WARNING {
		log.PrintLog("", "badger", "info", f, v)
	}
}

func (l *badgerLog) Infof(f string, v ...interface{}) {
	if l.level <= INFO {
		log.PrintLog("", "badger", "info", f, v)
	}
}

func (l *badgerLog) Debugf(f string, v ...interface{}) {
	if l.level <= DEBUG {
		log.PrintLog("", "badger", "debug", f, v)
	}
}
