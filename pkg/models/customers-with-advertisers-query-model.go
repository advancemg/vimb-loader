package models

import (
	"bytes"
	"encoding/json"
	"github.com/advancemg/badgerhold"
	"github.com/advancemg/vimb-loader/pkg/storage/badger"
)

type CustomersWithAdvertisersBadgerQuery struct {
}

func (query *CustomersWithAdvertisersBadgerQuery) connect() *badgerhold.Store {
	return badger.Open(DbCustomersWithAdvertisers)
}

func (query *CustomersWithAdvertisersBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *CustomersWithAdvertisersBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return query.Find(result, filterNetworks)
}

type CustomersWithAdvertisersDataBadgerQuery struct {
}

func (query *CustomersWithAdvertisersDataBadgerQuery) connect() *badgerhold.Store {
	return badger.Open(DbCustomersWithAdvertisersData)
}

func (query *CustomersWithAdvertisersDataBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *CustomersWithAdvertisersDataBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return query.Find(result, filterNetworks)
}
