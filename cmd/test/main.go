package main

import (
	"encoding/json"
	"fmt"
	vimb "github.com/advancemg/vimb-loader"
	"github.com/advancemg/vimb-loader/models"
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

	var jsStruct models.GetProgramBreaksRequest
	json.Unmarshal(js, &jsStruct)
	sorted, _ := jsStruct.Sorted()
	fmt.Println(string(sorted))
	res, err := vimb.Request(sorted)
	if err != nil {
		println(err.Error())
	}
	fmt.Println("post:\n", string(res))
}
