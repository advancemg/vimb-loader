package models

import (
	"encoding/json"
	"errors"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/usecase"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type GetProgramBreaksLight struct {
	Path int
	goConvert.UnsortedMap
}

type ProgramBreaksLightConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

func (cfg *ProgramBreaksLightConfiguration) StartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := GetProgramBreaksLightModeType
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
			var bodyJson GetProgramBreaksLight
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			s3Message, err := bodyJson.UploadToS3()
			if err != nil {
				errorCh <- err
			}
			err = amqpConfig.PublishJson(ProgramBreaksLightModeUpdateQueue, s3Message)
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (cfg *ProgramBreaksLightConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetProgramBreaksLightModeType
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "ProgramBreaksLightMode InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "ProgramBreaksLightMode InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		type Channels struct {
			Cnl string `json:"Cnl"`
		}
		channels := map[int64]struct{}{}
		var budgets []Budget
		var cnl []Channels
		months := map[int64][]time.Time{}
		db, table := utils.SplitDbAndTable(DbBudgets)
		repo := usecase.OpenDb(db, table)
		err = repo.FindWhereGe(&budgets, "Month", int64(-1))
		if err != nil {
			log.PrintLog("vimb-loader", "ProgramBreaksLightMode InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		for _, budget := range budgets {
			channels[*budget.CnlID] = struct{}{}
			days, err := utils.GetDaysFromYearMonthInt(*budget.Month)
			if err != nil {
				panic(err)
			}
			months[*budget.Month] = days
		}
		for channel, _ := range channels {
			cnl = append(cnl, Channels{
				Cnl: fmt.Sprintf("%d", channel),
			})
		}
		for month, days := range months {
			for _, day := range days {
				chunk := 50
				chunkCount := 0
				var j int
				for i := 0; i < len(cnl); i += chunk {
					chunkCount++
					j += chunk
					if j > len(cnl) {
						j = len(cnl)
					}
					var startEndDay string
					if day.Day() <= 9 {
						startEndDay = fmt.Sprintf("%d0%d", month, day.Day())
					} else {
						startEndDay = fmt.Sprintf("%d%d", month, day.Day())
					}
					request := goConvert.New()
					request.Set("SellingDirectionID", cfg.SellingDirection)
					request.Set("InclProgAttr", "0")
					request.Set("InclForecast", "0")
					request.Set("AudRatDec", "8")
					request.Set("StartDate", startEndDay)
					request.Set("EndDate", startEndDay)
					request.Set("LightMode", "1")
					request.Set("CnlList", cnl[i:j])
					request.Set("ProtocolVersion", "2")
					request.Set("Path", chunkCount)
					err := amqpConfig.PublishJson(qName, request)
					if err != nil {
						log.PrintLog("vimb-loader", "ProgramBreaksLightMode InitJob", "error", "Q:", qName, "err:", err.Error())
						return
					}
				}
			}
		}
	}
}

func (request *GetProgramBreaksLight) GetDataJson() (*JsonResponse, error) {
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

func (request *GetProgramBreaksLight) GetDataXmlZip() (*StreamResponse, error) {
	for {
		var isTimeout utils.Timeout
		db, table := utils.SplitDbAndTable(DbTimeout)
		repo := usecase.OpenDb(db, table)
		err := repo.Get("_id", &isTimeout)
		if err != nil {
			if errors.Is(err, usecase.ErrNotFound) {
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

func (request *GetProgramBreaksLight) UploadToS3() (*MqUpdateMessage, error) {
	for {
		typeName := GetProgramBreaksLightModeType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("GetProgramBreaksLight")
				continue
			}
			return nil, err
		}
		sellingDirectionID, _ := request.Get("SellingDirectionID")
		startDate, _ := request.Get("StartDate")
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%s/%s/%s/%d/%s-%s.gz", sellingDirectionID, utils.Actions.Client, typeName, startDate.(string)[4:6], startDate, request.Path, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return nil, err
		}
		return &MqUpdateMessage{
			Key: newS3Key,
		}, nil
	}
}

func (request *GetProgramBreaksLight) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	path, exist := request.Get("Path")
	if exist {
		request.Path = int(path.(float64))
	}
	sellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", sellingDirectionID)
	}
	inclProgAttr, exist := request.Get("InclProgAttr")
	if exist {
		body.Set("InclProgAttr", inclProgAttr)
	}
	inclForecast, exist := request.Get("InclForecast")
	if exist {
		body.Set("InclForecast", inclForecast)
	}
	audRatDec, exist := request.Get("AudRatDec")
	if exist {
		body.Set("AudRatDec", audRatDec)
	}
	startDate, exist := request.Get("StartDate")
	if exist {
		body.Set("StartDate", startDate)
	}
	endDate, exist := request.Get("EndDate")
	if exist {
		body.Set("EndDate", endDate)
	}
	lightMode, exist := request.Get("LightMode")
	if exist {
		body.Set("LightMode", lightMode)
	}
	cnlList, exist := request.Get("CnlList")
	if exist {
		body.Set("CnlList", cnlList)
	}
	protocolVersion, exist := request.Get("ProtocolVersion")
	if exist {
		body.Set("ProtocolVersion", protocolVersion)
	}
	xmlRequestHeader.Set("GetProgramBreaks", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}
