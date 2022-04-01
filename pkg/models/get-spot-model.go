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
			var bodyJson GetSpots
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

func (cfg *SpotsConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetSpotsType
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
		type Cnl struct {
			Cnl  int
			Main int
		}
		var cnl []int
		var advertisers []int
		var allChannels []Cnl
		var budgets []Budget
		var channels []Channel
		channelList := map[int]Cnl{}
		months := map[int][]string{}
		badgerBudgets := storage.NewBadger(DbBudgets)
		badgerBudgets.Iterate(func(key []byte, value []byte) {
			var budget Budget
			json.Unmarshal(value, &budget)
			budgets = append(budgets, budget)
		})
		for _, budget := range budgets {
			advertisers = append(advertisers, *budget.AdtID)
			cnl = append(cnl, *budget.CnlID)
			monthStr := fmt.Sprintf("%d", *budget.Month)
			month, err := strconv.Atoi(monthStr[4:6])
			if err != nil {
				panic(err)
			}
			year, err := strconv.Atoi(monthStr[0:4])
			if err != nil {
				panic(err)
			}
			days, err := utils.GetDaysFromMonth(year, time.Month(month))
			if err != nil {
				panic(err)
			}
			months[month] = days
		}
		channelsBadger := storage.NewBadger(DbChannels)
		if channelsBadger.Count() > 0 {
			channelsBadger.Iterate(func(key []byte, value []byte) {
				var channel Channel
				json.Unmarshal(value, &channel)
				channels = append(channels, channel)
			})
			for _, channel := range channels {
				allChannels = append(allChannels, Cnl{
					Cnl:  *channel.ID,
					Main: *channel.MainChnl,
				})
			}
			for _, channel := range cnl {
				for _, cnlIdMain := range allChannels {
					if cnlIdMain.Cnl == channel {
						channelList[channel] = cnlIdMain
					}
				}
			}
		}
		for month, days := range months {
			request := goConvert.New()
			startDay := fmt.Sprintf("%d%s", month, days[0])
			endDay := fmt.Sprintf("%d%s", month, days[len(days)-1])
			request.Set("SellingDirectionID", cfg.SellingDirection)
			request.Set("StartDate", startDay)
			request.Set("EndDate", endDay)
			request.Set("InclOrdBlocks", "1")
			request.Set("ChannelList", channelList)
			request.Set("AdtList", advertisers)
			err := amqpConfig.PublishJson(qName, request)
			if err != nil {
				fmt.Printf("Q:%s - err:%s", qName, err.Error())
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

func (request *GetSpots) UploadToS3() error {
	for {
		typeName := GetSpotsType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout()
				continue
			}
			return err
		}
		month, _ := request.Get("StartDate")
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%v/%s-%s.gz", utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
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
