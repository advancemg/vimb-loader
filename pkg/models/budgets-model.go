package models

import (
	goConvert "github.com/advancemg/go-convert"
)

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

func (request GetBudgets) GetData() (*StreamResponse, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	sellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", sellingDirectionID)
	}
	startMonth, exist := request.Get("StartMonth")
	if exist {
		body.Set("StartMonth", startMonth)
	}
	endMonth, exist := request.Get("EndMonth")
	if exist {
		body.Set("EndMonth", endMonth)
	}
	advertiserList, exist := request.Get("AdvertiserList")
	if exist {
		body.Set("AdvertiserList", advertiserList)
	}
	channelList, exist := request.Get("ChannelList")
	if exist {
		body.Set("ChannelList", channelList)
	}
	xmlRequestHeader.Set("GetBudgets", body)
	req, err := xmlRequestHeader.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(req),
	}, nil
}
