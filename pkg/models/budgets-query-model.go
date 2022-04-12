package models

import (
	"encoding/json"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
)

type BudgetsBadgerQuery struct {
}

func (query *BudgetsBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbBudgets)
}

func (query *BudgetsBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

type months struct {
	Month int `json:"Month"`
}

func (query *BudgetsBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var m months
	err := json.Unmarshal(filter, &m)
	if err != nil {
		return err
	}
	filterBudgets := badgerhold.Where("Month").Eq(m.Month)
	return query.Find(result, filterBudgets)
}
