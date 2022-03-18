package models

import goConvert "github.com/advancemg/go-convert"

type SwaggerDeleteMPlanFilmRequest struct {
	CommInMplID string `json:"CommInMplID"`
}

type DeleteMPlanFilm struct {
	goConvert.UnsortedMap
}

func (request *DeleteMPlanFilm) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}
