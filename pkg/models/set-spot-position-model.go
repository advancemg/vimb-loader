package models

import (
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerSetSpotPositionRequest struct {
	SpotID   string `json:"SpotID"`
	Distance string `json:"Distance"`
}

type SetSpotPosition struct {
	goConvert.UnsortedMap
}

func (request *SetSpotPosition) GetData() (*StreamResponse, error) {
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
