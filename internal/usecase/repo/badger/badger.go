package badger

import (
	"bytes"
	"encoding/json"
	"github.com/advancemg/badgerhold"
	"strconv"
	"strings"
)

type DbRepo struct {
	*badgerhold.Store
}

func New(db *badgerhold.Store) *DbRepo {
	return &DbRepo{db}
}

func (r *DbRepo) Find(result interface{}, filter interface{}) error {
	return r.Find(result, filter)
}

func (r *DbRepo) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks := HandleBadgerRequest(request)
	return r.Find(result, filterNetworks)
}

func (r *DbRepo) AddOrUpdate(key interface{}, data interface{}) error {
	return r.Upsert(key, data)
}

func (r *DbRepo) Get(key interface{}, result interface{}) error {
	return r.Get(key, result)
}

func (r *DbRepo) Delete(key interface{}, dataType interface{}) error {
	return r.Delete(key, dataType)
}

func HandleBadgerRequest(request map[string]interface{}) *badgerhold.Query {
	var query *badgerhold.Query
	once := true
	for field, value := range request {
		for key, val := range value.(map[string]interface{}) {
			if once {
				query = switchBadgerFilterWhere(query, key, field, val)
				once = false
			} else {
				query = switchBadgerFilterAnd(query, key, field, val)
			}
		}
	}
	return query
}

func switchBadgerFilterAnd(filter *badgerhold.Query, key, filed string, value interface{}) *badgerhold.Query {
	value = jsonNumber(value)
	switch key {
	case "eq":
		filter = filter.And(filed).Eq(value)
	case "ne":
		filter = filter.And(filed).Ne(value)
	case "gt":
		filter = filter.And(filed).Gt(value)
	case "lt":
		filter = filter.And(filed).Lt(value)
	case "ge":
		filter = filter.And(filed).Ge(value)
	case "le":
		filter = filter.And(filed).Le(value)
	case "in":
		filter = filter.And(filed).In(value)
	case "isnil":
		filter = filter.And(filed).IsNil()
	}
	return filter
}

func switchBadgerFilterWhere(filter *badgerhold.Query, key, filed string, value interface{}) *badgerhold.Query {
	value = jsonNumber(value)
	switch key {
	case "eq":
		filter = badgerhold.Where(filed).Eq(value)
	case "ne":
		filter = badgerhold.Where(filed).Ne(value)
	case "gt":
		filter = badgerhold.Where(filed).Gt(value)
	case "lt":
		filter = badgerhold.Where(filed).Lt(value)
	case "ge":
		filter = badgerhold.Where(filed).Ge(value)
	case "le":
		filter = badgerhold.Where(filed).Le(value)
	case "in":
		filter = badgerhold.Where(filed).In(value)
	case "isnil":
		filter = badgerhold.Where(filed).IsNil()
	}
	return filter
}

func jsonNumber(value interface{}) interface{} {
	if number, ok := value.(json.Number); ok {
		strconv.ParseInt(string(number), 10, 64)
		dot := strings.Contains(number.String(), ".")
		if dot {
			i, err := number.Float64()
			if err != nil {
				panic(err)
			}
			return i
		} else {
			i, err := number.Int64()
			if err != nil {
				panic(err)
			}
			return i
		}
	}
	return value
}
