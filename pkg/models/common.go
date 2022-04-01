package models

import (
	"encoding/json"
	"os"
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
	SpotsUpdateQueue                    = "spots-update"
	RanksUpdateQueue                    = "ranks-update"
	ChannelsUpdateQueue                 = "channels-update"
	DbCustomConfigMonth                 = "db/custom-config-month"
	DbCustomConfigAdvertisers           = "db/custom-config-advertisers "
	DbCustomConfigChannels              = "db/custom-config-channels"
	DbChannels                          = "db/channels"
	DbBudgets                           = "db/budgets"
	DbMediaplans                        = "db/mediaplans"
	DbAggMediaplans                     = "db/agg-mediaplans"
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

type WeekInfo struct {
	Number int       `json:"Number"`
	Close  bool      `json:"Close"`
	Date   time.Time `json:"Date"`
}

type MqUpdateMessage struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
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
