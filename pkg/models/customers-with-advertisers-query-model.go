package models

import (
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
)

type CustomersWithAdvertisersBadgerQuery struct {
}

func (query *CustomersWithAdvertisersBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbCustomersWithAdvertisers)
}

func (query *CustomersWithAdvertisersBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}