package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
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
	err := amqpConfig.DeclareSimpleQueue(qName)
	if err != nil {
		return nil, err
	}
	result := CommonResponse{}
	req := goConvert.New()
	req.Set("SellingDirectionID", request.SellingDirectionID)
	err = amqpConfig.PublishJson(qName, req)
	if err != nil {
		fmt.Printf("Q:%s - err:%s", qName, err.Error())
		return nil, err
	}
	result["status"] = "ok"
	return result, nil
}

func (request *Any) QueryCustomersWithAdvertisers() ([]CustomersWithAdvertisers, error) {
	var result []CustomersWithAdvertisers
	query := CustomersWithAdvertisersBadgerQuery{}
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

func (request *Any) QueryCustomersWithAdvertisersData() ([]CustomersWithAdvertisersData, error) {
	var result []CustomersWithAdvertisersData
	query := CustomersWithAdvertisersDataBadgerQuery{}
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
