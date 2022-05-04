package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"testing"
)

var js = `{
   "ChangeFilms": [
     {
       "FakeSpotIDs": "1",
       "CommInMplIDs": "1"
     },
     {
       "FakeSpotIDs": "2",
       "CommInMplIDs": "2"
     },
     {
       "FakeSpotIDs": "3",
       "CommInMplIDs": "3"
     },
     {
       "FakeSpotIDs": "4",
       "CommInMplIDs": "4"
     },
     {
       "FakeSpotIDs": "5",
       "CommInMplIDs": "5"
     },
     {
       "FakeSpotIDs": "6",
       "CommInMplIDs": "6"
     },
     {
       "FakeSpotIDs": "7",
       "CommInMplIDs": "7"
     },
     {
       "FakeSpotIDs": "8",
       "CommInMplIDs": "8"
     }
   ]
 }`

func TestChangeFilms_GetDataXmlZip(t *testing.T) {
	request := goConvert.New()
	request.UnmarshalJSON([]byte(js))
	//x := ChangeFilms{}
	xml, err := request.ToXml()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(xml))
	//x.GetDataXmlZip()
}
