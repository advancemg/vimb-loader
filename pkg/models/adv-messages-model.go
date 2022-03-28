package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"strconv"
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
	if !cfg.Loading {
		return nil
	}
	errorCh := make(chan error)
	go func() {
		qName := GetAdvMessagesType
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
			var bodyJson GetAdvMessages
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

func (cfg *AdvMessagesConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetAdvMessagesType
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
		type AdtID struct {
			AdtID string `json:"AdtID"`
		}
		badgerAdvertisers := storage.NewBadger(DbCustomConfigAdvertisers)
		badgerMonth := storage.NewBadger(DbCustomConfigMonth)
		defer badgerAdvertisers.Close()
		defer badgerMonth.Close()
		months := map[string][]string{}
		advertisers := map[string]AdtID{}
		badgerAdvertisers.Iterate(func(key []byte, value []byte) {
			advertisers[string(key)] = AdtID{AdtID: string(value)}
		})
		badgerMonth.Iterate(func(key []byte, value []byte) {
			month, err := strconv.Atoi(string(value)[4:6])
			if err != nil {
				panic(err)
			}
			year, err := strconv.Atoi(string(value)[0:4])
			if err != nil {
				panic(err)
			}
			days, err := utils.GetDaysFromMonth(year, time.Month(month))
			if err != nil {
				panic(err)
			}
			months[string(key)] = days
		})
		var adtList []AdtID
		for _, adt := range advertisers {
			adtList = append(adtList, adt)
		}
		for month, days := range months {
			startDay := fmt.Sprintf("%s%s", month, days[0])
			endDay := fmt.Sprintf("%s%s", month, days[len(days)-1])
			request := goConvert.New()
			request.Set("CreationDateStart", fmt.Sprintf("%s-%s-%s", startDay[0:4], startDay[4:6], startDay[6:8]))
			request.Set("CreationDateEnd", fmt.Sprintf("%s-%s-%s", endDay[0:4], endDay[4:6], endDay[6:8]))
			//request.Set("Advertisers", adtList)
			request.Set("Aspects", "2")
			//request.Set("AdvertisingMessageIDs", "")
			request.Set("FillMaterialTags", "true")
			err := amqpConfig.PublishJson(qName, request)
			if err != nil {
				fmt.Printf("Q:%s - err:%s", qName, err.Error())
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

func (request *GetAdvMessages) UploadToS3() error {
	for {
		typeName := GetAdvMessagesType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout()
				continue
			}
			return err
		}
		month, _ := request.Get("CreationDateStart")
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%v/%s-%s.gz", utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
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
	xml, _ := xmlRequestHeader.ToXml()
	fmt.Println(string(xml))
	return xmlRequestHeader.ToXml()
}
