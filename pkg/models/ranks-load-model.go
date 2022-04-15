package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
)

type RanksLoadRequest struct{}

type RankQuery struct {
	ID struct {
		Eq int `json:"eq" example:"1"`
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
		fmt.Printf("Q:%s - err:%s", qName, err.Error())
		return nil, err
	}
	result["status"] = "ok"
	return result, nil
}

func (request *Any) QueryRanks() ([]Ranks, error) {
	var result []Ranks
	query := RanksBadgerQuery{}
	marshal, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}
	err = query.FindJson(&result, marshal)
	if err != nil {
		return nil, err
	}
	return result, nil
}
