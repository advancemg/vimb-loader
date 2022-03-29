package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type BudgetLoadRequest struct {
	StartMonth         string `json:"StartMonth" example:"201903"`
	EndMonth           string `json:"EndMonth" example:"201912"`
	SellingDirectionID string `json:"SellingDirectionID" example:"23"`
}

func (request *BudgetLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetBudgetsType
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
		err := amqpConfig.PublishJson(qName, req)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return nil, err
		}
	}
	result["months"] = months
	result["status"] = "ok"
	return result, nil
}

func (request *BudgetLoadRequest) getMonths() ([]utils.YearMonth, error) {
	return utils.GetPeriodFromYearMonths(request.StartMonth, request.EndMonth)
}
