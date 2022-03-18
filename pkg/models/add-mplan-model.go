package models

import goConvert "github.com/advancemg/go-convert"

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

func (request *AddMPlan) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}
