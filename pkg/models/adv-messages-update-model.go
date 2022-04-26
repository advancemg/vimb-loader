package models

import (
	"encoding/json"
	"fmt"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type AdvertiserUpdateRequest struct {
	S3Key string
}

type Advertiser struct {
	AdtID                 *string   `json:"AdtID"`
	GroupID               *string   `json:"GroupID"`
	AdtName               *string   `json:"AdtName"`
	FilmName              *string   `json:"FilmName"`
	FilmVersion           *string   `json:"FilmVersion"`
	BrandName             *string   `json:"BrandName"`
	AspectName            *string   `json:"AspectName"`
	FilmID                *int64    `json:"FilmID"`
	FilmDur               *int64    `json:"FilmDur"`
	BrandID               *int64    `json:"BrandID"`
	ProdClassID           *int64    `json:"ProdClassID"`
	FfoaAllocated         *int64    `json:"FfoaAllocated"`
	FfoaLawAcc            *int64    `json:"FfoaLawAcc"`
	AspectID              *int64    `json:"AspectID"`
	DoubleAdvertiser      *bool     `json:"DoubleAdvertiser"`
	HasSpots              *bool     `json:"HasSpots"`
	HasBroadcastMaterials *bool     `json:"HasBroadcastMaterials"`
	HasPreviewMaterials   *bool     `json:"HasPreviewMaterials"`
	Timestamp             time.Time `json:"Timestamp"`
}

func (adv *Advertiser) Key() string {
	return fmt.Sprintf("%d", *adv.FilmID)
}

func (r *internalRow) ConvertAdvertiser() (*Advertiser, error) {
	timestamp := time.Now()
	rows := &Advertiser{
		AdtID:                 utils.StringI(r.Row["AdtID"]),
		GroupID:               utils.StringI(r.Row["GroupID"]),
		AdtName:               utils.StringI(r.Row["AdtName"]),
		FilmName:              utils.StringI(r.Row["FilmName"]),
		FilmVersion:           utils.StringI(r.Row["FilmVersion"]),
		BrandName:             utils.StringI(r.Row["BrandName"]),
		AspectName:            utils.StringI(r.Row["AspectName"]),
		FilmID:                utils.Int64I(r.Row["FilmID"]),
		FilmDur:               utils.Int64I(r.Row["FilmDur"]),
		BrandID:               utils.Int64I(r.Row["BrandID"]),
		ProdClassID:           utils.Int64I(r.Row["ProdClassID"]),
		FfoaAllocated:         utils.Int64I(r.Row["FfoaAllocated"]),
		FfoaLawAcc:            utils.Int64I(r.Row["FfoaLawAcc"]),
		AspectID:              utils.Int64I(r.Row["AspectID"]),
		DoubleAdvertiser:      utils.BoolI(r.Row["DoubleAdvertiser"]),
		HasSpots:              utils.BoolI(r.Row["HasSpots"]),
		HasBroadcastMaterials: utils.BoolI(r.Row["HasBroadcastMaterials"]),
		HasPreviewMaterials:   utils.BoolI(r.Row["HasPreviewMaterials"]),
		Timestamp:             timestamp,
	}
	return rows, nil
}

func AdvertiserStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := AdvMessagesUpdateQueue
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
			req := AdvertiserUpdateRequest{
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

func (request *AdvertiserUpdateRequest) Update() error {
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

func (request *AdvertiserUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	convertData, err := resp.Convert("AdvertisingMessagesData")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalRow
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	badgerAdvertiser := storage.Open(DbAdvertisers)
	for _, dataRow := range internalData {
		advertiser, err := dataRow.ConvertAdvertiser()
		if err != nil {
			return err
		}
		err = badgerAdvertiser.Upsert(advertiser.Key(), advertiser)
		if err != nil {
			return err
		}
	}
	return nil
}
