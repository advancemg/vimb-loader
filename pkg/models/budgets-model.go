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
				vimbError.CheckTimeout()
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
		/*update data from gz file*/
		/*err = request.DataConfiguration(newS3Key)
		if err != nil {
			return err
		}*/
		return &MqUpdateMessage{
			Key: newS3Key,
		}, nil
	}
}

func (request *GetBudgets) DataConfiguration(s3Key string) error {
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
	badgerMonth := storage.NewBadger(DbCustomConfigMonth)
	badgerAdvertisers := storage.NewBadger(DbCustomConfigAdvertisers)
	badgerChannels := storage.NewBadger(DbCustomConfigChannels)
	defer badgerMonth.Close()
	defer badgerAdvertisers.Close()
	defer badgerChannels.Close()
	var budgetMap map[string]interface{}
	err = json.Unmarshal(toJson, &budgetMap)
	if err != nil {
		return err
	}
	for _, v := range budgetMap {
		if map1, ok := v.(map[string]interface{}); ok {
			for _, v1 := range map1 {
				if arr, ok := v1.([]interface{}); ok {
					for _, arrVal := range arr {
						for _, value := range arrVal.(map[string]interface{}) {
							if mpL4, ok := value.(map[string]interface{}); ok {
								adtId := mpL4["AdtID"].(string)
								cnlId := mpL4["CnlID"].(string)
								month := mpL4["Month"].(string)
								badgerAdvertisers.Set(adtId, []byte(adtId))
								badgerChannels.Set(cnlId, []byte(cnlId))
								badgerMonth.Set(month, []byte(month))
							}
						}
					}
				}
			}
		}
	}
	return nil
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
