package models

import (
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
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
