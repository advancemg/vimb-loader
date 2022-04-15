package models

import (
	"encoding/json"
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

type SpotsQuery struct {
	Rating struct {
		Eq float64 `json:"eq" example:"0.58321"`
	} `json:"Rating"`
	SpotID struct {
		Eq int64 `json:"eq" example:"451118797"`
	} `json:"SpotID"`
}

type QuerySpotsOrderBlockQuery struct {
	OrdID struct {
		Eq int `json:"eq" example:"319260"`
	} `json:"OrdID"`
	BlockID struct {
		Eq int64 `json:"eq" example:"451118797"`
	} `json:"BlockID"`
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

func (request *Any) QuerySpots() ([]Spot, error) {
	var result []Spot
	query := SpotBadgerQuery{}
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

func (request *Any) QuerySpotsOrderBlock() ([]SpotOrderBlock, error) {
	var result []SpotOrderBlock
	query := SpotsOrderBlockQuery{}
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
