package models

import (
	"encoding/json"
	"github.com/advancemg/badgerhold"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	GetProgramBreaksType                = "GetProgramBreaks"
	GetProgramBreaksLightModeType       = "GetProgramBreaksLightMode"
	AddMPlanFilmType                    = "AddMPlanFilm"
	AddMPlanType                        = "AddMPlan"
	AddSpotType                         = "AddSpot"
	GetAdvMessagesType                  = "GetAdvMessages"
	GetBudgetsType                      = "GetBudgets"
	ChangeMPlanFilmPlannedInventoryType = "ChangeMPlanFilmPlannedInventory"
	ChangeSpotType                      = "ChangeSpot"
	GetCustomersWithAdvertisersType     = "GetCustomersWithAdvertisers"
	DeleteMPlanFilmType                 = "DeleteMPlanFilm"
	DeleteSpotType                      = "DeleteSpot"
	GetDeletedSpotInfoType              = "GetDeletedSpotInfo"
	GetMPLansType                       = "GetMPLans"
	GetSpotsType                        = "GetSpots"
	GetRanksType                        = "GetRanks"
	SetSpotPositionType                 = "SetSpotPosition"
	GetChannelsType                     = "GetChannels"
	ProgramBreaksUpdateQueue            = "program-breaks-update"
	ProgramBreaksLightModeUpdateQueue   = "program-breaks-light-mode-update"
	AdvMessagesUpdateQueue              = "adv-messages-update"
	BudgetsUpdateQueue                  = "budgets-update"
	CustomersWithAdvertisersUpdateQueue = "customers-with-advertisers-update"
	DeletedSpotInfoUpdateQueue          = "deleted-spot-info-update"
	MPLansUpdateQueue                   = "mediaplans-update"
	MediaplanAggUpdateQueue             = "mediaplans-agg-update"
	ProgramsUpdateQueue                 = "programs-update"
	SpotsUpdateQueue                    = "spots-update"
	RanksUpdateQueue                    = "ranks-update"
	ChannelsUpdateQueue                 = "channels-update"
	DbChannels                          = "db/channels"
	DbBudgets                           = "db/budgets"
	DbRanks                             = "db/ranks"
	DbAdvertisers                       = "db/advertisers"
	DbCustomersWithAdvertisers          = "db/customers-with-advertisers"
	DbCustomersWithAdvertisersData      = "db/customers-with-advertisers-data"
	DbSpots                             = "db/spots"
	DbSpotsOrderBlock                   = "db/spots-order-block"
	DbProgramBreaksLightMode            = "db/program-breaks-light-mode"
	DbProgramBreaks                     = "db/program-breaks"
	DbPrograms                          = "db/programs"
	DbProgramBreaksProMaster            = "db/program-breaks-pro-master"
	DbProgramBreaksBlockForecast        = "db/program-breaks-block-forecast"
	DbProgramBreaksBlockForecastTgr     = "db/program-breaks-block-forecast-tgr"
	DbDeletedSpotInfo                   = "db/deleted-spot-info"
	DbMediaplans                        = "db/mediaplans"
	DbAggMediaplans                     = "db/agg-mediaplans"
	DbTimeout                           = "db/timeout"
)

var QueueNames = []string{
	GetProgramBreaksType,
	GetProgramBreaksLightModeType,
	AddMPlanFilmType,
	AddMPlanType,
	AddSpotType,
	GetAdvMessagesType,
	GetBudgetsType,
	ChangeMPlanFilmPlannedInventoryType,
	ChangeSpotType,
	GetCustomersWithAdvertisersType,
	DeleteMPlanFilmType,
	DeleteSpotType,
	GetDeletedSpotInfoType,
	GetMPLansType,
	GetSpotsType,
	GetRanksType,
	SetSpotPositionType,
	GetChannelsType,
	ProgramBreaksUpdateQueue,
	ProgramBreaksLightModeUpdateQueue,
	AdvMessagesUpdateQueue,
	BudgetsUpdateQueue,
	CustomersWithAdvertisersUpdateQueue,
	DeletedSpotInfoUpdateQueue,
	MPLansUpdateQueue,
	MediaplanAggUpdateQueue,
	SpotsUpdateQueue,
	RanksUpdateQueue,
	ChannelsUpdateQueue,
}

type Configuration struct {
	Mediaplan                MediaplanConfiguration                `json:"mediaplan"`
	Budget                   BudgetConfiguration                   `json:"budget"`
	Channel                  ChannelConfiguration                  `json:"channel"`
	AdvMessages              AdvMessagesConfiguration              `json:"advMessages"`
	CustomersWithAdvertisers CustomersWithAdvertisersConfiguration `json:"customersWithAdvertisers"`
	DeletedSpotInfo          DeletedSpotInfoConfiguration          `json:"deletedSpotInfo"`
	Rank                     RanksConfiguration                    `json:"rank"`
	ProgramBreaks            ProgramBreaksConfiguration            `json:"programBreaks"`
	ProgramBreaksLight       ProgramBreaksLightConfiguration       `json:"programBreaksLight"`
	Spots                    SpotsConfiguration                    `json:"spots"`
}

type internalM struct {
	M map[string]interface{} `json:"m"`
}
type internalI struct {
	I map[string]interface{} `json:"i"`
}
type internalItem struct {
	Item map[string]interface{} `json:"item"`
}
type internalB struct {
	B map[string]interface{} `json:"b"`
}
type internalBB struct {
	B map[string]interface{} `json:"B"`
}
type internalTgr struct {
	Tgr map[string]interface{} `json:"tgr"`
}
type internalAttributes struct {
	Attributes map[string]interface{} `json:"attributes"`
}
type internalP struct {
	P map[string]interface{} `json:"p"`
}
type internalRow struct {
	Row map[string]interface{} `json:"Row"`
}
type internalS struct {
	S map[string]interface{} `json:"s"`
}
type internalObl struct {
	Obl map[string]interface{} `json:"obl"`
}
type internalChannel struct {
	Channel map[string]interface{} `json:"Channel"`
}

type Any struct {
	Body map[string]interface{}
}

type WeekInfo struct {
	Number int       `json:"Number"`
	Close  bool      `json:"Close"`
	Date   time.Time `json:"Date"`
}

type MqUpdateMessage struct {
	Bucket             string `json:"bucket"`
	Key                string `json:"key"`
	Month              string `json:"month"`
	SellingDirectionID string `json:"sellingDirectionID"`
}

type CommonResponse map[string]interface{}

type StreamResponse struct {
	Body    []byte `json:"body"`
	Request string `json:"request"`
}

type JsonResponse struct {
	Body    interface{} `json:"body"`
	Request string      `json:"request"`
}

var (
	Config = Configuration{}
)

func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	configFile, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config, nil
}

func HandleBadgerRequest(request map[string]interface{}) *badgerhold.Query {
	var query *badgerhold.Query
	once := true
	for field, value := range request {
		for key, val := range value.(map[string]interface{}) {
			if once {
				query = switchBadgerFilterWhere(query, key, field, val)
				once = false
			} else {
				query = switchBadgerFilterAnd(query, key, field, val)
			}
		}
	}
	return query
}

func switchBadgerFilterAnd(filter *badgerhold.Query, key, filed string, value interface{}) *badgerhold.Query {
	value = jsonNumber(value)
	switch key {
	case "eq":
		filter = filter.And(filed).Eq(value)
	case "ne":
		filter = filter.And(filed).Ne(value)
	case "gt":
		filter = filter.And(filed).Gt(value)
	case "lt":
		filter = filter.And(filed).Lt(value)
	case "ge":
		filter = filter.And(filed).Ge(value)
	case "le":
		filter = filter.And(filed).Le(value)
	case "in":
		filter = filter.And(filed).In(value)
	case "isnil":
		filter = filter.And(filed).IsNil()
	}
	return filter
}

func switchBadgerFilterWhere(filter *badgerhold.Query, key, filed string, value interface{}) *badgerhold.Query {
	value = jsonNumber(value)
	switch key {
	case "eq":
		filter = badgerhold.Where(filed).Eq(value)
	case "ne":
		filter = badgerhold.Where(filed).Ne(value)
	case "gt":
		filter = badgerhold.Where(filed).Gt(value)
	case "lt":
		filter = badgerhold.Where(filed).Lt(value)
	case "ge":
		filter = badgerhold.Where(filed).Ge(value)
	case "le":
		filter = badgerhold.Where(filed).Le(value)
	case "in":
		filter = badgerhold.Where(filed).In(value)
	case "isnil":
		filter = badgerhold.Where(filed).IsNil()
	}
	return filter
}

func jsonNumber(value interface{}) interface{} {
	if number, ok := value.(json.Number); ok {
		strconv.ParseInt(string(number), 10, 64)
		dot := strings.Contains(number.String(), ".")
		if dot {
			i, err := number.Float64()
			if err != nil {
				panic(err)
			}
			return i
		} else {
			i, err := number.Int64()
			if err != nil {
				panic(err)
			}
			return i
		}
	}
	return value
}
