package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerChangeFilmsRequest struct {
	ChangeFilms []struct {
		FakeSpotIDs  string `json:"FakeSpotIDs"`
		CommInMplIDs string `json:"CommInMplIDs"`
	} `json:"ChangeFilms"`
}

type ChangeFilms struct {
	goConvert.UnsortedMap
}

func (request *ChangeFilms) GetDataJson() (*JsonResponse, error) {
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.Actions.RequestJson(req)
	if err != nil {
		return nil, err
	}
	var body = map[string]interface{}{}
	err = json.Unmarshal(resp, &body)
	if err != nil {
		return nil, err
	}
	return &JsonResponse{
		Body:    body,
		Request: string(req),
	}, nil
}

func (request *ChangeFilms) GetDataXmlZip() (*StreamResponse, error) {
	//for {
	//	var isTimeout utils.Timeout
	//	err := storage.Open(DbTimeout).Get("vimb-timeout", &isTimeout)
	//	if err != nil {
	//		if errors.Is(err, badgerhold.ErrNotFound) {
	//			isTimeout.IsTimeout = false
	//		} else {
	//			return nil, err
	//		}
	//	}
	//	if isTimeout.IsTimeout {
	//		time.Sleep(1 * time.Second)
	//		continue
	//	}
	//	if !isTimeout.IsTimeout {
	//		break
	//	}
	//}
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.Actions.Request(req)
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    resp,
		Request: string(req),
	}, nil
}

func (request *ChangeFilms) UploadToS3() (*MqUpdateMessage, error) {
	for {
		typeName := GetBudgetsType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("GetBudgets")
				continue
			}
			return nil, err
		}
		month, _ := request.Get("StartMonth")
		if err != nil {
			return nil, err
		}
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%v/%s-%s.gz", utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return nil, err
		}
		return &MqUpdateMessage{
			Key: newS3Key,
		}, nil
	}
}

func (request *ChangeFilms) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	FakeSpotIDs, exist := request.Get("FakeSpotIDs")
	if exist {
		body.Set("FakeSpotIDs", FakeSpotIDs)
	}
	CommInMplIDs, exist := request.Get("CommInMplIDs")
	if exist {
		body.Set("CommInMplIDs", CommInMplIDs)
	}
	xmlRequestHeader.Set("ChangeFilms", body)
	return xmlRequestHeader.ToXml()
}
