package models

import (
	"github.com/advancemg/badgerhold"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"testing"
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
			fields:  fields{"../../dev-test-data/budgets201902.gz"},
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
	if testing.Short() {
		t.SkipNow()
	}
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
	_ = badgerBudgets.Find(&budgets, badgerhold.Where("Month").Ge(int64(-1)))
	badgerChannels := storage.Open(DbChannels)
	_ = badgerChannels.Find(&channels, badgerhold.Where("ID").Ge(int64(-1)))
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
