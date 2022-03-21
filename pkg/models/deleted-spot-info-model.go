package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerGetDeletedSpotInfoRequest struct {
	DateStart  string `json:"DateStart"`
	DateEnd    string `json:"DateEnd"`
	Agreements []struct {
		Id string `json:"ID"`
	} `json:"Agreements"`
}

type GetDeletedSpotInfo struct {
	goConvert.UnsortedMap
}

func (request *GetDeletedSpotInfo) GetData() (*StreamResponse, error) {
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

func (request *GetDeletedSpotInfo) GetRawData() (*StreamResponse, error) {
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

func (request *GetDeletedSpotInfo) GetDataToS3() error {
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

func (request *GetDeletedSpotInfo) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	dateStart, exist := request.Get("DateStart")
	if exist {
		body.Set("DateStart", dateStart)
	}
	dateEnd, exist := request.Get("DateEnd")
	if exist {
		body.Set("DateEnd", dateEnd)
	}
	agreements, exist := request.Get("Agreements")
	if exist {
		body.Set("Agreements", agreements)
	}
	xmlRequestHeader.Set("GetDeletedSpotInfo", body)
	return xmlRequestHeader.ToXml()
}
