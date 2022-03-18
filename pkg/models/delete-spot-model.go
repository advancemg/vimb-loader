package models

import (
	goConvert "github.com/advancemg/go-convert"
)

type SwaggerDeleteSpotRequest struct {
	SpotID string `json:"SpotID"`
}

type DeleteSpot struct {
	goConvert.UnsortedMap
}

func (request *DeleteSpot) GetData() (*StreamResponse, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	SpotID, exist := request.Get("SpotID")
	if exist {
		body.Set("SpotID", SpotID)
	}
	xmlRequestHeader.Set("DeleteSpot", body)
	req, err := xmlRequestHeader.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(req),
	}, nil
}
