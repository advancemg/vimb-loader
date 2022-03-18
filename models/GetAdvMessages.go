package models

import (
	convert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/models/types"
)

type GetAdvMessagesRequest struct {
	convert.UnsortedMap
}

type GetAdvMessages struct {
	CreationDateStart     string         `json:"creationDateStart"`
	CreationDateEnd       string         `json:"creationDateEnd"`
	Advertisers           []types.ItemId `json:"advertisers"`
	Aspects               []types.ItemId `json:"aspects"`
	AdvertisingMessageIDs []types.ItemId `json:"advertisingMessageIDs"`
	FillMaterialTags      string         `json:"fillMaterialTags"`
}

func (request GetAdvMessagesRequest) Sorted() ([]byte, error) {
	attributes := convert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := convert.New()
	body := convert.New()
	creationDateStart, exist := request.Get("CreationDateStart")
	if exist {
		body.Set("CreationDateStart", creationDateStart)
	}
	creationDateEnd, exist := request.Get("CreationDateEnd")
	if exist {
		body.Set("CreationDateEnd", creationDateEnd)
	}
	advertisers, exist := request.Get("Advertisers")
	if exist {
		body.Set("Advertisers", advertisers)
	}
	aspects, exist := request.Get("Aspects")
	if exist {
		body.Set("Aspects", aspects)
	}
	advertisingMessageIDs, exist := request.Get("AdvertisingMessageIDs")
	if exist {
		body.Set("AdvertisingMessageIDs", advertisingMessageIDs)
	}
	fillMaterialTags, exist := request.Get("FillMaterialTags")
	if exist {
		body.Set("FillMaterialTags", fillMaterialTags)
	}
	xmlRequestHeader.Set("GetAdvMessages", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}
