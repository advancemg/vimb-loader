package models

import (
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
)

type DeletedSpotInfoBadgerQuery struct {
}

func (query *DeletedSpotInfoBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbDeletedSpotInfo)
}

func (query *DeletedSpotInfoBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}
