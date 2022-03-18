package main

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/models"
	vimb "github.com/advancemg/vimb-loader/pkg/utils"
)

func main() {
	js := []byte(`{
    	"SellingDirectionID": "21",
    	"InclProgAttr": "1",
    	"InclForecast": "1",
    	"AudRatDec": "9",
    	"StartDate": "20210309",
    	"EndDate": "20210309",
    	"LightMode": "0",
    	"CnlList": {"Cnl": "1018566"},
    	"ProtocolVersion": "2"
	}`)

	var jsStruct models.GetProgramBreaks
	json.Unmarshal(js, &jsStruct)
	sorted, _ := jsStruct.GetData()
	fmt.Println(sorted.Request)
	res, err := vimb.Request([]byte(sorted.Request))
	if err != nil {
		println(err.Error())
	}
	fmt.Println("post:\n", string(res))
}
