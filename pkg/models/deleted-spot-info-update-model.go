package models

import (
	"encoding/json"
	"fmt"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage/badger"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type DeletedSpotInfoUpdateRequest struct {
	S3Key string
}

type DeletedSpotInfo struct {
	AgrID                  *int64     `json:"AgrID"`
	OrdID                  *int64     `json:"OrdID"`
	FilmID                 *int64     `json:"FilmID"`
	FilmDur                *int64     `json:"FilmDur"`
	BlockDate              *int64     `json:"BlockDate"`
	BlockTime              *int64     `json:"BlockTime"`
	Position               *int64     `json:"Position"`
	BlockNumber            *int64     `json:"BlockNumber"`
	Reason                 *int64     `json:"Reason"`
	AffiliationType        *int64     `json:"AffiliationType"`
	CurrentAuctionBidValue *int64     `json:"CurrentAuctionBidValue"`
	BlockID                *int64     `json:"BlockID"`
	SpotID                 *int64     `json:"SpotID"`
	AgrName                *string    `json:"AgrName"`
	CnlName                *string    `json:"CnlName"`
	PrgName                *string    `json:"PrgName"`
	OrdName                *string    `json:"OrdName"`
	FilmName               *string    `json:"FilmName"`
	FilmVersion            *string    `json:"FilmVersion"`
	DeleteDateTime         *time.Time `json:"DeleteDateTime"`
	Timestamp              time.Time  `json:"Timestamp"`
}

func (deletedSpotInfo *DeletedSpotInfo) Key() string {
	return fmt.Sprintf("%d", *deletedSpotInfo.SpotID)
}

func (i *internalI) ConvertDeletedSpotInfo() (*DeletedSpotInfo, error) {
	timestamp := time.Now()
	items := &DeletedSpotInfo{
		AgrID:                  utils.Int64I(i.I["AgrID"]),
		OrdID:                  utils.Int64I(i.I["OrdID"]),
		FilmID:                 utils.Int64I(i.I["FilmID"]),
		FilmDur:                utils.Int64I(i.I["FilmDur"]),
		BlockDate:              utils.Int64I(i.I["BlockDate"]),
		BlockTime:              utils.Int64I(i.I["BlockTime"]),
		Position:               utils.Int64I(i.I["Position"]),
		BlockNumber:            utils.Int64I(i.I["BlockNumber"]),
		Reason:                 utils.Int64I(i.I["Reason"]),
		AffiliationType:        utils.Int64I(i.I["AffiliationType"]),
		CurrentAuctionBidValue: utils.Int64I(i.I["CurrentAuctionBidValue"]),
		BlockID:                utils.Int64I(i.I["BlockID"]),
		SpotID:                 utils.Int64I(i.I["SpotID"]),
		AgrName:                utils.StringI(i.I["AgrName"]),
		CnlName:                utils.StringI(i.I["CnlName"]),
		PrgName:                utils.StringI(i.I["PrgName"]),
		OrdName:                utils.StringI(i.I["OrdName"]),
		FilmName:               utils.StringI(i.I["FilmName"]),
		FilmVersion:            utils.StringI(i.I["FilmVersion"]),
		DeleteDateTime:         utils.TimeI(i.I["DeleteDateTime"], `2006-01-02T15:04:05`),
		Timestamp:              timestamp,
	}
	return items, nil
}

func DeletedSpotInfoStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := DeletedSpotInfoUpdateQueue
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
			var bodyJson MqUpdateMessage
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			/*read from s3 by s3Key*/
			req := DeletedSpotInfoUpdateRequest{
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

func (request *DeletedSpotInfoUpdateRequest) Update() error {
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

func (request *DeletedSpotInfoUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	convertData, err := resp.Convert("Items")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalI
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	badgerDeletedSpotInfo := badger.Open(DbDeletedSpotInfo)
	for _, dataI := range internalData {
		deletedSpotInfo, err := dataI.ConvertDeletedSpotInfo()
		if err != nil {
			return err
		}
		err = badgerDeletedSpotInfo.Upsert(deletedSpotInfo.Key(), deletedSpotInfo)
		if err != nil {
			return err
		}
	}
	return nil
}
