package models

import (
	"encoding/json"
	"errors"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/store"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type SwaggerGetDeletedSpotInfoRequest struct {
	DateStart  string `json:"DateStart"`
	DateEnd    string `json:"DateEnd"`
	Agreements []struct {
		Id string `json:"ID"`
	} `json:"Agreements"`
}

type GetDeletedSpotInfo struct {
	goConvert.UnsortedMap
}

type DeletedSpotInfoConfiguration struct {
	Cron    string `json:"cron"`
	Loading bool   `json:"loading"`
}

func (cfg *DeletedSpotInfoConfiguration) StartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := GetDeletedSpotInfoType
		amqpConfig := mq_broker.InitConfig()
		defer amqpConfig.Close()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			errorCh <- err
		}
		ch, err := amqpConfig.Channel()
		if err != nil {
			errorCh <- err
		}
		defer ch.Close()
		err = ch.Qos(1, 0, false)
		messages, err := ch.Consume(qName, "",
			false,
			false,
			false,
			false,
			nil)
		for msg := range messages {
			var bodyJson GetDeletedSpotInfo
			err = json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			s3Message, err := bodyJson.UploadToS3()
			if err != nil {
				errorCh <- err
			}
			err = amqpConfig.PublishJson(DeletedSpotInfoUpdateQueue, s3Message)
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (cfg *DeletedSpotInfoConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetDeletedSpotInfoType
		amqpConfig := mq_broker.InitConfig()
		defer amqpConfig.Close()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "DeletedSpotInfo InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "DeletedSpotInfo InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		type agreement struct {
			Id int64 `json:"ID"`
		}
		agr := map[int64]struct{}{}
		var agreements []agreement
		var budgets []Budget
		months := map[int64][]time.Time{}
		db, table := utils.SplitDbAndTable(DbBudgets)
		repo := store.OpenDb(db, table)
		err = repo.FindWhereGe(&budgets, "Month", int64(-1))
		if err != nil {
			log.PrintLog("vimb-loader", "DeletedSpotInfo InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		for _, budget := range budgets {
			agr[*budget.AgrID] = struct{}{}
			days, err := utils.GetDaysFromYearMonthInt(*budget.Month)
			if err != nil {
				panic(err)
			}
			months[*budget.Month] = days
		}
		for agrId, _ := range agr {
			agreements = append(agreements, agreement{agrId})
		}
		for _, days := range months {
			startDay := fmt.Sprintf("%v", days[0].Format(time.RFC3339))
			endDay := fmt.Sprintf("%v", days[len(days)-1].Format(time.RFC3339))
			startDay = startDay[0 : len(startDay)-1]
			endDay = endDay[0 : len(endDay)-1]
			request := goConvert.New()
			request.Set("DateStart", startDay)
			request.Set("DateEnd", endDay)
			request.Set("Agreements", agreements)
			err = amqpConfig.PublishJson(qName, request)
			if err != nil {
				log.PrintLog("vimb-loader", "DeletedSpotInfo InitJob", "error", "Q:", qName, "err:", err.Error())
				return
			}
		}
	}
}

func (request *GetDeletedSpotInfo) GetDataJson() (*JsonResponse, error) {
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.Actions.RequestJson(req)
	if err != nil {
		return nil, err
	}
	var body = map[string]interface{}{}
	err = json.Unmarshal(resp, &body)
	if err != nil {
		return nil, err
	}
	return &JsonResponse{
		Body:    body,
		Request: string(req),
	}, nil
}

func (request *GetDeletedSpotInfo) GetDataXmlZip() (*StreamResponse, error) {
	for {
		var isTimeout utils.Timeout
		db, table := utils.SplitDbAndTable(DbTimeout)
		repo := store.OpenDb(db, table)
		err := repo.Get("_id", &isTimeout)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				isTimeout.IsTimeout = false
			} else {
				return nil, err
			}
		}
		if isTimeout.IsTimeout {
			time.Sleep(1 * time.Second)
			continue
		}
		if !isTimeout.IsTimeout {
			break
		}
	}
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.Actions.Request(req)
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    resp,
		Request: string(req),
	}, nil
}

func (request *GetDeletedSpotInfo) UploadToS3() (*MqUpdateMessage, error) {
	for {
		typeName := GetDeletedSpotInfoType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("GetDeletedSpotInfo")
				continue
			}
			return nil, err
		}
		month, _ := request.Get("DateStart")
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%v/%s-%s.gz", utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return nil, err
		}
		return &MqUpdateMessage{
			Key: newS3Key,
		}, nil
	}
}

func (request *GetDeletedSpotInfo) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	dateStart, exist := request.Get("DateStart")
	if exist {
		body.Set("DateStart", dateStart)
	}
	dateEnd, exist := request.Get("DateEnd")
	if exist {
		body.Set("DateEnd", dateEnd)
	}
	agreements, exist := request.Get("Agreements")
	if exist {
		body.Set("Agreements", agreements)
	}
	xmlRequestHeader.Set("GetDeletedSpotInfo", body)
	return xmlRequestHeader.ToXml()
}
