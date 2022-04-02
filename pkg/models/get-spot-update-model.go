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

type SpotsUpdateRequest struct {
	S3Key string
}

type Spot struct {
	SptChnlPTR             *int      `json:"SptChnlPTR"`
	CommInMplID            *int      `json:"CommInMplID"`
	Positioning            *int      `json:"Positioning"`
	AgrID                  *int      `json:"AgrID"`
	MplID                  *int      `json:"MplID"`
	OrdID                  *int      `json:"OrdID"`
	AtpID                  *int      `json:"AtpID"`
	DtpID                  *int      `json:"DtpID"`
	TgrID                  *int      `json:"TgrID"`
	SptDateL               *int      `json:"SptDateL"`
	FloatPriority          *int      `json:"FloatPriority"`
	CurrentAuctionBidValue *int      `json:"CurrentAuctionBidValue"`
	SpotOrderNo            *int      `json:"SpotOrderNo"`
	SpotBroadcastTime      *int      `json:"SpotBroadcastTime"`
	SpotFactBroadcastTime  *int      `json:"SpotFactBroadcastTime"`
	AllocationType         *int      `json:"AllocationType"`
	FixPriority            *int      `json:"FixPriority"`
	SpotReserve            *int      `json:"SpotReserve"`
	RankID                 *int      `json:"RankID"`
	SpotID                 *int64    `json:"SpotID"`
	BlockID                *int64    `json:"BlockID"`
	TNSSpotsID             *int64    `json:"TNSSpotsID"`
	TNSBlockID             *int64    `json:"TNSBlockID"`
	Rating                 *float64  `json:"Rating"`
	BaseRating             *float64  `json:"BaseRating"`
	OTS                    *float64  `json:"OTS"`
	IRating                *float64  `json:"IRating"`
	IBaseRating            *float64  `json:"IBaseRating"`
	IsHumanBeing           *bool     `json:"IsHumanBeing"`
	Timestamp              time.Time `json:"Timestamp"`
}

type SpotOrderBlock struct {
	OrdID     *int      `json:"OrdID"`
	BlockID   *int64    `json:"BlockID"`
	Rate      *float64  `json:"Rate"`
	IRate     *float64  `json:"IRate"`
	Timestamp time.Time `json:"Timestamp"`
}

func (orderBlock *SpotOrderBlock) Key() string {
	return fmt.Sprintf("%d-%d", *orderBlock.OrdID, *orderBlock.BlockID)
}

func (spot *Spot) Key() string {
	return fmt.Sprintf("%d", *spot.SpotID)
}

func (s *internalS) ConvertSpot() (*Spot, error) {
	timestamp := time.Now()
	spot := &Spot{
		SptChnlPTR:             utils.IntI(s.S["SptChnlPTR"]),
		CommInMplID:            utils.IntI(s.S["CommInMplID"]),
		Positioning:            utils.IntI(s.S["Positioning"]),
		AgrID:                  utils.IntI(s.S["AgrID"]),
		MplID:                  utils.IntI(s.S["MplID"]),
		OrdID:                  utils.IntI(s.S["OrdID"]),
		AtpID:                  utils.IntI(s.S["AtpID"]),
		DtpID:                  utils.IntI(s.S["DtpID"]),
		TgrID:                  utils.IntI(s.S["TgrID"]),
		SptDateL:               utils.IntI(s.S["SptDateL"]),
		FloatPriority:          utils.IntI(s.S["FloatPriority"]),
		CurrentAuctionBidValue: utils.IntI(s.S["CurrentAuctionBidValue"]),
		SpotOrderNo:            utils.IntI(s.S["SpotOrderNo"]),
		SpotBroadcastTime:      utils.IntI(s.S["SpotBroadcastTime"]),
		SpotFactBroadcastTime:  utils.IntI(s.S["SpotFactBroadcastTime"]),
		AllocationType:         utils.IntI(s.S["AllocationType"]),
		FixPriority:            utils.IntI(s.S["FixPriority"]),
		SpotReserve:            utils.IntI(s.S["SpotReserve"]),
		RankID:                 utils.IntI(s.S["RankID"]),
		SpotID:                 utils.Int64I(s.S["SpotID"]),
		BlockID:                utils.Int64I(s.S["BlockID"]),
		TNSSpotsID:             utils.Int64I(s.S["TNSSpotsID"]),
		TNSBlockID:             utils.Int64I(s.S["TNSBlockID"]),
		Rating:                 utils.FloatI(s.S["BaseRating"]),
		BaseRating:             utils.FloatI(s.S["BaseRating"]),
		OTS:                    utils.FloatI(s.S["OTS"]),
		IRating:                utils.FloatI(s.S["IRating"]),
		IBaseRating:            utils.FloatI(s.S["IBaseRating"]),
		IsHumanBeing:           utils.BoolI(s.S["IsHumanBeing"]),
		Timestamp:              timestamp,
	}
	return spot, nil
}

func (o *internalObl) ConvertOrderBlock() (*SpotOrderBlock, error) {
	timestamp := time.Now()
	orderBlock := &SpotOrderBlock{
		OrdID:     utils.IntI(o.Obl["OrdID"]),
		BlockID:   utils.Int64I(o.Obl["BlockID"]),
		Rate:      utils.FloatI(o.Obl["Rate"]),
		IRate:     utils.FloatI(o.Obl["IRate"]),
		Timestamp: timestamp,
	}
	return orderBlock, nil
}

func SpotStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := SpotsUpdateQueue
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
			req := SpotsUpdateRequest{
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

func (request *SpotsUpdateRequest) Update() error {
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

func (request *SpotsUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	convertData, err := resp.Convert("SpotList")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalS
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	badgerSpots := storage.Open(DbSpots)
	for _, dataM := range internalData {
		spot, err := dataM.ConvertSpot()
		if err != nil {
			return err
		}
		err = badgerSpots.Upsert(spot.Key(), spot)
		if err != nil {
			return err
		}
	}
	convertDataOrderBlock, err := resp.Convert("OrdBlocks")
	if err != nil {
		return err
	}
	marshalDataOrderBlock, err := json.Marshal(convertDataOrderBlock)
	if err != nil {
		return err
	}
	var internalDataOrderBlock []internalObl
	err = json.Unmarshal(marshalDataOrderBlock, &internalDataOrderBlock)
	if err != nil {
		return err
	}
	badgerSpotsOrderBlock := storage.Open(DbSpotsOrderBlock)
	for _, dataM := range internalDataOrderBlock {
		spot, err := dataM.ConvertOrderBlock()
		if err != nil {
			return err
		}
		err = badgerSpotsOrderBlock.Upsert(spot.Key(), spot)
		if err != nil {
			return err
		}
	}
	return nil
}
