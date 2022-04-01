package storage

import (
	"bytes"
	"encoding/json"
	"github.com/dgraph-io/badger"
	"github.com/timshannon/badgerhold"
	"os"
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
