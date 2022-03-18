package models

import goConvert "github.com/advancemg/go-convert"

type SwaggerChangeSpotRequest struct {
	FirstSpotID  string `json:"FirstSpotID"`
	SecondSpotID string `json:"FilmID"`
}

type ChangeSpot struct {
	goConvert.UnsortedMap
}

func (request *ChangeSpot) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}
