package storage

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"github.com/timshannon/badgerhold"
	"testing"
	"time"
)

type mplTest struct {
	AdtID                 *int       `json:"AdtID"`
	AdtName               *string    `json:"AdtName"`
	AdvID                 *int       `json:"AdvID"`
	AgrID                 *int       `json:"AgrID"`
	AgrName               *string    `json:"AgrName"`
	AllocationType        *string    `json:"AllocationType"`
	AmountFact            *float64   `json:"AmountFact"`
	AmountPlan            *float64   `json:"AmountPlan"`
	BrandID               *int       `json:"BrandID"`
	BrandName             *string    `json:"BrandName"`
	CPPoffprime           *float64   `json:"CPPoffprime"`
	CPPprime              *float64   `json:"CPPprime"`
	ComplimentId          *int       `json:"CommInMplID"`
	ContractBeg           *int       `json:"ContractBeg"`
	ContractEnd           *int       `json:"ContractEnd"`
	DateFrom              *int       `json:"DateFrom"`
	DateTo                *int       `json:"DateTo"`
	DealChannelStatus     *int       `json:"DealChannelStatus"`
	Discounts             []struct{} `json:"Discounts"`
	DoubleAdvertiser      *bool      `json:"DoubleAdvertiser"`
	DtpID                 *int       `json:"DtpID"`
	DublSpot              *int       `json:"DublSpot"`
	FbrName               *string    `json:"FbrName"`
	FilmDur               *int       `json:"FilmDur"`
	FilmDurKoef           *float64   `json:"FilmDurKoef"`
	FilmID                *int       `json:"FilmID"`
	FilmName              *string    `json:"FilmName"`
	FilmVersion           *string    `json:"FilmVersion"`
	FixPriceAsFloat       *int       `json:"FixPriceAsFloat"`
	FixPriority           *int       `json:"FixPriority"`
	GRP                   *float64   `json:"GRP"`
	GRPShift              *float64   `json:"GRPShift"`
	GrpFact               *float64   `json:"GrpFact"`
	GrpPlan               *float64   `json:"GrpPlan"`
	GrpTotal              *float64   `json:"GrpTotal"`
	GrpTotalPrime         *float64   `json:"GrpTotalPrime"`
	HasReserve            *int       `json:"HasReserve"`
	InventoryUnitDuration *int       `json:"InventoryUnitDuration"`
	MplCbrID              *int       `json:"MplCbrID"`
	MplCbrName            *string    `json:"MplCbrName"`
	MplCnlID              *int       `json:"MplCnlID"`
	MplID                 *string    `json:"MplID"`
	MplMonth              *string    `json:"MplMonth"`
	MplName               *string    `json:"MplName"`
	MplState              *int       `json:"MplState"`
	Multiple              *int       `json:"Multiple"`
	OBDPos                *int       `json:"OBDPos"`
	OrdFrID               *int       `json:"OrdFrID"`
	OrdID                 *int       `json:"OrdID"`
	OrdIsTriggered        *int       `json:"OrdIsTriggered"`
	OrdName               *string    `json:"OrdName"`
	PBACond               *string    `json:"PBACond"`
	PBAObjID              *int       `json:"PBAObjID"`
	ProdClassID           *int       `json:"ProdClassID"`
	SellingDirection      *int       `json:"SellingDirection"`
	SplitMessageGroupID   *int       `json:"SplitMessageGroupID"`
	SumShift              *float64   `json:"SumShift"`
	TPName                *string    `json:"TPName"`
	TgrID                 *string    `json:"TgrID"`
	TgrName               *string    `json:"TgrName"`
	Timestamp             *time.Time `json:"Timestamp"`
	AdvertiserId          *int       `json:"advertiserId"`
	AgreementId           *int       `json:"agreementId"`
	ChannelId             *int       `json:"channelId"`
	FfoaAllocated         *int       `json:"ffoaAllocated"`
	FfoaLawAcc            *int       `json:"ffoaLawAcc"`
	FilmId                *int       `json:"filmId"`
	MediaplanId           *int       `json:"mediaplanId"`
	Month                 *int       `json:"month"`
	OrdBegDate            *int       `json:"ordBegDate"`
	OrdEndDate            *int       `json:"ordEndDate"`
	OrdManager            *string    `json:"ordManager"`
}

func (mediaplan *mplTest) Key() string {
	return fmt.Sprintf("%v-%v", *mediaplan.MediaplanId, *mediaplan.ComplimentId)
}

func BenchmarkBadger_AddValues(b *testing.B) {
	if testing.Short() {
		return
	}
	badgerBenchmark := NewBadger("db/benchmark/storage/add-stream")
	b.ResetTimer()
	values := map[string]interface{}{}
	for i := 0; i < 100000; i++ {
		mplTest := mplTest{
			AdtID:                 utils.IntI(i + 1),
			AdtName:               nil,
			AdvID:                 utils.IntI(i + 2),
			AgrID:                 utils.IntI(i + 3),
			AgrName:               nil,
			AllocationType:        nil,
			AmountFact:            nil,
			AmountPlan:            nil,
			BrandID:               nil,
			BrandName:             nil,
			CPPoffprime:           nil,
			CPPprime:              nil,
			ComplimentId:          utils.IntI(i),
			ContractBeg:           nil,
			ContractEnd:           nil,
			DateFrom:              nil,
			DateTo:                nil,
			DealChannelStatus:     nil,
			Discounts:             nil,
			DoubleAdvertiser:      nil,
			DtpID:                 nil,
			DublSpot:              nil,
			FbrName:               nil,
			FilmDur:               nil,
			FilmDurKoef:           nil,
			FilmID:                nil,
			FilmName:              nil,
			FilmVersion:           nil,
			FixPriceAsFloat:       nil,
			FixPriority:           nil,
			GRP:                   nil,
			GRPShift:              nil,
			GrpFact:               nil,
			GrpPlan:               nil,
			GrpTotal:              nil,
			GrpTotalPrime:         nil,
			HasReserve:            nil,
			InventoryUnitDuration: nil,
			MplCbrID:              nil,
			MplCbrName:            nil,
			MplCnlID:              utils.IntI(i + 4),
			MplID:                 nil,
			MplMonth:              utils.StringI(201903),
			MplName:               nil,
			MplState:              nil,
			Multiple:              nil,
			OBDPos:                nil,
			OrdFrID:               nil,
			OrdID:                 nil,
			OrdIsTriggered:        nil,
			OrdName:               nil,
			PBACond:               nil,
			PBAObjID:              nil,
			ProdClassID:           nil,
			SellingDirection:      utils.IntI(23),
			SplitMessageGroupID:   nil,
			SumShift:              nil,
			TPName:                nil,
			TgrID:                 nil,
			TgrName:               nil,
			Timestamp:             nil,
			AdvertiserId:          utils.IntI(i + 1),
			AgreementId:           utils.IntI(i + 5),
			ChannelId:             utils.IntI(i + 4),
			FfoaAllocated:         nil,
			FfoaLawAcc:            nil,
			FilmId:                nil,
			MediaplanId:           utils.IntI(i),
			Month:                 utils.IntI(201903),
			OrdBegDate:            nil,
			OrdEndDate:            nil,
			OrdManager:            nil,
		}
		values[mplTest.Key()] = mplTest
	}
	badgerBenchmark.AddValues(values)
}

func TestBadger_FilterValue(t *testing.T) {
	if testing.Short() {
		return
	}
	badgerBenchmark := NewBadger("db/benchmark/storage/add-stream")
	var result []mplTest
	search := mplTest{}
	search.MediaplanId = utils.IntI(500)
	badgerBenchmark.Iterate(func(key []byte, value []byte) {
		var val mplTest
		json.Unmarshal(value, &val)
		if *search.MediaplanId == *val.MediaplanId {
			result = append(result, val)
		}
	})
	println(result)
}

func TestBadger_Count(t *testing.T) {
	if testing.Short() {
		return
	}
	badgerBenchmark := NewBadger("db/benchmark/storage/add-stream")
	println(badgerBenchmark.Count())
}

func TestQuery(t *testing.T) {
	var result []mplTest
	store := Open("db/benchmark/storage/add-stream")
	err := store.Find(&result, badgerhold.Where("AdtID").Le(2))
	if err != nil {
		panic(err)
	}
	println(result)
}
