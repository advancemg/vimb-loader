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

type SwaggerGetCustomersWithAdvertisersRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
}

type GetCustomersWithAdvertisers struct {
	goConvert.UnsortedMap
}

type CustomersWithAdvertisersConfiguration struct {
	Cron               string `json:"cron"`
	SellingDirectionID string `json:"sellingDirectionID"`
	Loading            bool   `json:"loading"`
}

func (cfg *CustomersWithAdvertisersConfiguration) StartJob() error {
	if !cfg.Loading {
		return nil
	}
	qName := GetCustomersWithAdvertisersType
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
		var bodyJson GetCustomersWithAdvertisers
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

func (cfg *CustomersWithAdvertisersConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetCustomersWithAdvertisersType
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

func (request *GetCustomersWithAdvertisers) GetDataJson() (*JsonResponse, error) {
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

func (request *GetCustomersWithAdvertisers) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *GetCustomersWithAdvertisers) UploadToS3() error {
	for {
		typeName := GetCustomersWithAdvertisersType
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

func (request *GetCustomersWithAdvertisers) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	sellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", sellingDirectionID)
	}
	xmlRequestHeader.Set("GetCustomersWithAdvertisers", body)
	return xmlRequestHeader.ToXml()
}
