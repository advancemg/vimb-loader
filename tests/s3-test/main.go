package main

import (
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/s3"
)

func main() {
	//var input = []byte(`{"SellingDirectionID": "21","InclProgAttr": "1","InclForecast": "1","AudRatDec": "9","StartDate": "20210309","EndDate": "20210309","LightMode": "0","CnlList": {"Cnl": "1018566"},"ProtocolVersion": "2"}`)
	//var js GetProgramBreaks
	//json.Unmarshal(input, &js)
	//err := js.UploadToS3()
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	err := s3.CreateBucket("test")
	if err != nil {
		fmt.Println(err.Error())
	}
}
