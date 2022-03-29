package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type ProgramBreaksLoadRequest struct {
	StartDate          string `json:"StartDate" example:"20190301"`
	EndDate            string `json:"EndDate" example:"20190331"`
	SellingDirectionID string `json:"SellingDirectionID" example:"23"`
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
		req := goConvert.New()
		req.Set("SellingDirectionID", request.SellingDirectionID)
		req.Set("StartDate", day.Format(`20060102`))
		req.Set("EndDate", day.Format(`20060102`))
		req.Set("InclProgAttr", "1")
		req.Set("InclForecast", "1")
		req.Set("LightMode", "0")
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

func (request *ProgramBreaksLoadRequest) getDays() ([]time.Time, error) {
	return utils.GetDaysByPeriod(request.StartDate, request.EndDate)
}
