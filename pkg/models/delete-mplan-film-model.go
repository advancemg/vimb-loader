package models

import (
	goConvert "github.com/advancemg/go-convert"
)

type SwaggerDeleteMPlanFilmRequest struct {
	CommInMplID string `json:"CommInMplID"`
}

type DeleteMPlanFilm struct {
	goConvert.UnsortedMap
}

func (request *DeleteMPlanFilm) GetData() (*StreamResponse, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	CommInMplID, exist := request.Get("CommInMplID")
	if exist {
		body.Set("CommInMplID", CommInMplID)
	}
	xmlRequestHeader.Set("DeleteMPlanFilm", body)
	req, err := xmlRequestHeader.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(req),
	}, nil
}
