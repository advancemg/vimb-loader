package models

import goConvert "github.com/advancemg/go-convert"

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

func (request *AddSpot) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}
