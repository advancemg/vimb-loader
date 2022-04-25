package storage

import (
	"bytes"
	"encoding/json"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	"github.com/dgraph-io/badger"
	"github.com/timshannon/badgerhold"
	"os"
	"time"
)

var queryBadger = map[string]*badgerhold.Store{}

type BadgerLog struct {
}

func (l *BadgerLog) Errorf(f string, v ...interface{}) {
	log.PrintLog("", "BADGER", "BADGER ERROR: "+f, v)
}

func (l *BadgerLog) Warningf(f string, v ...interface{}) {
	log.PrintLog("", "BADGER", "BADGER WARNING: "+f, v)
}

func (l *BadgerLog) Infof(f string, v ...interface{}) {
	log.PrintLog("", "BADGER", "BADGER INFO: "+f, v)
}

func (l *BadgerLog) Debugf(f string, v ...interface{}) {
	log.PrintLog("", "BADGER", "BADGER DEBUG: "+f, v)
}

var badgerLogger = &BadgerLog{}

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
	opts.MaxTableSize = 256 << 20
	opts.NumVersionsToKeep = 0
	opts.NumVersionsToKeep = 1
	opts.ValueLogFileSize = 1<<30 - 1
	opts.Logger = badgerLogger
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
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		for dbName, db := range queryBadger {
		again:
			err := db.Badger().RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
			log.PrintLog("vimb-loader", "badger", "info", "Clean BadgerGC:"+dbName)
		}
	}
}
