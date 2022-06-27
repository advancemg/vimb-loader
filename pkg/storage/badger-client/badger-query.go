package badger_client

import (
	"bytes"
	"encoding/json"
	"github.com/advancemg/badgerhold"
	log "github.com/advancemg/vimb-loader/pkg/logging/logrus"
	"github.com/dgraph-io/badger/v3"
	"os"
	"time"
)

var queryBadger = map[string]*badgerhold.Store{}

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

func DefaultEncode(value interface{}) ([]byte, error) {
	var buff bytes.Buffer
	en := json.NewEncoder(&buff)
	err := en.Encode(value)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func DefaultDecode(data []byte, value interface{}) error {
	var buff bytes.Buffer
	de := json.NewDecoder(&buff)
	_, err := buff.Write(data)
	if err != nil {
		return err
	}
	return de.Decode(value)
}

type Badger struct {
	Store *badgerhold.Store
}

func New(storageDir string) *badgerhold.Store {
	return Open(storageDir)
}

func Open(storageDir string) *badgerhold.Store {
	if queryBadger[storageDir] != nil {
		return queryBadger[storageDir]
	}
	var err error
	err = os.MkdirAll(storageDir, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	opts := badger.DefaultOptions(storageDir)
	opts.SyncWrites = true
	opts.Dir = storageDir
	opts.ValueDir = storageDir
	opts.MemTableSize = 256 << 20
	opts.Logger = defaultLogger(INFO)
	options := badgerhold.Options{
		Encoder:          DefaultEncode,
		Decoder:          DefaultDecode,
		SequenceBandwith: 1,
		Options:          opts,
	}
	store, err := badgerhold.Open(options)
	if err != nil {
		return nil
	}
	queryBadger[storageDir] = store
	return queryBadger[storageDir]
}

func CleanGC() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		for dbName, db := range queryBadger {
		again:
			err := db.Badger().RunValueLogGC(0.5)
			if err == nil {
				goto again
			}
			log.PrintLog("vimb-loader", "badger", "info", "Clean BadgerGC:"+dbName)
		}
	}
}
