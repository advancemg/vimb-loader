package models

import (
	"encoding/json"
	"errors"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/store"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type SwaggerGetBudgetsRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
	StartMonth         string `json:"StartMonth"`
	EndMonth           string `json:"EndMonth"`
	AdvertiserList     []struct {
		Id string `json:"ID"`
	} `json:"AdvertiserList"`
	ChannelList []struct {
		Id string `json:"ID"`
	} `json:"ChannelList"`
}

type GetBudgets struct {
	goConvert.UnsortedMap
}

type BudgetConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

func (cfg *BudgetConfiguration) StartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := GetBudgetsType
		amqpConfig := mq_broker.InitConfig()
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
			var bodyJson GetBudgets
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			s3Message, err := bodyJson.UploadToS3()
			if err != nil {
				errorCh <- err
			}
			err = amqpConfig.PublishJson(BudgetsUpdateQueue, s3Message)
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (cfg *BudgetConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetBudgetsType
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "Budget InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "Budget InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		months, err := utils.GetActualMonths()
		if err != nil {
			log.PrintLog("vimb-loader", "Budget InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		for _, month := range months {
			request := goConvert.New()
			request.Set("SellingDirectionID", cfg.SellingDirection)
			request.Set("StartMonth", month.ValueString)
			request.Set("EndMonth", month.ValueString)
			err := amqpConfig.PublishJson(qName, request)
			if err != nil {
				log.PrintLog("vimb-loader", "Budget InitJob", "error", "Q:", qName, "err:", err.Error())
				return
			}
		}
	}
}

func (request *GetBudgets) GetDataJson() (*JsonResponse, error) {
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

func (request *GetBudgets) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *GetBudgets) UploadToS3() (*MqUpdateMessage, error) {
	for {
		typeName := GetBudgetsType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("GetBudgets")
				continue
			}
			return nil, err
		}
		month, _ := request.Get("StartMonth")
		if err != nil {
			return nil, err
		}
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

func (request *GetBudgets) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	SellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", SellingDirectionID)
	}
	StartMonth, exist := request.Get("StartMonth")
	if exist {
		body.Set("StartMonth", StartMonth)
	}
	EndMonth, exist := request.Get("EndMonth")
	if exist {
		body.Set("EndMonth", EndMonth)
	}
	AdvertiserList, exist := request.Get("AdvertiserList")
	if exist {
		body.Set("AdvertiserList", AdvertiserList)
	}
	ChannelList, exist := request.Get("ChannelList")
	if exist {
		body.Set("ChannelList", ChannelList)
	}
	xmlRequestHeader.Set("GetBudgets", body)
	return xmlRequestHeader.ToXml()
}
