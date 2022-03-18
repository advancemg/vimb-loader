package models

import goConvert "github.com/advancemg/go-convert"

type SwaggerGetBudgetsRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
	StartMonth         string `json:"StartMonth"`
	EndMonth           string `json:"EndMonth"`
	AdvertiserList     []struct {
		Id string `json:"ID"`
	} `json:"AdvertiserList"`
	ChannelList []struct {
		Id string `json:"ID"`
	} `json:"ChannelList"`
}

type GetBudgets struct {
	goConvert.UnsortedMap
}

func (request *GetBudgets) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}
