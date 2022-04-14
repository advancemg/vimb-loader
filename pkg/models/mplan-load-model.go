package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type MediaplanLoadRequest struct {
	StartMonth         string `json:"StartMonth" example:"201903"`
	EndMonth           string `json:"EndMonth" example:"201912"`
	SellingDirectionID string `json:"SellingDirectionID" example:"23"`
}

func (request *MediaplanLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetMPLansType
	amqpConfig := mq_broker.InitConfig()
	err := amqpConfig.DeclareSimpleQueue(qName)
	if err != nil {
		return nil, err
	}
	months, err := request.getMonths()
	if err != nil {
		return nil, err
	}
	result := CommonResponse{}
	for _, month := range months {
		req := goConvert.New()
		req.Set("SellingDirectionID", request.SellingDirectionID)
		req.Set("StartMonth", month.ValueString)
		req.Set("EndMonth", month.ValueString)
		err = amqpConfig.PublishJson(qName, req)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return nil, err
		}
	}
	result["months"] = months
	result["status"] = "ok"
	return result, nil
}

func (request *MediaplanLoadRequest) getMonths() ([]utils.YearMonth, error) {
	return utils.GetPeriodFromYearMonths(request.StartMonth, request.EndMonth)
}

type MediaplanQuery struct {
	AdtID struct {
		Eq int `json:"eq" example:"700068653"`
	} `json:"AdtID"`
	AgrID struct {
		Eq int `json:"eq" example:"81024"`
	} `json:"AgrID"`
	MplID struct {
		Eq int `json:"eq" example:"14824608"`
	} `json:"MplID"`
}

func (request *Any) QueryMediaplans() ([]Mediaplan, error) {
	var result []Mediaplan
	query := MediaplanBadgerQuery{}
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
