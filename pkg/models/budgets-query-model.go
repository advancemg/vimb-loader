package models

import (
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
