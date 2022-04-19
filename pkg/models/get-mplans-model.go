package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"github.com/timshannon/badgerhold"
)

type SwaggerGetMPLansRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
	StartMonth         string `json:"StartMonth"`
	EndMonth           string `json:"EndMonth"`
	AdtList            []struct {
		AdtID string `json:"AdtID"`
	} `json:"AdtList"`
	ChannelList []struct {
		Cnl string `json:"Cnl"`
	} `json:"ChannelList"`
	IncludeEmpty string `json:"IncludeEmpty"`
}

type GetMPLans struct {
	goConvert.UnsortedMap
}

type MediaplanConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

func (cfg *MediaplanConfiguration) StartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := GetMPLansType
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			errorCh <- err
		}
		ch, err := amqpConfig.Channel()
		if err != nil {
			errorCh <- err
		}
		err = ch.Qos(1, 0, false)
		messages, err := ch.Consume(qName, "",
			false,
			false,
			false,
			false,
			nil)
		for msg := range messages {
			var bodyJson GetMPLans
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			s3Message, err := bodyJson.UploadToS3()
			if err != nil {
				errorCh <- err
			}
			err = amqpConfig.PublishJson(MPLansUpdateQueue, s3Message)
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (cfg *MediaplanConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetMPLansType
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "Mediaplans InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "Mediaplans InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		months := map[int64]struct{}{}
		var budgets []Budget
		badgerBudgets := storage.Open(DbBudgets)
		err = badgerBudgets.Find(&budgets, badgerhold.Where("Month").Ge(int64(-1)))
		if err != nil {
			log.PrintLog("vimb-loader", "Mediaplans InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		for _, budget := range budgets {
			months[*budget.Month] = struct{}{}
		}
		for month, _ := range months {
			startMonth := fmt.Sprintf("%d", month)
			request := goConvert.New()
			request.Set("SellingDirectionID", cfg.SellingDirection)
			request.Set("StartMonth", startMonth)
			request.Set("EndMonth", startMonth)
			request.Set("AdtList", []string{})
			request.Set("ChannelList", []string{})
			request.Set("IncludeEmpty", "false")
			err := amqpConfig.PublishJson(qName, request)
			if err != nil {
				log.PrintLog("vimb-loader", "Mediaplans InitJob", "error", "Q:", qName, "err:", err.Error())
				return
			}
		}
	}
}

func (request *GetMPLans) GetDataJson() (*JsonResponse, error) {
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

func (request *GetMPLans) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *GetMPLans) UploadToS3() (*MqUpdateMessage, error) {
	for {
		typeName := GetMPLansType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout()
				continue
			}
			return nil, err
		}
		sellingDirectionID, _ := request.Get("SellingDirectionID")
		month, _ := request.Get("StartMonth")
		if err != nil {
			return nil, err
		}
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%s/%v/%s-%s.gz", sellingDirectionID, utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return nil, err
		}
		return &MqUpdateMessage{
			Key: newS3Key,
		}, nil
	}
}

func (request *GetMPLans) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
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
	AdtList, exist := request.Get("AdtList")
	if exist {
		body.Set("AdtList", AdtList)
	}
	ChannelList, exist := request.Get("ChannelList")
	if exist {
		body.Set("ChannelList", ChannelList)
	}
	IncludeEmpty, exist := request.Get("IncludeEmpty")
	if exist {
		body.Set("IncludeEmpty", IncludeEmpty)
	}
	xmlRequestHeader.Set("GetMPlans", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}
