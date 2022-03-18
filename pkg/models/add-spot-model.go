package models


import (
	goConvert "github.com/advancemg/go-convert"
)


type SwaggerAddSpotRequest struct {
	BlockID         string `json:"BlockID"`
	FilmID          string `json:"FilmID"`
	Position        string `json:"Position"`
	FixedPosition   string `json:"FixedPosition"`
	AuctionBidValue string `json:"AuctionBidValue"`
}

type AddSpot struct {
	goConvert.UnsortedMap
}

func (request AddSpot) GetData() (*StreamResponse, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	BlockID, exist := request.Get("BlockID")
	if exist {
		body.Set("BlockID", BlockID)
	}
	FilmID, exist := request.Get("FilmID")
	if exist {
		body.Set("FilmID", FilmID)
	}
	Position, exist := request.Get("Position")
	if exist {
		body.Set("Position", Position)
	}
	FixedPosition, exist := request.Get("FixedPosition")
	if exist {
		body.Set("FixedPosition", FixedPosition)
	}
	AuctionBidValue, exist := request.Get("AuctionBidValue")
	if exist {
		body.Set("AuctionBidValue", AuctionBidValue)
	}
	xmlRequestHeader.Set("AddSpot", body)
	xmlRequestHeader.Set("attributes", attributes)
	req, err := xmlRequestHeader.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(req),
	}, nil
}
