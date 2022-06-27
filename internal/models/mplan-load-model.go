package models

import (
	"encoding/json"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/store"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
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
			log.PrintLog("vimb-loader", "Mediaplan InitTasks", "error", "Q:", qName, "err:", err.Error())
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
		Eq int64 `json:"eq" example:"700068653"`
	} `json:"AdtID"`
	AgrID struct {
		Eq int64 `json:"eq" example:"81024"`
	} `json:"AgrID"`
	MplID struct {
		Eq int64 `json:"eq" example:"14824608"`
	} `json:"MplID"`
}

type AggMediaplanQuery struct {
	MplMonth struct {
		Eq int64 `json:"eq" example:"201902"`
	} `json:"MplMonth"`
	CppPrime struct {
		Eq float64 `json:"eq" example:"2205.20143600803"`
	} `json:"CppPrime"`
	CppOffPrime struct {
		Eq float64 `json:"ge" example:"100"`
	} `json:"CppOffPrime"`
}

func (request *Any) QueryMediaplans() ([]Mediaplan, error) {
	var result []Mediaplan
	marshal, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}
	db, table := utils.SplitDbAndTable(DbMediaplans)
	dbMediaplans := store.OpenDb(db, table)
	err = dbMediaplans.FindJson(&result, marshal)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (request *Any) QueryAggMediaplans() ([]MediaplanAgg, error) {
	var result []MediaplanAgg
	marshal, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}
	db, table := utils.SplitDbAndTable(DbAggMediaplans)
	dbAggMediaplans := store.OpenDb(db, table)
	err = dbAggMediaplans.FindJson(&result, marshal)
	if err != nil {
		return nil, err
	}
	return result, nil
}
