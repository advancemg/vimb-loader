package models


import (
	goConvert "github.com/advancemg/go-convert"
)

type SwaggerAddMPlanRequest struct {
	OrdID             string `json:"OrdID"`
	MplCnlID          string `json:"MplCnlID"`
	DateFrom          string `json:"DateFrom"`
	DateTo            string `json:"DateTo"`
	MplName           string `json:"MplName"`
	BrandID           string `json:"BrandID"`
	MultiSpotsInBlock string `json:"MultiSpotsInBlock"`
}

type AddMPlan struct {
	goConvert.UnsortedMap
}


func (request AddMPlan) GetData() (*StreamResponse, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	OrdID, exist := request.Get("OrdID")
	if exist {
		body.Set("OrdID", OrdID)
	}
	MplCnlID, exist := request.Get("MplCnlID")
	if exist {
		body.Set("MplCnlID", MplCnlID)
	}
	DateFrom, exist := request.Get("DateFrom")
	if exist {
		body.Set("DateFrom", DateFrom)
	}
	MplName, exist := request.Get("MplName")
	if exist {
		body.Set("MplName", MplName)
	}
	BrandID, exist := request.Get("BrandID")
	if exist {
		body.Set("BrandID", BrandID)
	}
	MultiSpotsInBlock, exist := request.Get("MultiSpotsInBlock")
	if exist {
		body.Set("MultiSpotsInBlock", MultiSpotsInBlock)
	}
	xmlRequestHeader.Set("AddMPlan", body)
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
