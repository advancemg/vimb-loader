package models

import (
	"bytes"
	"encoding/json"
	"github.com/advancemg/badgerhold"
	"github.com/advancemg/vimb-loader/pkg/storage"
)

type RanksBadgerQuery struct {
}

func (query *RanksBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbRanks)
}

func (query *RanksBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *RanksBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return query.Find(result, filterNetworks)
}
