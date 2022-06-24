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

type SwaggerGetRanksRequest struct {
}

type GetRanks struct {
	goConvert.UnsortedMap
}

type RanksConfiguration struct {
	Cron    string `json:"cron"`
	Loading bool   `json:"loading"`
}

func (cfg *RanksConfiguration) StartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := GetRanksType
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
			var bodyJson GetRanks
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			s3Message, err := bodyJson.UploadToS3()
			if err != nil {
				errorCh <- err
			}
			err = amqpConfig.PublishJson(RanksUpdateQueue, s3Message)
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (cfg *RanksConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetRanksType
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "Ranks InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "Ranks InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		request := goConvert.New()
		request.Set("GetRanks", "")
		err = amqpConfig.PublishJson(qName, request)
		if err != nil {
			log.PrintLog("vimb-loader", "Ranks InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
	}
}

func (request *GetRanks) GetDataJson() (*JsonResponse, error) {
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

func (request *GetRanks) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *GetRanks) UploadToS3() (*MqUpdateMessage, error) {
	for {
		typeName := GetRanksType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("GetRanks")
				continue
			}
			return nil, err
		}
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%s-%s.gz", utils.Actions.Client, typeName, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return nil, err
		}
		return &MqUpdateMessage{
			Key: newS3Key,
		}, nil
	}
}

func (request *GetRanks) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	xmlRequestHeader.Set("GetRanks", "")
	return xmlRequestHeader.ToXml()
}
