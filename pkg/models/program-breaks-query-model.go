package models

import (
	"bytes"
	"encoding/json"
	"github.com/advancemg/badgerhold"
	"github.com/advancemg/vimb-loader/pkg/storage/badger"
)

type ProgramBreaksBadgerQuery struct {
}

func (query *ProgramBreaksBadgerQuery) connect() *badgerhold.Store {
	return badger.Open(DbProgramBreaks)
}

func (query *ProgramBreaksBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *ProgramBreaksBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return query.Find(result, filterNetworks)
}

type ProgramBreaksProMasterBadgerQuery struct {
}

func (query *ProgramBreaksProMasterBadgerQuery) connect() *badgerhold.Store {
	return badger.Open(DbProgramBreaksProMaster)
}

func (query *ProgramBreaksProMasterBadgerQuery) Find(result interface{}, filter *badgerhold.Query) error {
	store := query.connect()
	return store.Find(result, filter)
}

func (query *ProgramBreaksProMasterBadgerQuery) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return query.Find(result, filterNetworks)
}
