package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type AdvMessagesLoadRequest struct {
	CreationDateStart string `json:"CreationDateStart" example:"20190301"`
	CreationDateEnd   string `json:"CreationDateEnd" example:"20191201"`
	FillMaterialTags  string `json:"FillMaterialTags" example:"true"`
	Advertisers       []struct {
		Id string `json:"ID"`
	} `json:"Advertisers"`
	Aspects []struct {
		Id string `json:"ID"`
	} `json:"Aspects"`
	AdvertisingMessageIDs []struct {
		Id string `json:"ID"`
	} `json:"AdvertisingMessageIDs"`
}

type AdvMessageQuery struct {
	AdtID struct {
		Eq int `json:"eq" example:"700061957"`
	} `json:"AdtID"`
	BrandID struct {
		Eq int `json:"eq" example:"44362"`
	} `json:"BrandID"`
}

func (request *AdvMessagesLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetAdvMessagesType
	amqpConfig := mq_broker.InitConfig()
	err := amqpConfig.DeclareSimpleQueue(qName)
	if err != nil {
		return nil, err
	}
	days, err := request.getDays()
	if err != nil {
		return nil, err
	}
	result := CommonResponse{}
	for _, day := range days {
		req := goConvert.New()
		req.Set("CreationDateStart", day.Format(`20060102`))
		req.Set("CreationDateEnd", day.Format(`20060102`))
		if request.FillMaterialTags != "" {
			req.Set("FillMaterialTags", request.FillMaterialTags)
		} else {
			req.Set("FillMaterialTags", "true")
		}
		req.Set("Advertisers", []struct{}{})
		req.Set("Aspects", []struct{}{})
		req.Set("AdvertisingMessageIDs", []struct{}{})
		err := amqpConfig.PublishJson(qName, req)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return nil, err
		}
	}
	result["days"] = days
	result["status"] = "ok"
	return result, nil
}

func (request *AdvMessagesLoadRequest) getDays() ([]time.Time, error) {
	return utils.GetDaysByPeriod(request.CreationDateStart, request.CreationDateEnd)
}

func (request *Any) QueryAdvMessages() ([]Advertiser, error) {
	var result []Advertiser
	query := AdvertiserBadgerQuery{}
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
