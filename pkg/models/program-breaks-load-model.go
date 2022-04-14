package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type ProgramBreaksLoadRequest struct {
	SellingDirectionID string `json:"SellingDirectionID" example:"21"` //ID направления продаж
	InclProgAttr       string `json:"InclProgAttr" example:"1"`        //Флаг "Заполнять секцию ProMaster". 1 - да, 0 - нет. (int, not nillable)
	InclForecast       string `json:"InclForecast" example:"1"`        //Признак "Как заполнять секцию прогнозных рейтингов". 0 - Не заполнять,  1 - Заполнять только ЦА программатика, 2 - Заполнять всеми возможными ЦА
	StartDate          string `json:"StartDate" example:"20170701"`    //Дата начала периода в формате YYYYMMDD
	EndDate            string `json:"EndDate" example:"20170702"`      //Дата окончания периода (включительно) в формате YYYYMMDD
	CnlList            []struct {
		Cnl string `json:"Cnl" example:"1018566"` //ID канала (int, not nillable)
	} `json:"CnlList"`
	ProtocolVersion string `json:"ProtocolVersion" example:"2"`
}

type ProgramBreaksQuery struct {
	Month struct {
		Eq int `json:"eq" example:"201902"`
	} `json:"Month"`
	CnlID struct {
		Eq int `json:"eq" example:"1020232"`
	} `json:"CnlID"`
	WeekDay struct {
		Ee int `json:"ge" example:"1"`
	} `json:"WeekDay"`
}

func (request *ProgramBreaksLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetProgramBreaksType
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
		chunk := 50
		chunkCount := 0
		var j int
		for i := 0; i < len(request.CnlList); i += chunk {
			chunkCount++
			j += chunk
			if j > len(request.CnlList) {
				j = len(request.CnlList)
			}
			req := goConvert.New()
			req.Set("SellingDirectionID", request.SellingDirectionID)
			req.Set("InclProgAttr", request.InclProgAttr)
			req.Set("InclForecast", request.InclForecast)
			req.Set("AudRatDec", "8")
			req.Set("StartDate", day.Format(`20060102`))
			req.Set("EndDate", day.Format(`20060102`))
			req.Set("LightMode", "0")
			req.Set("CnlList", request.CnlList[i:j])
			req.Set("ProtocolVersion", "2")
			req.Set("Path", chunkCount)
			err := amqpConfig.PublishJson(qName, req)
			if err != nil {
				fmt.Printf("Q:%s - err:%s", qName, err.Error())
				return nil, err
			}
		}
	}
	result["days"] = days
	result["status"] = "ok"
	return result, nil
}

func (request *ProgramBreaksLoadRequest) getDays() ([]time.Time, error) {
	return utils.GetDaysByPeriod(request.StartDate, request.EndDate)
}

func (request *Any) QueryProgramBreaks() ([]ProgramBreaks, error) {
	var result []ProgramBreaks
	query := ProgramBreaksBadgerQuery{}
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
