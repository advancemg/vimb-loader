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

type SwaggerAddSpotRequest struct {
	BlockID         string `json:"BlockID"`
	FilmID          string `json:"FilmID"`
	Position        string `json:"Position"`
	FixedPosition   string `json:"FixedPosition"`
	AuctionBidValue string `json:"AuctionBidValue"`
}

type AddSpot struct {
	goConvert.UnsortedMap
}

type AddSpotConfiguration struct {
	Cron    string `json:"cron"`
	Loading bool   `json:"loading"`
}

func (cfg *AddSpotConfiguration) StartJob() error {
	if !cfg.Loading {
		return nil
	}
	qName := AddSpotType
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
		var bodyJson AddSpot
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

func (cfg *AddSpotConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := AddSpotType
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

func (request *AddSpot) GetDataJson() (*JsonResponse, error) {
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

func (request *AddSpot) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *AddSpot) UploadToS3() error {
	for {
		typeName := AddSpotType
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

func (request *AddSpot) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	BlockID, exist := request.Get("BlockID")
	if exist {
		body.Set("BlockID", BlockID)
	}
	FilmID, exist := request.Get("FilmID")
	if exist {
		body.Set("FilmID", FilmID)
	}
	Position, exist := request.Get("Position")
	if exist {
		body.Set("Position", Position)
	}
	FixedPosition, exist := request.Get("FixedPosition")
	if exist {
		body.Set("FixedPosition", FixedPosition)
	}
	AuctionBidValue, exist := request.Get("AuctionBidValue")
	if exist {
		body.Set("AuctionBidValue", AuctionBidValue)
	}
	xmlRequestHeader.Set("AddSpot", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}
