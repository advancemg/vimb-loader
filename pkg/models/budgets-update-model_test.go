package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
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

func TestBudgetsUpdateRequest_readBudgets(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	var cnl []int
	var result []Budget
	months := map[int][]string{}
	badgerBudgets := storage.NewBadger(DbBudgets)
	badgerBudgets.Iterate(func(key []byte, value []byte) {
		var budget Budget
		json.Unmarshal(value, &budget)
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
		months[month] = days
	}
	fmt.Println(cnl)
	fmt.Println(months)
}

func TestBudgetsUpdateRequest_readBudgetsAndChannels(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type Cnl struct {
		Cnl  int
		Main int
	}
	var cnl []int
	var advertisers []int
	var allchannels []Cnl
	var budgets []Budget
	var channels []Channel
	channelList := map[int]Cnl{}
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
