package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/badgerhold"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"strconv"
	"time"
)

type SpotsUpdateRequest struct {
	S3Key string
	Month string
}

type Spot struct {
	Rating30               *float64   `json:"Rating30"`
	IsPrime                *int64     `json:"IsPrime"`
	FilmID                 *int64     `json:"FilmID"`
	FilmVersion            *string    `json:"FilmVersion"`
	FilmName               *string    `json:"FilmName"`
	FilmDur                *int64     `json:"FilmDur"`
	SpotPullRating         *float64   `json:"SpotPullRating"`
	DLDate                 *time.Time `json:"DLDate"`
	SptChnlPTR             *int64     `json:"sptChnlPTR"`
	CommInMplID            *int64     `json:"CommInMplID"`
	Positioning            *int64     `json:"Positioning"`
	AgrID                  *int64     `json:"AgrID"`
	MplID                  *int64     `json:"MplID"`
	OrdID                  *int64     `json:"OrdID"`
	AtpID                  *int64     `json:"AtpID"`
	DtpID                  *int64     `json:"DtpID"`
	TgrID                  *int64     `json:"TgrID"`
	SptDateL               *int64     `json:"SptDateL"`
	FloatPriority          *int64     `json:"FloatPriority"`
	CurrentAuctionBidValue *int64     `json:"CurrentAuctionBidValue"`
	SpotOrderNo            *int64     `json:"SpotOrderNo"`
	SpotBroadcastTime      *int64     `json:"SpotBroadcastTime"`
	SpotFactBroadcastTime  *int64     `json:"SpotFactBroadcastTime"`
	AllocationType         *int64     `json:"AllocationType"`
	FixPriority            *int64     `json:"FixPriority"`
	SpotReserve            *int64     `json:"SpotReserve"`
	RankID                 *int64     `json:"RankID"`
	SpotID                 *int64     `json:"SpotID"`
	BlockID                *int64     `json:"BlockID"`
	TNSSpotsID             *int64     `json:"TNSSpotsID"`
	TNSBlockID             *int64     `json:"TNSBlockID"`
	Rating                 *float64   `json:"Rating"`
	BaseRating             *float64   `json:"BaseRating"`
	OTS                    *float64   `json:"OTS"`
	IRating                *float64   `json:"IRating"`
	IBaseRating            *float64   `json:"IBaseRating"`
	IsHumanBeing           *bool      `json:"IsHumanBeing"`
	Timestamp              time.Time  `json:"Timestamp"`
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
		SptChnlPTR:             utils.Int64I(s.S["sptChnlPTR"]),
		CommInMplID:            utils.Int64I(s.S["CommInMplID"]),
		Positioning:            utils.Int64I(s.S["Positioning"]),
		AgrID:                  utils.Int64I(s.S["AgrID"]),
		MplID:                  utils.Int64I(s.S["MplID"]),
		OrdID:                  utils.Int64I(s.S["OrdID"]),
		AtpID:                  utils.Int64I(s.S["AtpID"]),
		DtpID:                  utils.Int64I(s.S["DtpID"]),
		TgrID:                  utils.Int64I(s.S["TgrID"]),
		SptDateL:               utils.Int64I(s.S["SptDateL"]),
		FloatPriority:          utils.Int64I(s.S["FloatPriority"]),
		CurrentAuctionBidValue: utils.Int64I(s.S["CurrentAuctionBidValue"]),
		SpotOrderNo:            utils.Int64I(s.S["SpotOrderNo"]),
		SpotBroadcastTime:      utils.Int64I(s.S["SpotBroadcastTime"]),
		SpotFactBroadcastTime:  utils.Int64I(s.S["SpotFactBroadcastTime"]),
		AllocationType:         utils.Int64I(s.S["AllocationType"]),
		FixPriority:            utils.Int64I(s.S["FixPriority"]),
		SpotReserve:            utils.Int64I(s.S["SpotReserve"]),
		RankID:                 utils.Int64I(s.S["RankID"]),
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
				Month: bodyJson.Month,
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
	/*OrdBlocks*/
	convertDataOrderBlock, err := resp.Convert("OrdBlocks")
	if err != nil {
		return err
	}
	marshalDataOrderBlock, err := json.Marshal(convertDataOrderBlock)
	if err != nil {
		return err
	}
	var orderBlocks []internalObl
	err = json.Unmarshal(marshalDataOrderBlock, &orderBlocks)
	if err != nil {
		return err
	}
	badgerSpotsOrderBlock := storage.Open(DbSpotsOrderBlock)
	for _, dataO := range orderBlocks {
		spot, err := dataO.ConvertOrderBlock()
		if err != nil {
			return err
		}
		var orders []SpotOrderBlock
		err = badgerSpotsOrderBlock.Find(&orders, badgerhold.Where("OrdID").Eq(*spot.OrdID))
		if err != nil {
			return err
		}
		for _, item := range orders {
			err = badgerSpotsOrderBlock.Delete(item.Key(), item)
			if err != nil {
				return err
			}
		}
		err = badgerSpotsOrderBlock.Upsert(spot.Key(), spot)
		if err != nil {
			return err
		}
	}
	/*Spots*/
	convertData, err := resp.Convert("SpotList")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var spots []internalS
	err = json.Unmarshal(marshalData, &spots)
	if err != nil {
		return err
	}
	badgerSpots := storage.Open(DbSpots)
	month, err := strconv.Atoi(request.Month)
	if err != nil {
		return err
	}
	for _, dataS := range spots {
		spot, err := dataS.ConvertSpot()
		if err != nil {
			return err
		}
		if spot.SptChnlPTR == nil {
			continue
		}
		/*load from networks*/
		networksQuery := ProgramBreaksBadgerQuery{}
		var networks []ProgramBreaks
		filterNetwork := badgerhold.Where("CnlID").Eq(*spot.SptChnlPTR).And("Month").Eq(int64(month))
		err = networksQuery.Find(&networks, filterNetwork)
		if err != nil {
			return err
		}
		if networks != nil {
			for _, network := range networks {
				if spot.BlockID == nil {
					*spot.IsPrime = 0
					continue
				}
				spot.DLDate = network.DLDate
				spot.IsPrime = network.IsPrime
			}
		}
		/*load from mediaplans*/
		mediaplanQuery := MediaplanBadgerQuery{}
		var mediaplans []Mediaplan
		filterMediplan := badgerhold.Where("ComplimentId").Eq(*spot.CommInMplID)
		err = mediaplanQuery.Find(&mediaplans, filterMediplan)
		if err != nil {
			return err
		}
		if spot.CommInMplID == nil {
			continue
		}
		if mediaplans == nil {
			continue
		}
		var inventoryDur float64
		if mediaplans[0].InventoryUnitDuration == nil {
			inventoryDur = 30.0
		} else {
			inventoryDur = float64(*mediaplans[0].InventoryUnitDuration)
		}
		spot.FilmID = mediaplans[0].FilmID
		spot.FilmName = mediaplans[0].FilmName
		spot.FilmVersion = mediaplans[0].FilmVersion
		spot.FilmDur = mediaplans[0].FilmDur
		if mediaplans[0].FilmDur != nil {
			if orderBlocks != nil {
				for _, block := range orderBlocks {
					orderBlock, _ := block.ConvertOrderBlock()
					rate := *orderBlock.Rate
					spotRating := float64(*mediaplans[0].FilmDur) * rate / inventoryDur
					spot.Rating30 = &spotRating
					spot.SpotPullRating = &rate
				}
			}
		}
		var spotItems []Spot
		err = badgerSpots.Find(&spotItems, badgerhold.Where("SpotID").Eq(*spot.SpotID))
		if err != nil {
			return err
		}
		for _, item := range spotItems {
			err = badgerSpots.Delete(item.Key(), item)
			if err != nil {
				return err
			}
		}
		err = badgerSpots.Upsert(spot.Key(), spot)
		if err != nil {
			return err
		}
	}
	return nil
}
