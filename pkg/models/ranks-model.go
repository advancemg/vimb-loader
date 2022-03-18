package models

import goConvert "github.com/advancemg/go-convert"

type SwaggerGetRanksRequest struct {
}

type GetRanks struct {
	goConvert.UnsortedMap
}

func (request *GetRanks) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}
