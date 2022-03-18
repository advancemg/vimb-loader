package models

import goConvert "github.com/advancemg/go-convert"

type SwaggerDeleteSpotRequest struct {
	SpotID string `json:"SpotID"`
}

type DeleteSpot struct {
	goConvert.UnsortedMap
}

func (request *DeleteSpot) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}
