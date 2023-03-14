package models

import (
	"encoding/json"
	"errors"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/store"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type SwaggerGetSpotsRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
	StartDate          string `json:"StartDate"`
	EndDate            string `json:"EndDate"`
	InclOrdBlocks      string `json:"InclOrdBlocks"`
	ChannelList        []struct {
		Cnl  string `json:"Cnl"`
		Main string `json:"Main"`
	} `json:"ChannelList"`
	AdtList []struct {
		AdtID string `json:"AdtID"`
	} `json:"AdtList"`
}

type GetSpots struct {
	goConvert.UnsortedMap
}

type SpotsConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

func (cfg *SpotsConfiguration) StartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := GetSpotsType
		amqpConfig := mq_broker.InitConfig()
		defer amqpConfig.Close()
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
			var bodyJson GetSpots
			err = json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			s3Message, err := bodyJson.UploadToS3()
			if err != nil {
				errorCh <- err
			}
			err = amqpConfig.PublishJson(SpotsUpdateQueue, s3Message)
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (cfg *SpotsConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetSpotsType
		amqpConfig := mq_broker.InitConfig()
		defer amqpConfig.Close()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "Spots InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			log.PrintLog("vimb-loader", "Spots InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		type Cnl struct {
			Cnl  string `json:"Cnl"`
			Main string `json:"Main"`
		}
		type Adv struct {
			AdtID string `json:"AdtID"`
		}
		var allChannels []Cnl
		var allAdvertisers []Adv
		var budgets []Budget
		var channels []Channel
		channelList := map[int64]Cnl{}
		months := map[int64][]time.Time{}
		advertisers := map[int64]int64{}
		db, table := utils.SplitDbAndTable(DbBudgets)
		repoBudgets, err := store.OpenDb(db, table)
		if err != nil {
			log.PrintLog("vimb-loader", "Spots InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		err = repoBudgets.FindWhereGe(&budgets, "Month", int64(-1))
		if err != nil {
			log.PrintLog("vimb-loader", "Spots InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		db, table = utils.SplitDbAndTable(DbChannels)
		repoChannels, err := store.OpenDb(db, table)
		if err != nil {
			log.PrintLog("vimb-loader", "Spots InitJob", "error", "Q:", qName, "err:", err.Error())
			return
		}
		err = repoChannels.FindWhereGe(&channels, "ID", int64(-1))
		for _, budget := range budgets {
			advertisers[*budget.AdtID] = *budget.AdtID
			channelList[*budget.CnlID] = Cnl{
				Cnl:  fmt.Sprintf("%d", *budget.CnlID),
				Main: "",
			}
			days, err := utils.GetDaysFromYearMonthInt(*budget.Month)
			if err != nil {
				log.PrintLog("vimb-loader", "Spots InitJob", "error", "Q:", qName, "err:", err.Error())
				return
			}
			months[*budget.Month] = days
		}
		for _, channel := range channels {
			if channelItem, ok := channelList[*channel.ID]; ok {
				channelItem.Main = fmt.Sprintf("%d", *channel.MainChnl)
				allChannels = append(allChannels, channelItem)
			}
		}
		for _, adv := range advertisers {
			allAdvertisers = append(allAdvertisers, Adv{
				AdtID: fmt.Sprintf("%v", adv),
			})
		}
		for month, days := range months {
			request := goConvert.New()
			var startDay, endDay string
			if days[0].Day() <= 9 {
				startDay = fmt.Sprintf("%d0%d", month, days[0].Day())
			} else {
				startDay = fmt.Sprintf("%d%d", month, days[0].Day())
			}
			if days[len(days)-1].Day() <= 9 {
				endDay = fmt.Sprintf("%d0%d", month, days[len(days)-1].Day())
			} else {
				endDay = fmt.Sprintf("%d%d", month, days[len(days)-1].Day())
			}
			request.Set("SellingDirectionID", cfg.SellingDirection)
			request.Set("StartDate", startDay)
			request.Set("EndDate", endDay)
			request.Set("InclOrdBlocks", "1")
			request.Set("ChannelList", allChannels)
			request.Set("AdtList", allAdvertisers)
			err = amqpConfig.PublishJson(qName, request)
			if err != nil {
				log.PrintLog("vimb-loader", "Spots InitJob", "error", "Q:", qName, "err:", err.Error())
				return
			}
		}
	}
}

func (request *GetSpots) GetDataJson() (*JsonResponse, error) {
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

func (request *GetSpots) GetDataXmlZip() (*StreamResponse, error) {
	for {
		var isTimeout utils.Timeout
		db, table := utils.SplitDbAndTable(DbTimeout)
		repo, err := store.OpenDb(db, table)
		if err != nil {
			return nil, err
		}
		err = repo.Get("_id", &isTimeout)
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

func (request *GetSpots) UploadToS3() (*MqUpdateMessage, error) {
	for {
		typeName := GetSpotsType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				err = vimbError.CheckTimeout("GetSpots")
				if err != nil {
					return nil, err
				}
				continue
			}
			return nil, err
		}
		month, _ := request.Get("StartDate")
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%v/%s-%s.gz", utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return nil, err
		}
		return &MqUpdateMessage{
			Key:   newS3Key,
			Month: month.(string)[:6],
		}, nil
	}
}

func (request *GetSpots) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	SellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", SellingDirectionID)
	}
	StartDate, exist := request.Get("StartDate")
	if exist {
		body.Set("StartDate", StartDate)
	}
	EndDate, exist := request.Get("EndDate")
	if exist {
		body.Set("EndDate", EndDate)
	}
	InclOrdBlocks, exist := request.Get("InclOrdBlocks")
	if exist {
		body.Set("InclOrdBlocks", InclOrdBlocks)
	}
	ChannelList, exist := request.Get("ChannelList")
	if exist {
		body.Set("ChannelList", ChannelList)
	}
	AdtList, exist := request.Get("AdtList")
	if exist {
		body.Set("AdtList", AdtList)
	}
	xmlRequestHeader.Set("GetSpots", body)
	return xmlRequestHeader.ToXml()
}
