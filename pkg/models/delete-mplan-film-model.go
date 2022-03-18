package models

import (
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerDeleteMPlanFilmRequest struct {
	CommInMplID string `json:"CommInMplID"`
}

type DeleteMPlanFilm struct {
	goConvert.UnsortedMap
}

func (request *DeleteMPlanFilm) GetData() (*StreamResponse, error) {
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

func (request *DeleteMPlanFilm) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	CommInMplID, exist := request.Get("CommInMplID")
	if exist {
		body.Set("CommInMplID", CommInMplID)
	}
	xmlRequestHeader.Set("DeleteMPlanFilm", body)
	return xmlRequestHeader.ToXml()
}
