package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type SpotsLoadRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
	StartDate          string `json:"StartDate"`
	EndDate            string `json:"EndDate"`
	InclOrdBlocks      string `json:"InclOrdBlocks"`
	ChannelList        []struct {
		Cnl  string `json:"Cnl"`
		Main string `json:"Main"`
	} `json:"ChannelList"`
	AdtList []struct {
		AdtID string `json:"AdtID"`
	} `json:"AdtList"`
}

func (request *SpotsLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetSpotsType
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
		req.Set("SellingDirectionID", request.SellingDirectionID)
		req.Set("StartDate", day.Format(`20060102`))
		req.Set("EndDate", day.Format(`20060102`))
		req.Set("InclOrdBlocks", request.InclOrdBlocks)
		req.Set("ChannelList", request.ChannelList)
		req.Set("AdtList", request.AdtList)
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

func (request *SpotsLoadRequest) getDays() ([]time.Time, error) {
	return utils.GetDaysByPeriod(request.StartDate, request.EndDate)
}
