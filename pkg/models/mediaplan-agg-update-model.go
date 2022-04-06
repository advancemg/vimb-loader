package models

import (
	"encoding/json"
	"fmt"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/timshannon/badgerhold"
	"time"
)

const MediaplanAggTable = "mediaplans-agg"

type MediaplanAggUpdateRequest struct {
	Month        int `json:"month"`
	ChannelId    int `json:"channelId"`
	MediaplanId  int `json:"mediaplanId"`
	AdvertiserId int `json:"advertiserId"`
	AgreementId  int `json:"agreementId"`
}

type MediaplanAgg struct {
	AdtID                   *int       `json:"AdtID"`
	AdtName                 *string    `json:"AdtName"`
	AgreementId             *int       `json:"AgreementId"`
	AllocationType          *string    `json:"AllocationType"`
	AmountFact              *float64   `json:"AmountFact"`
	AmountPlan              *float64   `json:"AmountPlan"`
	BrandName               *string    `json:"BrandName"`
	Budget                  *float64   `json:"Budget"`
	ChannelId               *int       `json:"ChannelId"`
	ChannelName             *string    `json:"ChannelName"`
	CppOffPrime             *float64   `json:"CppOffPrime"`
	CppOffPrimeWithDiscount *float64   `json:"CppOffPrimeWithDiscount"`
	CppPrime                *float64   `json:"CppPrime"`
	CppPrimeWithDiscount    *float64   `json:"CppPrimeWithDiscount"`
	DealChannelStatus       *int       `json:"DealChannelStatus"`
	FactOff                 *float64   `json:"FactOff"`
	FactPrime               *float64   `json:"FactPrime"`
	FixPriority             *int       `json:"FixPriority"`
	GrpPlan                 *float64   `json:"GrpPlan"`
	GrpTotal                *float64   `json:"GrpTotal"`
	InventoryUnitDuration   *int       `json:"InventoryUnitDuration"`
	MediaplanState          *int       `json:"MediaplanState"`
	MplID                   *int       `json:"MplID"`
	MplMonth                *int       `json:"MplMonth"`
	MplName                 *string    `json:"MplName"`
	SpotsPrimePercent       *float64   `json:"SpotsPrimePercent"`
	SuperFix                *string    `json:"SuperFix"`
	UpdateDate              *time.Time `json:"UpdateDate"`
	Timestamp               *time.Time `json:"Timestamp"`
	UserGrpPlan             *string    `json:"UserGrpPlan"`
	WeeksInfo               []WeekInfo `json:"WeeksInfo"`
	BcpCentralID            *int       `json:"bcpCentralID"`
	BcpName                 *string    `json:"bcpName"`
}

func (mediaplan *MediaplanAgg) Key() string {
	return fmt.Sprintf("%d", *mediaplan.MplID)
}

func MediaplanAggStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := MediaplanAggUpdateQueue
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
			var bodyJson MediaplanAggUpdateRequest
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			err = bodyJson.Update()
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (request *MediaplanAggUpdateRequest) Update() error {
	timestamp := time.Now()
	/*create week info*/
	/*load from mediaplans*/
	mediaplanQuery := MediaplanBadgerQuery{}
	var mediaplans []Mediaplan
	filter := badgerhold.Where("MediaplanId").Eq(request.MediaplanId)
	err := mediaplanQuery.Find(&mediaplans, filter)
	if err != nil {
		return err
	}
	/*load from budgets*/
	budgetQuery := BudgetsBadgerQuery{}
	var budgets []Budget
	filterBudgets := badgerhold.Where("Month").Eq(request.Month).
		And("CnlID").Eq(request.ChannelId).
		And("AdtID").Eq(request.AdvertiserId).
		And("AgrID").Eq(request.AgreementId)
	err = budgetQuery.Find(&budgets, filterBudgets)
	if err != nil {
		return err
	}
	/*load from channels*/
	channelQuery := ChannelBadgerQuery{}
	var channels []Channel
	filterChannels := badgerhold.Where("ID").Eq(request.ChannelId)
	err = channelQuery.Find(&channels, filterChannels)
	if err != nil {
		return err
	}
	/*load from spots*/
	spotQuery := ChannelBadgerQuery{}
	var spots []Spot
	filterSpots := badgerhold.Where("MplID").Eq(request.MediaplanId)
	err = spotQuery.Find(&spots, filterSpots)
	if err != nil {
		return err
	}
	var budget Budget
	var channel Channel
	var spot Spot
	if len(spots) == 1 {
		spot = spots[0]
	}
	if len(budgets) == 1 {
		budget = budgets[0]
	}
	if len(channels) == 1 {
		channel = channels[0]
	}
	badgerAggMediaplans := storage.Open(DbAggMediaplans)
	for _, mediaplan := range mediaplans {
		var cppOffPrimeWithDiscount float64
		var cppPrimeWithDiscount float64
		var discountFactor float64
		if len(mediaplan.Discounts) > 0 {
			if len(mediaplan.Discounts) == 1 {
				discountFactor = *mediaplan.Discounts[0].DiscountFactor
				if discountFactor != 0 && *mediaplan.CPPprime != 0 {
					cppOffPrimeWithDiscount = *mediaplan.CPPprime * discountFactor
					if *mediaplan.CPPoffprime != 0 {
						cppOffPrimeWithDiscount = *mediaplan.CPPoffprime * discountFactor
					}
				}
			}
			if len(mediaplan.Discounts) > 1 {
				discountFactor = 1.0000000000000
				for _, item := range mediaplan.Discounts {
					discountFactor *= *item.DiscountFactor
				}
				cppPrimeWithDiscount = discountFactor * *mediaplan.CPPprime
				cppOffPrimeWithDiscount = discountFactor * *mediaplan.CPPoffprime
			}
		}
		aggMediaplan := &MediaplanAgg{
			AdtID:                   &request.AdvertiserId,
			AdtName:                 mediaplan.AdtName,
			AgreementId:             &request.AgreementId,
			AllocationType:          mediaplan.AllocationType,
			AmountFact:              mediaplan.AmountFact,
			AmountPlan:              mediaplan.AmountPlan,
			BrandName:               mediaplan.BrandName,
			Budget:                  budget.Budget,
			ChannelId:               &request.ChannelId,
			ChannelName:             channel.ShortName,
			CppOffPrime:             mediaplan.CPPoffprime,
			CppOffPrimeWithDiscount: &cppOffPrimeWithDiscount,
			CppPrimeWithDiscount:    &cppPrimeWithDiscount,
			CppPrime:                mediaplan.CPPprime,
			DealChannelStatus:       budget.DealChannelStatus,
			FactOff:                 spot.BaseRating, //
			FactPrime:               nil,
			FixPriority:             mediaplan.FixPriority,
			GrpPlan:                 mediaplan.GrpPlan,
			GrpTotal:                mediaplan.GrpTotal,
			InventoryUnitDuration:   mediaplan.InventoryUnitDuration,
			MediaplanState:          mediaplan.MplState,
			MplID:                   &request.MediaplanId,
			MplMonth:                &request.Month,
			MplName:                 mediaplan.MplName,
			SpotsPrimePercent:       nil,
			SuperFix:                nil,
			Timestamp:               &timestamp,
			UserGrpPlan:             nil,
			WeeksInfo:               []WeekInfo{},
			BcpCentralID:            channel.BcpCentralID,
			BcpName:                 channel.BcpName,
		}
		err = badgerAggMediaplans.Upsert(aggMediaplan.Key(), aggMediaplan)
		if err != nil {
			return err
		}
	}
	return nil
}
