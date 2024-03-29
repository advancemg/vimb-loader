package models

import (
	"encoding/json"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/store"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type CustomersWithAdvertisersLoadRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
}

type CustomersWithAdvertiserQuery struct {
	ID struct {
		Eq int64 `json:"eq" example:"1"`
	} `json:"ID"`
}

type CustomersWithAdvertiserDataQuery struct {
	CustID struct {
		Eq int64 `json:"eq" example:"1"`
	} `json:"CustID"`
}

func (request *CustomersWithAdvertisersLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetCustomersWithAdvertisersType
	amqpConfig := mq_broker.InitConfig()
	defer amqpConfig.Close()
	err := amqpConfig.DeclareSimpleQueue(qName)
	if err != nil {
		return nil, err
	}
	result := CommonResponse{}
	req := goConvert.New()
	req.Set("SellingDirectionID", request.SellingDirectionID)
	err = amqpConfig.PublishJson(qName, req)
	if err != nil {
		log.PrintLog("vimb-loader", "CustomersWithAdvertisers InitTasks", "error", "Q:", qName, "err:", err.Error())
		return nil, err
	}
	result["status"] = "ok"
	return result, nil
}

func (request *Any) QueryCustomersWithAdvertisers() ([]CustomersWithAdvertisers, error) {
	var result []CustomersWithAdvertisers
	db, table := utils.SplitDbAndTable(DbCustomersWithAdvertisers)
	repo, err := store.OpenDb(db, table)
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}
	err = repo.FindJson(&result, marshal)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (request *Any) QueryCustomersWithAdvertisersData() ([]CustomersWithAdvertisersData, error) {
	var result []CustomersWithAdvertisersData
	db, table := utils.SplitDbAndTable(DbCustomersWithAdvertisersData)
	repo, err := store.OpenDb(db, table)
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}
	err = repo.FindJson(&result, marshal)
	if err != nil {
		return nil, err
	}
	return result, nil
}
