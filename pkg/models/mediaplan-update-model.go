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

const MediaplanTable = "mediaplans"

type internalMDiscount struct {
	Item map[string]interface{} `json:"item"`
}

type Mediaplan struct {
	AdtID                 *int                    `json:"AdtID"`
	AdtName               *string                 `json:"AdtName"`
	AdvID                 *int                    `json:"AdvID"`
	AgrID                 *int                    `json:"AgrID"`
	AgrName               *string                 `json:"AgrName"`
	AllocationType        *string                 `json:"AllocationType"`
	AmountFact            *float64                `json:"AmountFact"`
	AmountPlan            *float64                `json:"AmountPlan"`
	BrandID               *int                    `json:"BrandID"`
	BrandName             *string                 `json:"BrandName"`
	CPPoffprime           *float64                `json:"CPPoffprime"`
	CPPprime              *float64                `json:"CPPprime"`
	ComplimentId          *int                    `json:"CommInMplID"`
	ContractBeg           *int                    `json:"ContractBeg"`
	ContractEnd           *int                    `json:"ContractEnd"`
	DateFrom              *int                    `json:"DateFrom"`
	DateTo                *int                    `json:"DateTo"`
	DealChannelStatus     *int                    `json:"DealChannelStatus"`
	Discounts             []MediaplanDiscountItem `json:"Discounts"`
	DoubleAdvertiser      *bool                   `json:"DoubleAdvertiser"`
	DtpID                 *int                    `json:"DtpID"`
	DublSpot              *int                    `json:"DublSpot"`
	FbrName               *string                 `json:"FbrName"`
	FilmDur               *int                    `json:"FilmDur"`
	FilmDurKoef           *float64                `json:"FilmDurKoef"`
	FilmID                *int                    `json:"FilmID"`
	FilmName              *string                 `json:"FilmName"`
	FilmVersion           *string                 `json:"FilmVersion"`
	FixPriceAsFloat       *int                    `json:"FixPriceAsFloat"`
	FixPriority           *int                    `json:"FixPriority"`
	GRP                   *float64                `json:"GRP"`
	GRPShift              *float64                `json:"GRPShift"`
	GrpFact               *float64                `json:"GrpFact"`
	GrpPlan               *float64                `json:"GrpPlan"`
	GrpTotal              *float64                `json:"GrpTotal"`
	GrpTotalPrime         *float64                `json:"GrpTotalPrime"`
	HasReserve            *int                    `json:"HasReserve"`
	InventoryUnitDuration *int                    `json:"InventoryUnitDuration"`
	MplCbrID              *int                    `json:"MplCbrID"`
	MplCbrName            *string                 `json:"MplCbrName"`
	MplCnlID              *int                    `json:"MplCnlID"`
	MplID                 *int                    `json:"MplID"`
	MplMonth              *string                 `json:"MplMonth"`
	MplName               *string                 `json:"MplName"`
	MplState              *int                    `json:"MplState"`
	Multiple              *int                    `json:"Multiple"`
	OBDPos                *int                    `json:"OBDPos"`
	OrdFrID               *int                    `json:"OrdFrID"`
	OrdID                 *int                    `json:"OrdID"`
	OrdIsTriggered        *int                    `json:"OrdIsTriggered"`
	OrdName               *string                 `json:"OrdName"`
	PBACond               *string                 `json:"PBACond"`
	PBAObjID              *int                    `json:"PBAObjID"`
	ProdClassID           *int                    `json:"ProdClassID"`
	SellingDirection      *int                    `json:"SellingDirection"`
	SplitMessageGroupID   *int                    `json:"SplitMessageGroupID"`
	SumShift              *float64                `json:"SumShift"`
	TPName                *string                 `json:"TPName"`
	TgrID                 *string                 `json:"TgrID"`
	TgrName               *string                 `json:"TgrName"`
	Timestamp             *time.Time              `json:"Timestamp"`
	AdvertiserId          *int                    `json:"AdvertiserId"`
	AgreementId           *int                    `json:"AgreementId"`
	ChannelId             *int                    `json:"ChannelId"`
	FfoaAllocated         *int                    `json:"ffoaAllocated"`
	FfoaLawAcc            *int                    `json:"ffoaLawAcc"`
	FilmId                *int                    `json:"FilmId"`
	MediaplanId           *int                    `json:"MediaplanId"`
	Month                 *int                    `json:"Month"`
	OrdBegDate            *int                    `json:"ordBegDate"`
	OrdEndDate            *int                    `json:"ordEndDate"`
	OrdManager            *string                 `json:"ordManager"`
}
type MediaplanDiscountItem struct {
	DiscountFactor          *float64   `json:"DiscountFactor"`
	TypeID                  *int       `json:"TypeID"`
	IsManual                *bool      `json:"IsManual"`
	DicountEndDate          *time.Time `json:"DicountEndDate"`
	IsSpotPositionDependent *bool      `json:"IsSpotPositionDependent"`
	DiscountTypeName        *string    `json:"DiscountTypeName"`
	DiscountStartDate       *time.Time `json:"DiscountStartDate"`
	ValueID                 *int       `json:"ValueID"`
	ApplicableToDeals       *bool      `json:"ApplicableToDeals"`
	ApplyingTypeID          *int       `json:"ApplyingTypeID"`
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
		TypeID:                  utils.IntI(m.Item["TypeID"]),
		IsManual:                utils.BoolI(m.Item["IsManual"]),
		DicountEndDate:          endTime,
		IsSpotPositionDependent: utils.BoolI(m.Item["IsSpotPositionDependent"]),
		DiscountTypeName:        utils.StringI(m.Item["DiscountTypeName"]),
		AggregationMethodName:   utils.StringI(m.Item["AggregationMethodName"]),
		DiscountStartDate:       startTime,
		ValueID:                 utils.IntI(m.Item["ValueID"]),
		ApplicableToDeals:       utils.BoolI(m.Item["ApplicableToDeals"]),
		ApplyingTypeID:          utils.IntI(m.Item["ApplyingTypeID"]),
		IsDiscountAggregate:     utils.BoolI(m.Item["IsDiscountAggregate"]),
	}
	return item, nil
}

func (m *internalM) ConvertMediaplan() (*Mediaplan, error) {
	timestamp := time.Now()
	month := utils.IntI(m.M["MplMonth"])
	channelId := utils.IntI(m.M["MplCnlID"])
	mediaplanId := utils.IntI(m.M["MplID"])
	advertiserId := utils.IntI(m.M["AdtID"])
	filmId := utils.IntI(m.M["FilmID"])
	agreementId := utils.IntI(m.M["AgrID"])
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
		AdtID:                 utils.IntI(m.M["AdtID"]),
		AdtName:               utils.StringI(m.M["AdtName"]),
		AdvID:                 utils.IntI(m.M["AdvID"]),
		AgrID:                 utils.IntI(m.M["AgrID"]),
		AgrName:               utils.StringI(m.M["AgrName"]),
		AllocationType:        utils.StringI(m.M["AllocationType"]),
		AmountFact:            utils.FloatI(m.M["AmountFact"]),
		AmountPlan:            utils.FloatI(m.M["AmountPlan"]),
		BrandID:               utils.IntI(m.M["BrandID"]),
		BrandName:             utils.StringI(m.M["BrandName"]),
		CPPoffprime:           utils.FloatI(m.M["CPPoffprime"]),
		CPPprime:              utils.FloatI(m.M["CPPprime"]),
		ComplimentId:          utils.IntI(m.M["CommInMplID"]),
		ContractBeg:           utils.IntI(m.M["ContractBeg"]),
		ContractEnd:           utils.IntI(m.M["ContractEnd"]),
		DateFrom:              utils.IntI(m.M["DateFrom"]),
		DateTo:                utils.IntI(m.M["DateTo"]),
		DealChannelStatus:     utils.IntI(m.M["DealChannelStatus"]),
		Discounts:             discounts,
		DoubleAdvertiser:      utils.BoolI(m.M["DoubleAdvertiser"]),
		DtpID:                 utils.IntI(m.M["DtpID"]),
		DublSpot:              utils.IntI(m.M["DublSpot"]),
		FbrName:               utils.StringI(m.M["FbrName"]),
		FilmDur:               utils.IntI(m.M["FilmDur"]),
		FilmDurKoef:           utils.FloatI(m.M["FilmDurKoef"]),
		FilmID:                utils.IntI(m.M["FilmID"]),
		FilmName:              utils.StringI(m.M["FilmName"]),
		FilmVersion:           utils.StringI(m.M["FilmVersion"]),
		FixPriceAsFloat:       utils.IntI(m.M["FixPriceAsFloat"]),
		FixPriority:           utils.IntI(m.M["FixPriority"]),
		GRP:                   utils.FloatI(m.M["GRP"]),
		GRPShift:              utils.FloatI(m.M["GRPShift"]),
		GrpFact:               utils.FloatI(m.M["GRPShift"]),
		GrpPlan:               utils.FloatI(m.M["GrpPlan"]),
		GrpTotal:              utils.FloatI(m.M["GrpTotal"]),
		GrpTotalPrime:         utils.FloatI(m.M["GrpTotalPrime"]),
		HasReserve:            utils.IntI(m.M["HasReserve"]),
		InventoryUnitDuration: utils.IntI(m.M["InventoryUnitDuration"]),
		MplCbrID:              utils.IntI(m.M["MplCbrID"]),
		MplCbrName:            utils.StringI(m.M["MplCbrName"]),
		MplCnlID:              utils.IntI(m.M["MplCnlID"]),
		MplID:                 utils.IntI(m.M["MplID"]),
		MplMonth:              utils.StringI(m.M["MplMonth"]),
		MplName:               utils.StringI(m.M["MplName"]),
		MplState:              utils.IntI(m.M["MplState"]),
		Multiple:              utils.IntI(m.M["Multiple"]),
		OBDPos:                utils.IntI(m.M["OBDPos"]),
		OrdFrID:               utils.IntI(m.M["OrdFrID"]),
		OrdID:                 utils.IntI(m.M["OrdID"]),
		OrdIsTriggered:        utils.IntI(m.M["OrdIsTriggered"]),
		OrdName:               utils.StringI(m.M["OrdName"]),
		PBACond:               utils.StringI(m.M["PBACond"]),
		PBAObjID:              utils.IntI(m.M["PBAObjID"]),
		ProdClassID:           utils.IntI(m.M["ProdClassID"]),
		SellingDirection:      utils.IntI(m.M["SellingDirection"]),
		SplitMessageGroupID:   utils.IntI(m.M["SplitMessageGroupID"]),
		SumShift:              utils.FloatI(m.M["SumShift"]),
		TPName:                utils.StringI(m.M["TPName"]),
		TgrID:                 utils.StringI(m.M["TgrID"]),
		TgrName:               utils.StringI(m.M["TgrName"]),
		Timestamp:             &timestamp,
		AdvertiserId:          advertiserId,
		AgreementId:           agreementId,
		ChannelId:             channelId,
		FfoaAllocated:         utils.IntI(m.M["ffoaAllocated"]),
		FfoaLawAcc:            utils.IntI(m.M["ffoaLawAcc"]),
		FilmId:                filmId,
		MediaplanId:           mediaplanId,
		Month:                 month,
		OrdBegDate:            utils.IntI(m.M["ordBegDate"]),
		OrdEndDate:            utils.IntI(m.M["ordEndDate"]),
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
	badgerMediaplans := storage.Open(DbMediaplans)
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
		err = badgerMediaplans.Upsert(mediaplan.Key(), mediaplan)
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
	badgerMediaplans := storage.Open(DbMediaplans)
	defer badgerMediaplans.Close()
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
		err = badgerMediaplans.Upsert(mediaplan.Key(), mediaplan)
		if err != nil {
			return err
		}
	}
	return nil
}
