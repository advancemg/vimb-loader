package badger_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/advancemg/badgerhold"
	log "github.com/advancemg/vimb-loader/pkg/logging/logrus"
	"github.com/dgraph-io/badger/v3"
	"os"
	"time"
)

var queryBadger = map[string]*badgerhold.Store{}

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

func New(storageDir string) (*badgerhold.Store, error) {
	open, err := Open(storageDir)
	if err != nil {
		return nil, fmt.Errorf("badger new: %w", err)
	}
	return open, nil
}

func Open(storageDir string) (*badgerhold.Store, error) {
	if queryBadger[storageDir] != nil {
		return queryBadger[storageDir], nil
	}
	var err error
	err = os.MkdirAll(storageDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("badger open: %w", err)
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
		return nil, fmt.Errorf("badger open: %w", err)
	}
	queryBadger[storageDir] = store
	return queryBadger[storageDir], nil
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
