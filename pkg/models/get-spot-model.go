package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerGetSpotsRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
	StartDate          string `json:"StartDate"`
	EndDate            string `json:"EndDate"`
	InclOrdBlocks      string `json:"InclOrdBlocks"`
	ChannelList        []struct {
		Cnl  string `json:"Cnl"`
		Main string `json:"Main"`
	} `json:"ChannelList"`
	AdtList []struct {
		AdtID string `json:"AdtID"`
	} `json:"AdtList"`
}

type GetSpots struct {
	goConvert.UnsortedMap
}

func (request *GetSpots) GetData() (*StreamResponse, error) {
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

func (request *GetSpots) GetRawData() (*StreamResponse, error) {
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

func (request *GetSpots) GetDataToS3() error {
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

func (request *GetSpots) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	SellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", SellingDirectionID)
	}
	StartDate, exist := request.Get("StartDate")
	if exist {
		body.Set("StartDate", StartDate)
	}
	EndDate, exist := request.Get("EndDate")
	if exist {
		body.Set("EndDate", EndDate)
	}
	InclOrdBlocks, exist := request.Get("InclOrdBlocks")
	if exist {
		body.Set("InclOrdBlocks", InclOrdBlocks)
	}
	ChannelList, exist := request.Get("ChannelList")
	if exist {
		body.Set("ChannelList", ChannelList)
	}
	AdtList, exist := request.Get("AdtList")
	if exist {
		body.Set("AdtList", AdtList)
	}
	xmlRequestHeader.Set("GetSpots", body)
	return xmlRequestHeader.ToXml()
}
