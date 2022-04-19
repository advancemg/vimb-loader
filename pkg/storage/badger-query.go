package storage

import (
	"bytes"
	"encoding/json"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	"github.com/dgraph-io/badger"
	"github.com/timshannon/badgerhold"
	"os"
)

var queryBadger = map[string]*badgerhold.Store{}

var badgerLogger = &log.BadgerLog{}

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
	opts.MaxTableSize = 128 << 20 //128 Mb
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
