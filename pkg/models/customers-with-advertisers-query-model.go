package models

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (query *CustomersWithAdvertisersBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request, true)
	err := query.Find(result, filterNetworks)
	if err != nil {
		var typeError *badgerhold.ErrTypeMismatch
		if errors.As(err, &typeError) {
			filterNetworks = HandleBadgerRequest(request, false)
			return query.Find(result, filterNetworks)
		} else {
			return err
		}
	}
	return err
}
