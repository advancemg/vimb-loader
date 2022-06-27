package models

import (
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
