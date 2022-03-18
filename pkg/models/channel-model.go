package models

import goConvert "github.com/advancemg/go-convert"

type SwaggerGetChannelsRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
}

type GetChannels struct {
	goConvert.UnsortedMap
}

func (request *GetChannels) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}
