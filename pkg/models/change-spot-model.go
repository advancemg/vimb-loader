package models

import (
	goConvert "github.com/advancemg/go-convert"
)

type SwaggerChangeSpotRequest struct {
	FirstSpotID  string `json:"FirstSpotID"`
	SecondSpotID string `json:"SecondSpotID"`
}

type ChangeSpot struct {
	goConvert.UnsortedMap
}

func (request ChangeSpot) GetData() (*StreamResponse, error) {
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
	xmlRequestHeader.Set("ChangeSpot", body)
	req, err := xmlRequestHeader.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(req),
	}, nil
}
