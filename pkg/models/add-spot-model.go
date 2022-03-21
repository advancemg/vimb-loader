package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerAddSpotRequest struct {
	BlockID         string `json:"BlockID"`
	FilmID          string `json:"FilmID"`
	Position        string `json:"Position"`
	FixedPosition   string `json:"FixedPosition"`
	AuctionBidValue string `json:"AuctionBidValue"`
}

type AddSpot struct {
	goConvert.UnsortedMap
}

func (request *AddSpot) GetData() (*StreamResponse, error) {
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

func (request *AddSpot) GetRawData() (*StreamResponse, error) {
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

func (request *AddSpot) GetDataToS3() error {
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

func (request *AddSpot) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	BlockID, exist := request.Get("BlockID")
	if exist {
		body.Set("BlockID", BlockID)
	}
	FilmID, exist := request.Get("FilmID")
	if exist {
		body.Set("FilmID", FilmID)
	}
	Position, exist := request.Get("Position")
	if exist {
		body.Set("Position", Position)
	}
	FixedPosition, exist := request.Get("FixedPosition")
	if exist {
		body.Set("FixedPosition", FixedPosition)
	}
	AuctionBidValue, exist := request.Get("AuctionBidValue")
	if exist {
		body.Set("AuctionBidValue", AuctionBidValue)
	}
	xmlRequestHeader.Set("AddSpot", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}
