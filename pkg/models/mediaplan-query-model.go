package models

import (
	"encoding/json"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
)

type MediaplanBadgerQuery struct {
}

func (query *MediaplanBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbMediaplans)
}

func (query *MediaplanBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *MediaplanBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var m MediaplanLoadBadgerRequest
	err := json.Unmarshal(filter, &m)
	if err != nil {
		return err
	}
	filterMediaplan := badgerhold.Where("Month").Eq(m.Month)
	return query.Find(result, filterMediaplan)
}
