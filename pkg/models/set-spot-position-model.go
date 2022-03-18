package models

import (
	goConvert "github.com/advancemg/go-convert"
)

type SwaggerSetSpotPositionRequest struct {
	SpotID   string `json:"SpotID"`
	Distance string `json:"Distance"`
}

type SetSpotPosition struct {
	goConvert.UnsortedMap
}

func (request *SetSpotPosition) GetData() (*StreamResponse, error) {
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
	req, err := xmlRequestHeader.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(req),
	}, nil
}
