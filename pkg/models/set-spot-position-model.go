package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerSetSpotPositionRequest struct {
	SpotID   string `json:"SpotID"`
	Distance string `json:"Distance"`
}

type SetSpotPosition struct {
	goConvert.UnsortedMap
}

func (request *SetSpotPosition) GetDataJson() (*JsonResponse, error) {
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

func (request *SetSpotPosition) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *SetSpotPosition) UploadToS3() error {
	typeName := SetSpotPositionType
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

func (request *SetSpotPosition) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	SpotID, exist := request.Get("SpotID")
	if exist {
		body.Set("SpotID", SpotID)
	}
	Distance, exist := request.Get("Distance")
	if exist {
		body.Set("Distance", Distance)
	}
	xmlRequestHeader.Set("SetSpotPosition", body)
	return xmlRequestHeader.ToXml()
}
