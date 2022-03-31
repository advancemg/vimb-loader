package models

import (
	"encoding/json"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"reflect"
	"time"
)

const MediaplanTable = "mediaplans"
const MediaplanAggTable = "mediaplans-agg"

type internalM struct {
	M map[string]interface{} `json:"m"`
}

type internalMDiscount struct {
	Item map[string]interface{} `json:"item"`
}

type Mediaplan struct {
	AdtID                 int                     `json:"AdtID"`
	AdtName               string                  `json:"AdtName"`
	AdvID                 int                     `json:"AdvID"`
	AgrID                 int                     `json:"AgrID"`
	AgrName               string                  `json:"AgrName"`
	AllocationType        string                  `json:"AllocationType"`
	AmountFact            float64                 `json:"AmountFact"`
	AmountPlan            float64                 `json:"AmountPlan"`
	BrandID               int                     `json:"BrandID"`
	BrandName             string                  `json:"BrandName"`
	CPPoffprime           float64                 `json:"CPPoffprime"`
	CPPprime              float64                 `json:"CPPprime"`
	ComplimentId          int                     `json:"CommInMplID"`
	ContractBeg           int                     `json:"ContractBeg"`
	ContractEnd           int                     `json:"ContractEnd"`
	DateFrom              int                     `json:"DateFrom"`
	DateTo                int                     `json:"DateTo"`
	DealChannelStatus     int                     `json:"DealChannelStatus"`
	Discounts             []MediaplanDiscountItem `json:"Discounts"`
	DoubleAdvertiser      string                  `json:"DoubleAdvertiser"`
	DtpID                 string                  `json:"DtpID"`
	DublSpot              string                  `json:"DublSpot"`
	FbrName               string                  `json:"FbrName"`
	FilmDur               string                  `json:"FilmDur"`
	FilmDurKoef           string                  `json:"FilmDurKoef"`
	FilmID                int                     `json:"FilmID"`
	FilmName              string                  `json:"FilmName"`
	FilmVersion           string                  `json:"FilmVersion"`
	FixPriceAsFloat       string                  `json:"FixPriceAsFloat"`
	FixPriority           string                  `json:"FixPriority"`
	GRP                   string                  `json:"GRP"`
	GRPShift              string                  `json:"GRPShift"`
	GrpFact               string                  `json:"GrpFact"`
	GrpPlan               string                  `json:"GrpPlan"`
	GrpTotal              string                  `json:"GrpTotal"`
	GrpTotalPrime         string                  `json:"GrpTotalPrime"`
	HasReserve            string                  `json:"HasReserve"`
	InventoryUnitDuration string                  `json:"InventoryUnitDuration"`
	MplCbrID              string                  `json:"MplCbrID"`
	MplCbrName            string                  `json:"MplCbrName"`
	MplCnlID              string                  `json:"MplCnlID"`
	MplID                 string                  `json:"MplID"`
	MplMonth              string                  `json:"MplMonth"`
	MplName               string                  `json:"MplName"`
	MplState              string                  `json:"MplState"`
	Multiple              string                  `json:"Multiple"`
	OBDPos                string                  `json:"OBDPos"`
	OrdFrID               string                  `json:"OrdFrID"`
	OrdID                 string                  `json:"OrdID"`
	OrdIsTriggered        string                  `json:"OrdIsTriggered"`
	OrdName               string                  `json:"OrdName"`
	PBACond               string                  `json:"PBACond"`
	PBAObjID              string                  `json:"PBAObjID"`
	ProdClassID           string                  `json:"ProdClassID"`
	SellingDirection      int                     `json:"SellingDirection"`
	SplitMessageGroupID   string                  `json:"SplitMessageGroupID"`
	SumShift              string                  `json:"SumShift"`
	TPName                string                  `json:"TPName"`
	TgrID                 string                  `json:"TgrID"`
	TgrName               string                  `json:"TgrName"`
	Timestamp             time.Time               `json:"Timestamp"`
	AdvertiserId          int                     `json:"advertiserId"`
	AgreementId           int                     `json:"agreementId"`
	ChannelId             int                     `json:"channelId"`
	FfoaAllocated         string                  `json:"ffoaAllocated"`
	FfoaLawAcc            string                  `json:"ffoaLawAcc"`
	FilmId                int                     `json:"filmId"`
	MediaplanId           int                     `json:"mediaplanId"`
	Month                 int                     `json:"month"`
	OrdBegDate            string                  `json:"ordBegDate"`
	OrdEndDate            string                  `json:"ordEndDate"`
	OrdManager            string                  `json:"ordManager"`
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
type MediaplanAgg struct {
	AdtID                   int        `json:"AdtID"`
	AdtName                 string     `json:"AdtName"`
	AgreementId             int        `json:"AgreementId"`
	AllocationType          int        `json:"AllocationType"`
	AmountFact              float64    `json:"AmountFact"`
	AmountPlan              float64    `json:"AmountPlan"`
	BrandName               string     `json:"BrandName"`
	Budget                  float64    `json:"Budget"`
	ChannelId               int        `json:"ChannelId"`
	ChannelName             string     `json:"ChannelName"`
	CppOffPrime             float64    `json:"CppOffPrime"`
	CppOffPrimeWithDiscount float64    `json:"CppOffPrimeWithDiscount"`
	CppPrime                float64    `json:"CppPrime"`
	CppPrimeWithDiscount    float64    `json:"CppPrimeWithDiscount"`
	DealChannelStatus       string     `json:"DealChannelStatus"`
	FactOff                 float64    `json:"FactOff"`
	FactPrime               float64    `json:"FactPrime"`
	FixPriority             int        `json:"FixPriority"`
	GrpPlan                 string     `json:"GrpPlan"`
	GrpTotal                float64    `json:"GrpTotal"`
	InventoryUnitDuration   int        `json:"InventoryUnitDuration"`
	MediaplanState          int        `json:"MediaplanState"`
	MplID                   int        `json:"MplID"`
	MplMonth                int        `json:"MplMonth"`
	MplName                 string     `json:"MplName"`
	SpotsPrimePercent       float64    `json:"SpotsPrimePercent"`
	SuperFix                string     `json:"SuperFix"`
	UpdateDate              time.Time  `json:"UpdateDate"`
	Timestamp               time.Time  `json:"Timestamp"`
	UserGrpPlan             string     `json:"UserGrpPlan"`
	WeeksInfo               []WeekInfo `json:"WeeksInfo"`
	BcpCentralID            string     `json:"bcpCentralID"`
	BcpName                 string     `json:"bcpName"`
}

type MediaplanUpdateRequest struct {
	S3Key string
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

func (m *internalM) Convert() (*Mediaplan, *MediaplanAgg, error) {
	timestamp := time.Now()
	month := *utils.IntI(m.M["MplMonth"])
	channelId := *utils.IntI(m.M["MplCnlID"])
	mediaplanId := *utils.IntI(m.M["MplID"])
	advertiserId := *utils.IntI(m.M["AdtID"])
	filmId := *utils.IntI(m.M["FilmID"])
	agreementId := *utils.IntI(m.M["AgrID"])
	/*discounts*/
	var discounts []MediaplanDiscountItem
	if _, ok := m.M["Discounts"]; ok {
		marshalData, err := json.Marshal(m.M["Discounts"])
		if err != nil {
			return nil, nil, err
		}
		switch reflect.TypeOf(m.M["Discounts"]).Kind() {
		case reflect.Array, reflect.Slice:
			var internalDiscountData []internalMDiscount
			err = json.Unmarshal(marshalData, &internalDiscountData)
			if err != nil {
				return nil, nil, err
			}
			for _, discountItem := range internalDiscountData {
				discount, err := discountItem.Convert()
				if err != nil {
					return nil, nil, err
				}
				discounts = append(discounts, *discount)
			}
		case reflect.Map, reflect.Struct:
			var internalDiscountData internalMDiscount
			err = json.Unmarshal(marshalData, &internalDiscountData)
			if err != nil {
				return nil, nil, err
			}
			discount, err := internalDiscountData.Convert()
			if err != nil {
				return nil, nil, err
			}
			discounts = append(discounts, *discount)
		}
	}
	mediaplan := &Mediaplan{
		AdtID:                 utils.Int(m.M["AdtID"].(string)),
		AdtName:               m.M["AdtName"].(string),
		AdvID:                 utils.Int(m.M["AdvID"].(string)),
		AgrID:                 utils.Int(m.M["AgrID"].(string)),
		AgrName:               m.M["AgrName"].(string),
		AllocationType:        m.M["AllocationType"].(string),
		AmountFact:            utils.Float(m.M["AmountFact"].(string)),
		AmountPlan:            utils.Float(m.M["AmountPlan"].(string)),
		BrandID:               utils.Int(m.M["BrandID"].(string)),
		BrandName:             m.M["BrandName"].(string),
		CPPoffprime:           utils.Float(m.M["CPPoffprime"].(string)),
		CPPprime:              utils.Float(m.M["CPPprime"].(string)),
		ComplimentId:          utils.Int(m.M["CommInMplID"].(string)),
		ContractBeg:           utils.Int(m.M["ContractBeg"].(string)),
		ContractEnd:           utils.Int(m.M["ContractEnd"].(string)),
		DateFrom:              utils.Int(m.M["DateFrom"].(string)),
		DateTo:                utils.Int(m.M["DateTo"].(string)),
		DealChannelStatus:     utils.Int(m.M["DealChannelStatus"].(string)),
		Discounts:             discounts,
		DoubleAdvertiser:      "",
		DtpID:                 "",
		DublSpot:              "",
		FbrName:               "",
		FilmDur:               "",
		FilmDurKoef:           "",
		FilmID:                utils.Int(m.M["FilmID"].(string)),
		FilmName:              "",
		FilmVersion:           "",
		FixPriceAsFloat:       "",
		FixPriority:           "",
		GRP:                   "",
		GRPShift:              "",
		GrpFact:               "",
		GrpPlan:               "",
		GrpTotal:              "",
		GrpTotalPrime:         "",
		HasReserve:            "",
		InventoryUnitDuration: "",
		MplCbrID:              "",
		MplCbrName:            "",
		MplCnlID:              "",
		MplID:                 "",
		MplMonth:              "",
		MplName:               "",
		MplState:              "",
		Multiple:              "",
		OBDPos:                "",
		OrdFrID:               "",
		OrdID:                 "",
		OrdIsTriggered:        "",
		OrdName:               "",
		PBACond:               "",
		PBAObjID:              "",
		ProdClassID:           "",
		SellingDirection:      0,
		SplitMessageGroupID:   "",
		SumShift:              "",
		TPName:                "",
		TgrID:                 "",
		TgrName:               "",
		Timestamp:             timestamp,
		AdvertiserId:          advertiserId,
		AgreementId:           0,
		ChannelId:             channelId,
		FfoaAllocated:         "",
		FfoaLawAcc:            "",
		FilmId:                filmId,
		MediaplanId:           mediaplanId,
		Month:                 month,
		OrdBegDate:            "",
		OrdEndDate:            "",
		OrdManager:            "",
	}
	mediaplanAgg := &MediaplanAgg{
		Timestamp:   timestamp,
		MplMonth:    month,
		ChannelId:   channelId,
		MplID:       mediaplanId,
		AdtID:       advertiserId,
		AgreementId: agreementId,
	}
	return mediaplan, mediaplanAgg, nil
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
	_, err = resp.Convert("MPlansList")
	if err != nil {
		return err
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
	for _, dataM := range internalData {
		mediaplan, agg, err := dataM.Convert()
		if err != nil {
			return err
		}
		println(mediaplan, agg)
	}
	return nil
}
