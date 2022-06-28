package models

import (
	"encoding/json"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/store"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type BudgetLoadRequest struct {
	StartMonth         string `json:"StartMonth" example:"201903"`
	EndMonth           string `json:"EndMonth" example:"201912"`
	SellingDirectionID string `json:"SellingDirectionID" example:"23"`
}

type BudgetQuery struct {
	Month struct {
		Eq int64 `json:"eq" example:"201902"`
	} `json:"Month"`
	CnlID struct {
		Eq int64 `json:"eq" example:"1020335"`
	} `json:"CnlID"`
	AdtID struct {
		Ee int64 `json:"eq" example:"700068653"`
	} `json:"AdtID"`
}

func (request *BudgetLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetBudgetsType
	amqpConfig := mq_broker.InitConfig()
	defer amqpConfig.Close()
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
			log.PrintLog("vimb-loader", "Budget InitTasks", "error", "Q:", qName, "err:", err.Error())
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

func (request *Any) QueryBudgets() ([]Budget, error) {
	var result []Budget
	db, table := utils.SplitDbAndTable(DbBudgets)
	repo := store.OpenDb(db, table)
	marshal, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}
	err = repo.FindJson(&result, marshal)
	if err != nil {
		return nil, err
	}
	return result, nil
}
