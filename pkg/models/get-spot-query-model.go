package models

import (
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
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
