package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerGetMPLansRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
	StartMonth         string `json:"StartMonth"`
	EndMonth           string `json:"EndMonth"`
	AdtList            []struct {
		AdtID string `json:"AdtID"`
	} `json:"AdtList"`
	ChannelList []struct {
		Cnl string `json:"Cnl"`
	} `json:"ChannelList"`
	IncludeEmpty string `json:"IncludeEmpty"`
}

type GetMPLans struct {
	goConvert.UnsortedMap
}

func (request *GetMPLans) GetData() (*StreamResponse, error) {
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

func (request *GetMPLans) GetRawData() (*StreamResponse, error) {
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

func (request *GetMPLans) GetDataToS3() error {
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

func (request *GetMPLans) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	SellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", SellingDirectionID)
	}
	StartMonth, exist := request.Get("StartMonth")
	if exist {
		body.Set("StartMonth", StartMonth)
	}
	EndMonth, exist := request.Get("EndMonth")
	if exist {
		body.Set("EndMonth", EndMonth)
	}
	AdtList, exist := request.Get("AdtList")
	if exist {
		body.Set("AdtList", AdtList)
	}
	ChannelList, exist := request.Get("ChannelList")
	if exist {
		body.Set("ChannelList", ChannelList)
	}
	IncludeEmpty, exist := request.Get("IncludeEmpty")
	if exist {
		body.Set("IncludeEmpty", IncludeEmpty)
	}
	xmlRequestHeader.Set("GetMPLans", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}
