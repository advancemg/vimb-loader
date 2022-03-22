package models

const (
	GetProgramBreaksType                = "GetProgramBreaks"
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
)

type StreamResponse struct {
	Body    []byte `json:"body"`
	Request string `json:"request"`
}
