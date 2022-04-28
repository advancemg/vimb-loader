package models

import (
	"bytes"
	"encoding/json"
	"github.com/advancemg/badgerhold"
	"github.com/advancemg/vimb-loader/pkg/storage"
)

type SpotBadgerQuery struct {
}

func (query *SpotBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbSpots)
}

func (query *SpotBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *SpotBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return query.Find(result, filterNetworks)
}

type SpotsOrderBlockQuery struct {
}

func (query *SpotsOrderBlockQuery) connect() *badgerhold.Store {
	return storage.Open(DbSpotsOrderBlock)
}

func (query *SpotsOrderBlockQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *SpotsOrderBlockQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return query.Find(result, filterNetworks)
}
