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
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "AddSpot":
		var js AddSpot
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "GetAdvMessages":
		var js GetAdvMessages
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "GetBudgets":
		var js GetBudgets
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "ChangeMPlanFilmPlannedInventory":
		var js ChangeMPlanFilmPlannedInventory
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "ChangeSpot":
		var js ChangeSpot
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "GetChannels":
		var js GetChannels
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "GetCustomersWithAdvertisers":
		var js GetCustomersWithAdvertisers
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "DeleteMPlanFilm":
		var js DeleteMPlanFilm
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "DeleteSpot":
		var js DeleteSpot
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "GetDeletedSpotInfo":
		var js GetDeletedSpotInfo
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "GetMPLans":
		var js GetMPLans
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "GetSpots":
		var js GetSpots
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "GetProgramBreaks":
		var js GetProgramBreaks
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "GetRanks":
		var js GetRanks
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "SetSpotPosition":
		var js SetSpotPosition
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	case "AddMPlanFilm":
		var js AddMPlanFilm
		json.Unmarshal(input, &js)
		data, err := js.GetData()
		if err != nil {
			return nil, err
		}
		return data.Body, nil
	}
	return nil, nil
}

func main() {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "GetProgramBreaks",
			input:   []byte(`{"SellingDirectionID": "21","InclProgAttr": "1","InclForecast": "1","AudRatDec": "9","StartDate": "20210309","EndDate": "20210309","LightMode": "0","CnlList": {"Cnl": "1018566"},"ProtocolVersion": "2"}`),
			wantErr: false,
		},
		{
			name:    "GetRanks",
			input:   []byte(`{"GetRanks": ""}`),
			wantErr: false,
		},
		{
			name:    "AddSpot",
			input:   []byte(`{"AuctionBidValue": "nil","BlockID": "1","FilmID": "1238114","FixedPosition": "true","Position": "nil"}`),
			wantErr: false,
		},
		{
			name:    "GetAdvMessages",
			input:   []byte(`{"Advertisers": {"ID": "1"},"AdvertisingMessageIDs": [{"ID": "1"},{"ID": "2"}],"Aspects": {"ID": "2"},"CreationDateEnd": "2019-03-02","CreationDateStart": "2019-03-01","FillMaterialTags": "true"}`),
			wantErr: false,
		},
		{
			name:    "GetBudgets",
			input:   []byte(`{"AdvertiserList": {"ID": "700064621"},"ChannelList": {"ID": "1018574"},"EndMonth": "20180115","SellingDirectionID": "21","StartMonth": "201801"}`),
			wantErr: false,
		},
		{
			name:    "ChangeMPlanFilmPlannedInventory",
			input:   []byte(`{"Data": [{"CommInMpl": {"ID": "3342386","Inventory": "10"}},{"CommInMpl": {"ID": "3342386","Inventory": "20"}},{"CommInMpl": {"ID": "93342386","Inventory": "10"}}]}`),
			wantErr: false,
		},
		{
			name:    "ChangeSpot",
			input:   []byte(`{"FirstSpotID": "1","SecondSpotID": "2"}`),
			wantErr: false,
		},
		{
			name:    "GetChannels",
			input:   []byte(`{"SellingDirectionID": "21"}`),
			wantErr: false,
		},
		{
			name:    "GetCustomersWithAdvertisers",
			input:   []byte(`{"SellingDirectionID": "23"}`),
			wantErr: false,
		},
		{
			name:    "DeleteMPlanFilm",
			input:   []byte(`{"CommInMplID": "9249649"}`),
			wantErr: false,
		},
		{
			name:    "DeleteSpot",
			input:   []byte(`{"SpotID": "123"}`),
			wantErr: false,
		},
		{
			name:    "GetDeletedSpotInfo",
			input:   []byte(`{"Agreements": {"ID": "1"},"DateEnd": "2019-12-11T12:00:00","DateStart": "2019-12-11T00:00:00"}`),
			wantErr: false,
		},
		{
			name:    "GetMPLans",
			input:   []byte(`{"AdtList": {"AdtID": "1"},"ChannelList": {"Cnl": "1018578"},"EndMonth": "201707","IncludeEmpty": "true","SellingDirectionID": "21","StartMonth": "201707"}`),
			wantErr: false,
		},
		{
			name:    "GetSpots",
			input:   []byte(`{"AdtList": {"AdtID": "1"},"ChannelList": {"Cnl": "2","Main": "1"},"EndDate": "20161016","InclOrdBlocks": "1","SellingDirectionID": "21","StartDate": "20161010"}`),
			wantErr: false,
		},
		{
			name:    "SetSpotPosition",
			input:   []byte(`{"Distance": "2","SpotID": "637435949"}`),
			wantErr: false,
		},
		{
			name:    "AddMPlanFilm",
			input:   []byte(`{"EndDate": "nil","FilmID": "751900","MplID": "14396424","StartDate": "2018-12-24"}`),
			wantErr: false,
		},
		{
			name:    "AddMPlan",
			input:   []byte(`{"BrandID": "nil","DateFrom": "2020-06-01","DateTo": "2020-06-30","MplCnlID": "1019882","MplName": "Тест","MultiSpotsInBlock": "false","OrdID": "499499"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := changeModel(tt.input, tt.name)
		if (err != nil) != tt.wantErr {
			fmt.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		fmt.Println(string(got))
		fmt.Println()
	}
}
