package models

import (
	"bytes"
	"encoding/json"
	"github.com/advancemg/badgerhold"
	"github.com/advancemg/vimb-loader/pkg/storage"
)

type AdvertiserBadgerQuery struct {
}

func (query *AdvertiserBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbAdvertisers)
}

func (query *AdvertiserBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *AdvertiserBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return query.Find(result, filterNetworks)
}
