package models

import (
	"bytes"
	"encoding/json"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
)

type MediaplanAggBadgerQuery struct {
}

func (query *MediaplanAggBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbAggMediaplans)
}

func (query *MediaplanAggBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *MediaplanAggBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return query.Find(result, filterNetworks)
}
