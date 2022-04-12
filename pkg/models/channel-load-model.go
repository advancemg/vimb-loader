package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
)

type ChannelLoadRequest struct {
	SellingDirectionID string `json:"SellingDirectionID" example:"23"`
}

func (request *ChannelLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetChannelsType
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

func (request *ChannelLoadRequest) LoadChannels() ([]Channel, error) {
	var result []Channel
	query := ChannelBadgerQuery{}
	marshal, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	query.FindJson(&result, marshal)
	return result, nil
}
