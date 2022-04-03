package models

import (
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
)

type ProgramBreaksLightBadgerQuery struct {
}

func (query *ProgramBreaksLightBadgerQuery) connect() *badgerhold.Store {
	return storage.Open(DbProgramBreaksLightMode)
}

func (query *ProgramBreaksLightBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}
