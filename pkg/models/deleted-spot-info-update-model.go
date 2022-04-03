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

type DeletedSpotInfoUpdateRequest struct {
	S3Key string
}

type DeletedSpotInfo struct {
	AgrID                  *int       `json:"AgrID"`
	OrdID                  *int       `json:"OrdID"`
	FilmID                 *int       `json:"FilmID"`
	FilmDur                *int       `json:"FilmDur"`
	BlockDate              *int       `json:"BlockDate"`
	BlockTime              *int       `json:"BlockTime"`
	Position               *int       `json:"Position"`
	BlockNumber            *int       `json:"BlockNumber"`
	Reason                 *int       `json:"Reason"`
	AffiliationType        *int       `json:"AffiliationType"`
	CurrentAuctionBidValue *int       `json:"CurrentAuctionBidValue"`
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
		AgrID:                  utils.IntI(i.I["AgrID"]),
		OrdID:                  utils.IntI(i.I["OrdID"]),
		FilmID:                 utils.IntI(i.I["FilmID"]),
		FilmDur:                utils.IntI(i.I["FilmDur"]),
		BlockDate:              utils.IntI(i.I["BlockDate"]),
		BlockTime:              utils.IntI(i.I["BlockTime"]),
		Position:               utils.IntI(i.I["Position"]),
		BlockNumber:            utils.IntI(i.I["BlockNumber"]),
		Reason:                 utils.IntI(i.I["Reason"]),
		AffiliationType:        utils.IntI(i.I["AffiliationType"]),
		CurrentAuctionBidValue: utils.IntI(i.I["CurrentAuctionBidValue"]),
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
	badgerDeletedSpotInfo := storage.Open(DbDeletedSpotInfo)
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
