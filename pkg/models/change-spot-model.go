package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerChangeSpotRequest struct {
	FirstSpotID  string `json:"FirstSpotID"`
	SecondSpotID string `json:"SecondSpotID"`
}

type ChangeSpot struct {
	goConvert.UnsortedMap
}

func (request *ChangeSpot) GetData() (*StreamResponse, error) {
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.RequestJson(req)
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    resp,
		Request: string(req),
	}, nil
}

func (request *ChangeSpot) GetRawData() (*StreamResponse, error) {
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.Request(req)
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    resp,
		Request: string(req),
	}, nil
}

func (request *ChangeSpot) GetDataToS3() error {
	data, err := request.GetRawData()
	if err != nil {
		return err
	}
	var newS3Key = fmt.Sprintf("%s.zip", "GetProgramBreaks")
	_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
	if err != nil {
		return err
	}
	return nil
}

func (request *ChangeSpot) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	FirstSpotID, exist := request.Get("FirstSpotID")
	if exist {
		body.Set("FirstSpotID", FirstSpotID)
	}
	SecondSpotID, exist := request.Get("SecondSpotID")
	if exist {
		body.Set("SecondSpotID", SecondSpotID)
	}
	xmlRequestHeader.Set("ChangeSpots", body)
	return xmlRequestHeader.ToXml()
}
