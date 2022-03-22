package main

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/models"
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
	fmt.Println()
	fmt.Println(string(sorted.Body))
}
