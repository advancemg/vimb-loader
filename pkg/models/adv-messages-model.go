package models

import (
	"encoding/json"
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

func (request *GetAdvMessages) GetDataJson() (*JsonResponse, error) {
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

func (request *GetAdvMessages) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *GetAdvMessages) UploadToS3() error {
	typeName := GetAdvMessagesType
	data, err := request.GetDataXmlZip()
	if err != nil {
		return err
	}
	var newS3Key = fmt.Sprintf("vimb/%s/%s/%s-%s.gz", utils.Actions.Client, typeName, utils.DateTimeNowInt(), typeName)
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
