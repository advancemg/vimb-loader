package models

import (
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
)

type ChannelBadgerQuery struct {
}

func (query *ChannelBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbChannels)
}

func (query *ChannelBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}
