package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
)

type CustomersWithAdvertisersLoadRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
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
