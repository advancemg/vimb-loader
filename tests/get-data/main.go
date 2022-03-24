package main

import (
	"encoding/json"
	"fmt"
	. "github.com/advancemg/vimb-loader/pkg/models"
)

func changeModel(input []byte, model string) ([]byte, error) {
	switch model {
	case "AddMPlan":
		var js AddMPlan
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "AddSpot":
		var js AddSpot
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "GetAdvMessages":
		var js GetAdvMessages
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "GetBudgets":
		var js GetBudgets
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "ChangeMPlanFilmPlannedInventory":
		var js ChangeMPlanFilmPlannedInventory
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "ChangeSpot":
		var js ChangeSpot
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "GetChannels":
		var js GetChannels
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "GetCustomersWithAdvertisers":
		var js GetCustomersWithAdvertisers
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "DeleteMPlanFilm":
		var js DeleteMPlanFilm
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "DeleteSpot":
		var js DeleteSpot
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "GetDeletedSpotInfo":
		var js GetDeletedSpotInfo
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "GetMPLans":
		var js GetMPLans
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "GetSpots":
		var js GetSpots
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "GetProgramBreaks":
		var js GetProgramBreaks
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "GetRanks":
		var js GetRanks
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "SetSpotPosition":
		var js SetSpotPosition
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	case "AddMPlanFilm":
		var js AddMPlanFilm
		json.Unmarshal(input, &js)
		data, err := js.GetDataJson()
		if err != nil {
			return nil, err
		}
		marshal, err := json.Marshal(data.Body)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	}
	return nil, nil
}

func main() {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "AddMPlanFilm",
			input: []byte(`{"EndDate": "nil","FilmID": "751900","MplID": "14396424","StartDate": "2022-03-21"}`),
		},
		{
			name:  "GetProgramBreaks",
			input: []byte(`{"SellingDirectionID": "21","InclProgAttr": "1","InclForecast": "1","AudRatDec": "9","StartDate": "20210309","EndDate": "20210309","LightMode": "0","CnlList": {"Cnl": "1018566"},"ProtocolVersion": "2"}`),
		},
		{
			name:  "GetRanks",
			input: []byte(`{"GetRanks": ""}`),
		},
		{
			name:  "GetDeletedSpotInfo",
			input: []byte(`{"Agreements": {"ID": "1"},"DateEnd": "2019-12-11T12:00:00","DateStart": "2019-12-11T00:00:00"}`),
		},
		{
			name:  "GetAdvMessages",
			input: []byte(`{"Advertisers": {"ID": "1"},"AdvertisingMessageIDs": [{"ID": "1"},{"ID": "2"}],"Aspects": {"ID": "2"},"CreationDateEnd": "2019-03-02","CreationDateStart": "2019-03-01","FillMaterialTags": "true"}`),
		},
		{
			name:  "GetBudgets",
			input: []byte(`{"AdvertiserList": {"ID": "700064621"},"ChannelList": {"ID": "1018574"},"EndMonth": "20210309","SellingDirectionID": "21","StartMonth": "202103"}`),
		},
		{
			name:  "ChangeMPlanFilmPlannedInventory",
			input: []byte(`{"Data": [{"CommInMpl": {"ID": "3342386","Inventory": "10"}},{"CommInMpl": {"ID": "3342387","Inventory": "20"}},{"CommInMpl": {"ID": "93342389","Inventory": "10"}}]}`),
		},
		{
			name:  "ChangeSpot",
			input: []byte(`{"FirstSpotID": "1","SecondSpotID": "2"}`),
		},
		{
			name:  "GetChannels",
			input: []byte(`{"SellingDirectionID": "21"}`),
		},
		{
			name:  "GetCustomersWithAdvertisers",
			input: []byte(`{"SellingDirectionID": "23"}`),
		},
		{
			name:  "DeleteMPlanFilm",
			input: []byte(`{"CommInMplID": "9249649"}`),
		},
		{
			name:  "DeleteSpot",
			input: []byte(`{"SpotID": "123"}`),
		},
		{
			name:  "GetMPLans",
			input: []byte(`{"AdtList": {"AdtID": "1"},"ChannelList": {"Cnl": "1018578"},"EndMonth": "20210309","IncludeEmpty": "true","SellingDirectionID": "21","StartMonth": "20210309"}`),
		},
		{
			name:  "GetSpots",
			input: []byte(`{"AdtList": {"AdtID": "1"},"ChannelList": {"Cnl": "2","Main": "1"},"EndDate": "20210309","InclOrdBlocks": "1","SellingDirectionID": "21","StartDate": "20210309"}`),
		},
		{
			name:  "SetSpotPosition",
			input: []byte(`{"Distance": "2","SpotID": "637435949"}`),
		},
		{
			name:  "AddSpot",
			input: []byte(`{"AuctionBidValue": "nil","BlockID": "1","FilmID": "1238114","FixedPosition": "true","Position": "nil"}`),
		},
		{
			name:  "AddMPlan",
			input: []byte(`{"BrandID": "nil","DateFrom": "2020-06-01","DateTo": "2020-06-30","MplCnlID": "1019882","MplName": "Тест","MultiSpotsInBlock": "false","OrdID": "499499"}`),
		}}
	for _, tt := range tests {
		got, err := changeModel(tt.input, tt.name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(tt.name)
		fmt.Println(string(got))
		fmt.Println()
	}
}
