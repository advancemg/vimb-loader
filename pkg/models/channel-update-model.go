package models

import (
	"encoding/json"
	"fmt"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"reflect"
	"time"
)

type ChannelsUpdateRequest struct {
	S3Key string
}

func (channel *Channel) Key() string {
	return fmt.Sprintf("%d", *channel.ID)
}

type internalChanelAspect struct {
	Aspect map[string]interface{} `json:"Aspect"`
}

type Channel struct {
	ID                 *int           `json:"ID"`
	MainChnl           *int           `json:"MainChnl"`
	SellingDirectionID *int           `json:"SellingDirectionID"`
	CnlOrderNo         *int           `json:"CnlOrderNo"`
	CnlCentralID       *int           `json:"CnlCentralID"`
	IsDisabled         *int           `json:"IsDisabled"`
	BcpCentralID       *int           `json:"BcpCentralID"`
	ShortName          *string        `json:"ShortName"`
	BcpName            *string        `json:"BcpName"`
	StartTime          *string        `json:"StartTime"`
	EndTime            *string        `json:"EndTime"`
	TotalOffset        *string        `json:"TotalOffset"`
	Timestamp          time.Time      `json:"Timestamp"`
	Aspects            []ChanelAspect `json:"Aspects"`
}

type ChanelAspect struct {
	StartDate *time.Time `json:"StartDate"`
	EndDate   *time.Time `json:"EndDate"`
	ID        *int       `json:"ID"`
}

func (m *internalChanelAspect) Convert() (*ChanelAspect, error) {
	aspect := &ChanelAspect{
		StartDate: utils.TimeI(m.Aspect["StartDate"], `2006-01-02T15:04:05`),
		EndDate:   utils.TimeI(m.Aspect["EndDate"], `2006-01-02T15:04:05`),
		ID:        utils.IntI(m.Aspect["ID"]),
	}
	return aspect, nil
}

func (m *internalChannel) Convert() (*Channel, error) {
	timestamp := time.Now()
	var aspects []ChanelAspect
	if _, ok := m.Channel["Aspects"]; ok {
		marshalData, err := json.Marshal(m.Channel["Aspects"])
		if err != nil {
			return nil, err
		}
		switch reflect.TypeOf(m.Channel["Aspects"]).Kind() {
		case reflect.Array, reflect.Slice:
			var internalChanelAspectData []internalChanelAspect
			err = json.Unmarshal(marshalData, &internalChanelAspectData)
			if err != nil {
				return nil, err
			}
			for _, aspectItem := range internalChanelAspectData {
				aspect, err := aspectItem.Convert()
				if err != nil {
					return nil, err
				}
				aspects = append(aspects, *aspect)
			}
		case reflect.Map, reflect.Struct:
			var internalChanelAspectData internalChanelAspect
			err = json.Unmarshal(marshalData, &internalChanelAspectData)
			if err != nil {
				return nil, err
			}
			aspect, err := internalChanelAspectData.Convert()
			if err != nil {
				return nil, err
			}
			aspects = append(aspects, *aspect)
		}
	}
	channel := &Channel{
		ID:                 utils.IntI(m.Channel["ID"]),
		MainChnl:           utils.IntI(m.Channel["MainChnl"]),
		SellingDirectionID: utils.IntI(m.Channel["SellingDirectionID"]),
		CnlOrderNo:         utils.IntI(m.Channel["CnlOrderNo"]),
		CnlCentralID:       utils.IntI(m.Channel["CnlCentralID"]),
		IsDisabled:         utils.IntI(m.Channel["IsDisabled"]),
		BcpCentralID:       utils.IntI(m.Channel["BcpCentralID"]),
		ShortName:          utils.StringI(m.Channel["ShortName"]),
		BcpName:            utils.StringI(m.Channel["BcpName"]),
		StartTime:          utils.StringI(m.Channel["StartTime"]),
		EndTime:            utils.StringI(m.Channel["EndTime"]),
		TotalOffset:        utils.StringI(m.Channel["TotalOffset"]),
		Timestamp:          timestamp,
		Aspects:            aspects,
	}
	return channel, nil
}

func ChannelsStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := ChannelsUpdateQueue
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
			var bodyJson MqUpdateMessage
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			/*read from s3 by s3Key*/
			req := ChannelsUpdateRequest{
				S3Key: bodyJson.Key,
			}
			err = req.Update()
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (request *ChannelsUpdateRequest) Update() error {
	var err error
	request.S3Key, err = s3.Download(request.S3Key)
	if err != nil {
		return err
	}
	err = request.loadFromFile()
	if err != nil {
		return err
	}
	return nil
}

func (request *ChannelsUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	convertData, err := resp.Convert("responseGetChannels")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalChannel
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	for _, dataM := range internalData {
		channel, err := dataM.Convert()
		if err != nil {
			return err
		}
		badgerChannels := storage.Open(DbChannels)
		err = badgerChannels.Upsert(channel.Key(), channel)
		if err != nil {
			return err
		}
	}
	return nil
}
