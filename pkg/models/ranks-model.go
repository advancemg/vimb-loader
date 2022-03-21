package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerGetRanksRequest struct {
}

type GetRanks struct {
	goConvert.UnsortedMap
}

func (request *GetRanks) GetData() (*StreamResponse, error) {
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

func (request *GetRanks) GetRawData() (*StreamResponse, error) {
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

func (request *GetRanks) GetDataToS3() error {
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

func (request *GetRanks) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	xmlRequestHeader.Set("GetRanks", "")
	return xmlRequestHeader.ToXml()
}
