package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
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
	if !cfg.Loading {
		return nil
	}
	errorCh := make(chan error)
	go func() {
		qName := GetDeletedSpotInfoType
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
			var bodyJson GetDeletedSpotInfo
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			err = bodyJson.UploadToS3()
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
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		months, err := utils.GetActualMonths()
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
		for _, month := range months {
			request := goConvert.New()
			request.Set("DateStart", month.ValueString)
			request.Set("DateEnd", month.ValueString)
			err := amqpConfig.PublishJson(qName, request)
			if err != nil {
				fmt.Printf("Q:%s - err:%s", qName, err.Error())
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

func (request *GetDeletedSpotInfo) UploadToS3() error {
	for {
		typeName := GetDeletedSpotInfoType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout()
				continue
			}
			return err
		}
		month, _ := request.Get("DateStart")
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%v/%s-%s.gz", utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
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
