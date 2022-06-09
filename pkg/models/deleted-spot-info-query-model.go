package models

import (
	"github.com/advancemg/badgerhold"
	"github.com/advancemg/vimb-loader/pkg/storage/badger"
)

type DeletedSpotInfoBadgerQuery struct {
}

func (query *DeletedSpotInfoBadgerQuery) connect() *badgerhold.Store {
	return badger.Open(DbDeletedSpotInfo)
}

func (query *DeletedSpotInfoBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}
