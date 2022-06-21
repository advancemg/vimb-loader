package app

import (
	"fmt"
	"github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/internal/usecase"
	"github.com/advancemg/vimb-loader/internal/usecase/repo/badger"
	"github.com/advancemg/vimb-loader/internal/usecase/repo/mongo"
	badgerClient "github.com/advancemg/vimb-loader/pkg/storage/badger-client"
	mongodbClient "github.com/advancemg/vimb-loader/pkg/storage/mongodb-client"
)

func Run(cfg *config.Configuration) {
	db_mongo, err := mongodbClient.New(cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.DB, cfg.Mongo.Password, cfg.Mongo.Password, cfg.Mongo.Debug)
	checkErr(err)
	usecase.New(mongo.New(db_mongo))

	db_badger := badgerClient.Open("database/table")
	repository := usecase.New(badger.New(db_badger))
	fmt.Println(repository)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
