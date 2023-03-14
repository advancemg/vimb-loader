package badger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/advancemg/badgerhold"
	"strconv"
	"strings"
	"time"
)

type DbRepo struct {
	*badgerhold.Store
}

func New(db *badgerhold.Store) *DbRepo {
	return &DbRepo{db}
}

func (r *DbRepo) FindWhereEq(result interface{}, filed string, value interface{}) error {
	var filter *badgerhold.Query
	filter = badgerhold.Where(filed).Eq(value)
	return r.Store.Find(result, filter)
}

func (r *DbRepo) FindWhereAnd2Eq(result interface{}, filed1 string, value1 interface{}, filed2 string, value2 interface{}) error {
	var filter *badgerhold.Query
	filter = badgerhold.Where(filed1).Eq(value1).And(filed2).Eq(value2)
	return r.Store.Find(result, filter)
}

func (r *DbRepo) FindWhereAnd4Eq(result interface{}, filed1 string, value1 interface{}, filed2 string, value2 interface{}, filed3 string, value3 interface{}, filed4 string, value4 interface{}) error {
	var filter *badgerhold.Query
	filter = badgerhold.Where(filed1).Eq(value1).
		And(filed2).Eq(value2).
		And(filed3).Eq(value3).
		And(filed4).Eq(value4)
	return r.Store.Find(result, filter)
}

func (r *DbRepo) FindWhereNe(result interface{}, filed string, value interface{}) error {
	var filter *badgerhold.Query
	filter = badgerhold.Where(filed).Ne(value)
	return r.Store.Find(result, filter)
}

func (r *DbRepo) FindWhereGt(result interface{}, filed string, value interface{}) error {
	var filter *badgerhold.Query
	filter = badgerhold.Where(filed).Gt(value)
	return r.Store.Find(result, filter)
}
func (r *DbRepo) FindWhereLt(result interface{}, filed string, value interface{}) error {
	var filter *badgerhold.Query
	filter = badgerhold.Where(filed).Lt(value)
	return r.Store.Find(result, filter)
}
func (r *DbRepo) FindWhereGe(result interface{}, filed string, value interface{}) error {
	var filter *badgerhold.Query
	filter = badgerhold.Where(filed).Ge(value)
	return r.Store.Find(result, filter)
}
func (r *DbRepo) FindWhereLe(result interface{}, filed string, value interface{}) error {
	var filter *badgerhold.Query
	filter = badgerhold.Where(filed).Le(value)
	return r.Store.Find(result, filter)
}
func (r *DbRepo) FindWhereIn(result interface{}, filed string, value interface{}) error {
	var filter *badgerhold.Query
	filter = badgerhold.Where(filed).In(value)
	return r.Store.Find(result, filter)
}

func (r *DbRepo) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return fmt.Errorf("FindJson Decode: %w", err)
	}
	filterNetworks, err := HandleBadgerRequest(request)
	if err != nil {
		return fmt.Errorf("FindJson: %w", err)
	}
	return r.Store.Find(result, filterNetworks)
}

func (r *DbRepo) AddOrUpdate(key interface{}, data interface{}) error {
	return r.Upsert(key, data)
}

func (r *DbRepo) AddWithTTL(key, value interface{}, ttl time.Duration) error {
	return r.UpsertTTL(key, value, ttl)
}

func (r *DbRepo) Get(key interface{}, result interface{}) error {
	return r.Store.Get(key, result)
}

func (r *DbRepo) Delete(key interface{}, dataType interface{}) error {
	return r.Store.Delete(key, dataType)
}

func HandleBadgerRequest(request map[string]interface{}) (*badgerhold.Query, error) {
	var query *badgerhold.Query
	once := true
	for field, value := range request {
		for key, val := range value.(map[string]interface{}) {
			if once {
				where, err := switchBadgerFilterWhere(query, key, field, val)
				query = where
				if err != nil {
					return nil, fmt.Errorf("HandleBadgerRequest: %w", err)
				}
				once = false
			} else {
				and, err := switchBadgerFilterAnd(query, key, field, val)
				query = and
				if err != nil {
					return nil, fmt.Errorf("HandleBadgerRequest: %w", err)
				}
			}
		}
	}
	return query, nil
}

func switchBadgerFilterAnd(filter *badgerhold.Query, key, filed string, value interface{}) (*badgerhold.Query, error) {
	value, err := jsonNumber(value)
	if err != nil {
		return nil, fmt.Errorf("switchBadgerFilterAnd: %w", err)
	}
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
	return filter, nil
}

func switchBadgerFilterWhere(filter *badgerhold.Query, key, filed string, value interface{}) (*badgerhold.Query, error) {
	value, err := jsonNumber(value)
	if err != nil {
		return nil, fmt.Errorf("switchBadgerFilterWhere: %w", err)
	}
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
	return filter, nil
}

func jsonNumber(value interface{}) (interface{}, error) {
	if number, ok := value.(json.Number); ok {
		strconv.ParseInt(string(number), 10, 64)
		dot := strings.Contains(number.String(), ".")
		if dot {
			i, err := number.Float64()
			if err != nil {
				return nil, fmt.Errorf("jsonNumber: %w", err)
			}
			return i, err
		} else {
			i, err := number.Int64()
			if err != nil {
				return nil, fmt.Errorf("jsonNumber: %w", err)
			}
			return i, err
		}
	}
	return value, nil
}
