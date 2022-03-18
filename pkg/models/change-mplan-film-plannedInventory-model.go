package models

import (
	goConvert "github.com/advancemg/go-convert"
)

type SwaggerChangeMPlanFilmPlannedInventoryRequest struct {
	CommInMpl [][]struct {
		ID        string `json:"ID"`
		Inventory string `json:"Inventory"`
	} `json:"CommInMpl"`
}

type ChangeMPlanFilmPlannedInventory struct {
	goConvert.UnsortedMap
}

func (request ChangeMPlanFilmPlannedInventory) GetData() (*StreamResponse, error) {
	xmlRequestHeader := goConvert.New()
	data := goConvert.New()
	body := goConvert.New()
	CommInMpl, exist := request.Get("CommInMpl")
	if exist {
		body.Set("CommInMpl", CommInMpl)
	}
	data.Set("Data", body)
	xmlRequestHeader.Set("ChangeMPlanFilmPlannedInventory", data)
	req, err := xmlRequestHeader.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(req),
	}, nil
}
