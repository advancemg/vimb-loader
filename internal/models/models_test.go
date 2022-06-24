package models

import (
	"bytes"
	"encoding/json"
	"github.com/ajankovic/xdiff"
	"github.com/ajankovic/xdiff/parser"
	"testing"
)

func changeModel(input []byte, model string) ([]byte, error) {
	switch model {
	case "AddMPlan":
		var js AddMPlan
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "AddSpot":
		var js AddSpot
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "GetAdvMessages":
		var js GetAdvMessages
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "GetBudgets":
		var js GetBudgets
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "ChangeMPlanFilmPlannedInventory":
		var js ChangeMPlanFilmPlannedInventory
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "ChangeSpot":
		var js ChangeSpot
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "GetChannels":
		var js GetChannels
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "GetCustomersWithAdvertisers":
		var js GetCustomersWithAdvertisers
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "DeleteMPlanFilm":
		var js DeleteMPlanFilm
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "DeleteSpot":
		var js DeleteSpot
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "GetDeletedSpotInfo":
		var js GetDeletedSpotInfo
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "GetMPLans":
		var js GetMPLans
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "GetSpots":
		var js GetSpots
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "GetProgramBreaks":
		var js GetProgramBreaks
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "GetRanks":
		var js GetRanks
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "SetSpotPosition":
		var js SetSpotPosition
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	case "AddMPlanFilm":
		var js AddMPlanFilm
		json.Unmarshal(input, &js)
		data, err := js.getXml()
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

func TestAddMPlan_GetData(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "AddMPlan",
			input:   []byte(`{"BrandID": "nil","DateFrom": "2020-06-01","DateTo": "2020-06-30","MplCnlID": "1019882","MplName": "Тест","MultiSpotsInBlock": "false","OrdID": "499499"}`),
			want:    []byte(`<AddMPlan xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><OrdID>499499</OrdID><MplCnlID>1019882</MplCnlID><DateFrom>2020-06-01</DateFrom><DateTo>2020-06-30</DateTo><MplName>Тест</MplName><BrandID xsi:nil="true"/><MultiSpotsInBlock>false</MultiSpotsInBlock></AddMPlan>`),
			wantErr: false,
		},
		{
			name:    "AddSpot",
			input:   []byte(`{"AuctionBidValue": "nil","BlockID": "1","FilmID": "1238114","FixedPosition": "true","Position": "nil"}`),
			want:    []byte(`<AddSpot xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><BlockID>1</BlockID><FilmID>1238114</FilmID><Position xsi:nil="true"></Position><FixedPosition>true</FixedPosition><AuctionBidValue xsi:nil="true"></AuctionBidValue></AddSpot>`),
			wantErr: false,
		},
		{
			name:    "GetAdvMessages",
			input:   []byte(`{"Advertisers": {"ID": "1"},"AdvertisingMessageIDs": [{"ID": "1"},{"ID": "2"}],"Aspects": {"ID": "2"},"CreationDateEnd": "2019-03-02","CreationDateStart": "2019-03-01","FillMaterialTags": "true"}`),
			want:    []byte(`<GetAdvMessages xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><CreationDateStart>2019-03-01</CreationDateStart><CreationDateEnd>2019-03-02</CreationDateEnd><Advertisers><ID>1</ID></Advertisers><Aspects><ID>2</ID></Aspects><AdvertisingMessageIDs><ID>1</ID><ID>2</ID></AdvertisingMessageIDs><FillMaterialTags>true</FillMaterialTags></GetAdvMessages>`),
			wantErr: false,
		},
		{
			name:    "GetBudgets",
			input:   []byte(`{"AdvertiserList": {"ID": "700064621"},"ChannelList": {"ID": "1018574"},"EndMonth": "20180115","SellingDirectionID": "21","StartMonth": "201801"}`),
			want:    []byte(`<GetBudgets><SellingDirectionID>21</SellingDirectionID><StartMonth>201801</StartMonth><EndMonth>20180115</EndMonth><AdvertiserList><ID>700064621</ID></AdvertiserList><ChannelList><ID>1018574</ID></ChannelList></GetBudgets>`),
			wantErr: false,
		},
		{
			name:    "ChangeMPlanFilmPlannedInventory",
			input:   []byte(`{"Data": [{"CommInMpl": {"ID": "3342386","Inventory": "10"}},{"CommInMpl": {"ID": "3342386","Inventory": "20"}},{"CommInMpl": {"ID": "93342386","Inventory": "10"}}]}`),
			want:    []byte(`<ChangeMPlanFilmPlannedInventory><Data><CommInMpl><ID>3342386</ID><Inventory>10</Inventory></CommInMpl><CommInMpl><ID>3342386</ID><Inventory>20</Inventory></CommInMpl><CommInMpl><ID>93342386</ID><Inventory>10</Inventory></CommInMpl></Data></ChangeMPlanFilmPlannedInventory>`),
			wantErr: false,
		},
		{
			name:    "ChangeSpot",
			input:   []byte(`{"FirstSpotID": "1","SecondSpotID": "2"}`),
			want:    []byte(`<ChangeSpots><FirstSpotID>1</FirstSpotID><SecondSpotID>2</SecondSpotID></ChangeSpots>`),
			wantErr: false,
		},
		{
			name:    "GetChannels",
			input:   []byte(`{"SellingDirectionID": "21"}`),
			want:    []byte(`<GetChannels><SellingDirectionID>21</SellingDirectionID></GetChannels>`),
			wantErr: false,
		},
		{
			name:    "GetCustomersWithAdvertisers",
			input:   []byte(`{"SellingDirectionID": "23"}`),
			want:    []byte(`<GetCustomersWithAdvertisers><SellingDirectionID>23</SellingDirectionID></GetCustomersWithAdvertisers>`),
			wantErr: false,
		},
		{
			name:    "DeleteMPlanFilm",
			input:   []byte(`{"CommInMplID": "9249649"}`),
			want:    []byte(`<DeleteMPlanFilm><CommInMplID>9249649</CommInMplID></DeleteMPlanFilm>`),
			wantErr: false,
		},
		{
			name:    "DeleteSpot",
			input:   []byte(`{"SpotID": "123"}`),
			want:    []byte(`<DeleteSpot><SpotID>123</SpotID></DeleteSpot>`),
			wantErr: false,
		},
		{
			name:    "GetDeletedSpotInfo",
			input:   []byte(`{"Agreements": {"ID": "1"},"DateEnd": "2019-12-11T12:00:00","DateStart": "2019-12-11T00:00:00"}`),
			want:    []byte(`<GetDeletedSpotInfo><DateStart>2019-12-11T00:00:00</DateStart><DateEnd>2019-12-11T12:00:00</DateEnd><Agreements><ID>1</ID></Agreements></GetDeletedSpotInfo>`),
			wantErr: false,
		},
		{
			name:    "GetMPLans",
			input:   []byte(`{"AdtList": {"AdtID": "1"},"ChannelList": {"Cnl": "1018578"},"EndMonth": "201707","IncludeEmpty": "true","SellingDirectionID": "21","StartMonth": "201707"}`),
			want:    []byte(`<GetMPLansxmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><SellingDirectionID>21</SellingDirectionID><StartMonth>201707</StartMonth><EndMonth>201707</EndMonth><AdtList><AdtID>1</AdtID></AdtList><ChannelList><Cnl>1018578</Cnl></ChannelList><IncludeEmpty>true</IncludeEmpty></GetMPLans>`),
			wantErr: false,
		},
		{
			name:    "GetSpots",
			input:   []byte(`{"AdtList": {"AdtID": "1"},"ChannelList": {"Cnl": "2","Main": "1"},"EndDate": "20161016","InclOrdBlocks": "1","SellingDirectionID": "21","StartDate": "20161010"}`),
			want:    []byte(`<GetSpots><SellingDirectionID>21</SellingDirectionID><StartDate>20161010</StartDate><EndDate>20161016</EndDate><InclOrdBlocks>1</InclOrdBlocks><ChannelList><Cnl>2</Cnl><Main>1</Main></ChannelList><AdtList><AdtID>1</AdtID></AdtList></GetSpots>`),
			wantErr: false,
		},
		{
			name:    "GetProgramBreaks",
			input:   []byte(`{"SellingDirectionID": "21","InclProgAttr": "1","InclForecast": "1","AudRatDec": "9","StartDate": "20210309","EndDate": "20210309","LightMode": "0","CnlList": {"Cnl": "1018566"},"ProtocolVersion": "2"}`),
			want:    []byte(`<GetProgramBreaks xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><SellingDirectionID>21</SellingDirectionID><InclProgAttr>1</InclProgAttr><InclForecast>1</InclForecast><AudRatDec>9</AudRatDec><StartDate>20210309</StartDate><EndDate>20210309</EndDate><LightMode>0</LightMode><CnlList><Cnl>1018566</Cnl></CnlList><ProtocolVersion>2</ProtocolVersion></GetProgramBreaks>`),
			wantErr: false,
		},
		{
			name:    "GetRanks",
			input:   []byte(`{"GetRanks": ""}`),
			want:    []byte(`<GetRanks></GetRanks>`),
			wantErr: false,
		},
		{
			name:    "SetSpotPosition",
			input:   []byte(`{"Distance": "2","SpotID": "637435949"}`),
			want:    []byte(`<SetSpotPosition><SpotID>637435949</SpotID><Distance>2</Distance></SetSpotPosition>`),
			wantErr: false,
		},
		{
			name:    "AddMPlanFilm",
			input:   []byte(`{"EndDate": "nil","FilmID": "751900","MplID": "14396424","StartDate": "2018-12-24"}`),
			want:    []byte(`<AddMPlanFilm xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><MplID>14396424</MplID><FilmID>751900</FilmID><StartDate>2018-12-24</StartDate><EndDate xsi:nil="true"/></AddMPlanFilm>`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := changeModel(tt.input, tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff, err := xmlDiff(got, tt.want)
			if (err != nil) != !diff {
				t.Errorf("GetData() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func xmlDiff(a, b []byte) (bool, error) {
	p := parser.New()
	left, err := p.ParseBytes(a)
	if err != nil {
		return false, err
	}
	right, err := p.ParseBytes(b)
	if err != nil {
		return false, err
	}
	diff, err := xdiff.Compare(left, right)
	if err != nil {
		return false, err
	}
	buf := new(bytes.Buffer)
	enc := xdiff.NewTextEncoder(buf)
	if err := enc.Encode(diff); err != nil {
		return false, err
	}
	if buf.String() == "No difference.\n" {
		return true, nil
	}
	return false, nil
}
