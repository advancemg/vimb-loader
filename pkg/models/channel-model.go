package models


import (
	goConvert "github.com/advancemg/go-convert"
)

type SwaggerGetChannelsRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
}

type GetChannels struct {
	goConvert.UnsortedMap
}


func (request GetChannels) GetData() (*StreamResponse, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	sellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", sellingDirectionID)
	}
	xmlRequestHeader.Set("GetChannels", body)
	req, err := xmlRequestHeader.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(req),
	}, nil
}
