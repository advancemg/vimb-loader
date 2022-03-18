package models

import goConvert "github.com/advancemg/go-convert"

type SwaggerChangeFilmsRequest struct {
	ChangeFilms []struct {
		FakeSpotIDs  string `json:"FakeSpotIDs"`
		CommInMplIDs string `json:"CommInMplIDs"`
	} `json:"ChangeFilms"`
}

type ChangeFilms struct {
	goConvert.UnsortedMap
}

func (request *ChangeFilms) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}
