package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerGetAdvMessagesRequest struct {
	CreationDateStart string `json:"CreationDateStart"`
	CreationDateEnd   string `json:"CreationDateEnd"`
	Advertisers       []struct {
		Id string `json:"ID"`
	} `json:"Advertisers"`
	Aspects []struct {
		Id string `json:"ID"`
	} `json:"Aspects"`
	AdvertisingMessageIDs []struct {
		Id string `json:"ID"`
	} `json:"AdvertisingMessageIDs"`
	FillMaterialTags string `json:"FillMaterialTags"`
}

type GetAdvMessages struct {
	goConvert.UnsortedMap
}

func (request *GetAdvMessages) GetData() (*StreamResponse, error) {
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

func (request *GetAdvMessages) GetRawData() (*StreamResponse, error) {
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

func (request *GetAdvMessages) GetDataToS3() error {
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

func (request *GetAdvMessages) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	creationDateStart, exist := request.Get("CreationDateStart")
	if exist {
		body.Set("CreationDateStart", creationDateStart)
	}
	creationDateEnd, exist := request.Get("CreationDateEnd")
	if exist {
		body.Set("CreationDateEnd", creationDateEnd)
	}
	advertisers, exist := request.Get("Advertisers")
	if exist {
		body.Set("Advertisers", advertisers)
	}
	aspects, exist := request.Get("Aspects")
	if exist {
		body.Set("Aspects", aspects)
	}
	advertisingMessageIDs, exist := request.Get("AdvertisingMessageIDs")
	if exist {
		body.Set("AdvertisingMessageIDs", advertisingMessageIDs)
	}
	fillMaterialTags, exist := request.Get("FillMaterialTags")
	if exist {
		body.Set("FillMaterialTags", fillMaterialTags)
	}
	xmlRequestHeader.Set("GetAdvMessages", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}
