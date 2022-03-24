package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type SwaggerChangeSpotRequest struct {
	FirstSpotID  string `json:"FirstSpotID"`
	SecondSpotID string `json:"SecondSpotID"`
}

type ChangeSpot struct {
	goConvert.UnsortedMap
}

type ChangeSpotConfiguration struct {
	Cron    string `json:"cron"`
	Loading bool   `json:"loading"`
}

func (cfg *ChangeSpotConfiguration) StartJob() error {
	if !cfg.Loading {
		return nil
	}
	qName := ChangeSpotType
	amqpConfig := mq_broker.InitConfig()
	err := amqpConfig.DeclareSimpleQueue(qName)
	if err != nil {
		return err
	}
	ch, err := amqpConfig.Channel()
	if err != nil {
		return err
	}
	err = ch.Qos(1, 0, false)
	messages, err := ch.Consume(qName, "",
		false,
		false,
		false,
		false,
		nil)
	for msg := range messages {
		var bodyJson ChangeSpot
		err := json.Unmarshal(msg.Body, &bodyJson)
		if err != nil {
			return err
		}
		err = bodyJson.UploadToS3()
		if err != nil {
			return err
		}
		msg.Ack(false)
	}
	return nil
}

func (cfg *ChangeSpotConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := ChangeSpotType
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
	}
}

func (request *ChangeSpot) GetDataJson() (*JsonResponse, error) {
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

func (request *ChangeSpot) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *ChangeSpot) UploadToS3() error {
	for {
		typeName := ChangeSpotType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				code := vimbError.Code
				switch code {
				case 1001:
					fmt.Printf("Vimb code %v timeout...", code)
					time.Sleep(time.Minute * 1)
					continue
				case 1003:
					fmt.Printf("Vimb code %v timeout...", code)
					time.Sleep(time.Minute * 2)
					continue
				default:
					fmt.Printf("Vimb code %v - not implemented timeout...", code)
					time.Sleep(time.Minute * 1)
					continue
				}
			}
			return err
		}
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%s-%s.gz", utils.Actions.Client, typeName, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
	}
}

func (request *ChangeSpot) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	FirstSpotID, exist := request.Get("FirstSpotID")
	if exist {
		body.Set("FirstSpotID", FirstSpotID)
	}
	SecondSpotID, exist := request.Get("SecondSpotID")
	if exist {
		body.Set("SecondSpotID", SecondSpotID)
	}
	xmlRequestHeader.Set("ChangeSpots", body)
	return xmlRequestHeader.ToXml()
}
