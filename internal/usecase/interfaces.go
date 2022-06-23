package usecase

import "time"

type DbInterface interface {
	FindJson(result interface{}, filter []byte) error
	AddOrUpdate(key interface{}, data interface{}) error
	AddWithTTL(key, value interface{}, ttl time.Duration) error
	Get(key interface{}, result interface{}) error
	Delete(key interface{}, dataType interface{}) error
	FindWhereEq(result interface{}, filed string, value interface{}) error
	FindWhereNe(result interface{}, filed string, value interface{}) error
	FindWhereGt(result interface{}, filed string, value interface{}) error
	FindWhereLt(result interface{}, filed string, value interface{}) error
	FindWhereGe(result interface{}, filed string, value interface{}) error
	FindWhereLe(result interface{}, filed string, value interface{}) error
	FindWhereAnd2Eq(result interface{}, filed1 string, value1 interface{}, filed2 string, value2 interface{}) error
	FindWhereAnd4Eq(result interface{}, filed1 string, value1 interface{}, filed2 string, value2 interface{}, filed3 string, value3 interface{}, filed4 string, value4 interface{}) error
}
