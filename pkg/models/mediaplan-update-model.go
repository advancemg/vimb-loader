package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/usecase"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"reflect"
	"time"
)

type internalMDiscount struct {
	Item map[string]interface{} `json:"item"`
}

type Mediaplan struct {
	AdtID                 *int64                  `json:"AdtID"`
	AdtName               *string                 `json:"AdtName"`
	AdvID                 *int64                  `json:"AdvID"`
	AgrID                 *int64                  `json:"AgrID"`
	AgrName               *string                 `json:"AgrName"`
	AllocationType        *int64                  `json:"AllocationType"`
	AmountFact            *float64                `json:"AmountFact"`
	AmountPlan            *float64                `json:"AmountPlan"`
	BrandID               *int64                  `json:"BrandID"`
	BrandName             *string                 `json:"BrandName"`
	CPPoffprime           *float64                `json:"CPPoffprime"`
	CPPprime              *float64                `json:"CPPprime"`
	ComplimentId          *int64                  `json:"CommInMplID"`
	ContractBeg           *int64                  `json:"ContractBeg"`
	ContractEnd           *int64                  `json:"ContractEnd"`
	DateFrom              *int64                  `json:"DateFrom"`
	DateTo                *int64                  `json:"DateTo"`
	DealChannelStatus     *int64                  `json:"DealChannelStatus"`
	Discounts             []MediaplanDiscountItem `json:"Discounts"`
	DoubleAdvertiser      *bool                   `json:"DoubleAdvertiser"`
	DtpID                 *int64                  `json:"DtpID"`
	DublSpot              *int64                  `json:"DublSpot"`
	FbrName               *string                 `json:"FbrName"`
	FilmDur               *int64                  `json:"FilmDur"`
	FilmDurKoef           *float64                `json:"FilmDurKoef"`
	FilmID                *int64                  `json:"FilmID"`
	FilmName              *string                 `json:"FilmName"`
	FilmVersion           *string                 `json:"FilmVersion"`
	FixPriceAsFloat       *int64                  `json:"FixPriceAsFloat"`
	FixPriority           *int64                  `json:"FixPriority"`
	GRP                   *float64                `json:"GRP"`
	GRPShift              *float64                `json:"GRPShift"`
	GrpFact               *float64                `json:"GrpFact"`
	GrpPlan               *float64                `json:"GrpPlan"`
	GrpTotal              *float64                `json:"GrpTotal"`
	GrpTotalPrime         *float64                `json:"GrpTotalPrime"`
	HasReserve            *int64                  `json:"HasReserve"`
	InventoryUnitDuration *int64                  `json:"InventoryUnitDuration"`
	MplCbrID              *int64                  `json:"MplCbrID"`
	MplCbrName            *string                 `json:"MplCbrName"`
	MplCnlID              *int64                  `json:"MplCnlID"`
	MplID                 *int64                  `json:"MplID"`
	MplMonth              *string                 `json:"MplMonth"`
	MplName               *string                 `json:"MplName"`
	MplState              *int64                  `json:"MplState"`
	Multiple              *int64                  `json:"Multiple"`
	OBDPos                *int64                  `json:"OBDPos"`
	OrdFrID               *int64                  `json:"OrdFrID"`
	OrdID                 *int64                  `json:"OrdID"`
	OrdIsTriggered        *int64                  `json:"OrdIsTriggered"`
	OrdName               *string                 `json:"OrdName"`
	PBACond               *string                 `json:"PBACond"`
	PBAObjID              *int64                  `json:"PBAObjID"`
	ProdClassID           *int64                  `json:"ProdClassID"`
	SellingDirection      *int64                  `json:"SellingDirection"`
	SplitMessageGroupID   *int64                  `json:"SplitMessageGroupID"`
	SumShift              *float64                `json:"SumShift"`
	TPName                *string                 `json:"TPName"`
	TgrID                 *string                 `json:"TgrID"`
	TgrName               *string                 `json:"TgrName"`
	Timestamp             *time.Time              `json:"Timestamp"`
	AdvertiserId          *int64                  `json:"AdvertiserId"`
	AgreementId           *int64                  `json:"AgreementId"`
	ChannelId             *int64                  `json:"ChannelId"`
	FfoaAllocated         *int64                  `json:"ffoaAllocated"`
	FfoaLawAcc            *int64                  `json:"ffoaLawAcc"`
	FilmId                *int64                  `json:"FilmId"`
	MediaplanId           *int64                  `json:"MediaplanId"`
	Month                 *int64                  `json:"Month"`
	OrdBegDate            *int64                  `json:"ordBegDate"`
	OrdEndDate            *int64                  `json:"ordEndDate"`
	OrdManager            *string                 `json:"ordManager"`
}
type MediaplanDiscountItem struct {
	DiscountFactor          *float64   `json:"DiscountFactor"`
	TypeID                  *int64     `json:"TypeID"`
	IsManual                *bool      `json:"IsManual"`
	DicountEndDate          *time.Time `json:"DicountEndDate"`
	IsSpotPositionDependent *bool      `json:"IsSpotPositionDependent"`
	DiscountTypeName        *string    `json:"DiscountTypeName"`
	DiscountStartDate       *time.Time `json:"DiscountStartDate"`
	ValueID                 *int64     `json:"ValueID"`
	ApplicableToDeals       *bool      `json:"ApplicableToDeals"`
	ApplyingTypeID          *int64     `json:"ApplyingTypeID"`
	IsDiscountAggregate     *bool      `json:"IsDiscountAggregate"`
	AggregationMethodName   *string    `json:"AggregationMethodName"`
}

type MediaplanUpdateRequest struct {
	S3Key string
}

func (mediaplan *Mediaplan) Key() string {
	return fmt.Sprintf("%d-%d", *mediaplan.MediaplanId, *mediaplan.ComplimentId)
}
func (mediaplan *Mediaplan) AggKey() string {
	return fmt.Sprintf("%d-%d-%d-%d-%d", *mediaplan.MediaplanId, *mediaplan.AdvertiserId, *mediaplan.ChannelId, *mediaplan.AgreementId, *mediaplan.Month)
}

func (m *internalMDiscount) Convert() (*MediaplanDiscountItem, error) {
	startTime := utils.TimeI(m.Item["DiscountStartDate"], `2006-01-02T15:04:05`)
	endTime := utils.TimeI(m.Item["DicountEndDate"], `2006-01-02T15:04:05`)
	item := &MediaplanDiscountItem{
		DiscountFactor:          utils.FloatI(m.Item["DiscountFactor"]),
		TypeID:                  utils.Int64I(m.Item["TypeID"]),
		IsManual:                utils.BoolI(m.Item["IsManual"]),
		DicountEndDate:          endTime,
		IsSpotPositionDependent: utils.BoolI(m.Item["IsSpotPositionDependent"]),
		DiscountTypeName:        utils.StringI(m.Item["DiscountTypeName"]),
		AggregationMethodName:   utils.StringI(m.Item["AggregationMethodName"]),
		DiscountStartDate:       startTime,
		ValueID:                 utils.Int64I(m.Item["ValueID"]),
		ApplicableToDeals:       utils.BoolI(m.Item["ApplicableToDeals"]),
		ApplyingTypeID:          utils.Int64I(m.Item["ApplyingTypeID"]),
		IsDiscountAggregate:     utils.BoolI(m.Item["IsDiscountAggregate"]),
	}
	return item, nil
}

func (m *internalM) ConvertMediaplan() (*Mediaplan, error) {
	timestamp := time.Now()
	month := utils.Int64I(m.M["MplMonth"])
	channelId := utils.Int64I(m.M["MplCnlID"])
	mediaplanId := utils.Int64I(m.M["MplID"])
	advertiserId := utils.Int64I(m.M["AdtID"])
	filmId := utils.Int64I(m.M["FilmID"])
	agreementId := utils.Int64I(m.M["AgrID"])
	var discounts []MediaplanDiscountItem
	if _, ok := m.M["Discounts"]; ok {
		marshalData, err := json.Marshal(m.M["Discounts"])
		if err != nil {
			return nil, err
		}
		switch reflect.TypeOf(m.M["Discounts"]).Kind() {
		case reflect.Array, reflect.Slice:
			var internalDiscountData []internalMDiscount
			err = json.Unmarshal(marshalData, &internalDiscountData)
			if err != nil {
				return nil, err
			}
			for _, discountItem := range internalDiscountData {
				discount, err := discountItem.Convert()
				if err != nil {
					return nil, err
				}
				discounts = append(discounts, *discount)
			}
		case reflect.Map, reflect.Struct:
			var internalDiscountData internalMDiscount
			err = json.Unmarshal(marshalData, &internalDiscountData)
			if err != nil {
				return nil, err
			}
			discount, err := internalDiscountData.Convert()
			if err != nil {
				return nil, err
			}
			discounts = append(discounts, *discount)
		}
	}
	mediaplan := &Mediaplan{
		AdtID:                 utils.Int64I(m.M["AdtID"]),
		AdtName:               utils.StringI(m.M["AdtName"]),
		AdvID:                 utils.Int64I(m.M["AdvID"]),
		AgrID:                 utils.Int64I(m.M["AgrID"]),
		AgrName:               utils.StringI(m.M["AgrName"]),
		AllocationType:        utils.Int64I(m.M["AllocationType"]),
		AmountFact:            utils.FloatI(m.M["AmountFact"]),
		AmountPlan:            utils.FloatI(m.M["AmountPlan"]),
		BrandID:               utils.Int64I(m.M["BrandID"]),
		BrandName:             utils.StringI(m.M["BrandName"]),
		CPPoffprime:           utils.FloatI(m.M["CPPoffprime"]),
		CPPprime:              utils.FloatI(m.M["CPPprime"]),
		ComplimentId:          utils.Int64I(m.M["CommInMplID"]),
		ContractBeg:           utils.Int64I(m.M["ContractBeg"]),
		ContractEnd:           utils.Int64I(m.M["ContractEnd"]),
		DateFrom:              utils.Int64I(m.M["DateFrom"]),
		DateTo:                utils.Int64I(m.M["DateTo"]),
		DealChannelStatus:     utils.Int64I(m.M["DealChannelStatus"]),
		Discounts:             discounts,
		DoubleAdvertiser:      utils.BoolI(m.M["DoubleAdvertiser"]),
		DtpID:                 utils.Int64I(m.M["DtpID"]),
		DublSpot:              utils.Int64I(m.M["DublSpot"]),
		FbrName:               utils.StringI(m.M["FbrName"]),
		FilmDur:               utils.Int64I(m.M["FilmDur"]),
		FilmDurKoef:           utils.FloatI(m.M["FilmDurKoef"]),
		FilmID:                utils.Int64I(m.M["FilmID"]),
		FilmName:              utils.StringI(m.M["FilmName"]),
		FilmVersion:           utils.StringI(m.M["FilmVersion"]),
		FixPriceAsFloat:       utils.Int64I(m.M["FixPriceAsFloat"]),
		FixPriority:           utils.Int64I(m.M["FixPriority"]),
		GRP:                   utils.FloatI(m.M["GRP"]),
		GRPShift:              utils.FloatI(m.M["GRPShift"]),
		GrpFact:               utils.FloatI(m.M["GRPShift"]),
		GrpPlan:               utils.FloatI(m.M["GrpPlan"]),
		GrpTotal:              utils.FloatI(m.M["GrpTotal"]),
		GrpTotalPrime:         utils.FloatI(m.M["GrpTotalPrime"]),
		HasReserve:            utils.Int64I(m.M["HasReserve"]),
		InventoryUnitDuration: utils.Int64I(m.M["InventoryUnitDuration"]),
		MplCbrID:              utils.Int64I(m.M["MplCbrID"]),
		MplCbrName:            utils.StringI(m.M["MplCbrName"]),
		MplCnlID:              utils.Int64I(m.M["MplCnlID"]),
		MplID:                 utils.Int64I(m.M["MplID"]),
		MplMonth:              utils.StringI(m.M["MplMonth"]),
		MplName:               utils.StringI(m.M["MplName"]),
		MplState:              utils.Int64I(m.M["MplState"]),
		Multiple:              utils.Int64I(m.M["Multiple"]),
		OBDPos:                utils.Int64I(m.M["OBDPos"]),
		OrdFrID:               utils.Int64I(m.M["OrdFrID"]),
		OrdID:                 utils.Int64I(m.M["OrdID"]),
		OrdIsTriggered:        utils.Int64I(m.M["OrdIsTriggered"]),
		OrdName:               utils.StringI(m.M["OrdName"]),
		PBACond:               utils.StringI(m.M["PBACond"]),
		PBAObjID:              utils.Int64I(m.M["PBAObjID"]),
		ProdClassID:           utils.Int64I(m.M["ProdClassID"]),
		SellingDirection:      utils.Int64I(m.M["SellingDirection"]),
		SplitMessageGroupID:   utils.Int64I(m.M["SplitMessageGroupID"]),
		SumShift:              utils.FloatI(m.M["SumShift"]),
		TPName:                utils.StringI(m.M["TPName"]),
		TgrID:                 utils.StringI(m.M["TgrID"]),
		TgrName:               utils.StringI(m.M["TgrName"]),
		Timestamp:             &timestamp,
		AdvertiserId:          advertiserId,
		AgreementId:           agreementId,
		ChannelId:             channelId,
		FfoaAllocated:         utils.Int64I(m.M["ffoaAllocated"]),
		FfoaLawAcc:            utils.Int64I(m.M["ffoaLawAcc"]),
		FilmId:                filmId,
		MediaplanId:           mediaplanId,
		Month:                 month,
		OrdBegDate:            utils.Int64I(m.M["ordBegDate"]),
		OrdEndDate:            utils.Int64I(m.M["ordEndDate"]),
		OrdManager:            utils.StringI(m.M["ordManager"]),
	}
	return mediaplan, nil
}

func MediaplanStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := MPLansUpdateQueue
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
			req := MediaplanUpdateRequest{
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

func (request *MediaplanUpdateRequest) Update() error {
	var err error
	filePath, err := s3.Download(request.S3Key)
	if err != nil {
		return err
	}
	resp := utils.VimbResponse{FilePath: filePath}
	convertData, err := resp.Convert("MPlansList")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalM
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	db, table := utils.SplitDbAndTable(DbMediaplans)
	dbMediaplans := usecase.OpenDb(db, table)
	aggTasks := map[string]MediaplanAggUpdateRequest{}
	for _, dataM := range internalData {
		mediaplan, err := dataM.ConvertMediaplan()
		if err != nil {
			return err
		}
		id := mediaplan.MediaplanId
		advertiserId := mediaplan.AdvertiserId
		channelId := mediaplan.ChannelId
		agreementId := mediaplan.AgreementId
		month := mediaplan.Month
		aggTasks[mediaplan.AggKey()] = MediaplanAggUpdateRequest{
			Month:        *month,
			ChannelId:    *channelId,
			MediaplanId:  *id,
			AdvertiserId: *advertiserId,
			AgreementId:  *agreementId,
		}
		err = dbMediaplans.AddOrUpdate(mediaplan.Key(), mediaplan)
		if err != nil {
			return err
		}
	}
	amqpConfig := mq_broker.InitConfig()
	for _, aggMessage := range aggTasks {
		err := amqpConfig.PublishJson(MediaplanAggUpdateQueue, aggMessage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (request *MediaplanUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	convertData, err := resp.Convert("MPlansList")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalM
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	aggTasks := map[string]MediaplanAggUpdateRequest{}
	db, table := utils.SplitDbAndTable(DbMediaplans)
	dbMediaplans := usecase.OpenDb(db, table)
	for _, dataM := range internalData {
		mediaplan, err := dataM.ConvertMediaplan()
		if err != nil {
			return err
		}
		id := mediaplan.MediaplanId
		advertiserId := mediaplan.AdvertiserId
		channelId := mediaplan.ChannelId
		agreementId := mediaplan.AgreementId
		month := mediaplan.Month
		key := fmt.Sprintf("%d-%d-%d-%d-%d", id, advertiserId, channelId, agreementId, month)
		aggTasks[key] = MediaplanAggUpdateRequest{
			Month:        *month,
			ChannelId:    *channelId,
			MediaplanId:  *id,
			AdvertiserId: *advertiserId,
			AgreementId:  *agreementId,
		}
		var mediaplans []Mediaplan
		if mediaplan.MediaplanId != nil {
			err = dbMediaplans.FindWhereEq(&mediaplans, "MediaplanId", *mediaplan.MediaplanId)
			if err != nil {
				return err
			}
		}
		for _, item := range mediaplans {
			err = dbMediaplans.Delete(item.Key(), item)
			if err != nil {
				return err
			}
		}
		err = dbMediaplans.AddOrUpdate(mediaplan.Key(), mediaplan)
		if err != nil {
			return err
		}
	}
	return nil
}
