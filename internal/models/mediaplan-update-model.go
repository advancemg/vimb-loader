package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/store"
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
	AdtID                 *int64                  `json:"AdtID" bson:"AdtID"`
	AdtName               *string                 `json:"AdtName" bson:"AdtName"`
	AdvID                 *int64                  `json:"AdvID" bson:"AdvID"`
	AgrID                 *int64                  `json:"AgrID" bson:"AgrID"`
	AgrName               *string                 `json:"AgrName" bson:"AgrName"`
	AllocationType        *int64                  `json:"AllocationType" bson:"AllocationType"`
	AmountFact            *float64                `json:"AmountFact" bson:"AmountFact"`
	AmountPlan            *float64                `json:"AmountPlan" bson:"AmountPlan"`
	BrandID               *int64                  `json:"BrandID" bson:"BrandID"`
	BrandName             *string                 `json:"BrandName" bson:"BrandName"`
	CPPoffprime           *float64                `json:"CPPoffprime" bson:"CPPoffprime"`
	CPPprime              *float64                `json:"CPPprime" bson:"CPPprime"`
	ComplimentId          *int64                  `json:"CommInMplID" bson:"ComplimentId"`
	ContractBeg           *int64                  `json:"ContractBeg" bson:"ContractBeg"`
	ContractEnd           *int64                  `json:"ContractEnd" bson:"ContractEnd"`
	DateFrom              *int64                  `json:"DateFrom" bson:"DateFrom"`
	DateTo                *int64                  `json:"DateTo" bson:"DateTo"`
	DealChannelStatus     *int64                  `json:"DealChannelStatus" bson:"DealChannelStatus"`
	Discounts             []MediaplanDiscountItem `json:"Discounts" bson:"Discounts"`
	DoubleAdvertiser      *bool                   `json:"DoubleAdvertiser" bson:"DoubleAdvertiser"`
	DtpID                 *int64                  `json:"DtpID" bson:"DtpID"`
	DublSpot              *int64                  `json:"DublSpot" bson:"DublSpot"`
	FbrName               *string                 `json:"FbrName" bson:"FbrName"`
	FilmDur               *int64                  `json:"FilmDur" bson:"FilmDur"`
	FilmDurKoef           *float64                `json:"FilmDurKoef" bson:"FilmDurKoef"`
	FilmID                *int64                  `json:"FilmID" bson:"FilmID"`
	FilmName              *string                 `json:"FilmName" bson:"FilmName"`
	FilmVersion           *string                 `json:"FilmVersion" bson:"FilmVersion"`
	FixPriceAsFloat       *int64                  `json:"FixPriceAsFloat" bson:"FixPriceAsFloat"`
	FixPriority           *int64                  `json:"FixPriority" bson:"FixPriority"`
	GRP                   *float64                `json:"GRP" bson:"GRP"`
	GRPShift              *float64                `json:"GRPShift" bson:"GRPShift"`
	GrpFact               *float64                `json:"GrpFact" bson:"GrpFact"`
	GrpPlan               *float64                `json:"GrpPlan" bson:"GrpPlan"`
	GrpTotal              *float64                `json:"GrpTotal" bson:"GrpTotal"`
	GrpTotalPrime         *float64                `json:"GrpTotalPrime" bson:"GrpTotalPrime"`
	HasReserve            *int64                  `json:"HasReserve" bson:"HasReserve"`
	InventoryUnitDuration *int64                  `json:"InventoryUnitDuration" bson:"InventoryUnitDuration"`
	MplCbrID              *int64                  `json:"MplCbrID" bson:"MplCbrID"`
	MplCbrName            *string                 `json:"MplCbrName" bson:"MplCbrName"`
	MplCnlID              *int64                  `json:"MplCnlID" bson:"MplCnlID"`
	MplID                 *int64                  `json:"MplID" bson:"MplID"`
	MplMonth              *string                 `json:"MplMonth" bson:"MplMonth"`
	MplName               *string                 `json:"MplName" bson:"MplName"`
	MplState              *int64                  `json:"MplState" bson:"MplState"`
	Multiple              *int64                  `json:"Multiple" bson:"Multiple"`
	OBDPos                *int64                  `json:"OBDPos" bson:"OBDPos"`
	OrdFrID               *int64                  `json:"OrdFrID" bson:"OrdFrID"`
	OrdID                 *int64                  `json:"OrdID" bson:"OrdID"`
	OrdIsTriggered        *int64                  `json:"OrdIsTriggered" bson:"OrdIsTriggered"`
	OrdName               *string                 `json:"OrdName" bson:"OrdName"`
	PBACond               *string                 `json:"PBACond" bson:"PBACond"`
	PBAObjID              *int64                  `json:"PBAObjID" bson:"PBAObjID"`
	ProdClassID           *int64                  `json:"ProdClassID" bson:"ProdClassID"`
	SellingDirection      *int64                  `json:"SellingDirection" bson:"SellingDirection"`
	SplitMessageGroupID   *int64                  `json:"SplitMessageGroupID" bson:"SplitMessageGroupID"`
	SumShift              *float64                `json:"SumShift" bson:"SumShift"`
	TPName                *string                 `json:"TPName" bson:"TPName"`
	TgrID                 *string                 `json:"TgrID" bson:"TgrID"`
	TgrName               *string                 `json:"TgrName" bson:"TgrName"`
	Timestamp             *time.Time              `json:"Timestamp" bson:"Timestamp"`
	AdvertiserId          *int64                  `json:"AdvertiserId" bson:"AdvertiserId"`
	AgreementId           *int64                  `json:"AgreementId" bson:"AgreementId"`
	ChannelId             *int64                  `json:"ChannelId" bson:"ChannelId"`
	FfoaAllocated         *int64                  `json:"ffoaAllocated" bson:"ffoaAllocated"`
	FfoaLawAcc            *int64                  `json:"ffoaLawAcc" bson:"ffoaLawAcc"`
	FilmId                *int64                  `json:"FilmId" bson:"FilmId"`
	MediaplanId           *int64                  `json:"MediaplanId" bson:"MediaplanId"`
	Month                 *int64                  `json:"Month" bson:"Month"`
	OrdBegDate            *int64                  `json:"ordBegDate" bson:"ordBegDate"`
	OrdEndDate            *int64                  `json:"ordEndDate" bson:"ordEndDate"`
	OrdManager            *string                 `json:"ordManager" bson:"ordManager"`
}
type MediaplanDiscountItem struct {
	DiscountFactor          *float64   `json:"DiscountFactor" bson:"DiscountFactor"`
	TypeID                  *int64     `json:"TypeID" bson:"TypeID"`
	IsManual                *bool      `json:"IsManual" bson:"IsManual"`
	DicountEndDate          *time.Time `json:"DicountEndDate" bson:"DicountEndDate"`
	IsSpotPositionDependent *bool      `json:"IsSpotPositionDependent" bson:"IsSpotPositionDependent"`
	DiscountTypeName        *string    `json:"DiscountTypeName" bson:"DiscountTypeName"`
	DiscountStartDate       *time.Time `json:"DiscountStartDate" bson:"DiscountStartDate"`
	ValueID                 *int64     `json:"ValueID" bson:"ValueID"`
	ApplicableToDeals       *bool      `json:"ApplicableToDeals" bson:"ApplicableToDeals"`
	ApplyingTypeID          *int64     `json:"ApplyingTypeID" bson:"ApplyingTypeID"`
	IsDiscountAggregate     *bool      `json:"IsDiscountAggregate" bson:"IsDiscountAggregate"`
	AggregationMethodName   *string    `json:"AggregationMethodName" bson:"AggregationMethodName"`
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
	dbMediaplans, err := store.OpenDb(db, table)
	if err != nil {
		return err
	}
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
	defer amqpConfig.Close()
	for _, aggMessage := range aggTasks {
		err = amqpConfig.PublishJson(MediaplanAggUpdateQueue, aggMessage)
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
	dbMediaplans, err := store.OpenDb(db, table)
	if err != nil {
		return err
	}
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
