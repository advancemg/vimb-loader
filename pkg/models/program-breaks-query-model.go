package models

import (
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
)

type ProgramBreaksBadgerQuery struct {
}

func (query *ProgramBreaksBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbProgramBreaks)
}

func (query *ProgramBreaksBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}
