package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"github.com/timshannon/badgerhold"
	"strconv"
	"testing"
	"time"
)

func TestBudgetsUpdateRequest_loadFromFile(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type fields struct {
		S3Key string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "loadFromFile-Budgets",
			fields:  fields{"../../dev-test-data/budgets.gz"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &BudgetsUpdateRequest{
				S3Key: tt.fields.S3Key,
			}
			if err := request.loadFromFile(); (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryBudgets(t *testing.T) {
	type Cnl struct {
		Cnl  int64
		Main int64
	}
	var budgets []Budget
	var channels []Channel
	var allChannels []Cnl
	advertisers := map[int64]int64{}
	channelList := map[int64]Cnl{}
	badgerBudgets := storage.Open(DbBudgets)
	_ = badgerBudgets.Find(&budgets, badgerhold.Where("Month").Ge(-1))
	badgerChannels := storage.Open(DbChannels)
	_ = badgerChannels.Find(&channels, badgerhold.Where("ID").Ge(-1))
	for _, budget := range budgets {
		advertisers[*budget.AdtID] = *budget.AdtID
		channelList[*budget.CnlID] = Cnl{
			Cnl:  *budget.CnlID,
			Main: 0,
		}
	}
	for _, channel := range channels {
		if channelItem, ok := channelList[*channel.ID]; ok {
			channelItem.Main = *channel.MainChnl
			allChannels = append(allChannels, channelItem)
		}
	}
}

func TestBudgetsUpdateRequest_readBudgets(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	var cnl []int64
	var result []Budget
	months := map[int64][]string{}
	badgerBudgets := storage.NewBadger(DbBudgets)
	badgerBudgets.Iterate(func(key []byte, value []byte) {
		var budget Budget
		err := json.Unmarshal(value, &budget)
		if err != nil {
			return
		}
		result = append(result, budget)
	})
	for _, budget := range result {
		cnl = append(cnl, *budget.CnlID)
		monthStr := fmt.Sprintf("%d", *budget.Month)
		month, err := strconv.Atoi(monthStr[4:6])
		if err != nil {
			panic(err)
		}
		year, err := strconv.Atoi(monthStr[0:4])
		if err != nil {
			panic(err)
		}
		days, err := utils.GetDaysFromMonth(year, time.Month(month))
		if err != nil {
			panic(err)
		}
		months[int64(month)] = days
	}
	fmt.Println(cnl)
	fmt.Println(months)
}

func TestBudgetsUpdateRequest_readBudgetsAndChannels(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type Cnl struct {
		Cnl  int64
		Main int64
	}
	var cnl []int64
	var advertisers []int64
	var allchannels []Cnl
	var budgets []Budget
	var channels []Channel
	channelList := map[int64]Cnl{}
	months := map[int][]string{}
	badgerBudgets := storage.NewBadger(DbBudgets)
	badgerBudgets.Iterate(func(key []byte, value []byte) {
		var budget Budget
		json.Unmarshal(value, &budget)
		budgets = append(budgets, budget)
	})
	for _, budget := range budgets {
		advertisers = append(advertisers, *budget.AdtID)
		cnl = append(cnl, *budget.CnlID)
		monthStr := fmt.Sprintf("%d", *budget.Month)
		month, err := strconv.Atoi(monthStr[4:6])
		if err != nil {
			panic(err)
		}
		year, err := strconv.Atoi(monthStr[0:4])
		if err != nil {
			panic(err)
		}
		days, err := utils.GetDaysFromMonth(year, time.Month(month))
		if err != nil {
			panic(err)
		}
		months[month] = days
	}
	channelsBadger := storage.NewBadger(DbChannels)
	if channelsBadger.Count() > 0 {
		channelsBadger.Iterate(func(key []byte, value []byte) {
			var channel Channel
			json.Unmarshal(value, &channel)
			channels = append(channels, channel)
		})
		for _, channel := range channels {
			allchannels = append(allchannels, Cnl{
				Cnl:  *channel.ID,
				Main: *channel.MainChnl,
			})
		}
		for _, channel := range cnl {
			for _, cnlIdMain := range allchannels {
				if cnlIdMain.Cnl == channel {
					channelList[channel] = cnlIdMain
				}
			}
		}
	}
	fmt.Println(channelList)
	fmt.Println(advertisers)
}
