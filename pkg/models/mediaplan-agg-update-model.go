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
	query := MediaplanBadgerQuery{}
	var mediaplans []Mediaplan
	filter := badgerhold.Where("MediaplanId").Eq(request.MediaplanId)
	err := query.Find(&mediaplans, filter)
	if err != nil {
		return err
	}
	/*load from budgets*/
	budgetQuery := BudgetsBadgerQuery{}
	var budgets []Budget
	filterBudgets := badgerhold.Where("Month").Eq(request.Month)
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
	filterSpots := badgerhold.Where("AgrID").Eq(request.AgreementId)
	err = spotQuery.Find(&spots, filterSpots)
	if err != nil {
		return err
	}
	badgerAggMediaplans := storage.Open(DbAggMediaplans)
	for _, mediaplan := range mediaplans {
		for _, budget := range budgets {
			for _, channel := range channels {
				if *budget.Month == *mediaplan.Month && *mediaplan.MplCnlID == *channel.ID {
					aggMediaplan := &MediaplanAgg{
						AdtID:                   mediaplan.AdtID,
						AdtName:                 mediaplan.AdtName,
						AgreementId:             mediaplan.AgreementId,
						AllocationType:          mediaplan.AllocationType,
						AmountFact:              mediaplan.AmountFact,
						AmountPlan:              mediaplan.AmountPlan,
						BrandName:               mediaplan.BrandName,
						Budget:                  budget.Budget,
						ChannelId:               mediaplan.ChannelId,
						ChannelName:             budget.CnlName,
						CppOffPrime:             mediaplan.CPPoffprime,
						CppOffPrimeWithDiscount: nil,
						CppPrime:                mediaplan.CPPprime,
						CppPrimeWithDiscount:    nil,
						DealChannelStatus:       mediaplan.DealChannelStatus,
						FactOff:                 nil,
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
			}
		}
	}
	return nil
}
