package models

import (
	"encoding/json"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

const MediaplanTable = "mediaplans"
const MediaplanAggTable = "mediaplans-agg"

type internalM struct {
	M map[string]interface{} `json:"m"`
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
	CPPoffprime           string                  `json:"CPPoffprime"`
	CPPprime              string                  `json:"CPPprime"`
	ComplimentId          int                     `json:"CommInMplID"`
	ContractBeg           string                  `json:"ContractBeg"`
	ContractEnd           string                  `json:"ContractEnd"`
	DateFrom              string                  `json:"DateFrom"`
	DateTo                string                  `json:"DateTo"`
	DealChannelStatus     string                  `json:"DealChannelStatus"`
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
	DiscountFactor          string      `json:"DiscountFactor"`
	TypeID                  string      `json:"TypeID"`
	IsManual                string      `json:"IsManual"`
	DicountEndDate          interface{} `json:"DicountEndDate"`
	IsSpotPositionDependent string      `json:"IsSpotPositionDependent"`
	DiscountTypeName        string      `json:"DiscountTypeName"`
	DiscountStartDate       interface{} `json:"DiscountStartDate"`
	ValueID                 interface{} `json:"ValueID"`
	ApplicableToDeals       string      `json:"ApplicableToDeals"`
	ApplyingTypeID          string      `json:"ApplyingTypeID"`
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

func (m *internalM) Convert() (*Mediaplan, *MediaplanAgg, error) {
	timestamp := time.Now()
	month := utils.Int(m.M["MplMonth"].(string))
	channelId := utils.Int(m.M["MplCnlID"].(string))
	mediaplanId := utils.Int(m.M["MplID"].(string))
	advertiserId := utils.Int(m.M["AdtID"].(string))
	filmId := utils.Int(m.M["FilmID"].(string))
	agreementId := utils.Int(m.M["AgrID"].(string))
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
		BrandName:             "",
		CPPoffprime:           "",
		CPPprime:              "",
		ComplimentId:          utils.Int(m.M["CommInMplID"].(string)),
		ContractBeg:           "",
		ContractEnd:           "",
		DateFrom:              "",
		DateTo:                "",
		DealChannelStatus:     "",
		Discounts:             nil,
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
	println(internalData)
	return nil
}
