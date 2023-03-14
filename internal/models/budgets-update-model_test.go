package models

import (
	"fmt"
	"github.com/advancemg/badgerhold"
	"github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/pkg/logging/zap"
	"github.com/advancemg/vimb-loader/pkg/storage/badger-client"
	"os"
	"testing"
	"time"
)

func TestLoadFromFile(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	err := config.Load()
	if err != nil {
		panic(err)
	}
	zap.Init()
	request := &BudgetsUpdateRequest{S3Key: "../../dev-test-data/budgets201902.gz"}
	start := time.Now()
	err = request.loadFromFile()
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start))
	os.RemoveAll("logs")
}

func BenchmarkLoadFromFil(b *testing.B) {
	if testing.Short() {
		b.SkipNow()
	}
	err := config.Load()
	if err != nil {
		panic(err)
	}
	zap.Init()
	request := &BudgetsUpdateRequest{S3Key: "../../dev-test-data/budgets201902.gz"}
	for i := 0; i < b.N; i++ {
		err = request.loadFromFile()
		if err != nil {
			panic(err)
		}
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
	badgerBudgets, _ := badger_client.Open(DbBudgets)
	_ = badgerBudgets.Find(&budgets, badgerhold.Where("Month").Ge(int64(-1)))
	badgerChannels, _ := badger_client.Open(DbChannels)
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
