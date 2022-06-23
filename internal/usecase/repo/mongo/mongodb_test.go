package mongo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

type Budget struct {
	Month                 *int64       `json:"Month" bson:"Month"`
	CnlID                 *int64       `json:"CnlID" bson:"CnlID"`
	AdtID                 *int64       `json:"AdtID" bson:"AdtID"`
	AgrID                 *int64       `json:"AgrID" bson:"AgrID"`
	InventoryUnitDuration *int64       `json:"InventoryUnitDuration" bson:"InventoryUnitDuration"`
	DealChannelStatus     *int64       `json:"DealChannelStatus" bson:"DealChannelStatus"`
	FixPercent            *int64       `json:"FixPercent" bson:"FixPercent"`
	GRPFix                *int64       `json:"GRPFix" bson:"GRPFix"`
	AdtName               *string      `json:"AdtName" bson:"AdtName"`
	AgrName               *string      `json:"AgrName" bson:"AgrName"`
	CmpName               *string      `json:"CmpName" bson:"CmpName"`
	CnlName               *string      `json:"CnlName" bson:"CnlName"`
	TP                    *string      `json:"TP" bson:"TP"`
	Budget                *float64     `json:"Budget" bson:"Budget"`
	CoordCost             *float64     `json:"CoordCost" bson:"CoordCost"`
	Cost                  *float64     `json:"Cost" bson:"Cost"`
	FixPercentPrime       *float64     `json:"FixPercentPrime" bson:"FixPercentPrime"`
	FloatPercent          *float64     `json:"FloatPercent" bson:"FloatPercent"`
	FloatPercentPrime     *float64     `json:"FloatPercentPrime" bson:"FloatPercentPrime"`
	GRP                   *float64     `json:"GRP" bson:"GRP"`
	GRPWithoutKF          *float64     `json:"GRPWithoutKF" bson:"GRPWithoutKF"`
	Timestamp             time.Time    `json:"Timestamp" bson:"Timestamp"`
	Quality               []BudgetItem `json:"Quality" bson:"Quality"`
}

type BudgetItem struct {
	RankID            *int64   `json:"RankID" bson:"RankID"`
	Percent           *float64 `json:"Percent" bson:"Percent"`
	BudgetOffprime    *float64 `json:"BudgetOffprime" bson:"BudgetOffprime"`
	BudgetPrime       *float64 `json:"BudgetPrime" bson:"BudgetPrime"`
	InventoryOffprime *float64 `json:"InventoryOffprime" bson:"InventoryOffprime"`
	InventoryPrime    *float64 `json:"InventoryPrime" bson:"InventoryPrime"`
	PercentPrime      *float64 `json:"PercentPrime" bson:"PercentPrime"`
}

func TestGet(t *testing.T) {
	var db = DbRepo{
		Client:   nil,
		table:    "budgets",
		database: "db",
	}
	var budgets Budget
	err := db.Get("Month", &budgets)
	if err != nil {
		panic(err)
	}
	os.RemoveAll("logs")
	fmt.Println(*budgets.Month, *budgets.Budget, *budgets.AdtName)
}

func TestDbRepo_AddWithTTL(t *testing.T) {
	var db = DbRepo{
		Client:   nil,
		table:    "timeout",
		database: "db",
	}
	timeout := Timeout{
		IsTimeout: true,
	}
	err := db.AddWithTTL("_id", timeout, time.Second*30)
	if err != nil {
		panic(err)
	}
	os.RemoveAll("logs")
}

func TestGetTTL(t *testing.T) {
	var db = DbRepo{
		Client:   nil,
		table:    "timeout",
		database: "db",
	}
	var timeout Timeout
	err := db.Get("_id", &timeout)
	if err != nil {
		panic(err)
	}
	os.RemoveAll("logs")
	fmt.Println(timeout)
}
