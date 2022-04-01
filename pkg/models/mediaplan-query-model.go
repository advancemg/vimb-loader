package models

import (
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
