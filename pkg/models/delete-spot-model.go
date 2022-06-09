package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/advancemg/badgerhold"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage/badger"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type SwaggerDeleteSpotRequest struct {
	SpotID string `json:"SpotID"`
}

type DeleteSpot struct {
	goConvert.UnsortedMap
}

func (request *DeleteSpot) GetDataJson() (*JsonResponse, error) {
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

func (request *DeleteSpot) GetDataXmlZip() (*StreamResponse, error) {
	for {
		var isTimeout utils.Timeout
		err := badger.Open(DbTimeout).Get("vimb-timeout", &isTimeout)
		if err != nil {
			if errors.Is(err, badgerhold.ErrNotFound) {
				isTimeout.IsTimeout = false
			} else {
				return nil, err
			}
		}
		if isTimeout.IsTimeout {
			time.Sleep(1 * time.Second)
			continue
		}
		if !isTimeout.IsTimeout {
			break
		}
	}
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

func (request *DeleteSpot) UploadToS3() error {
	for {
		typeName := DeleteSpotType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("DeleteSpot")
				continue
			}
			return err
		}
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%s-%s.gz", utils.Actions.Client, typeName, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
	}
}

func (request *DeleteSpot) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	SpotID, exist := request.Get("SpotID")
	if exist {
		body.Set("SpotID", SpotID)
	}
	xmlRequestHeader.Set("DeleteSpot", body)
	return xmlRequestHeader.ToXml()
}
