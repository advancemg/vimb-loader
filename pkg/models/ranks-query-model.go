package models

import (
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
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
