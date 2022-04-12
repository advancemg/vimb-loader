package models

import (
	"encoding/json"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
	"strconv"
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

type sellingDirectionID struct {
	SellingDirectionID string `json:"SellingDirectionID"`
}

func (query *ChannelBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var s sellingDirectionID
	err := json.Unmarshal(filter, &s)
	if err != nil {
		return err
	}
	intSellingDirectionID, err := strconv.Atoi(s.SellingDirectionID)
	if err != nil {
		return err
	}
	filterBudgets := badgerhold.Where("SellingDirectionID").Eq(intSellingDirectionID)
	return query.Find(result, filterBudgets)
}
