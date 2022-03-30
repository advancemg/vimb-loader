package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"io"
	"os"
)

type SwaggerGetChannelsRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
}

type GetChannels struct {
	goConvert.UnsortedMap
}

type ChannelConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

func (cfg *ChannelConfiguration) StartJob() chan error {
	if !cfg.Loading {
		return nil
	}
	errorCh := make(chan error)
	go func() {
		qName := GetChannelsType
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
			var bodyJson GetChannels
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

func (cfg *ChannelConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetChannelsType
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
		request := goConvert.New()
		request.Set("SellingDirectionID", cfg.SellingDirection)
		err = amqpConfig.PublishJson(qName, request)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
	}
}

func (request *GetChannels) GetDataJson() (*JsonResponse, error) {
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

func (request *GetChannels) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *GetChannels) UploadToS3() error {
	for {
		typeName := GetChannelsType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout()
				continue
			}
			return err
		}
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%s-%s.gz", utils.Actions.Client, typeName, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		/*update data from gz file*/
		err = request.DataConfiguration(newS3Key)
		if err != nil {
			return err
		}
		return nil
	}
}

func (request *GetChannels) DataConfiguration(s3Key string) error {
	download, err := s3.Download(s3Key)
	if err != nil {
		return err
	}
	open, err := os.Open(download)
	if err != nil {
		return err
	}
	defer open.Close()
	zipBuffer := new(bytes.Buffer)
	_, err = io.Copy(zipBuffer, open)
	if err != nil {
		return fmt.Errorf("not copy zip data, %v", err)
	}
	toJson, err := goConvert.ZipXmlToJson(zipBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("not ZipXmlToJson, %v", err)
	}
	channels := storage.NewBadger(DbChannels)
	defer channels.Close()
	var channelMap map[string]interface{}
	json.Unmarshal(toJson, &channelMap)
	if err != nil {
		return err
	}
	for _, v := range channelMap {
		if mapLevel2, ok := v.(map[string]interface{}); ok {
			for _, v2 := range mapLevel2 {
				if array, ok := v2.([]interface{}); ok {
					for _, value := range array {
						id := value.(map[string]interface{})["ID"].(string)
						mainChan := value.(map[string]interface{})["MainChnl"].(string)
						channels.Set(id, []byte(mainChan))
					}
				}
			}
		}
	}
	return nil
}

func (request *GetChannels) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	sellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", sellingDirectionID)
	}
	xmlRequestHeader.Set("GetChannels", body)
	return xmlRequestHeader.ToXml()
}
