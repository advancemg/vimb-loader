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

func (cfg *MediaplanConfiguration) StartJob() error {
	if !cfg.Loading {
		return nil
	}
	qName := GetMPLansType
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
		var bodyJson GetMPLans
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

func (cfg *MediaplanConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetMPLansType
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
			request.Set("SellingDirectionID", cfg.SellingDirection)
			request.Set("StartMonth", month.ValueString)
			request.Set("EndMonth", month.ValueString)
			err := amqpConfig.PublishJson(qName, request)
			if err != nil {
				fmt.Printf("Q:%s - err:%s", qName, err.Error())
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

func (request *GetMPLans) UploadToS3() error {
	for {
		typeName := GetMPLansType
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
		month, _ := request.Get("StartMonth")
		if err != nil {
			return err
		}
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%v/%s-%s.gz", utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
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
	xmlRequestHeader.Set("GetMPLans", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}
