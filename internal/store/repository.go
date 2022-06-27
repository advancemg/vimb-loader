package store

import (
	"errors"
	"fmt"
	cfg "github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/internal/store/repo/badger"
	"github.com/advancemg/vimb-loader/internal/store/repo/mongo"
	badger_client "github.com/advancemg/vimb-loader/pkg/storage/badger-client"
	mongodb_client "github.com/advancemg/vimb-loader/pkg/storage/mongodb-client"
	"strings"
	"time"
)

const mongodb = "mongodb"

var ErrNotFound = errors.New("No data found for this key")

type Repository struct {
	repo DbInterface
}

func (r *Repository) FindWhereEq(result interface{}, filed string, value interface{}) error {
	err := r.repo.FindWhereEq(result, filed, value)
	if err != nil {
		return fmt.Errorf("repository - FindWhereEq - r.repo.FindWhereEq: %w", err)
	}
	return nil
}

func (r *Repository) FindWhereAnd2Eq(result interface{}, filed1 string, value1 interface{}, filed2 string, value2 interface{}) error {
	err := r.repo.FindWhereAnd2Eq(result, filed1, value1, filed2, value2)
	if err != nil {
		return fmt.Errorf("repository - FindWhereAnd2Eq - r.repo.FindWhereAnd2Eq: %w", err)
	}
	return nil
}

func (r *Repository) FindWhereAnd4Eq(result interface{}, filed1 string, value1 interface{}, filed2 string, value2 interface{}, filed3 string, value3 interface{}, filed4 string, value4 interface{}) error {
	err := r.repo.FindWhereAnd4Eq(result, filed1, value1, filed2, value2, filed3, value3, filed4, value4)
	if err != nil {
		return fmt.Errorf("repository - FindWhereAnd4Eq - r.repo.FindWhereAnd4Eq: %w", err)
	}
	return nil
}

func (r *Repository) FindWhereNe(result interface{}, filed string, value interface{}) error {
	err := r.repo.FindWhereNe(result, filed, value)
	if err != nil {
		return fmt.Errorf("repository - FindWhereNe - r.repo.FindWhereNe: %w", err)
	}
	return nil
}

func (r *Repository) FindWhereGt(result interface{}, filed string, value interface{}) error {
	err := r.repo.FindWhereGt(result, filed, value)
	if err != nil {
		return fmt.Errorf("repository - FindWhereGt - r.repo.FindWhereGt: %w", err)
	}
	return nil
}

func (r *Repository) FindWhereLt(result interface{}, filed string, value interface{}) error {
	err := r.repo.FindWhereLt(result, filed, value)
	if err != nil {
		return fmt.Errorf("repository - FindWhereLt - r.repo.FindWhereLt: %w", err)
	}
	return nil
}

func (r *Repository) FindWhereGe(result interface{}, filed string, value interface{}) error {
	err := r.repo.FindWhereGe(result, filed, value)
	if err != nil {
		return fmt.Errorf("repository - FindWhereGe - r.repo.FindWhereGe: %w", err)
	}
	return nil
}

func (r *Repository) FindWhereLe(result interface{}, filed string, value interface{}) error {
	err := r.repo.FindWhereLe(result, filed, value)
	if err != nil {
		return fmt.Errorf("repository - FindWhereLe - r.repo.FindWhereLe: %w", err)
	}
	return nil
}

func (r *Repository) FindJson(result interface{}, filter []byte) error {
	err := r.repo.FindJson(result, filter)
	if err != nil {
		return fmt.Errorf("repository - FindJson - r.repo.FindJson: %w", err)
	}
	return nil
}

func (r *Repository) AddOrUpdate(key interface{}, data interface{}) error {
	err := r.repo.AddOrUpdate(key, data)
	if err != nil {
		return fmt.Errorf("repository - AddOrUpdate - r.repo.AddOrUpdate: %w", err)
	}
	return nil
}

func (r *Repository) AddWithTTL(key, value interface{}, ttl time.Duration) error {
	err := r.repo.AddWithTTL(key, value, ttl)
	if err != nil {
		return fmt.Errorf("repository - AddWithTTL - r.repo.AddWithTTL: %w", err)
	}
	return nil
}

func (r *Repository) Get(key interface{}, result interface{}) error {
	err := r.repo.Get(key, result)
	if err != nil {
		if strings.Contains(err.Error(), "No data found for this key") || strings.Contains(err.Error(), "mongo: no documents in result") {
			return ErrNotFound
		}
		return fmt.Errorf("repository - Get - r.repo.Get: %w", err)
	}
	return nil
}

func (r *Repository) Delete(key interface{}, dataType interface{}) error {
	err := r.repo.Delete(key, dataType)
	if err != nil {
		return fmt.Errorf("repository - Delete - r.repo.Delete: %w", err)
	}
	return nil
}

func New(r DbInterface) *Repository {
	return &Repository{repo: r}
}

func OpenDb(database, table string) *Repository {
	var repository *Repository
	switch cfg.Config.Database {
	case mongodb:
		cfgMongo := mongodb_client.InitConfig()
		cfgMongo.DB = database
		mongoClient, err := cfgMongo.New()
		if err != nil {
			panic(err)
		}
		db := mongo.New(mongoClient, table, database)
		repository = New(db)
	default:
		db := badger.New(badger_client.Open(database + "/" + table))
		repository = New(db)
	}
	return repository
}
