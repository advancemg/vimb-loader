package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/store"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type MediaplanAggUpdateRequest struct {
	Month        int64 `json:"month" bson:"month"`
	ChannelId    int64 `json:"channelId" bson:"channelId"`
	MediaplanId  int64 `json:"mediaplanId" bson:"mediaplanId"`
	AdvertiserId int64 `json:"advertiserId" bson:"advertiserId"`
	AgreementId  int64 `json:"agreementId" bson:"agreementId"`
}

type MediaplanAgg struct {
	AdtID                   *int64     `json:"AdtID" bson:"AdtID"`
	AdtName                 *string    `json:"AdtName" bson:"AdtName"`
	AgreementId             *int64     `json:"AgreementId" bson:"AgreementId"`
	AllocationType          *int64     `json:"AllocationType" bson:"AllocationType"`
	AmountFact              *float64   `json:"AmountFact" bson:"AmountFact"`
	AmountPlan              *float64   `json:"AmountPlan" bson:"AmountPlan"`
	BrandName               *string    `json:"BrandName" bson:"BrandName"`
	Budget                  *float64   `json:"Budget" bson:"Budget"`
	ChannelId               *int64     `json:"ChannelId" bson:"ChannelId"`
	ChannelName             *string    `json:"ChannelName" bson:"ChannelName"`
	CppOffPrime             *float64   `json:"CppOffPrime" bson:"CppOffPrime"`
	CppOffPrimeWithDiscount *float64   `json:"CppOffPrimeWithDiscount" bson:"CppOffPrimeWithDiscount"`
	CppPrime                *float64   `json:"CppPrime" bson:"CppPrime"`
	CppPrimeWithDiscount    *float64   `json:"CppPrimeWithDiscount" bson:"CppPrimeWithDiscount"`
	DealChannelStatus       *int64     `json:"DealChannelStatus" bson:"DealChannelStatus"`
	FactOff                 *float64   `json:"FactOff" bson:"FactOff"`
	FactPrime               *float64   `json:"FactPrime" bson:"FactPrime"`
	FixPriority             *int64     `json:"FixPriority" bson:"FixPriority"`
	GrpPlan                 *float64   `json:"GrpPlan" bson:"GrpPlan"`
	GrpTotal                *float64   `json:"GrpTotal" bson:"GrpTotal"`
	InventoryUnitDuration   *int64     `json:"InventoryUnitDuration" bson:"InventoryUnitDuration"`
	MediaplanState          *int64     `json:"MediaplanState" bson:"MediaplanState"`
	MplID                   *int64     `json:"MplID" bson:"MplID"`
	MplMonth                *int64     `json:"MplMonth" bson:"MplMonth"`
	MplName                 *string    `json:"MplName" bson:"MplName"`
	SpotsPrimePercent       *float64   `json:"SpotsPrimePercent" bson:"SpotsPrimePercent"`
	SuperFix                *string    `json:"SuperFix" bson:"SuperFix"`
	UpdateDate              *time.Time `json:"UpdateDate" bson:"UpdateDate"`
	UserGrpPlan             *string    `json:"UserGrpPlan" bson:"UserGrpPlan"`
	WeeksInfo               []WeekInfo `json:"WeeksInfo" bson:"WeeksInfo"`
	BcpCentralID            *int64     `json:"bcpCentralID" bson:"bcpCentralID"`
	BcpName                 *string    `json:"bcpName" bson:"bcpName"`
}

func (mediaplan *MediaplanAgg) Key() string {
	return fmt.Sprintf("%d", *mediaplan.MplID)
}

func MediaplanAggStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := MediaplanAggUpdateQueue
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
			var bodyJson MediaplanAggUpdateRequest
			err = json.Unmarshal(msg.Body, &bodyJson)
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
	db, table := utils.SplitDbAndTable(DbMediaplans)
	dbMediaplan, err := store.OpenDb(db, table)
	if err != nil {
		return err
	}
	var mediaplans []Mediaplan
	err = dbMediaplan.FindWhereEq(&mediaplans, "MediaplanId", request.MediaplanId)
	if err != nil {
		return err
	}
	db, table = utils.SplitDbAndTable(DbAggMediaplans)
	dbAggMediaplans, err := store.OpenDb(db, table)
	if err != nil {
		return err
	}
	for _, mediaplan := range mediaplans {
		var cppOffPrimeWithDiscount float64
		var cppPrimeWithDiscount float64
		var discountFactor float64
		var cppOffPrime float64
		var cppPrime float64
		/*load from budgets*/
		var budget float64
		var dealChannelStatus int64
		db, table = utils.SplitDbAndTable(DbBudgets)
		dbBudgets, err := store.OpenDb(db, table)
		if err != nil {
			return err
		}
		var budgets []Budget
		err = dbBudgets.FindWhereAnd4Eq(&budgets,
			"Month", request.Month,
			"CnlID", *mediaplan.ChannelId,
			"AdtID", *mediaplan.AdtID,
			"AgrID", *mediaplan.AgreementId,
		)
		if err != nil {
			return err
		}
		if budgets != nil {
			budget = *budgets[0].Budget
			dealChannelStatus = *budgets[0].DealChannelStatus
		}
		/*load from channels*/
		var channelName string
		var bcpName string
		var bcpCentralID int64
		db, table = utils.SplitDbAndTable(DbChannels)
		dbChannels, err := store.OpenDb(db, table)
		if err != nil {
			return err
		}
		var channels []Channel
		if mediaplan.ChannelId != nil {
			err = dbChannels.FindWhereEq(&channels, "ID", *mediaplan.ChannelId)
			if err != nil {
				return err
			}
		}
		if channels != nil {
			channelName = *channels[0].ShortName
			bcpName = *channels[0].BcpName
			bcpCentralID = *channels[0].BcpCentralID
		}
		/*load from spots*/
		db, table = utils.SplitDbAndTable(DbSpots)
		dbSpots, err := store.OpenDb(db, table)
		if err != nil {
			return err
		}
		var spots []Spot
		if mediaplan.MplID != nil {
			err = dbSpots.FindWhereEq(&spots, "MplID", *mediaplan.MplID)
			if err != nil {
				return err
			}
		}
		/*update plans*/
		if len(mediaplan.Discounts) > 0 {
			if len(mediaplan.Discounts) == 1 {
				discountFactor = *mediaplan.Discounts[0].DiscountFactor
				cppOffPrime = *mediaplan.CPPoffprime
				cppPrime = *mediaplan.CPPprime
				if discountFactor != 0 && cppOffPrime != 0 {
					cppPrimeWithDiscount = cppPrime * discountFactor
					if *mediaplan.CPPoffprime != 0 {
						cppOffPrimeWithDiscount = cppOffPrime * discountFactor
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
		/*Spot facts*/
		var primeFactRating float64
		var offFactRating float64
		var primePercent float64
		var totalSpotsRatingSum float64
		var primeTimeRatingSum float64
		if spots != nil {
			for _, spot := range spots {
				var isPrime int64
				if (spot.DLDate == nil) || (time.Now().Unix() < spot.DLDate.Unix()) {
					break
				}
				totalSpotsRatingSum += *spot.Rating30
				if spot.IsPrime == nil {
					isPrime = 0
					continue
				}
				if *spot.IsPrime == 1 {
					/*prime*/
					primeTimeRatingSum += *spot.Rating30
					primeFactRating += *spot.Rating30
				}
				if isPrime == 0 {
					/*off*/
					offFactRating += *spot.Rating30
				}
			}
		}
		if totalSpotsRatingSum > 0 {
			primePercent = primeTimeRatingSum / totalSpotsRatingSum
		}
		/*get open/close weeks*/
		var resultWeekStatus []WeekInfo
		mondays, err := utils.GetWeekDayByYearMonth(*mediaplan.Month)
		if err != nil {
			return err
		}
		nowInt := int(time.Now().Weekday())
		for monday, dateTime := range mondays {
			if nowInt > monday {
				var weekItem WeekInfo
				weekItem.Number = monday
				weekItem.Date = dateTime
				weekItem.Close = true
				resultWeekStatus = append(resultWeekStatus, weekItem)
			} else {
				var weekItem WeekInfo
				weekItem.Number = monday
				weekItem.Date = dateTime
				weekItem.Close = false
				resultWeekStatus = append(resultWeekStatus, weekItem)
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
			Budget:                  &budget,
			ChannelId:               &request.ChannelId,
			ChannelName:             &channelName,
			CppOffPrime:             mediaplan.CPPoffprime,
			CppOffPrimeWithDiscount: &cppOffPrimeWithDiscount,
			CppPrimeWithDiscount:    &cppPrimeWithDiscount,
			CppPrime:                mediaplan.CPPprime,
			DealChannelStatus:       &dealChannelStatus,
			FactOff:                 &offFactRating,
			FactPrime:               &primeFactRating,
			FixPriority:             mediaplan.FixPriority,
			GrpPlan:                 mediaplan.GRP,
			GrpTotal:                mediaplan.GrpTotal,
			InventoryUnitDuration:   mediaplan.InventoryUnitDuration,
			MediaplanState:          mediaplan.MplState,
			MplID:                   &request.MediaplanId,
			MplMonth:                &request.Month,
			MplName:                 mediaplan.MplName,
			SpotsPrimePercent:       &primePercent,
			SuperFix:                nil,
			UpdateDate:              &timestamp,
			UserGrpPlan:             nil,
			WeeksInfo:               resultWeekStatus,
			BcpCentralID:            &bcpCentralID,
			BcpName:                 &bcpName,
		}
		err = dbAggMediaplans.AddOrUpdate(aggMediaplan.Key(), aggMediaplan)
		if err != nil {
			return err
		}
	}
	return nil
}
