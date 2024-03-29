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

type SwaggerGetAdvMessagesRequest struct {
	CreationDateStart string `json:"CreationDateStart"`
	CreationDateEnd   string `json:"CreationDateEnd"`
	Advertisers       []struct {
		Id string `json:"ID"`
	} `json:"Advertisers"`
	Aspects []struct {
		Id string `json:"ID"`
	} `json:"Aspects"`
	AdvertisingMessageIDs []struct {
		Id string `json:"ID"`
	} `json:"AdvertisingMessageIDs"`
	FillMaterialTags string `json:"FillMaterialTags"`
}

type GetAdvMessages struct {
	goConvert.UnsortedMap
}

type AdvMessagesConfiguration struct {
	Cron    string `json:"cron"`
	Loading bool   `json:"loading"`
}

func (cfg *AdvMessagesConfiguration) StartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := GetAdvMessagesType
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
			var bodyJson GetAdvMessages
			err = json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			s3Message, err := bodyJson.UploadToS3()
			if err != nil {
				errorCh <- err
			}
			err = amqpConfig.PublishJson(AdvMessagesUpdateQueue, s3Message)
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()

	return errorCh
}

func (cfg *AdvMessagesConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetAdvMessagesType
		amqpConfig := mq_broker.InitConfig()
		defer amqpConfig.Close()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "AdvMessages InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "AdvMessages InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		type adtID struct {
			AdtID string `json:"AdtID"`
		}
		type filed struct {
			Id string `json:"ID"`
		}
		var budgets []Budget
		months := map[int64][]time.Time{}
		db, table := utils.SplitDbAndTable(DbBudgets)
		repo, err := store.OpenDb(db, table)
		if err != nil {
			log.PrintLog("vimb-loader", "AdvMessages InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		err = repo.FindWhereGe(&budgets, "Month", int64(-1))
		if err != nil {
			log.PrintLog("vimb-loader", "AdvMessages InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		for _, budget := range budgets {
			days, err := utils.GetDaysFromYearMonthInt(*budget.Month)
			if err != nil {
				log.PrintLog("vimb-loader", "AdvMessages InitJob", "error", "Q:", qName, "err:", err.Error())
				return
			}
			months[*budget.Month] = days
		}
		for _, days := range months {
			request := goConvert.New()
			request.Set("CreationDateStart", days[0].String()[:10])
			request.Set("CreationDateEnd", days[len(days)-1].String()[:10])
			request.Set("Advertisers", []adtID{})
			request.Set("Aspects", []filed{})
			request.Set("AdvertisingMessageIDs", []filed{})
			request.Set("FillMaterialTags", "true")
			err = amqpConfig.PublishJson(qName, request)
			if err != nil {
				log.PrintLog("vimb-loader", "AdvMessages InitJob", "error", "Q:", qName, "err:", err.Error())
				return
			}
		}
	}
}

func (request *GetAdvMessages) GetDataJson() (*JsonResponse, error) {
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

func (request *GetAdvMessages) GetDataXmlZip() (*StreamResponse, error) {
	for {
		var isTimeout utils.Timeout
		db, table := utils.SplitDbAndTable(DbTimeout)
		repo, err := store.OpenDb(db, table)
		if err != nil {
			return nil, err
		}
		err = repo.Get("_id", &isTimeout)
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

func (request *GetAdvMessages) UploadToS3() (*MqUpdateMessage, error) {
	for {
		typeName := GetAdvMessagesType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				err = vimbError.CheckTimeout("GetAdvMessages")
				if err != nil {
					return nil, err
				}
				continue
			}
			return nil, err
		}
		month, _ := request.Get("CreationDateStart")
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

func (request *GetAdvMessages) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	creationDateStart, exist := request.Get("CreationDateStart")
	if exist {
		body.Set("CreationDateStart", creationDateStart)
	}
	creationDateEnd, exist := request.Get("CreationDateEnd")
	if exist {
		body.Set("CreationDateEnd", creationDateEnd)
	}
	advertisers, exist := request.Get("Advertisers")
	if exist {
		body.Set("Advertisers", advertisers)
	}
	aspects, exist := request.Get("Aspects")
	if exist {
		body.Set("Aspects", aspects)
	}
	advertisingMessageIDs, exist := request.Get("AdvertisingMessageIDs")
	if exist {
		body.Set("AdvertisingMessageIDs", advertisingMessageIDs)
	}
	fillMaterialTags, exist := request.Get("FillMaterialTags")
	if exist {
		body.Set("FillMaterialTags", fillMaterialTags)
	}
	xmlRequestHeader.Set("GetAdvMessages", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}
