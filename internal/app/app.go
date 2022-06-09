package app

import (
	"github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/internal/usecase"
	"github.com/advancemg/vimb-loader/internal/usecase/repo/mongo"
	"github.com/advancemg/vimb-loader/pkg/storage/mongodb"
)

func Run(cfg *config.Configuration) {
	db, err := mongodb.New(cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.DB, cfg.Mongo.Password, cfg.Mongo.Password, cfg.Mongo.Debug)
	checkErr(err)
	usecase.New(mongo.New(db))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
