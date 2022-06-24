package models

import (
	"encoding/json"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/store"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type RanksLoadRequest struct{}

type RankQuery struct {
	ID struct {
		Eq int64 `json:"eq" example:"1"`
	} `json:"ID"`
}

func (request *RanksLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetRanksType
	amqpConfig := mq_broker.InitConfig()
	err := amqpConfig.DeclareSimpleQueue(qName)
	if err != nil {
		return nil, err
	}
	result := CommonResponse{}
	req := goConvert.New()
	err = amqpConfig.PublishJson(qName, req)
	if err != nil {
		log.PrintLog("vimb-loader", "Ranks InitTasks", "error", "Q:", qName, "err:", err.Error())
		return nil, err
	}
	result["status"] = "ok"
	return result, nil
}

func (request *Any) QueryRanks() ([]Ranks, error) {
	var result []Ranks
	db, table := utils.SplitDbAndTable(DbRanks)
	dbRanks := store.OpenDb(db, table)
	marshal, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}
	err = dbRanks.FindJson(&result, marshal)
	if err != nil {
		return nil, err
	}
	return result, nil
}
